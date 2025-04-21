package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/smithy-go"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

type CreateAWSIntegrationInput struct {
	TeamId                           model.Id
	Name                             string
	RoleARN                          string
	GetAccountNamesFromOrganizations bool
	ManageSCPs                       bool
	CloudTrailTrail                  *CreateAWSIntegrationCloudTrailTrailInput
	QueueReportGeneration            bool
}

type CreateAWSIntegrationCloudTrailTrailInput struct {
	S3BucketName string
	S3KeyPrefix  string
}

func (s *Session) ValidateAWSIntegration(ctx context.Context, input CreateAWSIntegrationInput) UserFacingError {
	if input.ManageSCPs && !input.GetAccountNamesFromOrganizations {
		return NewUserError("To manage SCPs you must also allow the integration to get account info from AWS Organizations.")
	}

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
		accounts, err := orgsClient.ListAccounts(ctx, &organizations.ListAccountsInput{})
		if err != nil || len(accounts.Accounts) == 0 {
			return NewUserError("Unable to get account info. Please make sure the role has permission to perform the organizations:ListAccounts action.")
		}

		if input.ManageSCPs {
			// organizations:ListPoliciesForTarget
			{
				if _, err := orgsClient.ListPoliciesForTarget(ctx, &organizations.ListPoliciesForTargetInput{
					Filter:   organizationstypes.PolicyTypeServiceControlPolicy,
					TargetId: accounts.Accounts[0].Id,
				}); err != nil {
					return NewUserError("Unable to get policies. Please make sure the role has permission to perform the organizations:ListPoliciesForTarget action.")
				}
			}

			// organizations:ListParents
			{
				if _, err := orgsClient.ListParents(ctx, &organizations.ListParentsInput{
					ChildId: accounts.Accounts[0].Id,
				}); err != nil {
					return NewUserError("Unable to get account parents. Please make sure the role has permission to perform the organizations:ListParents action.")
				}
			}

			// organizations:ListRoots
			{
				output, err := orgsClient.ListRoots(ctx, &organizations.ListRootsInput{})
				if err != nil || len(output.Roots) == 0 {
					return NewUserError("Unable to get organization roots. Please make sure the role has permission to perform the organizations:ListRoots action.")
				}
				var hasSCPsEnabled bool
				for _, t := range output.Roots[0].PolicyTypes {
					if t.Status == organizationstypes.PolicyTypeStatusEnabled && t.Type == organizationstypes.PolicyTypeServiceControlPolicy {
						hasSCPsEnabled = true
						break
					}
				}
				if !hasSCPsEnabled {
					return NewUserError("SCPs are not enabled for the organization. Please enable them and try again.")
				}
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

	if err := s.ValidateAWSIntegration(ctx, input); err != nil {
		return nil, err
	}

	integration := &model.AWSIntegration{
		Id:                               model.NewAWSIntegrationId(),
		CreationTime:                     time.Now(),
		TeamId:                           input.TeamId,
		Name:                             input.Name,
		RoleARN:                          input.RoleARN,
		GetAccountNamesFromOrganizations: input.GetAccountNamesFromOrganizations,
		ManageSCPs:                       input.ManageSCPs,
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

	{
		today := time.Now().Truncate(24 * time.Hour)
		for i := 0; i < 7; i++ {
			if err := s.app.doReconAndQueueAWSIntegrationReportGeneration(ctx, doReconAndQueueAWSIntegrationReportGenerationInput{
				Integration: integration,
				StartTime:   today.Add(-time.Duration(i) * 24 * time.Hour),
				Duration:    24 * time.Hour,
				Retention:   team.Entitlements.ReportRetention(),
				ReconOnly:   !input.QueueReportGeneration,
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
	}
	integration, err := s.app.store.GetAWSIntegrationById(ctx, id)
	if integration == nil || err != nil {
		return s.SanitizedError(err)
	} else if err := s.RequireTeamAdministrator(ctx, integration.TeamId); err != nil {
		return err
	}

	if deleteAssociatedData {
		if reports, err := s.app.store.GetReportsByTeamId(ctx, integration.TeamId); err != nil {
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
	CanManageSCPs    bool
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
		CanManageSCPs:    input.CanManageSCPs,
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

func (s *Session) GetAWSIntegrationReconsByTeamId(ctx context.Context, teamId model.Id) ([]*model.AWSIntegrationRecon, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}
	ret, err := s.app.store.GetAWSIntegrationReconsByTeamId(ctx, teamId)
	return ret, s.SanitizedError(err)
}

func (a *App) assumeAWSIntegrationRole(ctx context.Context, integration *model.AWSIntegration) (*ststypes.Credentials, error) {
	output, err := a.sts.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &integration.RoleARN,
		RoleSessionName: aws.String("cloud_snitch"),
		ExternalId:      aws.String(integration.TeamId.String()),
	})
	if err != nil {
		return nil, err
	}
	return output.Credentials, nil
}

// Finds an integration for the given team which is capable of managing SCPs for the given account.
func (a *App) awsSCPManagementIntegrationByTeamAndAccountId(ctx context.Context, teamId model.Id, accountId string) (*model.AWSIntegration, error) {
	recons, err := a.store.GetAWSIntegrationReconsByTeamId(ctx, teamId)
	if err != nil {
		return nil, err
	}
	for _, recon := range recons {
		for _, account := range recon.Accounts {
			if account.Id == accountId {
				if integration, err := a.store.GetAWSIntegrationById(ctx, recon.AWSIntegrationId); err != nil {
					return nil, err
				} else if integration != nil && integration.ManageSCPs {
					return integration, nil
				}
			}
		}
	}
	return nil, nil
}

const ManagedAWSSCPNamePrefix = "CloudSnitchManagedSCP-"

const CloudSnitchManagedResourceTag = "CloudSnitchManaged"

func (a *App) findManagedAWSSCP(ctx context.Context, orgsClient AWSOrganizationsAPI, accountId string) (*organizationstypes.PolicySummary, error) {
	var nextToken *string
	for {
		output, err := orgsClient.ListPoliciesForTarget(ctx, &organizations.ListPoliciesForTargetInput{
			Filter:    organizationstypes.PolicyTypeServiceControlPolicy,
			TargetId:  aws.String(accountId),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list account scps: %w", err)
		}

		for _, policySummary := range output.Policies {
			if policySummary.Name == nil || *policySummary.Name != ManagedAWSSCPNamePrefix+accountId {
				continue
			}
			return &policySummary, nil
		}

		if output.NextToken == nil {
			return nil, nil
		}
		nextToken = output.NextToken
	}
}

func (s *Session) GetManagedAWSSCPByTeamAndAccountId(ctx context.Context, teamId model.Id, accountId string) (*model.AWSSCP, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}

	integration, err := s.app.awsSCPManagementIntegrationByTeamAndAccountId(ctx, teamId, accountId)
	if err != nil || integration == nil {
		return nil, s.SanitizedError(err)
	}

	creds, err := s.app.assumeAWSIntegrationRole(ctx, integration)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to assume role: %w", err))
	}

	orgsClient, err := s.app.organizationsFactory.NewFromSTSCredentials(ctx, creds)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to create organizations client: %w", err))
	}

	policySummary, err := s.app.findManagedAWSSCP(ctx, orgsClient, accountId)
	if err != nil || policySummary == nil {
		return nil, s.SanitizedError(err)
	}

	if policy, err := orgsClient.DescribePolicy(ctx, &organizations.DescribePolicyInput{
		PolicyId: policySummary.Id,
	}); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to describe policy: %w", err))
	} else if policy.Policy != nil && policy.Policy.Content != nil {
		return &model.AWSSCP{
			Content: *policy.Policy.Content,
		}, nil
	} else {
		return nil, nil
	}
}

type PutManagedAWSSCPInput struct {
	Content string
}

func (s *Session) PutManagedAWSSCPByTeamAndAccountId(ctx context.Context, teamId model.Id, accountId string, input PutManagedAWSSCPInput) (*model.AWSSCP, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}

	integration, err := s.app.awsSCPManagementIntegrationByTeamAndAccountId(ctx, teamId, accountId)
	if err != nil || integration == nil {
		return nil, s.SanitizedError(err)
	}

	creds, err := s.app.assumeAWSIntegrationRole(ctx, integration)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to assume role: %w", err))
	}

	orgsClient, err := s.app.organizationsFactory.NewFromSTSCredentials(ctx, creds)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to create organizations client: %w", err))
	}

	ret := &model.AWSSCP{
		Content: input.Content,
	}

	policySummary, err := s.app.findManagedAWSSCP(ctx, orgsClient, accountId)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	if policySummary != nil {
		// Existing policy found, just update it.

		if _, err := orgsClient.UpdatePolicy(ctx, &organizations.UpdatePolicyInput{
			PolicyId: policySummary.Id,
			Content:  aws.String(input.Content),
		}); err != nil {
			return nil, s.SanitizedError(fmt.Errorf("failed to update policy: %w", err))
		}
	} else {
		// Create a new policy and attach it.

		policy, err := orgsClient.CreatePolicy(ctx, &organizations.CreatePolicyInput{
			Name:        aws.String(ManagedAWSSCPNamePrefix + accountId),
			Description: aws.String("Managed by CloudSnitch (" + s.app.config.FrontendURL + "). Do not modify directly."),
			Type:        organizationstypes.PolicyTypeServiceControlPolicy,
			Content:     &input.Content,
			Tags: []organizationstypes.Tag{
				{
					Key:   aws.String(CloudSnitchManagedResourceTag),
					Value: aws.String("true"),
				},
			},
		})
		if err != nil {
			return nil, s.SanitizedError(fmt.Errorf("failed to create policy: %w", err))
		}

		if _, err := orgsClient.AttachPolicy(ctx, &organizations.AttachPolicyInput{
			PolicyId: policy.Policy.PolicySummary.Id,
			TargetId: aws.String(accountId),
		}); err != nil {
			return nil, s.SanitizedError(fmt.Errorf("failed to attach policy: %w", err))
		}
	}

	return ret, nil
}

func (s *Session) GetAWSAccessReportByTeamAndAccountId(ctx context.Context, teamId model.Id, accountId string) (*model.AWSAccessReport, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}

	integration, err := s.app.awsSCPManagementIntegrationByTeamAndAccountId(ctx, teamId, accountId)
	if err != nil || integration == nil {
		return nil, s.SanitizedError(err)
	}

	creds, err := s.app.assumeAWSIntegrationRole(ctx, integration)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to assume role: %w", err))
	}

	orgsClient, err := s.app.organizationsFactory.NewFromSTSCredentials(ctx, creds)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to create organizations client: %w", err))
	}

	var entityPath string

	// It sure takes a lot of effort to get the entity path... Is there an easier way to do this?
	{
		entityPathComponents := []string{accountId}
		for {
			output, err := orgsClient.ListParents(ctx, &organizations.ListParentsInput{
				ChildId: aws.String(entityPathComponents[0]),
			})
			if err != nil {
				return nil, s.SanitizedError(fmt.Errorf("failed to list parents: %w", err))
			}
			if len(output.Parents) == 0 {
				return nil, s.SanitizedError(fmt.Errorf("failed to find organization root"))
			}
			p := output.Parents[0]
			entityPathComponents = append([]string{*p.Id}, entityPathComponents...)
			if p.Type == organizationstypes.ParentTypeRoot {
				break
			}
		}

		output, err := orgsClient.ListRoots(ctx, &organizations.ListRootsInput{})
		if err != nil || len(output.Roots) == 0 {
			return nil, s.SanitizedError(fmt.Errorf("failed to list roots: %w", err))
		}
		rootArnParts := strings.Split(*output.Roots[0].Arn, "/")
		organizationId := rootArnParts[1]
		entityPathComponents = append([]string{organizationId}, entityPathComponents...)

		entityPath = strings.Join(entityPathComponents, "/")
	}

	iamClient, err := s.app.iamFactory.NewFromSTSCredentials(ctx, creds)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to create organizations client: %w", err))
	}

	output, err := iamClient.GenerateOrganizationsAccessReport(ctx, &iam.GenerateOrganizationsAccessReportInput{
		EntityPath: aws.String(entityPath),
	})
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("failed to generate access report: %w", err))
	}

	var ret model.AWSAccessReport

	jobId := *output.JobId
	var marker *string

	for {
		output, err := iamClient.GetOrganizationsAccessReport(ctx, &iam.GetOrganizationsAccessReportInput{
			JobId:  aws.String(jobId),
			Marker: marker,
		})
		if err != nil {
			return nil, s.SanitizedError(fmt.Errorf("failed to get access report: %w", err))
		}
		if output.JobStatus == iamtypes.JobStatusTypeInProgress {
			time.Sleep(time.Second)
			continue
		} else if output.JobStatus != iamtypes.JobStatusTypeCompleted {
			var message string
			if output.ErrorDetails != nil && output.ErrorDetails.Message != nil {
				message = *output.ErrorDetails.Message
			}
			return nil, s.SanitizedError(fmt.Errorf("unexpected job status for access report: %s, error message: %s", string(output.JobStatus), message))
		}

		for _, detail := range output.AccessDetails {
			service := model.AWSAccessReportService{
				Name:      *detail.ServiceName,
				Namespace: *detail.ServiceNamespace,
			}
			if detail.LastAuthenticatedTime != nil {
				service.LastAuthenticationTime = *detail.LastAuthenticatedTime
			}
			ret.Services = append(ret.Services, service)
		}

		if !output.IsTruncated {
			break
		}
		marker = output.Marker
	}

	return &ret, nil
}
