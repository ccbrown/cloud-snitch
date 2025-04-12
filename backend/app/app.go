package app

import (
	"context"
	"crypto"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"

	"github.com/ccbrown/cloud-snitch/backend/store"
)

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

func renderTemplateFunc(name string, params any) (template.HTML, error) {
	out := strings.Builder{}
	if err := templates.ExecuteTemplate(&out, name, params); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

func init() {
	templates = template.New("").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"render": renderTemplateFunc,
	})
	if _, err := templates.ParseFS(templatesFS, "templates/*.tmpl"); err != nil {
		panic(fmt.Sprintf("unable to parse templates: %s", err))
	}
}

type App struct {
	store                *store.Store
	emailer              Emailer
	config               Config
	webAuthn             *webauthn.WebAuthn
	organizationsFactory AWSOrganizationsAPIFactory
	sts                  AWSSTSAPI
	awsRegion            string
	sqs                  map[string]AmazonSQSAPI
	s3                   AmazonS3API
	s3Factory            AmazonS3APIFactory
	urlSigner            *sign.URLSigner
	stripe               *client.API
}

func New(cfg Config) (*App, error) {
	var stripeClient *client.API
	if cfg.StripeSecretKey != "" {
		var backends *stripe.Backends
		if cfg.StripeAPIBackend != nil {
			backends = &stripe.Backends{
				API: cfg.StripeAPIBackend,
			}
		}
		stripeClient = client.New(cfg.StripeSecretKey, backends)
	}

	store, err := store.New(cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize store: %w", err)
	}

	emailer, err := NewEmailer(cfg.Email)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize emailer: %w", err)
	}

	frontendURL, err := url.Parse(cfg.FrontendURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse frontend url: %w", err)
	}

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Cloud Snitch",
		RPID:          frontendURL.Hostname(),
		RPOrigins:     []string{cfg.FrontendURL},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce: true,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce: true,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize app webauthn: %w", err)
	}

	var awsConfig aws.Config
	if cfg.S3 == nil || cfg.STS == nil || cfg.SQSFactory == nil || cfg.OrganizationsFactory == nil {
		awsConfig, err = config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error loading default aws config: %w", err)
		}
		awsConfig.RetryMaxAttempts = 5
	}

	stsAPI := cfg.STS
	if stsAPI == nil {
		stsAPI = sts.NewFromConfig(awsConfig)
	}

	sqsFactory := cfg.SQSFactory
	if sqsFactory == nil {
		sqsFactory = LiveAmazonSQSAPIFactory{Config: awsConfig}
	}

	s3API := cfg.S3
	if s3API == nil {
		s3API = s3.NewFromConfig(awsConfig)
	}

	s3Factory := cfg.S3Factory
	if s3Factory == nil {
		s3Factory = LiveAmazonS3APIFactory{}
	}

	organizationsFactory := cfg.OrganizationsFactory
	if organizationsFactory == nil {
		organizationsFactory = LiveAWSOrganizationsAPIFactory{}
	}

	sqsAPI := make(map[string]AmazonSQSAPI, len(cfg.AWSRegions))
	for _, region := range cfg.AWSRegions {
		sqs, err := sqsFactory.NewWithRegion(context.Background(), region)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize sqs client for region %s: %w", region, err)
		}
		sqsAPI[region] = sqs
	}

	var urlSigner *sign.URLSigner
	if cfg.CloudFrontKeyId != "" && cfg.CloudFrontPrivateKey != "" {
		block, rest := pem.Decode([]byte(cfg.CloudFrontPrivateKey))
		if block == nil {
			return nil, fmt.Errorf("unable to decode cloudfront private key: %s", rest)
		}
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse cloudfront private key: %w", err)
		}
		urlSigner = sign.NewURLSigner(cfg.CloudFrontKeyId, key.(crypto.Signer))
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	return &App{
		store:                store,
		emailer:              emailer,
		config:               cfg,
		webAuthn:             webAuthn,
		organizationsFactory: organizationsFactory,
		awsRegion:            awsRegion,
		sts:                  stsAPI,
		sqs:                  sqsAPI,
		s3:                   s3API,
		s3Factory:            s3Factory,
		urlSigner:            urlSigner,
		stripe:               stripeClient,
	}, nil
}

func ValidateName(name string) UserFacingError {
	if len(name) == 0 {
		return NewUserError("A name is required.")
	}
	if len(name) > 250 {
		return NewUserError("Names must be 250 characters or less.")
	}
	return nil
}

type ContactUsInput struct {
	Name         string
	EmailAddress string
	Subject      string
	Message      string
}

func (input ContactUsInput) Validate() UserFacingError {
	if len(input.Name) == 0 {
		return NewUserError("A name is required.")
	} else if len(input.Name) > 250 {
		return NewUserError("Please provide a shorter name.")
	}

	if err := ValidateEmailAddress(input.EmailAddress); err != nil {
		return err
	}

	if len(input.Subject) == 0 {
		return NewUserError("A subject is required.")
	} else if len(input.Subject) > 500 {
		return NewUserError("Please provide a shorter subject.")
	}

	if len(input.Message) == 0 {
		return NewUserError("A message is required.")
	} else if len(input.Message) > 5000 {
		return NewUserError("Please provide a shorter message.")
	}

	return nil
}

func (s *Session) ContactUs(ctx context.Context, input ContactUsInput) UserFacingError {
	if err := input.Validate(); err != nil {
		return err
	}

	if err := s.app.Email(ctx, s.app.config.ContactEmailAddress, "Contact Form Submission: "+input.Subject, "contact_us_email.html.tmpl", map[string]any{
		"Name":         input.Name,
		"EmailAddress": input.EmailAddress,
		"Subject":      input.Subject,
		"Message":      input.Message,
	}); err != nil {
		return s.SanitizedError(fmt.Errorf("unable to send contact us email: %w", err))
	}

	return nil
}

func (a *App) Stripe() *client.API {
	return a.stripe
}

func (a *App) RenderTemplate(template string, data map[string]any) (string, error) {
	params := map[string]any{
		"FrontendURL": a.config.FrontendURL,
	}
	for k, v := range data {
		params[k] = v
	}

	html := strings.Builder{}
	if err := templates.ExecuteTemplate(&html, template, params); err != nil {
		return "", err
	}
	return html.String(), nil
}

func nilIfEmpty[T comparable](v T) *T {
	var empty T
	if v == empty {
		return nil
	}
	return &v
}

func emptyIfNil[T any](v *T) T {
	if v == nil {
		var empty T
		return empty
	}
	return *v
}
