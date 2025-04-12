package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	DynamoDB DynamoDBConfig
}

type DynamoDBConfig struct {
	TableName string

	// These overrides are typically only used for tests or local environments.
	Endpoint          string
	Region            string
	StaticCredentials *DynamoDBStaticCredentials
}

type DynamoDBStaticCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func (cfg *DynamoDBConfig) AWSConfig() (aws.Config, error) {
	ret, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return ret, fmt.Errorf("error loading default config: %w", err)
	}

	ret.RetryMaxAttempts = 5

	if cfg.Region != "" {
		ret.Region = cfg.Region
	}

	if creds := cfg.StaticCredentials; creds != nil {
		ret.Credentials = aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     creds.AccessKeyId,
				SecretAccessKey: creds.SecretAccessKey,
			}, nil
		})
	}

	if cfg.Endpoint != "" {
		ret.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		})
	}

	return ret, nil
}
