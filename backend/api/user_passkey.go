package api

import (
	"context"

	"github.com/go-webauthn/webauthn/protocol"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func (api *API) DeleteUserPasskeyById(ctx context.Context, request apispec.DeleteUserPasskeyByIdRequestObject) (apispec.DeleteUserPasskeyByIdResponseObject, error) {
	sess := ctxSession(ctx)

	if err := sess.DeleteUserPasskeyById(ctx, model.Id(request.PasskeyId)); err != nil {
		return nil, err
	} else {
		return apispec.DeleteUserPasskeyById200Response{}, nil
	}
}

func (api *API) UpdateUserPasskeyById(ctx context.Context, request apispec.UpdateUserPasskeyByIdRequestObject) (apispec.UpdateUserPasskeyByIdResponseObject, error) {
	sess := ctxSession(ctx)

	patch := app.UserPasskeyPatch{
		Name: request.Body.Name,
	}

	if result, err := sess.PatchUserPasskeyById(ctx, model.Id(request.PasskeyId), patch); err != nil {
		return nil, err
	} else if result == nil {
		return nil, app.NotFoundError("No such passkey.")
	} else {
		return apispec.UpdateUserPasskeyById200JSONResponse(UserPasskeyFromModel(result)), nil
	}
}

func UserPasskeyFromModel(passkey *model.UserPasskey) apispec.UserPasskey {
	return apispec.UserPasskey{
		Id:           passkey.Id.String(),
		CreationTime: passkey.CreationTime,
		Name:         passkey.Name,
		UserId:       passkey.UserId.String(),
	}
}

func (api *API) BeginUserPasskeyAuthentication(ctx context.Context, request apispec.BeginUserPasskeyAuthenticationRequestObject) (apispec.BeginUserPasskeyAuthenticationResponseObject, error) {
	sess := ctxSession(ctx)

	if output, err := sess.BeginUserPasskeyAuthentication(ctx); err != nil {
		return nil, err
	} else {
		return apispec.BeginUserPasskeyAuthentication200JSONResponse{
			SessionId:                  output.SessionId.String(),
			CredentialAssertionOptions: output.CredentialAssertion,
		}, nil
	}
}

func (api *API) BeginUserPasskeyRegistration(ctx context.Context, request apispec.BeginUserPasskeyRegistrationRequestObject) (apispec.BeginUserPasskeyRegistrationResponseObject, error) {
	sess := ctxSession(ctx)

	if !sess.HasUserId(UserIdFromRequest(sess, request.UserId)) {
		return nil, app.AuthorizationError{}
	}

	if output, err := sess.BeginUserPasskeyRegistration(ctx); err != nil {
		return nil, err
	} else {
		return apispec.BeginUserPasskeyRegistration200JSONResponse{
			SessionId:                 output.SessionId.String(),
			CredentialCreationOptions: output.CredentialCreation,
		}, nil
	}
}

func (api *API) CompleteUserPasskeyRegistration(ctx context.Context, request apispec.CompleteUserPasskeyRegistrationRequestObject) (apispec.CompleteUserPasskeyRegistrationResponseObject, error) {
	sess := ctxSession(ctx)

	if !sess.HasUserId(UserIdFromRequest(sess, request.UserId)) {
		return nil, app.AuthorizationError{}
	}

	if response, err := reshape[protocol.CredentialCreationResponse](request.Body.CredentialCreationResponse); err != nil {
		return nil, err
	} else if passkey, err := sess.CompleteUserPasskeyRegistration(ctx, app.CompleteUserPasskeyRegistrationInput{
		SessionId:                  model.Id(request.Body.SessionId),
		PasskeyName:                request.Body.PasskeyName,
		CredentialCreationResponse: &response,
	}); err != nil {
		return nil, err
	} else {
		return apispec.CompleteUserPasskeyRegistration200JSONResponse(UserPasskeyFromModel(passkey)), nil
	}
}

func (api *API) GetUserPasskeys(ctx context.Context, request apispec.GetUserPasskeysRequestObject) (apispec.GetUserPasskeysResponseObject, error) {
	sess := ctxSession(ctx)

	if passkeys, err := sess.GetUserPasskeysByUserId(ctx, UserIdFromRequest(sess, request.UserId)); err != nil {
		return nil, err
	} else {
		return apispec.GetUserPasskeys200JSONResponse(mapSlice(passkeys, UserPasskeyFromModel)), nil
	}
}
