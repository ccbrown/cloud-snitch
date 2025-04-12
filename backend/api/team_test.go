package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestAPI_Team(t *testing.T) {
	api := NewTestAPI(t)
	_, adminCtx := api.NewTestUser("admin@example.com", model.UserRoleAdministrator)
	alice, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)
	bob, bobCtx := api.NewTestUser("bob@example.com", model.UserRoleCustomer)

	resp, err := api.CreateTeam(aliceCtx, apispec.CreateTeamRequestObject{
		Body: &apispec.CreateTeamJSONRequestBody{
			Name: "Foo",
		},
	})
	require.NoError(t, err)
	team := resp.(apispec.CreateTeam200JSONResponse)
	assert.Equal(t, "Foo", team.Name)

	t.Run("GetTeamById", func(t *testing.T) {
		t.Run("Member", func(t *testing.T) {
			resp, err := api.GetTeam(aliceCtx, apispec.GetTeamRequestObject{TeamId: team.Id})
			require.NoError(t, err)
			assert.Equal(t, team.Id, resp.(apispec.GetTeam200JSONResponse).Id)
		})

		t.Run("Admin", func(t *testing.T) {
			resp, err := api.GetTeam(adminCtx, apispec.GetTeamRequestObject{TeamId: team.Id})
			require.NoError(t, err)
			assert.Equal(t, team.Id, resp.(apispec.GetTeam200JSONResponse).Id)
		})

		t.Run("Nonmember", func(t *testing.T) {
			_, err := api.GetTeam(bobCtx, apispec.GetTeamRequestObject{TeamId: team.Id})
			require.Error(t, err)
		})
	})

	t.Run("UpdateTeam", func(t *testing.T) {
		t.Run("Member", func(t *testing.T) {
			resp, err := api.UpdateTeam(aliceCtx, apispec.UpdateTeamRequestObject{
				TeamId: team.Id,
				Body: &apispec.UpdateTeamJSONRequestBody{
					Name: pointer("Bar"),
				},
			})
			require.NoError(t, err)
			assert.Equal(t, "Bar", resp.(apispec.UpdateTeam200JSONResponse).Name)
		})

		t.Run("Nonmember", func(t *testing.T) {
			_, err := api.UpdateTeam(bobCtx, apispec.UpdateTeamRequestObject{
				TeamId: team.Id,
				Body: &apispec.UpdateTeamJSONRequestBody{
					Name: pointer("Bar"),
				},
			})
			require.Error(t, err)
		})
	})

	t.Run("GetTeams", func(t *testing.T) {
		_, err := api.GetTeams(aliceCtx, apispec.GetTeamsRequestObject{})
		require.Error(t, err)

		resp, err := api.GetTeams(adminCtx, apispec.GetTeamsRequestObject{})
		require.NoError(t, err)
		teams := resp.(apispec.GetTeams200JSONResponse)
		assert.Len(t, teams, 1)
	})

	t.Run("JoinWithoutInvite", func(t *testing.T) {
		_, err = api.JoinTeam(bobCtx, apispec.JoinTeamRequestObject{
			TeamId: team.Id,
		})
		require.Error(t, err)
	})

	t.Run("Invite", func(t *testing.T) {
		t.Run("IndividualTeam", func(t *testing.T) {
			team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)

			_, err := api.CreateTeamInvite(aliceCtx, apispec.CreateTeamInviteRequestObject{
				TeamId: team.Id.String(),
				Body: &apispec.CreateTeamInviteJSONRequestBody{
					EmailAddress: "bob@example.com",
					Role:         apispec.TeamMembershipRoleMEMBER,
				},
			})
			require.Error(t, err)
		})

		team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierTeam)

		_, err := api.CreateTeamInvite(aliceCtx, apispec.CreateTeamInviteRequestObject{
			TeamId: team.Id.String(),
			Body: &apispec.CreateTeamInviteJSONRequestBody{
				EmailAddress: "bob@example.com",
				Role:         apispec.TeamMembershipRoleMEMBER,
			},
		})
		require.NoError(t, err)

		email := <-api.app.Emails()
		assert.Contains(t, email.Subject, "Invite")

		t.Run("Nonadmin", func(t *testing.T) {
			_, err := api.CreateTeamInvite(bobCtx, apispec.CreateTeamInviteRequestObject{
				TeamId: team.Id.String(),
				Body: &apispec.CreateTeamInviteJSONRequestBody{
					EmailAddress: "foo@example.com",
					Role:         apispec.TeamMembershipRoleMEMBER,
				},
			})
			require.Error(t, err)
		})

		t.Run("GetByTeamId", func(t *testing.T) {
			resp, err := api.GetTeamInvitesByTeamId(aliceCtx, apispec.GetTeamInvitesByTeamIdRequestObject{
				TeamId: team.Id.String(),
			})
			require.NoError(t, err)
			invites := resp.(apispec.GetTeamInvitesByTeamId200JSONResponse)
			assert.Len(t, invites, 1)
		})

		t.Run("GetByUserId", func(t *testing.T) {
			t.Run("OtherUser", func(t *testing.T) {
				_, err := api.GetTeamInvitesByUserId(aliceCtx, apispec.GetTeamInvitesByUserIdRequestObject{
					UserId: bob.Id.String(),
				})
				require.Error(t, err)
			})

			t.Run("Self", func(t *testing.T) {
				resp, err := api.GetTeamInvitesByUserId(bobCtx, apispec.GetTeamInvitesByUserIdRequestObject{
					UserId: bob.Id.String(),
				})
				require.NoError(t, err)
				invites := resp.(apispec.GetTeamInvitesByUserId200JSONResponse)
				assert.Len(t, invites, 1)
				assert.Equal(t, "alice@example.com", invites[0].Sender.EmailAddress)
			})
		})

		t.Run("Delete", func(t *testing.T) {
			_, err := api.DeleteTeamInvite(bobCtx, apispec.DeleteTeamInviteRequestObject{
				TeamId:       team.Id.String(),
				EmailAddress: "bob@example.com",
			})
			require.NoError(t, err)

			resp, err := api.GetTeamInvitesByUserId(bobCtx, apispec.GetTeamInvitesByUserIdRequestObject{
				UserId: bob.Id.String(),
			})
			require.NoError(t, err)
			invites := resp.(apispec.GetTeamInvitesByUserId200JSONResponse)
			assert.Empty(t, invites)

			_, err = api.CreateTeamInvite(aliceCtx, apispec.CreateTeamInviteRequestObject{
				TeamId: team.Id.String(),
				Body: &apispec.CreateTeamInviteJSONRequestBody{
					EmailAddress: "bob@example.com",
					Role:         apispec.TeamMembershipRoleMEMBER,
				},
			})
			require.NoError(t, err)

			email := <-api.app.Emails()
			assert.Contains(t, email.Subject, "Invite")
		})

		_, err = api.JoinTeam(bobCtx, apispec.JoinTeamRequestObject{
			TeamId: team.Id.String(),
		})
		require.NoError(t, err)

		t.Run("GetTeamMembershipsByTeamId", func(t *testing.T) {
			t.Run("Nonadmin", func(t *testing.T) {
				_, err := api.GetTeamMembershipsByTeamId(bobCtx, apispec.GetTeamMembershipsByTeamIdRequestObject{
					TeamId: team.Id.String(),
				})
				require.Error(t, err)
			})

			t.Run("Admin", func(t *testing.T) {
				resp, err := api.GetTeamMembershipsByTeamId(aliceCtx, apispec.GetTeamMembershipsByTeamIdRequestObject{
					TeamId: team.Id.String(),
				})
				require.NoError(t, err)
				memberships := resp.(apispec.GetTeamMembershipsByTeamId200JSONResponse)
				assert.Len(t, memberships, 2)
			})
		})

		t.Run("GetTeamMembershipsByUserId", func(t *testing.T) {
			t.Run("OtherUser", func(t *testing.T) {
				_, err := api.GetTeamMembershipsByUserId(aliceCtx, apispec.GetTeamMembershipsByUserIdRequestObject{
					UserId: bob.Id.String(),
				})
				require.Error(t, err)
			})

			t.Run("Self", func(t *testing.T) {
				resp, err := api.GetTeamMembershipsByUserId(bobCtx, apispec.GetTeamMembershipsByUserIdRequestObject{
					UserId: bob.Id.String(),
				})
				require.NoError(t, err)
				memberships := resp.(apispec.GetTeamMembershipsByUserId200JSONResponse)
				assert.Len(t, memberships, 1)
			})
		})

		t.Run("DeleteTeamMembership", func(t *testing.T) {
			t.Run("NonadminOtherMember", func(t *testing.T) {
				_, err := api.DeleteTeamMembership(bobCtx, apispec.DeleteTeamMembershipRequestObject{
					TeamId: team.Id.String(),
					UserId: alice.Id.String(),
				})
				require.Error(t, err)
			})

			t.Run("Self", func(t *testing.T) {
				_, err := api.DeleteTeamMembership(bobCtx, apispec.DeleteTeamMembershipRequestObject{
					TeamId: team.Id.String(),
					UserId: bob.Id.String(),
				})
				require.NoError(t, err)
			})
		})
	})
}

func TestAPI_QueueTeamReportGeneration(t *testing.T) {
	api := NewTestAPI(t)
	_, adminCtx := api.NewTestUser("admin@example.com", model.UserRoleAdministrator)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)

	// Create an integration so something actually gets queued.
	{
		_, err := api.CreateAWSIntegration(aliceCtx, apispec.CreateAWSIntegrationRequestObject{
			TeamId: team.Id.String(),
			Body: &apispec.CreateAWSIntegrationJSONRequestBody{
				Name:    "Foo",
				RoleArn: "arn:aws:iam::123456789012:role/MyRole",
				CloudtrailTrail: &apispec.CreateAWSIntegrationCloudTrailTrailInput{
					S3BucketName: "aws-cloudtrail-logs",
				},
			},
		})
		require.NoError(t, err)
	}

	req := apispec.QueueTeamReportGenerationRequestObject{
		TeamId: team.Id.String(),
		Body: &apispec.QueueTeamReportGenerationJSONRequestBody{
			StartTime:       time.Now(),
			DurationSeconds: 60 * 60,
			Retention:       apispec.ONEWEEK,
		},
	}

	t.Run("Nonadmin", func(t *testing.T) {
		_, err := api.QueueTeamReportGeneration(aliceCtx, req)
		require.Error(t, err)
	})

	ignoreSQSRequests := len(api.app.SQSRequests("us-east-1"))

	_, err := api.QueueTeamReportGeneration(adminCtx, req)
	require.NoError(t, err)

	sqsRequests := api.app.SQSRequests("us-east-1")[ignoreSQSRequests:]
	require.Len(t, sqsRequests, 1)
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
}

func TestAPI_TeamPrincipalSettings(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	team := api.NewTestTeamWithSubscription(aliceCtx, app.TeamSubscriptionTierIndividual)

	const principalKey = "foo"

	t.Run("Default", func(t *testing.T) {
		resp, err := api.GetTeamPrincipalSettings(aliceCtx, apispec.GetTeamPrincipalSettingsRequestObject{
			TeamId:       team.Id.String(),
			PrincipalKey: principalKey,
		})
		require.NoError(t, err)
		settings := resp.(apispec.GetTeamPrincipalSettings200JSONResponse)
		assert.Nil(t, settings.Description)
	})

	resp, err := api.UpdateTeamPrincipalSettings(aliceCtx, apispec.UpdateTeamPrincipalSettingsRequestObject{
		TeamId:       team.Id.String(),
		PrincipalKey: principalKey,
		Body: &apispec.UpdateTeamPrincipalSettingsJSONRequestBody{
			Description: pointer("This is a description."),
		},
	})
	require.NoError(t, err)
	settings := resp.(apispec.UpdateTeamPrincipalSettings200JSONResponse)
	assert.Equal(t, "This is a description.", *settings.Description)

	t.Run("Get", func(t *testing.T) {
		resp, err := api.GetTeamPrincipalSettings(aliceCtx, apispec.GetTeamPrincipalSettingsRequestObject{
			TeamId:       team.Id.String(),
			PrincipalKey: principalKey,
		})
		require.NoError(t, err)
		settings := resp.(apispec.GetTeamPrincipalSettings200JSONResponse)
		assert.Equal(t, "This is a description.", *settings.Description)
	})

	t.Run("Update", func(t *testing.T) {
		resp, err := api.UpdateTeamPrincipalSettings(aliceCtx, apispec.UpdateTeamPrincipalSettingsRequestObject{
			TeamId:       team.Id.String(),
			PrincipalKey: principalKey,
			Body: &apispec.UpdateTeamPrincipalSettingsJSONRequestBody{
				Description: pointer("This is a new description."),
			},
		})
		require.NoError(t, err)
		settings := resp.(apispec.UpdateTeamPrincipalSettings200JSONResponse)
		assert.Equal(t, "This is a new description.", *settings.Description)

		t.Run("Get", func(t *testing.T) {
			resp, err := api.GetTeamPrincipalSettings(aliceCtx, apispec.GetTeamPrincipalSettingsRequestObject{
				TeamId:       team.Id.String(),
				PrincipalKey: principalKey,
			})
			require.NoError(t, err)
			settings := resp.(apispec.GetTeamPrincipalSettings200JSONResponse)
			assert.Equal(t, "This is a new description.", *settings.Description)
		})
	})
}
