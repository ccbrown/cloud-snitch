package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestAPI_AWSIntegration(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)

	t.Run("CreateWithoutExternalId", func(t *testing.T) {
		_, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
			TeamId: team.Id.String(),
			Body: &apispec.CreateAWSIntegrationJSONRequestBody{
				Name:    "Foo",
				RoleArn: "arn:aws:iam::123456789012:role/MyRoleNoExternalIdRequired",
				CloudtrailTrail: &apispec.CreateAWSIntegrationCloudTrailTrailInput{
					S3BucketName: "my-bucket",
				},
			},
		})
		assert.Error(t, err)
	})

	resp, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
		TeamId: team.Id.String(),
		Body: &apispec.CreateAWSIntegrationJSONRequestBody{
			Name:    "Foo",
			RoleArn: "arn:aws:iam::123456789012:role/MyRole",
			CloudtrailTrail: &apispec.CreateAWSIntegrationCloudTrailTrailInput{
				S3BucketName: "my-bucket",
			},
		},
	})
	require.NoError(t, err)
	integration := resp.(apispec.CreateAWSIntegration200JSONResponse)
	assert.Equal(t, "Foo", integration.Name)

	t.Run("CreateDuplicate", func(t *testing.T) {
		_, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
			TeamId: team.Id.String(),
			Body: &apispec.CreateAWSIntegrationJSONRequestBody{
				Name:    "Foo",
				RoleArn: "arn:aws:iam::123456789012:role/MyRole",
				CloudtrailTrail: &apispec.CreateAWSIntegrationCloudTrailTrailInput{
					S3BucketName: "my-bucket",
				},
			},
		})
		assert.Error(t, err)
	})

	t.Run("GetTeamAWSIntegrations", func(t *testing.T) {
		resp, err := api.GetAWSIntegrationsByTeamId(aliceCtx, apispec.GetAWSIntegrationsByTeamIdRequestObject{
			TeamId: team.Id.String(),
		})
		require.NoError(t, err)
		integrations := resp.(apispec.GetAWSIntegrationsByTeamId200JSONResponse)
		assert.Len(t, integrations, 1)
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := api.DeleteAWSIntegration(aliceCtx, apispec.DeleteAWSIntegrationRequestObject{
			IntegrationId: integration.Id,
		})
		require.NoError(t, err)
	})
}

func TestAPI_CreateAWSIntegration_WithBackfill(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)
	ignoreSQSRequests := len(api.app.SQSRequests("us-east-1"))

	createResp, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
		TeamId: team.Id.String(),
		Body: &apispec.CreateAWSIntegrationJSONRequestBody{
			Name:    "Foo",
			RoleArn: "arn:aws:iam::123456789012:role/MyRole",
			CloudtrailTrail: &apispec.CreateAWSIntegrationCloudTrailTrailInput{
				S3BucketName: "aws-cloudtrail-logs",
			},
			QueueReportGeneration: pointer(true),
		},
	})
	require.NoError(t, err)
	integration := createResp.(apispec.CreateAWSIntegration200JSONResponse)

	sqsRequests := api.app.SQSRequests("us-east-1")[ignoreSQSRequests:]
	require.Len(t, sqsRequests, 7)
	sqsRequest := sqsRequests[0]
	assert.Len(t, sqsRequest.Entries, 2)

	t.Run("AWSAccounts", func(t *testing.T) {
		resp, err := api.GetAWSAccountsByTeamId(aliceCtx, apispec.GetAWSAccountsByTeamIdRequestObject{
			TeamId: team.Id.String(),
		})
		require.NoError(t, err)
		accounts := resp.(apispec.GetAWSAccountsByTeamId200JSONResponse)
		assert.Len(t, accounts, 2)
	})

	{
		_, err := api.DeleteAWSIntegration(aliceCtx, apispec.DeleteAWSIntegrationRequestObject{
			IntegrationId: integration.Id,
			Body: &apispec.DeleteAWSIntegrationJSONRequestBody{
				DeleteAssociatedData: pointer(true),
			},
		})
		require.NoError(t, err)
	}

	t.Run("AWSAccountsAfterDelete", func(t *testing.T) {
		resp, err := api.GetAWSAccountsByTeamId(aliceCtx, apispec.GetAWSAccountsByTeamIdRequestObject{
			TeamId: team.Id.String(),
		})
		require.NoError(t, err)
		accounts := resp.(apispec.GetAWSAccountsByTeamId200JSONResponse)
		assert.Len(t, accounts, 0)
	})
}

func TestAPI_AWSIntegration_SCP_Management(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)

	{
		resp, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
			TeamId: team.Id.String(),
			Body: &apispec.CreateAWSIntegrationJSONRequestBody{
				Name:                             "Foo",
				RoleArn:                          "arn:aws:iam::123456789012:role/MyRole",
				GetAccountNamesFromOrganizations: pointer(true),
				ManageScps:                       pointer(true),
			},
		})
		require.NoError(t, err)
		integration := resp.(apispec.CreateAWSIntegration200JSONResponse)
		assert.Equal(t, "Foo", integration.Name)
	}

	t.Run("NoSCP", func(t *testing.T) {
		_, err := api.GetManagedAWSSCP(aliceCtx, apispec.GetManagedAWSSCPRequestObject{
			TeamId:    team.Id.String(),
			AccountId: "123456789012",
		})
		require.Error(t, err)
	})

	t.Run("CreateSCP", func(t *testing.T) {
		resp, err := api.PutManagedAWSSCP(aliceCtx, apispec.PutManagedAWSSCPRequestObject{
			TeamId:    team.Id.String(),
			AccountId: "123456789012",
			Body: &apispec.PutManagedAWSSCPJSONRequestBody{
				Content: "foo",
			},
		})
		require.NoError(t, err)
		scp := resp.(apispec.PutManagedAWSSCP200JSONResponse)
		assert.Equal(t, "foo", scp.Content)

		t.Run("GetSCP", func(t *testing.T) {
			resp, err := api.GetManagedAWSSCP(aliceCtx, apispec.GetManagedAWSSCPRequestObject{
				TeamId:    team.Id.String(),
				AccountId: "123456789012",
			})
			require.NoError(t, err)
			scp := resp.(apispec.GetManagedAWSSCP200JSONResponse)
			assert.Equal(t, "foo", scp.Content)
		})

		t.Run("UpdateSCP", func(t *testing.T) {
			resp, err := api.PutManagedAWSSCP(aliceCtx, apispec.PutManagedAWSSCPRequestObject{
				TeamId:    team.Id.String(),
				AccountId: "123456789012",
				Body: &apispec.PutManagedAWSSCPJSONRequestBody{
					Content: "bar",
				},
			})
			require.NoError(t, err)
			scp := resp.(apispec.PutManagedAWSSCP200JSONResponse)
			assert.Equal(t, "bar", scp.Content)

			t.Run("GetSCP", func(t *testing.T) {
				resp, err := api.GetManagedAWSSCP(aliceCtx, apispec.GetManagedAWSSCPRequestObject{
					TeamId:    team.Id.String(),
					AccountId: "123456789012",
				})
				require.NoError(t, err)
				scp := resp.(apispec.GetManagedAWSSCP200JSONResponse)
				assert.Equal(t, "bar", scp.Content)
			})
		})
	})

	t.Run("AccessReport", func(t *testing.T) {
		resp, err := api.GetAWSAccessReport(aliceCtx, apispec.GetAWSAccessReportRequestObject{
			TeamId:    team.Id.String(),
			AccountId: "123456789012",
		})
		require.NoError(t, err)
		report := resp.(apispec.GetAWSAccessReport200JSONResponse)
		assert.NotEmpty(t, report.Services)
	})
}
