package app_test

import (
	"context"
	"encoding/base64"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/app/apptest"
)

var userRegistrationTokenRegexp = regexp.MustCompile(`token=([a-zA-Z0-9_\-]+)`)

func TestUserRegistration(t *testing.T) {
	a := apptest.NewTestApp(t)
	sess := a.NewAnonymousSession()

	t.Run("MissingAgreements", func(t *testing.T) {
		require.Error(t, sess.BeginUserRegistration(context.Background(), app.BeginUserRegistrationInput{
			EmailAddress: "foo@example.com",
		}))
	})

	require.NoError(t, sess.BeginUserRegistration(context.Background(), app.BeginUserRegistrationInput{
		EmailAddress:                    "foo@example.com",
		TermsOfServiceAgreementRevision: "2024.01.01",
		PrivacyPolicyAgreementRevision:  "2024.01.01",
		CookiePolicyAgreementRevision:   "2024.01.01",
	}))

	email := <-a.Emails()
	assert.Equal(t, "foo@example.com", email.To)

	token := userRegistrationTokenRegexp.FindStringSubmatch(email.HTML)[1]
	assert.NotEmpty(t, token)

	tokenBytes, err := base64.RawURLEncoding.DecodeString(token)
	require.NoError(t, err)

	user, accessToken, err := sess.CompleteUserRegistration(context.Background(), app.CompleteUserRegistrationInput{
		Token: tokenBytes,
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "foo@example.com", user.EmailAddress)

	userSess, err := sess.WithUserAccessToken(context.Background(), accessToken)
	require.NoError(t, err)
	require.NotNil(t, userSess)
	assert.Equal(t, user.Id, userSess.User().Id)

	t.Run("DoubleUse", func(t *testing.T) {
		_, _, err := sess.CompleteUserRegistration(context.Background(), app.CompleteUserRegistrationInput{
			Token: tokenBytes,
		})
		assert.Error(t, err)
	})

	t.Run("ExistingUser", func(t *testing.T) {
		require.NoError(t, sess.BeginUserRegistration(context.Background(), app.BeginUserRegistrationInput{
			EmailAddress:                    "foo@example.com",
			TermsOfServiceAgreementRevision: "2024.01.01",
			PrivacyPolicyAgreementRevision:  "2024.01.01",
			CookiePolicyAgreementRevision:   "2024.01.01",
		}))
	})
}
