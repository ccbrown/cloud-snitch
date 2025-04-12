package apptest

import (
	"context"
	"encoding/base64"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store/storetest"
)

const testFrontendURL = "http://localhost:3000"

type TestApp struct {
	*app.App
	T          *testing.T
	sqsFactory *TestAmazonSQSAPIFactory
}

func NewTestApp(t *testing.T) *TestApp {
	sqsFactory := &TestAmazonSQSAPIFactory{}
	cfg := app.Config{
		FrontendURL:           testFrontendURL,
		PasswordEncryptionKey: []byte("12345678901234567890123456789012"),
		Store:                 storetest.NewStoreConfig(t),
		STS:                   &TestAWSSTSAPI{},
		S3:                    &TestAmazonS3API{},
		S3Factory:             &TestAmazonS3APIFactory{},
		SQSFactory:            sqsFactory,
		StripeSecretKey:       "sk_test_12345678901234567890123456789012",
		Pricing: app.PricingConfig{
			IndividualSubscriptionStripePriceId: DummyStripePriceIndividualSubscription.ID,
			TeamSubscriptionStripePriceId:       DummyStripePriceTeamSubscription.ID,
		},
		StripeAPIBackend: &MockStripeBackend{},
		AWSRegions:       []string{"us-west-2", "us-east-1"},
		CloudFrontKeyId:  "MyTestKeyId",
		CloudFrontPrivateKey: `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDQPDski30uMMOg
xwNJTdgvJ6ok7GLYsYNj+x5LG+bKqdim+6N4yBp+9H8YatyGQ/qiIoNE5a3jIJL7
c817H9Z1I1ssgF222ucjXTe8HPER8nPMlVOIP5n7ZSnrIAQdUXhLQJkTy8NBchSz
cFAXaA951pTgztt+2MTfzzCer80bHipH0ZJahGYR0I/a6GygS62rWDw7hT/Diy/X
Vu6rL4o4hzO6QDsjP5da7fGbxmueLOXnX7YFu1eVbzJ+FLKXoarCAJZXjoScdu3e
gIlYCXGr5GW36dLfaVQpAzzrQpmQBFc2MbXG0LSfDoJ1vrItPm1aYU1OO+BSfJ0n
7nyrGrjvAgMBAAECggEAYEDvH6ngn7jHvKoxUTGL8+QUSEQCPuLy9oSt0I8ho91V
uX4I5rwsTzHKC+ndbKzAjjCA2BiIw7ubZWL2gOrLEVNaAhyF9Q+DlvuwVyJTpnWZ
ZGBD/+9SSHvPIBGpBTpS7gn6mEVwSHCos/b+9orR2IJBSBcmK6CchE27ziY6G0WF
Vytfg6yKpvXtlXZ5HB2SW7rtaDi855own4BgrlK7ZYmJmvdfoxiAW2Crp9GZUzKG
lYdZahh7NPXjeZ3SrQeL/kwaoVnMsfguowbE5P8IJLs29S9Hn+hLlPBnCFbA3E+a
Wwxn7WfCuA02bAvGLtDbx4Khv+jOjCdYeV/xSLH1DQKBgQDnnbFL8u4nzRC0wKZe
sCrPBkrR7J29UGLbSHIwdC/crWvJV6zCGVVspjqOSQZsEhM7DiyM1nMHQnf+soDN
b3aBlazsm5xoVenqcgbS+zVyg1wIVv00PY32wqkFsPHbhszfyIMlCoCyQCls1Cwr
kKr5TRVEves/fzT4fvbaxxeUgwKBgQDmKGcGzw5nlF7waoJb7r7EKGcdFbjDgRzt
r8vQ9xQl0+j6Ay1AtgyZc+ezlV1H2T03h8HXnpHwBEOLzfd70hnai5pZYryp8+vc
1ypxDnhSosxJCv754YDydpT/rfkQqFSG5FGHLwho9pfbvwVwrnC5e9DB3Xxi1w7m
clA6TSUWJQKBgQC2wM+HbNyLnRvEc4oShpCQn5djwn7IROhru+MV5mdpiZDV4o8W
7CRiQVtMr8QYe76ex1VRn1hN7x19Y12MO5nXL8GtRC+Kh9e1PFm0GbEXdRymG/VY
RgWBIiF5sx9zJw76uFz6Wion+8Zz74oWqeSuJytT/ULk7Dtgo9Wl1Jf/fwKBgG23
0PVz+3/xQRkHDUHaHGLSh+Vbl4rLoAjHBziEsLwfuy6EGSyCHJuCn3ACwkcgDojM
VOH/G775qWGaPGJwlBaU+m2mHh4+w6+xevAOX9m09gHUFhz1HU14rir8uoDwXes4
KI1zJGU1OMtu0p8W6XHizm+8sDFGKDo/QHOqlwVpAoGAH/6OvKhplozITHi7EpNL
7OTW1jMyF5XAuTJhCTwbWZr66eTicfpScWByvxg5u7+ooHs66ZFQaHgyjfG1/lSk
rSWtxxB19xvfLgkTclF4LJanKhKR+G6JuMRdtBHAoB/iQ6Thguuqy9X9QZYhZfyc
2qfIl5IsK7gjgYtHT0JkzXg=
-----END PRIVATE KEY-----`,
		S3CDNURL:     testFrontendURL,
		S3BucketName: "MyTestBucket",
		SQSQueueName: "MyTestQueue",
	}
	a, err := app.New(cfg)
	require.NoError(t, err)
	return &TestApp{
		App:        a,
		T:          t,
		sqsFactory: sqsFactory,
	}
}

var userRegistrationTokenRegexp = regexp.MustCompile(`token=([a-zA-Z0-9_\-]+)`)

func (a *TestApp) NewTestUser(emailAddress string, role model.UserRole) (*model.User, *app.Session) {
	var err error

	err = a.NewAnonymousSession().BeginUserRegistration(context.Background(), app.BeginUserRegistrationInput{
		EmailAddress:                    emailAddress,
		TermsOfServiceAgreementRevision: "2024.01.01",
		PrivacyPolicyAgreementRevision:  "2024.01.01",
		CookiePolicyAgreementRevision:   "2024.01.01",
	})
	require.NoError(a.T, err)

	var email app.Email
	for {
		email = <-a.Emails()
		if strings.Contains(email.Subject, "Registration") {
			break
		}
	}
	token := userRegistrationTokenRegexp.FindStringSubmatch(email.HTML)[1]
	tokenBytes, err := base64.RawURLEncoding.DecodeString(token)
	require.NoError(a.T, err)

	user, accessToken, err := a.NewAnonymousSession().CompleteUserRegistration(context.Background(), app.CompleteUserRegistrationInput{
		Token: tokenBytes,
	})
	require.NoError(a.T, err)

	patched, err := a.SetUserRole(context.Background(), user.Id, role)
	require.NoError(a.T, err)
	require.NotNil(a.T, patched)

	sess, err := a.NewAnonymousSession().WithUserAccessToken(context.Background(), accessToken)
	require.NoError(a.T, err)
	require.NotNil(a.T, sess)

	return user, sess
}

func (a *TestApp) NewTestTeamWithSubscription(sess *app.Session, tier app.TeamSubscriptionTier) *model.Team {
	var err error

	team, err := sess.CreateTeam(context.Background(), app.CreateTeamInput{
		Name: "Test Team",
	})
	require.NoError(a.T, err)

	_, err = sess.CreateTeamBillingProfileById(context.Background(), team.Id, app.CreateTeamBillingProfileInput{
		Name: "Test Team",
		Address: model.TeamBillingAddress{
			Line1:      aws.String("123 Main St"),
			City:       aws.String("Seattle"),
			State:      aws.String("WA"),
			PostalCode: "98101",
			Country:    "US",
		},
	})
	require.NoError(a.T, err)

	_, err = sess.PutTeamPaymentMethodById(context.Background(), team.Id, app.PutTeamPaymentMethodInput{
		StripePaymentMethodId: DummyStripeCard.ID,
	})
	require.NoError(a.T, err)

	_, err = sess.CreateTeamSubscriptionById(context.Background(), team.Id, app.CreateTeamSubscriptionInput{
		Tier: tier,
	})
	require.NoError(a.T, err)

	return team
}

func (a *TestApp) SQSRequests(region string) []*sqs.SendMessageBatchInput {
	return a.sqsFactory.Requests(region)
}
