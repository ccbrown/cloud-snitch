package app

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	jsoniter "github.com/json-iterator/go"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/report"
)

func (s *Session) GetReportsByTeamId(ctx context.Context, teamId model.Id) ([]*model.Report, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}
	reports, err := s.app.store.GetReportsByTeamId(ctx, teamId)
	return reports, s.SanitizedError(err)
}

type QueueReportGenerationInput struct {
	StartTime time.Time
	Duration  time.Duration
}

func (a *App) QueueReportGeneration(ctx context.Context, input QueueReportGenerationInput) error {
	teams, err := a.store.GetTeams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get teams: %w", err)
	}

	messages := make([]OutgoingQueueMessage, 0, len(teams))

	for _, team := range teams {
		if !team.Entitlements.IndividualFeatures {
			continue
		}
		messages = append(messages, OutgoingQueueMessage{
			Message: QueueMessage{
				QueueTeamReportGeneration: &QueueTeamReportGenerationInput{
					TeamId:                         team.Id,
					StartTime:                      input.StartTime,
					Duration:                       input.Duration,
					Retention:                      team.Entitlements.ReportRetention(),
					MaxSourceBytesPerAccountRegion: team.Entitlements.MaxSourceBytesPerAccountRegion(),
				},
			},
		})
	}

	return a.QueueMessages(ctx, map[string][]OutgoingQueueMessage{
		a.awsRegion: messages,
	})
}

type QueueTeamReportGenerationInput struct {
	TeamId                         model.Id
	StartTime                      time.Time
	Duration                       time.Duration
	Retention                      model.ReportRetention
	MaxSourceBytesPerAccountRegion int64
}

func (s *Session) QueueTeamReportGeneration(ctx context.Context, input QueueTeamReportGenerationInput) UserFacingError {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		return AuthorizationError{}
	}
	return s.SanitizedError(s.app.QueueTeamReportGeneration(ctx, input))
}

func (a *App) QueueTeamReportGeneration(ctx context.Context, input QueueTeamReportGenerationInput) error {
	integrations, err := a.store.GetAWSIntegrationsByTeamId(ctx, input.TeamId)
	if err != nil {
		return fmt.Errorf("failed to get aws integrations: %w", err)
	}
	for _, integration := range integrations {
		if err := a.queueAWSIntegrationReportGeneration(ctx, queueAWSIntegrationReportGenerationInput{
			Integration:                    integration,
			StartTime:                      input.StartTime,
			Duration:                       input.Duration,
			Retention:                      input.Retention,
			MaxSourceBytesPerAccountRegion: input.MaxSourceBytesPerAccountRegion,
		}); err != nil {
			return fmt.Errorf("failed to queue aws integration report generation: %w", err)
		}
	}
	return nil
}

type QueueAWSIntegrationReportGenerationInput struct {
	IntegrationId                  model.Id
	StartTime                      time.Time
	Duration                       time.Duration
	Retention                      model.ReportRetention
	MaxSourceBytesPerAccountRegion int64
}

func (s *Session) QueueAWSIntegrationReportGeneration(ctx context.Context, input QueueAWSIntegrationReportGenerationInput) UserFacingError {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		return AuthorizationError{}
	}
	integration, err := s.app.store.GetAWSIntegrationById(ctx, input.IntegrationId)
	if err != nil {
		return s.SanitizedError(err)
	} else if integration == nil {
		return NotFoundError("No such AWS integration.")
	}
	return s.SanitizedError(s.app.queueAWSIntegrationReportGeneration(ctx, queueAWSIntegrationReportGenerationInput{
		Integration:                    integration,
		StartTime:                      input.StartTime,
		Duration:                       input.Duration,
		Retention:                      input.Retention,
		MaxSourceBytesPerAccountRegion: input.MaxSourceBytesPerAccountRegion,
	}))
}

type queueAWSIntegrationReportGenerationInput struct {
	Integration                    *model.AWSIntegration
	StartTime                      time.Time
	Duration                       time.Duration
	Retention                      model.ReportRetention
	MaxSourceBytesPerAccountRegion int64
}

func (a *App) queueAWSIntegrationReportGeneration(ctx context.Context, input queueAWSIntegrationReportGenerationInput) error {
	output, err := a.sts.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &input.Integration.RoleARN,
		RoleSessionName: aws.String("cloud_snitch"),
		ExternalId:      aws.String(input.Integration.TeamId.String()),
	})
	if err != nil {
		return fmt.Errorf("failed to assume role: %w", err)
	}
	creds := output.Credentials

	accountRecon := map[string]PutAWSIntegrationReconAccountInput{}

	if input.Integration.GetAccountNamesFromOrganizations {
		orgsClient, err := a.organizationsFactory.NewFromSTSCredentials(ctx, creds)
		if err != nil {
			return fmt.Errorf("failed to create organizations client: %w", err)
		}

		var nextToken *string
		for {
			output, err := orgsClient.ListAccounts(ctx, &organizations.ListAccountsInput{
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("failed to list accounts: %w", err)
			}

			for _, account := range output.Accounts {
				accountRecon[*account.Id] = PutAWSIntegrationReconAccountInput{
					Id:   *account.Id,
					Name: emptyIfNil(account.Name),
				}
			}

			if output.NextToken == nil {
				break
			}
			nextToken = output.NextToken
		}
	}

	if trail := input.Integration.CloudTrailTrail; trail != nil {
		s3Client, err := a.s3Factory.NewFromSTSCredentials(ctx, creds)
		if err != nil {
			return fmt.Errorf("failed to create s3 client: %w", err)
		}

		locationOutput, err := s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: &trail.S3BucketName,
		})
		if err != nil {
			return fmt.Errorf("failed to get bucket location: %w", err)
		}
		bucketRegion := S3BucketLocationConstraintRegion(locationOutput.LocationConstraint)

		accountRegions, err := report.ScanAWSCloudTrailLogBucket(ctx, report.ScanAWSCloudTrailLogBucketConfig{
			S3:         s3Client,
			BucketName: trail.S3BucketName,
			KeyPrefix:  trail.S3KeyPrefix,
		})
		if err != nil {
			return fmt.Errorf("failed to scan aws cloudtrail log bucket: %w", err)
		}

		messagesByQueueRegion := map[string][]OutgoingQueueMessage{}

		for _, accountRegion := range accountRegions {
			if _, ok := accountRecon[accountRegion.AccountId]; !ok {
				accountRecon[accountRegion.AccountId] = PutAWSIntegrationReconAccountInput{
					Id: accountRegion.AccountId,
				}
			}

			queueRegion := a.ClosestAvailableAWSRegion(bucketRegion)
			messagesByQueueRegion[queueRegion] = append(messagesByQueueRegion[queueRegion], OutgoingQueueMessage{
				Message: QueueMessage{
					GenerateAWSCloudTrailReport: &GenerateAWSCloudTrailReportInput{
						FutureReportId:    model.NewReportId(),
						AWSIntegrationId:  input.Integration.Id,
						StartTime:         input.StartTime,
						Duration:          input.Duration,
						AccountsKeyPrefix: accountRegion.AccountsPrefix,
						AccountId:         accountRegion.AccountId,
						Region:            accountRegion.Region,
						Retention:         input.Retention,
						MaxSourceBytes:    input.MaxSourceBytesPerAccountRegion,
					},
				},
			})
		}

		if err := a.QueueMessages(ctx, messagesByQueueRegion); err != nil {
			return fmt.Errorf("failed to queue messages: %w", err)
		}
	}

	accountRecons := make([]PutAWSIntegrationReconAccountInput, 0, len(accountRecon))
	for _, accountRecon := range accountRecon {
		accountRecons = append(accountRecons, accountRecon)
	}
	if err := a.PutAWSIntegrationRecon(ctx, PutAWSIntegrationReconInput{
		AWSIntegrationId: input.Integration.Id,
		TeamId:           input.Integration.TeamId,
		Time:             time.Now(),
		Accounts:         accountRecons,
	}); err != nil {
		return fmt.Errorf("failed to put aws integration recon: %w", err)
	}

	return nil
}

type GenerateAWSCloudTrailReportInput struct {
	// The id is generated in advance so that report generation is idempotent.
	FutureReportId model.Id

	AWSIntegrationId  model.Id
	StartTime         time.Time
	Duration          time.Duration
	AccountsKeyPrefix string
	AccountId         string
	Region            string
	Retention         model.ReportRetention
	MaxSourceBytes    int64
}

// Synchronously generates and persists a report.
func (a *App) GenerateAWSCloudTrailReport(ctx context.Context, input GenerateAWSCloudTrailReportInput) (*model.Report, error) {
	startTime := time.Now()

	integration, err := a.store.GetAWSIntegrationById(ctx, input.AWSIntegrationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get aws integration: %w", err)
	}

	output, err := a.sts.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &integration.RoleARN,
		RoleSessionName: aws.String("cloud_snitch"),
		ExternalId:      aws.String(integration.TeamId.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to assume role: %w", err)
	}
	creds := output.Credentials

	s3Client, err := a.s3Factory.NewFromSTSCredentials(ctx, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	r := &report.Report{
		StartTime:       input.StartTime,
		DurationSeconds: int(input.Duration.Seconds()),
	}

	if err := r.ImportAWSCloudTrailLogsForAccountRegion(ctx, report.ImportAWSCloudTrailLogsForAccountRegionConfig{
		S3:             s3Client,
		BucketName:     integration.CloudTrailTrail.S3BucketName,
		AccountsPrefix: input.AccountsKeyPrefix,
		AccountId:      input.AccountId,
		Region:         input.Region,
		MaxSourceBytes: input.MaxSourceBytes,
	}); err != nil {
		return nil, fmt.Errorf("failed to import aws cloudtrail logs: %w", err)
	}

	if r.IsEmpty() {
		return nil, nil
	}

	buf, err := jsoniter.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report: %w", err)
	}

	key := "reports/" + input.FutureReportId.String() + ".json"

	if _, err := a.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:  &a.config.S3BucketName,
		Key:     &key,
		Body:    bytes.NewReader(buf),
		Tagging: aws.String("team_id=" + integration.TeamId.String() + "&retention=" + string(input.Retention)),
	}); err != nil {
		return nil, fmt.Errorf("failed to put report in s3: %w", err)
	}

	expirationTime := input.StartTime.Add(input.Duration + input.Retention.Duration())

	downloadURL := ""
	if a.urlSigner != nil {
		url := a.config.S3CDNURL + "/" + key
		downloadURL, err = a.urlSigner.Sign(url, expirationTime.Add(5*time.Minute))
		if err != nil {
			return nil, fmt.Errorf("failed to sign url: %w", err)
		}
	}

	ret := &model.Report{
		Id:               input.FutureReportId,
		CreationTime:     time.Now(),
		ExpirationTime:   expirationTime,
		TeamId:           integration.TeamId,
		AWSIntegrationId: input.AWSIntegrationId,
		Scope: model.ReportScope{
			StartTime: input.StartTime,
			Duration:  input.Duration,
			AWS: model.ReportScopeAWS{
				AccountId: input.AccountId,
				Region:    input.Region,
			},
		},
		Location: model.ReportLocation{
			AWSRegion: a.awsRegion,
			S3Bucket:  a.config.S3BucketName,
			Key:       key,
		},
		DownloadURL:        downloadURL,
		Size:               len(buf),
		SourceBytes:        int(r.SourceBytes),
		IsIncomplete:       r.IsIncomplete,
		GenerationDuration: time.Since(startTime),
	}

	if err := a.store.PutReport(ctx, ret); err != nil {
		return nil, fmt.Errorf("failed to put report in store: %w", err)
	}

	if err := a.store.PutTeamBillableAccount(ctx, &model.TeamBillableAccount{
		Id:             input.AccountId,
		TeamId:         integration.TeamId,
		ExpirationTime: time.Now().Add(72 * time.Hour),
	}); err != nil {
		return nil, fmt.Errorf("failed to put team billable account: %w", err)
	}

	return ret, nil
}

func (s *Session) DeleteReportById(ctx context.Context, id model.Id) UserFacingError {
	report, err := s.app.store.GetReportById(ctx, id)
	if err != nil || report == nil {
		return s.SanitizedError(err)
	} else if err := s.RequireTeamAdministrator(ctx, report.TeamId); err != nil {
		return err
	}
	if err := s.app.store.DeleteReportById(ctx, id); err != nil {
		return s.SanitizedError(err)
	}
	// TODO: should we also delete from S3?
	return nil
}
