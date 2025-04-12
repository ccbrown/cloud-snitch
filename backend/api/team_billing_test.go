package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v81"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app/apptest"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestAPI_TeamBillingProfile(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	resp, err := api.CreateTeam(aliceCtx, apispec.CreateTeamRequestObject{
		Body: &apispec.CreateTeamJSONRequestBody{
			Name: "Foo",
		},
	})
	require.NoError(t, err)
	team := resp.(apispec.CreateTeam200JSONResponse)
	assert.Equal(t, "Foo", team.Name)

	t.Run("NotFound", func(t *testing.T) {
		_, err := api.GetTeamBillingProfile(aliceCtx, apispec.GetTeamBillingProfileRequestObject{
			TeamId: team.Id,
		})
		require.Error(t, err)
	})

	_, err = api.CreateTeamBillingProfile(aliceCtx, apispec.CreateTeamBillingProfileRequestObject{
		TeamId: team.Id,
		Body: &apispec.CreateTeamBillingProfileJSONRequestBody{
			Name: "Foo",
			Address: apispec.TeamBillingAddress{
				Line1:      pointer("123 Main St"),
				City:       pointer("Seattle"),
				State:      pointer("WA"),
				PostalCode: "98101",
				Country:    "US",
			},
		},
	})
	require.NoError(t, err)

	t.Run("Get", func(t *testing.T) {
		resp, err := api.GetTeamBillingProfile(aliceCtx, apispec.GetTeamBillingProfileRequestObject{
			TeamId: team.Id,
		})
		require.NoError(t, err)
		billingProfile := resp.(apispec.GetTeamBillingProfile200JSONResponse)
		assert.Equal(t, "123 Main St", *billingProfile.Address.Line1)
		assert.Equal(t, "Seattle", *billingProfile.Address.City)
		assert.Equal(t, "WA", *billingProfile.Address.State)
		assert.Equal(t, "98101", billingProfile.Address.PostalCode)
		assert.Equal(t, "US", billingProfile.Address.Country)
	})

	t.Run("Update", func(t *testing.T) {
		_, err = api.UpdateTeamBillingProfile(aliceCtx, apispec.UpdateTeamBillingProfileRequestObject{
			TeamId: team.Id,
			Body: &apispec.UpdateTeamBillingProfileJSONRequestBody{
				Address: &apispec.TeamBillingAddress{
					Line1:      pointer("567 Main St"),
					City:       pointer("Seattle"),
					State:      pointer("WA"),
					PostalCode: "98101",
					Country:    "US",
				},
			},
		})
		require.NoError(t, err)

		resp, err := api.GetTeamBillingProfile(aliceCtx, apispec.GetTeamBillingProfileRequestObject{
			TeamId: team.Id,
		})
		require.NoError(t, err)
		billingProfile := resp.(apispec.GetTeamBillingProfile200JSONResponse)
		assert.Equal(t, "567 Main St", *billingProfile.Address.Line1)
		assert.Equal(t, "Seattle", *billingProfile.Address.City)
		assert.Equal(t, "WA", *billingProfile.Address.State)
		assert.Equal(t, "98101", billingProfile.Address.PostalCode)
		assert.Equal(t, "US", billingProfile.Address.Country)
	})
}

func TestAPI_TeamPaymentMethod(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	resp, err := api.CreateTeam(aliceCtx, apispec.CreateTeamRequestObject{
		Body: &apispec.CreateTeamJSONRequestBody{
			Name: "Foo",
		},
	})
	require.NoError(t, err)
	team := resp.(apispec.CreateTeam200JSONResponse)
	assert.Equal(t, "Foo", team.Name)

	_, err = api.CreateTeamBillingProfile(aliceCtx, apispec.CreateTeamBillingProfileRequestObject{
		TeamId: team.Id,
		Body: &apispec.CreateTeamBillingProfileJSONRequestBody{
			Name: "Foo",
			Address: apispec.TeamBillingAddress{
				Line1:      pointer("123 Main St"),
				City:       pointer("Seattle"),
				State:      pointer("WA"),
				PostalCode: "98101",
				Country:    "US",
			},
		},
	})
	require.NoError(t, err)

	t.Run("NotSet", func(t *testing.T) {
		_, err := api.GetTeamPaymentMethod(aliceCtx, apispec.GetTeamPaymentMethodRequestObject{
			TeamId: team.Id,
		})
		require.Error(t, err)
	})

	{
		resp, err := api.PutTeamPaymentMethod(aliceCtx, apispec.PutTeamPaymentMethodRequestObject{
			TeamId: team.Id,
			Body: &apispec.PutTeamPaymentMethodJSONRequestBody{
				StripePaymentMethodId: apptest.DummyStripeCard.ID,
			},
		})
		require.NoError(t, err)
		method := apispec.TeamPaymentMethod(resp.(apispec.PutTeamPaymentMethod200JSONResponse))
		_, err = method.AsTeamPaymentMethodCard()
		require.NoError(t, err)
	}

	t.Run("Get", func(t *testing.T) {
		resp, err := api.GetTeamPaymentMethod(aliceCtx, apispec.GetTeamPaymentMethodRequestObject{
			TeamId: team.Id,
		})
		require.NoError(t, err)
		method := apispec.TeamPaymentMethod(resp.(apispec.GetTeamPaymentMethod200JSONResponse))
		_, err = method.AsTeamPaymentMethodCard()
		require.NoError(t, err)
	})
}

func TestAPI_TeamSubscription(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	resp, err := api.CreateTeam(aliceCtx, apispec.CreateTeamRequestObject{
		Body: &apispec.CreateTeamJSONRequestBody{
			Name: "Foo",
		},
	})
	require.NoError(t, err)
	team := resp.(apispec.CreateTeam200JSONResponse)

	_, err = api.CreateTeamBillingProfile(aliceCtx, apispec.CreateTeamBillingProfileRequestObject{
		TeamId: team.Id,
		Body: &apispec.CreateTeamBillingProfileJSONRequestBody{
			Name: "Foo",
			Address: apispec.TeamBillingAddress{
				Line1:      pointer("123 Main St"),
				City:       pointer("Seattle"),
				State:      pointer("WA"),
				PostalCode: "98101",
				Country:    "US",
			},
		},
	})
	require.NoError(t, err)

	t.Run("NotSet", func(t *testing.T) {
		_, err := api.GetTeamSubscription(aliceCtx, apispec.GetTeamSubscriptionRequestObject{
			TeamId: team.Id,
		})
		require.Error(t, err)
	})

	_, err = api.PutTeamPaymentMethod(aliceCtx, apispec.PutTeamPaymentMethodRequestObject{
		TeamId: team.Id,
		Body: &apispec.PutTeamPaymentMethodJSONRequestBody{
			StripePaymentMethodId: apptest.DummyStripeCard.ID,
		},
	})
	require.NoError(t, err)

	_, err = api.CreateTeamSubscription(aliceCtx, apispec.CreateTeamSubscriptionRequestObject{
		TeamId: team.Id,
		Body: &apispec.CreateTeamSubscriptionJSONRequestBody{
			Tier: apispec.INDIVIDUAL,
		},
	})
	require.NoError(t, err)

	t.Run("Get", func(t *testing.T) {
		resp, err := api.GetTeamSubscription(aliceCtx, apispec.GetTeamSubscriptionRequestObject{
			TeamId: team.Id,
		})
		require.NoError(t, err)
		subscription := resp.(apispec.GetTeamSubscription200JSONResponse)
		assert.Equal(t, "Individual Subscription", subscription.Name)
		assert.EqualValues(t, 0, subscription.Accounts)
		assert.Equal(t, "$0.99", subscription.Price.AccountMonth.Text)
	})

	t.Run("Update", func(t *testing.T) {
		_, err = api.UpdateTeamSubscription(aliceCtx, apispec.UpdateTeamSubscriptionRequestObject{
			TeamId: team.Id,
			Body: &apispec.UpdateTeamSubscriptionJSONRequestBody{
				Tier: pointer(apispec.TEAM),
			},
		})
		require.NoError(t, err)

		resp, err := api.GetTeamSubscription(aliceCtx, apispec.GetTeamSubscriptionRequestObject{
			TeamId: team.Id,
		})
		require.NoError(t, err)
		subscription := resp.(apispec.GetTeamSubscription200JSONResponse)
		assert.Equal(t, "Team Subscription", subscription.Name)
		assert.EqualValues(t, 0, subscription.Accounts)
		assert.Equal(t, "$9.99", subscription.Price.AccountMonth.Text)
	})

	t.Run("Custom", func(t *testing.T) {
		// Update the subscription in Stripe
		{
			subscriptions := api.app.Stripe().Subscriptions.List(nil)
			require.True(t, subscriptions.Next())
			subscription := subscriptions.Subscription()
			require.Len(t, subscription.Items.Data, 1)
			item := subscription.Items.Data[0]
			_, err := api.app.Stripe().Subscriptions.Update(subscription.ID, &stripe.SubscriptionParams{
				Items: []*stripe.SubscriptionItemsParams{
					{
						ID:    &item.ID,
						Price: stripe.String(apptest.DummyStripePriceCustomSubscription.ID),
					},
				},
			})
			require.NoError(t, err)
		}

		t.Run("Get", func(t *testing.T) {
			resp, err := api.GetTeamSubscription(aliceCtx, apispec.GetTeamSubscriptionRequestObject{
				TeamId: team.Id,
			})
			require.NoError(t, err)
			subscription := resp.(apispec.GetTeamSubscription200JSONResponse)
			assert.Equal(t, "My Custom Subscription", subscription.Name)
			assert.EqualValues(t, 0, subscription.Accounts)
			assert.Nil(t, subscription.Price)
		})

		t.Run("Update", func(t *testing.T) {
			_, err = api.UpdateTeamSubscription(aliceCtx, apispec.UpdateTeamSubscriptionRequestObject{
				TeamId: team.Id,
				Body: &apispec.UpdateTeamSubscriptionJSONRequestBody{
					Tier: pointer(apispec.TEAM),
				},
			})
			require.Error(t, err)
		})
	})
}
