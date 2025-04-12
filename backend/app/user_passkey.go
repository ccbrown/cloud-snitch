package app

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

func (s *Session) validateUserPasskey(passkey *model.UserPasskey) UserFacingError {
	if len(passkey.Name) == 0 {
		return NewUserError("A name is required.")
	} else if len(passkey.Name) > 200 {
		return NewUserError("Please use a shorter name.")
	}
	return nil
}

type webAuthnUser struct {
	User     *model.User
	Passkeys []*model.UserPasskey
}

var _ webauthn.User = webAuthnUser{}

func (u webAuthnUser) WebAuthnID() []byte {
	return []byte(u.User.Id)
}

func (u webAuthnUser) WebAuthnName() string {
	return u.User.EmailAddress
}

func (u webAuthnUser) WebAuthnDisplayName() string {
	return u.User.EmailAddress
}

func (u webAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	ret := make([]webauthn.Credential, 0, len(u.Passkeys))
	for _, passkey := range u.Passkeys {
		ret = append(ret, passkey.Credential)
	}
	return ret
}

func (u webAuthnUser) WebAuthnIcon() string {
	return ""
}

func (s *Session) GetUserPasskeysByUserId(ctx context.Context, userId model.Id) ([]*model.UserPasskey, UserFacingError) {
	if !s.HasUserId(userId) {
		return nil, AuthorizationError{}
	}

	ret, err := s.app.store.GetUserPasskeysByUserId(ctx, userId)
	return ret, s.SanitizedError(err)
}

const maxUserPasskeys = 5

type BeginUserPasskeyRegistrationOutput struct {
	CredentialCreation *protocol.CredentialCreation
	SessionId          model.Id
}

func (s *Session) BeginUserPasskeyRegistration(ctx context.Context) (*BeginUserPasskeyRegistrationOutput, UserFacingError) {
	if s.user == nil {
		return nil, AuthorizationError{}
	}

	passkeys, serr := s.GetUserPasskeysByUserId(ctx, s.user.Id)
	if serr != nil {
		return nil, serr
	} else if len(passkeys) >= maxUserPasskeys {
		return nil, NewUserError(fmt.Sprintf("A maximum of %v passkeys are allowed per user. Please remove one to add a new one.", maxUserPasskeys))
	}

	user := webAuthnUser{
		User:     s.user,
		Passkeys: passkeys,
	}
	creation, sessionData, err := s.app.webAuthn.BeginRegistration(user)
	if err != nil {
		return nil, s.SanitizedError(err)
	}
	sessionId := model.NewUserPasskeySessionId()
	if err := s.app.store.PutUserPasskeySession(ctx, &model.UserPasskeySession{
		Id:   sessionId,
		Type: model.UserPasskeySessionTypeRegistration,
		Data: *sessionData,
	}); err != nil {
		return nil, s.SanitizedError(err)
	}
	return &BeginUserPasskeyRegistrationOutput{
		CredentialCreation: creation,
		SessionId:          sessionId,
	}, nil
}

type CompleteUserPasskeyRegistrationInput struct {
	SessionId                  model.Id
	PasskeyName                string
	CredentialCreationResponse *protocol.CredentialCreationResponse
}

func (s *Session) CompleteUserPasskeyRegistration(ctx context.Context, input CompleteUserPasskeyRegistrationInput) (*model.UserPasskey, UserFacingError) {
	if s.user == nil {
		return nil, AuthorizationError{}
	}

	passkeys, serr := s.GetUserPasskeysByUserId(ctx, s.user.Id)
	if serr != nil {
		return nil, serr
	} else if len(passkeys) >= maxUserPasskeys {
		return nil, NewUserError(fmt.Sprintf("A maximum of %v passkeys are allowed per user. Please remove one to add a new one.", maxUserPasskeys))
	}

	user := webAuthnUser{
		User:     s.user,
		Passkeys: passkeys,
	}

	session, err := s.app.store.GetUserPasskeySessionById(ctx, input.SessionId)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if session == nil || session.Type != model.UserPasskeySessionTypeRegistration {
		return nil, NewUserError("Session expired. Please try again.")
	}

	data, err := input.CredentialCreationResponse.Parse()
	if err != nil {
		return nil, NewUserError("Invalid credential creation response.")
	}

	credential, err := s.app.webAuthn.CreateCredential(user, session.Data, data)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	passkey := &model.UserPasskey{
		Id:           model.NewUserPasskeyId(),
		CreationTime: time.Now(),
		Name:         input.PasskeyName,
		UserId:       s.user.Id,
		Credential:   *credential,
	}
	if err := s.validateUserPasskey(passkey); err != nil {
		return nil, err
	}

	if err := s.app.store.DeleteUserPasskeySessionById(ctx, session.Id); err != nil {
		return nil, s.SanitizedError(err)
	}

	if err := s.app.store.PutUserPasskey(ctx, passkey); err != nil {
		return nil, s.SanitizedError(err)
	}
	return passkey, nil
}

type BeginUserPasskeyAuthenticationOutput struct {
	CredentialAssertion *protocol.CredentialAssertion
	SessionId           model.Id
}

func (s *Session) BeginUserPasskeyAuthentication(ctx context.Context) (*BeginUserPasskeyAuthenticationOutput, UserFacingError) {
	ret, sessionData, err := s.app.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, s.SanitizedError(err)
	}
	sessionId := model.NewUserPasskeySessionId()
	if err := s.app.store.PutUserPasskeySession(ctx, &model.UserPasskeySession{
		Id:   sessionId,
		Type: model.UserPasskeySessionTypeAuthentication,
		Data: *sessionData,
	}); err != nil {
		return nil, s.SanitizedError(err)
	}
	return &BeginUserPasskeyAuthenticationOutput{
		CredentialAssertion: ret,
		SessionId:           sessionId,
	}, nil
}

type CompleteUserPasskeyAuthenticationInput struct {
	SessionId                   model.Id
	CredentialAssertionResponse *protocol.CredentialAssertionResponse
}

// Completes passkey authentication and returns the user that was authenticated.
func (s *Session) CompleteUserPasskeyAuthentication(ctx context.Context, input CompleteUserPasskeyAuthenticationInput) (*model.User, UserFacingError) {
	session, err := s.app.store.GetUserPasskeySessionById(ctx, input.SessionId)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if session == nil || session.Type != model.UserPasskeySessionTypeAuthentication {
		return nil, NewUserError("Session expired. Please try again.")
	}

	data, err := input.CredentialAssertionResponse.Parse()
	if err != nil {
		return nil, NewUserError("Invalid credential assertion response.")
	}

	user, err := s.app.store.GetUserById(ctx, model.Id(data.Response.UserHandle), store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if user == nil {
		return nil, NewUserError("User not found.")
	}

	passkeys, err := s.app.store.GetUserPasskeysByUserId(ctx, user.Id)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	credential, err := s.app.webAuthn.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		if user.Id != model.Id(userHandle) {
			// This shouldn't be possible. `userHandle` should always be the same as
			// `data.Response.UserHandle` and `user.Id`.
			return nil, fmt.Errorf("unexpected user handle in webauthn discoverable user handler")
		}
		return webAuthnUser{
			User:     user,
			Passkeys: passkeys,
		}, nil
	}, session.Data, data)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	if err := s.app.store.DeleteUserPasskeySessionById(ctx, session.Id); err != nil {
		return nil, s.SanitizedError(err)
	}
	for _, passkey := range passkeys {
		if bytes.Equal(passkey.Credential.ID, credential.ID) {
			if _, err := s.app.store.PatchUserPasskeyById(ctx, passkey.Id, &store.UserPasskeyPatch{
				Credential: credential,
			}); err != nil {
				return nil, s.SanitizedError(err)
			}
			break
		}
	}

	return user, nil
}

func (s *Session) DeleteUserPasskeyById(ctx context.Context, id model.Id) UserFacingError {
	passkey, err := s.app.store.GetUserPasskeyById(ctx, id)
	if err != nil {
		return s.SanitizedError(err)
	} else if passkey == nil {
		return nil
	} else if !s.HasUserId(passkey.UserId) {
		return AuthorizationError{}
	}

	return s.SanitizedError(s.app.store.DeleteUserPasskeyById(ctx, id))
}

type UserPasskeyPatch struct {
	Name *string
}

func (s *Session) PatchUserPasskeyById(ctx context.Context, id model.Id, patch UserPasskeyPatch) (*model.UserPasskey, UserFacingError) {
	passkey, err := s.app.store.GetUserPasskeyById(ctx, id)
	if passkey == nil {
		return nil, s.SanitizedError(err)
	}

	if !s.HasUserId(passkey.UserId) {
		return nil, AuthorizationError{}
	}

	possibleOutcome := *passkey
	storePatch := &store.UserPasskeyPatch{}
	if patch.Name != nil {
		possibleOutcome.Name = *patch.Name
		storePatch.Name = patch.Name
	}

	if err := s.validateUserPasskey(&possibleOutcome); err != nil {
		return nil, err
	}

	ret, err := s.app.store.PatchUserPasskeyById(ctx, id, storePatch)
	return ret, s.SanitizedError(err)
}
