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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"

	"github.com/ccbrown/cloud-snitch/backend/app"
)

type TestAWSSTSAPI struct{}

func (api *TestAWSSTSAPI) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	if !strings.Contains(*params.RoleArn, "NoExternalIdRequired") && params.ExternalId == nil {
		return nil, fmt.Errorf("external id required")
	}
	return &sts.AssumeRoleOutput{
		Credentials: &ststypes.Credentials{},
	}, nil
}

// Paths are relative to the root of the Go module.
var bucketPaths = map[string]string{
	"aws-cloudtrail-logs": "report/testdata/aws-cloudtrail-logs",
}

type TestAmazonS3API struct {
}

func (api *TestAmazonS3API) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return &s3.GetBucketLocationOutput{}, nil
}

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

func (f TestAmazonS3APIFactory) NewFromSTSCredentials(ctx context.Context, credentials *ststypes.Credentials) (app.AmazonS3API, error) {
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
