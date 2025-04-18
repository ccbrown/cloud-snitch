package apptest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

type TestAWSSTSAPI struct{}

func (api *TestAWSSTSAPI) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	if !strings.Contains(*params.RoleArn, "NoExternalIdRequired") && params.ExternalId == nil {
		return nil, fmt.Errorf("external id required")
	}
	return &sts.AssumeRoleOutput{
		Credentials: &ststypes.Credentials{
			AccessKeyId: params.RoleArn,
		},
	}, nil
}

// Paths are relative to the root of the Go module.
var bucketPaths = map[string]string{
	"aws-cloudtrail-logs": "report/testdata/aws-cloudtrail-logs",
}

type TestAmazonS3API struct{}

func (api *TestAmazonS3API) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, nil
}

func (api *TestAmazonS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if bucketPath, ok := bucketPaths[*params.Bucket]; ok {
		_, b, _, _ := runtime.Caller(0)
		bucketPath = filepath.Dir(filepath.Dir(filepath.Dir(b))) + "/" + bucketPath
		objectPath := filepath.Join(bucketPath, *params.Key)
		buf, err := os.ReadFile(objectPath)
		if err != nil {
			return nil, err
		}
		return &s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader(buf)),
		}, nil
	}

	return &s3.GetObjectOutput{}, nil
}

func (api *TestAmazonS3API) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}

func (api *TestAmazonS3API) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if bucketPath, ok := bucketPaths[*params.Bucket]; ok {
		_, b, _, _ := runtime.Caller(0)
		bucketPath = filepath.Dir(filepath.Dir(filepath.Dir(b))) + "/" + bucketPath

		if params.Delimiter != nil {
			entries, err := os.ReadDir(bucketPath + "/" + *params.Prefix)
			if !os.IsNotExist(err) && err != nil {
				return nil, err
			}
			commonPrefixes := make([]s3types.CommonPrefix, 0, len(entries))
			for _, entry := range entries {
				if entry.IsDir() {
					commonPrefixes = append(commonPrefixes, s3types.CommonPrefix{
						Prefix: aws.String(*params.Prefix + entry.Name() + "/"),
					})
				}
			}
			return &s3.ListObjectsV2Output{
				CommonPrefixes: commonPrefixes,
			}, nil
		}

		prefixDir := *params.Prefix
		if !strings.HasSuffix(prefixDir, "/") {
			prefixDir = filepath.Dir(prefixDir)
		}
		var objects []s3types.Object
		err := filepath.Walk(bucketPath+"/"+prefixDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				key := strings.TrimPrefix(path, bucketPath+"/")
				if !info.IsDir() && strings.HasPrefix(key, *params.Prefix) {
					objects = append(objects, s3types.Object{
						Key: aws.String(key),
					})
				}
				return nil
			})
		if !os.IsNotExist(err) && err != nil {
			return nil, err
		}
		return &s3.ListObjectsV2Output{
			Contents: objects,
		}, nil
	}

	return &s3.ListObjectsV2Output{}, nil
}

type TestAmazonS3APIFactory struct{}

func (TestAmazonS3APIFactory) GetBucketRegion(ctx context.Context, bucketName string) (string, error) {
	return "us-east-1", nil
}

func (f TestAmazonS3APIFactory) NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials, region string) (app.AmazonS3API, error) {
	return &TestAmazonS3API{}, nil
}

type TestAmazonSQSAPI struct {
	m        sync.Mutex
	requests []*sqs.SendMessageBatchInput
}

func (api *TestAmazonSQSAPI) SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()
	api.requests = append(api.requests, params)
	return &sqs.SendMessageBatchOutput{}, nil
}

type TestAmazonSQSAPIFactory struct {
	m    sync.Mutex
	apis map[string]*TestAmazonSQSAPI
}

func (f *TestAmazonSQSAPIFactory) NewWithRegion(ctx context.Context, region string) (app.AmazonSQSAPI, error) {
	f.m.Lock()
	defer f.m.Unlock()
	if f.apis == nil {
		f.apis = map[string]*TestAmazonSQSAPI{}
	}
	if f.apis[region] == nil {
		f.apis[region] = &TestAmazonSQSAPI{}
	}
	return f.apis[region], nil
}

func (f *TestAmazonSQSAPIFactory) Requests(region string) []*sqs.SendMessageBatchInput {
	f.m.Lock()
	defer f.m.Unlock()
	return append([]*sqs.SendMessageBatchInput{}, f.apis[region].requests...)
}

type TestAWSIAMAPI struct{}

func (TestAWSIAMAPI) GenerateOrganizationsAccessReport(ctx context.Context, params *iam.GenerateOrganizationsAccessReportInput, optFns ...func(*iam.Options)) (*iam.GenerateOrganizationsAccessReportOutput, error) {
	if *params.EntityPath != "o-1234/r-1234/123456789012" {
		return nil, fmt.Errorf("invalid entity path")
	}
	return &iam.GenerateOrganizationsAccessReportOutput{
		JobId: aws.String("job-id"),
	}, nil
}

func (TestAWSIAMAPI) GetOrganizationsAccessReport(ctx context.Context, params *iam.GetOrganizationsAccessReportInput, optFns ...func(*iam.Options)) (*iam.GetOrganizationsAccessReportOutput, error) {
	var marker string
	if params.Marker != nil {
		marker = *params.Marker
	}
	switch marker {
	case "":
		return &iam.GetOrganizationsAccessReportOutput{
			JobStatus: iamtypes.JobStatusTypeCompleted,
			AccessDetails: []iamtypes.AccessDetail{
				{
					ServiceName:      aws.String("AWS App2Container"),
					ServiceNamespace: aws.String("a2c"),
				},
				{
					ServiceName:      aws.String("Alexa for Business"),
					ServiceNamespace: aws.String("a4b"),
				},
			},
			IsTruncated: true,
			Marker:      aws.String("p2"),
		}, nil
	case "p2":
		return &iam.GetOrganizationsAccessReportOutput{
			JobStatus: iamtypes.JobStatusTypeCompleted,
			AccessDetails: []iamtypes.AccessDetail{
				{
					ServiceName:           aws.String("AWS IAM Access Analyzer"),
					ServiceNamespace:      aws.String("access-analyzer"),
					LastAuthenticatedTime: aws.Time(time.Now()),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid marker")
	}
}

type TestAWSIAMAPIFactory struct{}

func (TestAWSIAMAPIFactory) NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials) (app.AWSIAMAPI, error) {
	return &TestAWSIAMAPI{}, nil
}

type TestAWSOrganizationsAPI struct {
	m                 sync.Mutex
	policiesById      map[string]*organizationstypes.Policy
	attachedPolicyIds map[string][]string
}

func (api *TestAWSOrganizationsAPI) ListAccounts(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error) {
	return &organizations.ListAccountsOutput{
		Accounts: []organizationstypes.Account{
			{
				Id:     aws.String("123456789012"),
				Name:   aws.String("Test Account"),
				Status: organizationstypes.AccountStatusActive,
			},
		},
	}, nil
}

func (api *TestAWSOrganizationsAPI) ListParents(ctx context.Context, params *organizations.ListParentsInput, optFns ...func(*organizations.Options)) (*organizations.ListParentsOutput, error) {
	return &organizations.ListParentsOutput{
		Parents: []organizationstypes.Parent{
			{
				Id:   aws.String("r-1234"),
				Type: organizationstypes.ParentTypeRoot,
			},
		},
	}, nil
}

func (api *TestAWSOrganizationsAPI) ListPoliciesForTarget(ctx context.Context, params *organizations.ListPoliciesForTargetInput, optFns ...func(*organizations.Options)) (*organizations.ListPoliciesForTargetOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()

	ret := &organizations.ListPoliciesForTargetOutput{}
	for _, id := range api.attachedPolicyIds[*params.TargetId] {
		ret.Policies = append(ret.Policies, *api.policiesById[id].PolicySummary)
	}
	return ret, nil
}

func (api *TestAWSOrganizationsAPI) DescribePolicy(ctx context.Context, params *organizations.DescribePolicyInput, optFns ...func(*organizations.Options)) (*organizations.DescribePolicyOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()

	if policy, ok := api.policiesById[*params.PolicyId]; ok {
		return &organizations.DescribePolicyOutput{
			Policy: policy,
		}, nil
	}
	return nil, fmt.Errorf("policy not found")
}

func (api *TestAWSOrganizationsAPI) AttachPolicy(ctx context.Context, params *organizations.AttachPolicyInput, optFns ...func(*organizations.Options)) (*organizations.AttachPolicyOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()

	if api.attachedPolicyIds == nil {
		api.attachedPolicyIds = map[string][]string{}
	}
	api.attachedPolicyIds[*params.TargetId] = append(api.attachedPolicyIds[*params.TargetId], *params.PolicyId)

	return &organizations.AttachPolicyOutput{}, nil
}

func (api *TestAWSOrganizationsAPI) CreatePolicy(ctx context.Context, params *organizations.CreatePolicyInput, optFns ...func(*organizations.Options)) (*organizations.CreatePolicyOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()

	if api.policiesById == nil {
		api.policiesById = map[string]*organizationstypes.Policy{}
	}

	policy := &organizationstypes.Policy{
		Content: params.Content,
		PolicySummary: &organizationstypes.PolicySummary{
			Id:   aws.String(model.NewId("p").String()),
			Name: params.Name,
		},
	}
	api.policiesById[*policy.PolicySummary.Id] = policy

	return &organizations.CreatePolicyOutput{
		Policy: policy,
	}, nil
}

func (api *TestAWSOrganizationsAPI) UpdatePolicy(ctx context.Context, params *organizations.UpdatePolicyInput, optFns ...func(*organizations.Options)) (*organizations.UpdatePolicyOutput, error) {
	api.m.Lock()
	defer api.m.Unlock()

	if policy, ok := api.policiesById[*params.PolicyId]; ok {
		policy.Content = params.Content

		return &organizations.UpdatePolicyOutput{
			Policy: policy,
		}, nil
	}
	return nil, fmt.Errorf("policy not found")
}

func (api *TestAWSOrganizationsAPI) ListRoots(ctx context.Context, params *organizations.ListRootsInput, optFns ...func(*organizations.Options)) (*organizations.ListRootsOutput, error) {
	return &organizations.ListRootsOutput{
		Roots: []organizationstypes.Root{
			{
				Id:   aws.String("r-1234"),
				Arn:  aws.String("arn:aws:organizations::123456789012:root/o-1234/r-1234"),
				Name: aws.String("Root"),
				PolicyTypes: []organizationstypes.PolicyTypeSummary{
					{
						Type:   organizationstypes.PolicyTypeServiceControlPolicy,
						Status: organizationstypes.PolicyTypeStatusEnabled,
					},
				},
			},
		},
	}, nil
}

type TestAWSOrganizationsAPIFactory struct {
	m    sync.Mutex
	orgs map[string]*TestAWSOrganizationsAPI
}

func (f *TestAWSOrganizationsAPIFactory) NewFromSTSCredentials(ctx context.Context, creds *ststypes.Credentials) (app.AWSOrganizationsAPI, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if org, ok := f.orgs[*creds.AccessKeyId]; ok {
		return org, nil
	} else {
		if f.orgs == nil {
			f.orgs = map[string]*TestAWSOrganizationsAPI{}
		}
		org := &TestAWSOrganizationsAPI{}
		f.orgs[*creds.AccessKeyId] = org
		return org, nil
	}
}
