package api

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func UserIdFromRequest(sess *app.Session, id string) model.Id {
	if id == "self" && sess.User() != nil {
		return sess.User().Id
	}
	return model.Id(id)
}

func UserAgreementFromModel(agreement model.UserAgreement) apispec.UserAgreement {
	return apispec.UserAgreement{
		Revision: string(agreement.Revision),
	}
}

func UserFromModel(user *model.User) apispec.User {
	ret := apispec.User{
		Id:           user.Id.String(),
		EmailAddress: user.EmailAddress,
	}
	if user.HasPassword() {
		ret.HasPassword = pointer(true)
	}
	if user.Role != model.UserRoleNone {
		ret.Role = pointer(UserRoleFromModel(user.Role))
	}
	if user.TermsOfServiceAgreement.Revision.IsValid() {
		ret.TermsOfServiceAgreement = pointer(UserAgreementFromModel(user.TermsOfServiceAgreement))
	}
	if user.PrivacyPolicyAgreement.Revision.IsValid() {
		ret.PrivacyPolicyAgreement = pointer(UserAgreementFromModel(user.PrivacyPolicyAgreement))
	}
	if user.CookiePolicyAgreement.Revision.IsValid() {
		ret.CookiePolicyAgreement = pointer(UserAgreementFromModel(user.CookiePolicyAgreement))
	}
	return ret
}

func UserRoleFromSpec(role apispec.UserRole) model.UserRole {
	switch role {
	case apispec.UserRoleADMINISTRATOR:
		return model.UserRoleAdministrator
	case apispec.UserRoleCUSTOMER:
		return model.UserRoleCustomer
	default:
		panic(fmt.Sprintf("unknown user role: %v", role))
	}
}

func UserRoleFromModel(role model.UserRole) apispec.UserRole {
	switch role {
	case model.UserRoleAdministrator:
		return apispec.UserRoleADMINISTRATOR
	case model.UserRoleCustomer:
		return apispec.UserRoleCUSTOMER
	default:
		panic(fmt.Sprintf("unknown user role: %v", role))
	}
}

func sessionWithCredentials(ctx context.Context, userCredentials *apispec.UserCredentials) (*app.Session, error) {
	sess := ctxSession(ctx)
	if creds, _ := userCredentials.AsUserEmailAddressAndPasswordCredentials(); creds.EmailAddress != "" && creds.Password != "" {
		return sess.WithUserCredentials(ctx, creds.EmailAddress, creds.Password)
	} else if creds, _ := userCredentials.AsUserPasskeyCredentials(); creds.SessionId != "" {
		if response, err := reshape[protocol.CredentialAssertionResponse](creds.CredentialAssertionResponse); err != nil {
			return nil, err
		} else {
			return sess.WithUserPasskey(ctx, model.Id(creds.SessionId), &response)
		}
	} else if creds, _ := userCredentials.AsUserEmailCredentials(); creds.Token != "" {
		token, _ := base64.RawURLEncoding.DecodeString(creds.Token)
		return sess.WithUserEmailAuthentication(ctx, token)
	} else {
		return nil, nil
	}
}

func (api *API) Authenticate(ctx context.Context, request apispec.AuthenticateRequestObject) (apispec.AuthenticateResponseObject, error) {
	sess, err := sessionWithCredentials(ctx, request.Body)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return apispec.Authenticate400JSONResponse{
			ErrorResponseJSONResponse: apispec.ErrorResponseJSONResponse{
				Message: "Invalid credentials.",
			},
		}, nil
	} else if token, err := sess.CreateUserAccessToken(ctx); err != nil {
		return nil, err
	} else {
		return apispec.Authenticate200JSONResponse{
			Token: base64.RawURLEncoding.EncodeToString(token),
			User:  UserFromModel(sess.User()),
		}, nil
	}
}

func (*API) SignOut(ctx context.Context, request apispec.SignOutRequestObject) (apispec.SignOutResponseObject, error) {
	sess := ctxSession(ctx)
	if err := sess.SignOut(ctx); err != nil {
		return nil, err
	} else {
		return apispec.SignOut200Response{}, nil
	}
}

func (api *API) GetUsers(ctx context.Context, request apispec.GetUsersRequestObject) (apispec.GetUsersResponseObject, error) {
	sess := ctxSession(ctx)

	if users, err := sess.GetUsers(ctx); err != nil {
		return nil, err
	} else {
		return apispec.GetUsers200JSONResponse(mapSlice(users, UserFromModel)), nil
	}
}

func (api *API) GetUser(ctx context.Context, request apispec.GetUserRequestObject) (apispec.GetUserResponseObject, error) {
	sess := ctxSession(ctx)
	userId := UserIdFromRequest(sess, request.UserId)

	if user, err := sess.GetUserById(ctx, userId); err != nil {
		return nil, err
	} else if user == nil {
		return nil, app.NotFoundError("No such user.")
	} else {
		return apispec.GetUser200JSONResponse(UserFromModel(user)), nil
	}
}

func (api *API) UpdateUser(ctx context.Context, request apispec.UpdateUserRequestObject) (apispec.UpdateUserResponseObject, error) {
	sess := ctxSession(ctx)

	patch := app.UserPatch{
		EmailAddress: request.Body.EmailAddress,
		Password:     request.Body.Password,
	}
	if request.Body.Role != nil {
		role := UserRoleFromSpec(*request.Body.Role)
		patch.Role = &role
	}
	if request.Body.TermsOfServiceAgreementRevision != nil {
		patch.TermsOfServiceAgreementRevision = pointer(model.UserAgreementRevision(*request.Body.TermsOfServiceAgreementRevision))
	}
	if request.Body.PrivacyPolicyAgreementRevision != nil {
		patch.PrivacyPolicyAgreementRevision = pointer(model.UserAgreementRevision(*request.Body.PrivacyPolicyAgreementRevision))
	}
	if request.Body.CookiePolicyAgreementRevision != nil {
		patch.CookiePolicyAgreementRevision = pointer(model.UserAgreementRevision(*request.Body.CookiePolicyAgreementRevision))
	}

	if user, err := sess.PatchUserById(ctx, UserIdFromRequest(sess, request.UserId), patch); err != nil {
		return nil, err
	} else if user == nil {
		return nil, app.NotFoundError("No such user.")
	} else {
		return apispec.UpdateUser200JSONResponse(UserFromModel(user)), nil
	}
}

func (api *API) BeginUserRegistration(ctx context.Context, request apispec.BeginUserRegistrationRequestObject) (apispec.BeginUserRegistrationResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.BeginUserRegistrationInput{
		EmailAddress:                    request.Body.EmailAddress,
		TermsOfServiceAgreementRevision: model.UserAgreementRevision(request.Body.TermsOfServiceAgreementRevision),
		PrivacyPolicyAgreementRevision:  model.UserAgreementRevision(request.Body.PrivacyPolicyAgreementRevision),
		CookiePolicyAgreementRevision:   model.UserAgreementRevision(request.Body.CookiePolicyAgreementRevision),
	}

	if err := sess.BeginUserRegistration(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.BeginUserRegistration200JSONResponse{}, nil
	}
}

func (api *API) CompleteUserRegistration(ctx context.Context, request apispec.CompleteUserRegistrationRequestObject) (apispec.CompleteUserRegistrationResponseObject, error) {
	sess := ctxSession(ctx)

	token, _ := base64.RawURLEncoding.DecodeString(request.Body.Token)

	input := app.CompleteUserRegistrationInput{
		Token: token,
	}

	if user, accessToken, err := sess.CompleteUserRegistration(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.CompleteUserRegistration200JSONResponse{
			User:        UserFromModel(user),
			AccessToken: base64.RawURLEncoding.EncodeToString(accessToken),
		}, nil
	}
}

func (api *API) BeginUserEmailAuthentication(ctx context.Context, request apispec.BeginUserEmailAuthenticationRequestObject) (apispec.BeginUserEmailAuthenticationResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.BeginUserEmailAuthenticationInput{
		EmailAddress: request.Body.EmailAddress,
	}

	if err := sess.BeginUserEmailAuthentication(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.BeginUserEmailAuthentication200JSONResponse{}, nil
	}
}
