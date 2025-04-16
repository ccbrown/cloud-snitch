package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

type CreateAWSIntegrationInput struct {
	TeamId                           model.Id
	Name                             string
	RoleARN                          string
	GetAccountNamesFromOrganizations bool
	CloudTrailTrail                  *CreateAWSIntegrationCloudTrailTrailInput
	QueueReportGeneration            bool
}

type CreateAWSIntegrationCloudTrailTrailInput struct {
	S3BucketName string
	S3KeyPrefix  string
}

func (s *Session) ValidateAWSIntegrationRole(ctx context.Context, input CreateAWSIntegrationInput) UserFacingError {
	// First, make sure we *can't* assume the role without the external id.
	if _, err := s.app.sts.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &input.RoleARN,
		RoleSessionName: aws.String("cloud_snitch_validation"),
	}); err == nil {
		return NewUserError("The provided role was able to be assumed without an external id. Please configure the role's trust relationship to require the external id.")
	}

	// Now try to assume the role with the external id.
	output, err := s.app.sts.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &input.RoleARN,
		RoleSessionName: aws.String("cloud_snitch_validation"),
		ExternalId:      aws.String(input.TeamId.String()),
	})
	if err != nil {
		return NewUserError("The provided role could not be assumed. Please double check the role's trust policy.")
	}
	creds := output.Credentials

	if input.GetAccountNamesFromOrganizations {
		orgsClient, err := s.app.organizationsFactory.NewFromSTSCredentials(ctx, creds)
		if err != nil {
			return s.SanitizedError(fmt.Errorf("failed to create organizations client: %w", err))
		}

		// organizations:ListAccounts
		{
			if _, err := orgsClient.ListAccounts(ctx, &organizations.ListAccountsInput{}); err != nil {
				return NewUserError("Unable to get account names. Please make sure the role has permission to perform the organizations:ListAccounts action.")
			}
		}
	}

	// Now check S3 permissions...
	if trail := input.CloudTrailTrail; trail != nil {
		bucketRegion, err := s.app.s3Factory.GetBucketRegion(ctx, trail.S3BucketName)
		if err != nil {
			return NewUserError("Unable to determine bucket location. Please make sure the given bucket exists.")
		}

		s3Client, err := s.app.s3Factory.NewFromSTSCredentials(ctx, creds, bucketRegion)
		if err != nil {
			return s.SanitizedError(fmt.Errorf("failed to create s3 client: %w", err))
		}

		// s3:ListBucket
		{
			_, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:  aws.String(trail.S3BucketName),
				Prefix:  aws.String(trail.S3KeyPrefix),
				MaxKeys: aws.Int32(1),
			})
			if err != nil {
				return NewUserError("Unable to list bucket contents. Please make sure the role has permission to perform the s3:ListBucket action.")
			}
		}

		// s3:GetObject
		{
			_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(trail.S3BucketName),
				Key:    aws.String(trail.S3KeyPrefix + "cloud-snitch-permission-check"),
			})
			if err != nil {
				var apiErr smithy.APIError
				if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "NotFound" {
					return NewUserError("Unable to request bucket objects. Please make sure the role has permission to perform the s3:GetObject action.")
				}
			}
		}
	}

	return nil
}

func (s *Session) CreateAWSIntegration(ctx context.Context, input CreateAWSIntegrationInput) (*model.AWSIntegration, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, input.TeamId); err != nil {
		return nil, err
	} else if err := ValidateName(input.Name); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, input.TeamId, store.ConsistencyEventual)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("No such team.")
	} else if !team.Entitlements.IndividualFeatures {
		return nil, NewUserError("An active subscription is required to create AWS integrations.")
	}

	if existing, err := s.app.store.GetAWSIntegrationsByTeamId(ctx, input.TeamId); err != nil {
		return nil, s.SanitizedError(err)
	} else if len(existing) > 50 {
		return nil, NewUserError("To prevent abuse, teams may only have up to 50 AWS integrations by default. If you need more, please contact support.")
	} else {
		for _, integration := range existing {
			trail := integration.CloudTrailTrail
			if trail.S3BucketName == input.CloudTrailTrail.S3BucketName && trail.S3KeyPrefix == input.CloudTrailTrail.S3KeyPrefix {
				return nil, NewUserError("An integration for that trail already exists.")
			}
		}
	}

	if err := s.ValidateAWSIntegrationRole(ctx, input); err != nil {
		return nil, err
	}

	integration := &model.AWSIntegration{
		Id:                               model.NewAWSIntegrationId(),
		CreationTime:                     time.Now(),
		TeamId:                           input.TeamId,
		Name:                             input.Name,
		RoleARN:                          input.RoleARN,
		GetAccountNamesFromOrganizations: input.GetAccountNamesFromOrganizations,
	}
	if trail := input.CloudTrailTrail; trail != nil {
		integration.CloudTrailTrail = &model.AWSIntegrationCloudTrailTrail{
			S3BucketName: trail.S3BucketName,
			S3KeyPrefix:  trail.S3KeyPrefix,
		}
	}

	if err := s.app.store.PutAWSIntegration(ctx, integration); err != nil {
		return nil, s.SanitizedError(err)
	}

	if input.QueueReportGeneration {
		today := time.Now().Truncate(24 * time.Hour)
		for i := 0; i < 7; i++ {
			if err := s.app.queueAWSIntegrationReportGeneration(ctx, queueAWSIntegrationReportGenerationInput{
				Integration: integration,
				StartTime:   today.Add(-time.Duration(i) * 24 * time.Hour),
				Duration:    24 * time.Hour,
				Retention:   team.Entitlements.ReportRetention(),
			}); err != nil {
				return nil, s.SanitizedError(fmt.Errorf("failed to queue report generation: %w", err))
			}
		}
	}

	return integration, nil
}

func (s *Session) GetAWSIntegrationsByTeamId(ctx context.Context, teamId model.Id) ([]*model.AWSIntegration, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}
	integrations, err := s.app.store.GetAWSIntegrationsByTeamId(ctx, teamId)
	return integrations, s.SanitizedError(err)
}

type AWSIntegrationPatch struct {
	Name *string
}

func (s *Session) PatchAWSIntegrationById(ctx context.Context, id model.Id, patch AWSIntegrationPatch) (*model.AWSIntegration, UserFacingError) {
	if err := s.RequireUser(); err != nil {
		return nil, err
	} else if integration, err := s.app.store.GetAWSIntegrationById(ctx, id); integration == nil || err != nil {
		return nil, s.SanitizedError(err)
	} else if err := s.RequireTeamAdministrator(ctx, integration.TeamId); err != nil {
		return nil, err
	}

	storePatch := &store.AWSIntegrationPatch{
		Name: patch.Name,
	}
	if patch.Name != nil {
		if err := ValidateName(*patch.Name); err != nil {
			return nil, err
		}
	}
	team, err := s.app.store.PatchAWSIntegrationById(ctx, id, storePatch)
	return team, s.SanitizedError(err)
}

func (s *Session) DeleteAWSIntegrationById(ctx context.Context, id model.Id, deleteAssociatedData bool) UserFacingError {
	if err := s.RequireUser(); err != nil {
		return err
	} else if integration, err := s.app.store.GetAWSIntegrationById(ctx, id); integration == nil || err != nil {
		return s.SanitizedError(err)
	} else if err := s.RequireTeamAdministrator(ctx, integration.TeamId); err != nil {
		return err
	}

	if deleteAssociatedData {
		if reports, err := s.app.store.GetReportsByTeamId(ctx, id); err != nil {
			return s.SanitizedError(err)
		} else {
			var toDelete []model.Id
			for _, report := range reports {
				if report.AWSIntegrationId == id {
					toDelete = append(toDelete, report.Id)
				}
			}
			if err := s.app.store.DeleteReportsByIds(ctx, toDelete...); err != nil {
				return s.SanitizedError(err)
			}
		}

		if err := s.app.store.DeleteAWSIntegrationReconByAWSIntegrationId(ctx, id); err != nil {
			return s.SanitizedError(err)
		}
	}

	return s.SanitizedError(s.app.store.DeleteAWSIntegrationById(ctx, id))
}

func BestAvailableAWSRegion(region string, available []string) string {
	best := ""
	bestCommonPrefixLength := -1

	// Find the region with the longest common prefix with the given region.
	// TODO: This could probably be more intelligent.
	for _, r := range available {
		commonPrefixLength := 0
		for i := 0; i < len(r) && i < len(region); i++ {
			if r[i] != region[i] {
				break
			}
			commonPrefixLength++
		}
		if commonPrefixLength > bestCommonPrefixLength {
			best = r
			bestCommonPrefixLength = commonPrefixLength
		}
	}

	return best
}

type PutAWSIntegrationReconInput struct {
	AWSIntegrationId model.Id
	TeamId           model.Id
	Time             time.Time
	Accounts         []PutAWSIntegrationReconAccountInput
}

type PutAWSIntegrationReconAccountInput struct {
	Id   string
	Name string
}

func (a *App) PutAWSIntegrationRecon(ctx context.Context, input PutAWSIntegrationReconInput) error {
	recon := &model.AWSIntegrationRecon{
		AWSIntegrationId: input.AWSIntegrationId,
		TeamId:           input.TeamId,
		Time:             input.Time,
		ExpirationTime:   input.Time.Add(3 * 24 * time.Hour),
		Accounts:         make([]model.AWSIntegrationAccountRecon, len(input.Accounts)),
	}
	for i, account := range input.Accounts {
		recon.Accounts[i] = model.AWSIntegrationAccountRecon{
			Id:   account.Id,
			Name: account.Name,
		}
	}

	if err := a.store.PutAWSIntegrationRecon(ctx, recon); err != nil {
		return fmt.Errorf("failed to put AWS integration recon: %w", err)
	}

	return nil
}

func (s *Session) GetAWSIntegrationReconsByTeamId(ctx context.Context, teamId model.Id) ([]*model.AWSIntegrationRecon, error) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}
	return s.app.store.GetAWSIntegrationReconsByTeamId(ctx, teamId)
}
