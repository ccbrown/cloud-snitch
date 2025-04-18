package app

import (
	"github.com/stripe/stripe-go/v81"

	"github.com/ccbrown/cloud-snitch/backend/store"
)

type Config struct {
	FrontendURL               string
	Store                     store.Config
	Email                     EmailConfig
	ContactEmailAddress       string
	PasswordEncryptionKey     []byte
	UserRegistrationAllowlist []string
	CloudFrontKeyId           string
	CloudFrontPrivateKey      string
	S3CDNURL                  string
	S3BucketName              string
	SQSQueueName              string
	AWSAccountId              string
	AWSRegions                []string
	StripeSecretKey           string
	Pricing                   PricingConfig

	// These can be overridden with mock implementations for testing.
	STS                  AWSSTSAPI
	SQSFactory           AmazonSQSAPIFactory
	S3                   AmazonS3API
	S3Factory            AmazonS3APIFactory
	OrganizationsFactory AWSOrganizationsAPIFactory
	IAMFactory           AWSIAMAPIFactory
	StripeAPIBackend     stripe.Backend
}

type PricingConfig struct {
	IndividualSubscriptionStripePriceId string
	TeamSubscriptionStripePriceId       string
}
