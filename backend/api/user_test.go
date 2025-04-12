package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestAPI_GetUser(t *testing.T) {
	api := NewTestAPI(t)
	alice, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)
	_, adminCtx := api.NewTestUser("admin@example.com", model.UserRoleAdministrator)

	user, err := api.GetUser(aliceCtx, apispec.GetUserRequestObject{UserId: "self"})
	require.NoError(t, err)
	assert.Equal(t, alice.Id.String(), user.(apispec.GetUser200JSONResponse).Id)

	t.Run("GetUser", func(t *testing.T) {
		user, err := api.GetUser(adminCtx, apispec.GetUserRequestObject{UserId: alice.Id.String()})
		require.NoError(t, err)
		assert.Equal(t, alice.Id.String(), user.(apispec.GetUser200JSONResponse).Id)
	})

	t.Run("GetUsers", func(t *testing.T) {
		_, err := api.GetUsers(aliceCtx, apispec.GetUsersRequestObject{})
		require.Error(t, err)

		resp, err := api.GetUsers(adminCtx, apispec.GetUsersRequestObject{})
		require.NoError(t, err)
		users := resp.(apispec.GetUsers200JSONResponse)
		assert.Len(t, users, 2)
	})

	t.Run("BadAuth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/users/self", nil)
		require.NoError(t, err)
		r.Header.Set("Authorization", "token foo")
		api.ServeHTTP(w, r)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "%v", string(body))
	})
}

func TestAPI_SignOut(t *testing.T) {
	api := NewTestAPI(t)

	_, err := api.BeginUserRegistration(api.AnonymousContext, apispec.BeginUserRegistrationRequestObject{
		Body: &apispec.BeginUserRegistrationJSONRequestBody{
			EmailAddress:                    "alice@example.com",
			TermsOfServiceAgreementRevision: "2024.01.01",
			PrivacyPolicyAgreementRevision:  "2024.01.01",
			CookiePolicyAgreementRevision:   "2024.01.01",
		},
	})
	require.NoError(t, err)

	email := <-api.app.Emails()
	userRegistrationTokenRegexp := regexp.MustCompile(`token=([a-zA-Z0-9_\-]+)`)
	token := userRegistrationTokenRegexp.FindStringSubmatch(email.HTML)[1]

	regResp, err := api.CompleteUserRegistration(api.AnonymousContext, apispec.CompleteUserRegistrationRequestObject{
		Body: &apispec.CompleteUserRegistrationJSONRequestBody{
			Token: token,
		},
	})
	require.NoError(t, err)
	output := regResp.(apispec.CompleteUserRegistration200JSONResponse)

	assert.NotEmpty(t, output.User.Id)
	assert.NotEmpty(t, output.AccessToken)

	// Make sure the token is good.
	{
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/users/self", nil)
		r.Header.Set("Authorization", "token "+output.AccessToken)
		require.NoError(t, err)
		api.ServeHTTP(w, r)

		resp := w.Result()
		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// Sign out.
	{
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/sign-out", nil)
		r.Header.Set("Authorization", "token "+output.AccessToken)
		require.NoError(t, err)
		api.ServeHTTP(w, r)

		resp := w.Result()
		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// Make sure the token is no longer good.
	{
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/users/self", nil)
		r.Header.Set("Authorization", "token "+output.AccessToken)
		require.NoError(t, err)
		api.ServeHTTP(w, r)

		resp := w.Result()
		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAPI_UserRegistration(t *testing.T) {
	api := NewTestAPI(t)

	_, err := api.BeginUserRegistration(api.AnonymousContext, apispec.BeginUserRegistrationRequestObject{
		Body: &apispec.BeginUserRegistrationJSONRequestBody{
			EmailAddress:                    "alice@example.com",
			TermsOfServiceAgreementRevision: "2024.01.01",
			PrivacyPolicyAgreementRevision:  "2024.01.01",
			CookiePolicyAgreementRevision:   "2024.01.01",
		},
	})
	require.NoError(t, err)

	email := <-api.app.Emails()
	userRegistrationTokenRegexp := regexp.MustCompile(`token=([a-zA-Z0-9_\-]+)`)
	token := userRegistrationTokenRegexp.FindStringSubmatch(email.HTML)[1]

	regResp, err := api.CompleteUserRegistration(api.AnonymousContext, apispec.CompleteUserRegistrationRequestObject{
		Body: &apispec.CompleteUserRegistrationJSONRequestBody{
			Token: token,
		},
	})
	require.NoError(t, err)
	output := regResp.(apispec.CompleteUserRegistration200JSONResponse)

	assert.NotEmpty(t, output.User.Id)
	assert.NotEmpty(t, output.AccessToken)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/users/self", nil)
	r.Header.Set("Authorization", "token "+output.AccessToken)
	require.NoError(t, err)
	api.ServeHTTP(w, r)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	var user apispec.User
	require.NoError(t, json.Unmarshal(body, &user))

	assert.Equal(t, output.User.Id, user.Id)
}

func TestAPI_UpdateUser(t *testing.T) {
	api := NewTestAPI(t)
	_, aliceCtx := api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	t.Run("Empty", func(t *testing.T) {
		_, err := api.UpdateUser(aliceCtx, apispec.UpdateUserRequestObject{UserId: "self", Body: &apispec.UpdateUserJSONRequestBody{}})
		require.NoError(t, err)
	})

	t.Run("Password", func(t *testing.T) {
		resp, err := api.UpdateUser(aliceCtx, apispec.UpdateUserRequestObject{UserId: "self", Body: &apispec.UpdateUserJSONRequestBody{
			Password: aws.String("mynewpassword!!"),
		}})
		require.NoError(t, err)
		user := resp.(apispec.UpdateUser200JSONResponse)
		assert.True(t, *user.HasPassword)

		resp, err = api.UpdateUser(aliceCtx, apispec.UpdateUserRequestObject{UserId: "self", Body: &apispec.UpdateUserJSONRequestBody{
			Password: aws.String(""),
		}})
		require.NoError(t, err)
		user = resp.(apispec.UpdateUser200JSONResponse)
		assert.False(t, user.HasPassword != nil && *user.HasPassword)
	})
}

func TestAPI_UserEmailAuthentication(t *testing.T) {
	api := NewTestAPI(t)

	api.NewTestUser("alice@example.com", model.UserRoleCustomer)

	_, err := api.BeginUserEmailAuthentication(api.AnonymousContext, apispec.BeginUserEmailAuthenticationRequestObject{
		Body: &apispec.BeginUserEmailAuthenticationJSONRequestBody{
			EmailAddress: "alice@example.com",
		},
	})
	require.NoError(t, err)

	email := <-api.app.Emails()

	var creds apispec.UserCredentials
	creds.FromUserEmailCredentials(apispec.UserEmailCredentials{
		Token: regexp.MustCompile(`token=([a-zA-Z0-9_\-]+)`).FindStringSubmatch(email.HTML)[1],
	})

	regResp, err := api.Authenticate(api.AnonymousContext, apispec.AuthenticateRequestObject{
		Body: &creds,
	})
	require.NoError(t, err)
	output := regResp.(apispec.Authenticate200JSONResponse)

	assert.NotEmpty(t, output.User.Id)
	assert.NotEmpty(t, output.Token)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/users/self", nil)
	r.Header.Set("Authorization", "token "+output.Token)
	require.NoError(t, err)
	api.ServeHTTP(w, r)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	var user apispec.User
	require.NoError(t, json.Unmarshal(body, &user))

	assert.Equal(t, "alice@example.com", user.EmailAddress)
}
