package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/jaytaylor/html2text"
)

type Emailer interface {
	Email(ctx context.Context, to, subject, html string) error
}

type EmailConfig struct {
	SES *EmailSESConfig
}

type EmailSESConfig struct {
	FromAddress string
	Region      string
}

func NewEmailer(e EmailConfig) (Emailer, error) {
	if cfg := e.SES; cfg != nil {
		awsConfig, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error loading default aws config: %w", err)
		}

		awsConfig.RetryMaxAttempts = 5

		if cfg.Region != "" {
			awsConfig.Region = cfg.Region
		}

		return &AmazonSESEmailer{
			client:      ses.NewFromConfig(awsConfig),
			fromAddress: cfg.FromAddress,
		}, nil
	} else {
		return &channelEmailer{
			Emails: make(chan Email, 100),
		}, nil
	}
}

type AmazonSESEmailer struct {
	client *ses.Client

	fromAddress string
}

func (e *AmazonSESEmailer) Email(ctx context.Context, to, subject, html string) error {
	text, err := html2text.FromString(html)
	if err != nil {
		return err
	}

	_, err = e.client.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    &html,
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    &text,
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    &subject,
			},
		},
		Source: aws.String("Cloud Snitch <" + e.fromAddress + ">"),
	})
	return err
}

type Email struct {
	To      string
	Subject string
	HTML    string
}

type channelEmailer struct {
	Emails chan Email
}

// If a real emailer is not configured, this will return a channel that can be used to read emails
// that would have been sent.
func (a *App) Emails() <-chan Email {
	return a.emailer.(*channelEmailer).Emails
}

func (e *channelEmailer) Email(ctx context.Context, to, subject, html string) error {
	email := Email{
		To:      to,
		Subject: subject,
		HTML:    html,
	}
	select {
	case e.Emails <- email:
	default:
	}
	return nil
}

func (a *App) Email(ctx context.Context, to, subject, template string, data map[string]any) error {
	html, err := a.RenderTemplate(template, data)
	if err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}
	return a.emailer.Email(ctx, to, subject, html)
}
