package app

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	geo "github.com/kellydunn/golang-geo"
)

type AWSIAMAPI interface {
	GenerateOrganizationsAccessReport(ctx context.Context, params *iam.GenerateOrganizationsAccessReportInput, optFns ...func(*iam.Options)) (*iam.GenerateOrganizationsAccessReportOutput, error)
	GetOrganizationsAccessReport(ctx context.Context, params *iam.GetOrganizationsAccessReportInput, optFns ...func(*iam.Options)) (*iam.GetOrganizationsAccessReportOutput, error)
}

type AWSIAMAPIFactory interface {
	NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials) (AWSIAMAPI, error)
}

type LiveAWSIAMAPIFactory struct{}

func (LiveAWSIAMAPIFactory) NewFromSTSCredentials(ctx context.Context, creds *ststypes.Credentials) (AWSIAMAPI, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)))
	if err != nil {
		return nil, err
	}
	// AWS IAM is a global service, but its control plane is in us-east-1. If us-east-1 is down,
	// there would be nothing we can do anyway, so we'll always us-east-1 in order to keep our
	// activity predictable.
	// See: https://docs.aws.amazon.com/whitepapers/latest/aws-fault-isolation-boundaries/global-services.html
	awsConfig.Region = "us-east-1"
	return iam.NewFromConfig(awsConfig), nil
}

type AWSOrganizationsAPI interface {
	ListAccounts(ctx context.Context, params *organizations.ListAccountsInput, optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error)
	ListParents(ctx context.Context, params *organizations.ListParentsInput, optFns ...func(*organizations.Options)) (*organizations.ListParentsOutput, error)
	ListPoliciesForTarget(ctx context.Context, params *organizations.ListPoliciesForTargetInput, optFns ...func(*organizations.Options)) (*organizations.ListPoliciesForTargetOutput, error)
	DescribePolicy(ctx context.Context, params *organizations.DescribePolicyInput, optFns ...func(*organizations.Options)) (*organizations.DescribePolicyOutput, error)
	AttachPolicy(ctx context.Context, params *organizations.AttachPolicyInput, optFns ...func(*organizations.Options)) (*organizations.AttachPolicyOutput, error)
	CreatePolicy(ctx context.Context, params *organizations.CreatePolicyInput, optFns ...func(*organizations.Options)) (*organizations.CreatePolicyOutput, error)
	UpdatePolicy(ctx context.Context, params *organizations.UpdatePolicyInput, optFns ...func(*organizations.Options)) (*organizations.UpdatePolicyOutput, error)
	ListRoots(ctx context.Context, params *organizations.ListRootsInput, optFns ...func(*organizations.Options)) (*organizations.ListRootsOutput, error)
}

type AWSOrganizationsAPIFactory interface {
	NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials) (AWSOrganizationsAPI, error)
}

type LiveAWSOrganizationsAPIFactory struct{}

func (LiveAWSOrganizationsAPIFactory) NewFromSTSCredentials(ctx context.Context, creds *ststypes.Credentials) (AWSOrganizationsAPI, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)))
	if err != nil {
		return nil, err
	}
	// AWS Organizations is a global service, but its control plane is in us-east-1. If us-east-1 is
	// down, there would be nothing we can do anyway, so we'll always us-east-1 in order to keep our
	// activity predictable.
	// See: https://docs.aws.amazon.com/whitepapers/latest/aws-fault-isolation-boundaries/global-services.html
	awsConfig.Region = "us-east-1"
	return organizations.NewFromConfig(awsConfig), nil
}

type AmazonS3API interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type AmazonS3APIFactory interface {
	GetBucketRegion(ctx context.Context, bucketName string) (string, error)
	NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials, region string) (AmazonS3API, error)
}

type LiveAmazonS3APIFactory struct{}

func (LiveAmazonS3APIFactory) GetBucketRegion(ctx context.Context, bucketName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+bucketName+".s3.amazonaws.com", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return "", err
	}
	if region := resp.Header.Get("x-amz-bucket-region"); region != "" {
		return region, nil
	}
	return "", fmt.Errorf("unable to determine bucket region")
}

func (LiveAmazonS3APIFactory) NewFromSTSCredentials(ctx context.Context, creds *ststypes.Credentials, region string) (AmazonS3API, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)))
	if err != nil {
		return nil, err
	}
	if region != "" {
		awsConfig.Region = region
	}
	return s3.NewFromConfig(awsConfig), nil
}

type AWSSTSAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

type AmazonSQSAPI interface {
	SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error)
}

type AmazonSQSAPIFactory interface {
	NewWithRegion(ctx context.Context, region string) (AmazonSQSAPI, error)
}

type LiveAmazonSQSAPIFactory struct {
	Config aws.Config
}

func (f LiveAmazonSQSAPIFactory) NewWithRegion(ctx context.Context, region string) (AmazonSQSAPI, error) {
	config := f.Config.Copy()
	config.Region = region
	return sqs.NewFromConfig(config), nil
}

type AWSRegion struct {
	Name               string
	GeolocationCountry string
	GeolocationRegion  string
	Partition          string
	Latitude           float64
	Longitude          float64
}

// Metadata on known AWS regions. It should be expected that this is not exhaustive as AWS regularly
// adds new regions.
var KnownAWSRegions = map[string]AWSRegion{
	// XXX: DO NOT EDIT MANUALLY. Data is generated by ../scripts/gather_aws_info.py.
	"ap-northeast-1": {"Asia Pacific (Tokyo)", "JP", "JP-13", "aws", 35.41, 139.42},
	"ap-southeast-5": {"Asia Pacific (Malaysia)", "MY", "MY-14", "aws", 4.2105, 101.9758},
	"ap-southeast-7": {"Asia Pacific (Thailand)", "TH", "TH-10", "aws", 15.87, 100.9925},
	"eu-central-1":   {"Europe (Frankfurt)", "DE", "DE-HE", "aws", 50, 8},
	"eu-north-1":     {"Europe (Stockholm)", "SE", "SE-AB", "aws", 59.25, 17.81},
	"me-south-1":     {"Middle East (Bahrain)", "BH", "BH-14", "aws", 26.1, 50.46},
	"sa-east-1":      {"South America (Sao Paulo)", "BR", "BR-SP", "aws", -23.34, -46.38},
	"us-gov-east-1":  {"AWS GovCloud (US-East)", "US", "US-OH", "aws-us-gov", 38.944, -77.455},
	"us-gov-west-1":  {"AWS GovCloud (US-West)", "US", "US-OR", "aws-us-gov", 37.618, -122.375},
	"us-west-1":      {"US West (N. California)", "US", "US-CA", "aws", 37.35, -121.96},
	"ap-northeast-2": {"Asia Pacific (Seoul)", "KR", "KR-28", "aws", 37.56, 126.98},
	"ap-northeast-3": {"Asia Pacific (Osaka)", "JP", "JP-27", "aws", 34.69, 135.49},
	"ap-south-1":     {"Asia Pacific (Mumbai)", "IN", "IN-MH", "aws", 19.08, 72.88},
	"ap-southeast-3": {"Asia Pacific (Jakarta)", "ID", "ID-JK", "aws", -6.125, 106.655},
	"cn-north-1":     {"China (Beijing)", "CN", "CN-BJ", "aws-cn", 40.08, 116.584},
	"eu-south-1":     {"Europe (Milan)", "IT", "IT-MI", "aws", 45.43, 9.29},
	"eu-west-1":      {"Europe (Ireland)", "IE", "IE-D", "aws", 53, -8},
	"eu-west-3":      {"Europe (Paris)", "FR", "FR-75C", "aws", 48.86, 2.35},
	"il-central-1":   {"Israel (Tel Aviv)", "IL", "IL-TA", "aws", 32.0853, 34.7818},
	"us-east-2":      {"US East (Ohio)", "US", "US-OH", "aws", 39.96, -83},
	"af-south-1":     {"Africa (Cape Town)", "ZA", "ZA-WC", "aws", -33.93, 18.42},
	"ap-east-1":      {"Asia Pacific (Hong Kong)", "CN", "CN-HK", "aws", 22.27, 114.16},
	"ap-south-2":     {"Asia Pacific (Hyderabad)", "IN", "IN-TG", "aws", 17.4065, 78.4772},
	"ap-southeast-1": {"Asia Pacific (Singapore)", "SG", "SG-01", "aws", 1.37, 103.8},
	"ap-southeast-2": {"Asia Pacific (Sydney)", "AU", "AU-NSW", "aws", -33.86, 151.2},
	"ca-west-1":      {"Canada West (Calgary)", "CA", "CA-AB", "aws", 51.0447, -114.0719},
	"cn-northwest-1": {"China (Ningxia)", "CN", "CN-NX", "aws-cn", 38.321667, 106.3925},
	"eu-central-2":   {"Europe (Zurich)", "CH", "CH-ZH", "aws", 47.3769, 8.5417},
	"me-central-1":   {"Middle East (UAE)", "AE", "AE-DU", "aws", 23.4241, 53.8478},
	"us-west-2":      {"US West (Oregon)", "US", "US-OR", "aws", 46.15, -123.88},
	"ap-southeast-4": {"Asia Pacific (Melbourne)", "AU", "AU-VIC", "aws", -37.8136, 144.9631},
	"ca-central-1":   {"Canada (Central)", "CA", "CA-QC", "aws", 45.5, -73.6},
	"eu-south-2":     {"Europe (Spain)", "ES", "ES-AR", "aws", 40.4637, -3.7492},
	"eu-west-2":      {"Europe (London)", "GB", "GB-LND", "aws", 51, -0.1},
	"mx-central-1":   {"Mexico (Central)", "MX", "MX-QUE", "aws", 19.4326, -99.1332},
	"us-east-1":      {"US East (N. Virginia)", "US", "US-VA", "aws", 38.13, -78.45},
}

// Returns the AWS region from `AWSRegions` whose identifier appears to be most similar to the given
// region.
//
// Returns "" if no good match is found.
func MostSimilarKnownAWSRegion(region string) string {
	best := ""
	bestCommonPrefixLength := -1

	// Find the region with the longest common prefix with the given region.
	for known := range KnownAWSRegions {
		commonPrefixLength := 0
		for i := 0; i < len(known) && i < len(region); i++ {
			if known[i] != region[i] {
				break
			}
			commonPrefixLength++
		}
		if commonPrefixLength > bestCommonPrefixLength {
			best = known
			bestCommonPrefixLength = commonPrefixLength
		}
	}

	return best
}

func (a *App) ClosestAvailableAWSRegion(region string) string {
	return ClosestAvailableAWSRegion(region, a.config.AWSRegions)
}

func ClosestAvailableAWSRegion(region string, availableRegions []string) string {
	region = MostSimilarKnownAWSRegion(region)
	if region == "" {
		region = "us-east-1"
	}

	info := KnownAWSRegions[region]
	loc := geo.NewPoint(info.Latitude, info.Longitude)

	ret := availableRegions[0]
	closestDistance := math.MaxFloat64
	for _, otherId := range availableRegions {
		if otherRegion, ok := KnownAWSRegions[otherId]; ok {
			dist := loc.GreatCircleDistance(geo.NewPoint(otherRegion.Latitude, otherRegion.Longitude))
			if dist < closestDistance {
				ret = otherId
				closestDistance = dist
			}
		}
	}
	return ret
}
