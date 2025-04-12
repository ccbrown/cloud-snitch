package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/app/apptest"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestGenerateAWSCloudTrailReport(t *testing.T) {
	a := apptest.NewTestApp(t)

	_, sess := a.NewTestUser("alice@example.com", model.UserRoleCustomer)

	var err error

	team := a.NewTestTeamWithSubscription(sess, app.TeamSubscriptionTierIndividual)

	integration, err := sess.CreateAWSIntegration(context.Background(), app.CreateAWSIntegrationInput{
		TeamId:  team.Id,
		Name:    "My Integration",
		RoleARN: "arn:aws:iam::123456789012:role/MyRole",
		CloudTrailTrail: &app.CreateAWSIntegrationCloudTrailTrailInput{
			S3BucketName: "aws-cloudtrail-logs",
		},
	})
	require.NoError(t, err)

	report, err := a.GenerateAWSCloudTrailReport(context.Background(), app.GenerateAWSCloudTrailReportInput{
		FutureReportId:    model.NewReportId(),
		AWSIntegrationId:  integration.Id,
		StartTime:         time.Date(2025, 3, 6, 2, 25, 0, 0, time.UTC),
		Duration:          60 * time.Minute,
		AccountsKeyPrefix: "AWSLogs/o-1234abcde/",
		AccountId:         "222222222222",
		Region:            "us-east-1",
		Retention:         model.ReportRetentionOneWeek,
	})
	require.NoError(t, err)

	assert.Equal(t, model.ReportScope{
		StartTime: time.Date(2025, 3, 6, 2, 25, 0, 0, time.UTC),
		Duration:  60 * time.Minute,
		AWS: model.ReportScopeAWS{
			AccountId: "222222222222",
			Region:    "us-east-1",
		},
	}, report.Scope)

	reports, err := sess.GetReportsByTeamId(context.Background(), team.Id)
	require.NoError(t, err)
	require.Len(t, reports, 1)
	assert.Equal(t, report.Id, reports[0].Id)
}
