package api

import (
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app/apptest"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestUserPasskey(t *testing.T) {
	api := NewTestAPI(t)
	alice, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	var passkey *apptest.Passkey
	var passkeyId string

	// Register the passkey.
	{
		beginResp, err := api.BeginUserPasskeyRegistration(aliceCtx, apispec.BeginUserPasskeyRegistrationRequestObject{
			UserId: alice.Id.String(),
		})
		require.NoError(t, err)
		beginOutput := beginResp.(apispec.BeginUserPasskeyRegistration200JSONResponse)

		credentialCreationOptions, err := reshape[protocol.CredentialCreation](beginOutput.CredentialCreationOptions)
		require.NoError(t, err)

		newPasskey, ccResp, err := apptest.NewPasskey(&credentialCreationOptions)
		require.NoError(t, err)
		passkey = newPasskey

		completeResp, err := api.CompleteUserPasskeyRegistration(aliceCtx, apispec.CompleteUserPasskeyRegistrationRequestObject{
			UserId: alice.Id.String(),
			Body: &apispec.CompleteUserPasskeyRegistrationJSONRequestBody{
				SessionId:                  beginOutput.SessionId,
				PasskeyName:                "My Passkey",
				CredentialCreationResponse: ccResp,
			},
		})
		require.NoError(t, err)
		completeOutput := completeResp.(apispec.CompleteUserPasskeyRegistration200JSONResponse)
		assert.Equal(t, "My Passkey", completeOutput.Name)
		passkeyId = completeOutput.Id
	}

	// Make sure we can authenticate with the passkey.
	t.Run("Authenticate", func(t *testing.T) {
		beginResp, err := api.BeginUserPasskeyAuthentication(api.AnonymousContext, apispec.BeginUserPasskeyAuthenticationRequestObject{})
		require.NoError(t, err)
		beginOutput := beginResp.(apispec.BeginUserPasskeyAuthentication200JSONResponse)

		credentialAssertionOptions, err := reshape[protocol.CredentialAssertion](beginOutput.CredentialAssertionOptions)
		require.NoError(t, err)

		caResp, err := passkey.Get(&credentialAssertionOptions)
		require.NoError(t, err)

		var creds apispec.UserCredentials
		creds.FromUserPasskeyCredentials(apispec.UserPasskeyCredentials{
			SessionId:                   beginOutput.SessionId,
			CredentialAssertionResponse: caResp,
		})

		authResp, err := api.Authenticate(api.AnonymousContext, apispec.AuthenticateRequestObject{
			Body: &creds,
		})
		require.NoError(t, err)
		auth := authResp.(apispec.Authenticate200JSONResponse)
		assert.NotEmpty(t, auth.Token)
		assert.Equal(t, alice.Id.String(), auth.User.Id)
	})

	// Make sure we can see the passkey.
	t.Run("GetUserPasskeys", func(t *testing.T) {
		resp, err := api.GetUserPasskeys(aliceCtx, apispec.GetUserPasskeysRequestObject{
			UserId: alice.Id.String(),
		})
		require.NoError(t, err)
		output := resp.(apispec.GetUserPasskeys200JSONResponse)
		require.Len(t, output, 1)
		got := output[0]
		assert.Equal(t, passkeyId, got.Id)
		assert.Equal(t, "My Passkey", got.Name)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		_, bobCtx := api.NewTestUser("bob@example.com", model.UserRoleCustomer)

		// Make sure Bob can't register Alice's passkeys.
		t.Run("BeginUserPasskeyRegistration", func(t *testing.T) {
			_, err := api.BeginUserPasskeyRegistration(bobCtx, apispec.BeginUserPasskeyRegistrationRequestObject{
				UserId: alice.Id.String(),
			})
			require.Error(t, err)
		})

		// Make sure Bob can't see Alice's passkeys.
		t.Run("GetUserPasskeys", func(t *testing.T) {
			_, err := api.GetUserPasskeys(bobCtx, apispec.GetUserPasskeysRequestObject{
				UserId: alice.Id.String(),
			})
			require.Error(t, err)
		})

		// Make sure Bob can't update Alice's passkeys.
		t.Run("UpdateUserPasskeyById", func(t *testing.T) {
			_, err := api.UpdateUserPasskeyById(bobCtx, apispec.UpdateUserPasskeyByIdRequestObject{
				PasskeyId: passkeyId,
				Body: &apispec.UpdateUserPasskeyByIdJSONRequestBody{
					Name: pointer("My New Passkey Name"),
				},
			})
			require.Error(t, err)
		})

		// Make sure Bob can't delete Alice's passkeys.
		t.Run("DeleteUserPasskeyById", func(t *testing.T) {
			_, err := api.DeleteUserPasskeyById(bobCtx, apispec.DeleteUserPasskeyByIdRequestObject{
				PasskeyId: passkeyId,
			})
			require.Error(t, err)
		})
	})

	// Make sure we can update the passkey.
	t.Run("UpdateUserPasskeyById", func(t *testing.T) {
		resp, err := api.UpdateUserPasskeyById(aliceCtx, apispec.UpdateUserPasskeyByIdRequestObject{
			PasskeyId: passkeyId,
			Body: &apispec.UpdateUserPasskeyByIdJSONRequestBody{
				Name: pointer("My New Passkey Name"),
			},
		})
		require.NoError(t, err)
		output := resp.(apispec.UpdateUserPasskeyById200JSONResponse)
		assert.Equal(t, "My New Passkey Name", output.Name)
	})

	t.Run("DeleteUserPasskeyById", func(t *testing.T) {
		// Make sure we can delete the passkey.
		{
			resp, err := api.DeleteUserPasskeyById(aliceCtx, apispec.DeleteUserPasskeyByIdRequestObject{
				PasskeyId: passkeyId,
			})
			require.NoError(t, err)
			_ = resp.(apispec.DeleteUserPasskeyById200Response)
		}

		// Make sure we can no longer see the passkey.
		{
			resp, err := api.GetUserPasskeys(aliceCtx, apispec.GetUserPasskeysRequestObject{
				UserId: alice.Id.String(),
			})
			require.NoError(t, err)
			output := resp.(apispec.GetUserPasskeys200JSONResponse)
			assert.Empty(t, output)
		}
	})
}
