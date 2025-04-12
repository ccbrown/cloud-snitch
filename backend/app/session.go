package app

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
	"github.com/go-webauthn/webauthn/protocol"
)

type Session struct {
	user            *model.User
	userAccessToken *model.UserAccessToken
	app             *App
	logger          *zap.Logger
}

func (s *Session) RequireUser() UserFacingError {
	if s.user != nil {
		return nil
	}
	return AuthorizationError{}
}

func (s *Session) RequireTeamMember(ctx context.Context, teamId model.Id) UserFacingError {
	if err := s.RequireUser(); err != nil {
		return err
	} else if s.user.Role == model.UserRoleAdministrator {
		// Just make sure the team exists.
		if team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyEventual); err != nil {
			return s.SanitizedError(err)
		} else if team == nil {
			return AuthorizationError{}
		}
		return nil
	} else if membership, err := s.app.store.GetTeamMembershipByTeamAndUserId(ctx, teamId, s.user.Id); err != nil {
		return s.SanitizedError(err)
	} else if membership == nil {
		return AuthorizationError{}
	}
	return nil
}

func (s *Session) RequireTeamAdministrator(ctx context.Context, teamId model.Id) UserFacingError {
	if err := s.RequireUser(); err != nil {
		return err
	} else if s.user.Role == model.UserRoleAdministrator {
		// Just make sure the team exists.
		if team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyEventual); err != nil {
			return s.SanitizedError(err)
		} else if team == nil {
			return AuthorizationError{}
		}
		return nil
	} else if membership, err := s.app.store.GetTeamMembershipByTeamAndUserId(ctx, teamId, s.user.Id); err != nil {
		return s.SanitizedError(err)
	} else if membership == nil || membership.Role != model.TeamMembershipRoleAdministrator {
		return AuthorizationError{}
	}
	return nil
}

func (s *Session) HasUserRole(role model.UserRole) bool {
	return s.user != nil && s.user.Role == role
}

func (s *Session) HasUserId(userId model.Id) bool {
	return s.user != nil && s.user.Id == userId
}

func (a *App) NewAnonymousSession() *Session {
	return &Session{
		app:    a,
		logger: zap.L(),
	}
}

func (sess *Session) Logger() *zap.Logger {
	return sess.logger
}

func (sess *Session) User() *model.User {
	return sess.user
}

func (sess Session) WithLogFields(fields ...zap.Field) *Session {
	sess.logger = sess.logger.With(fields...)
	return &sess
}

// NewUserSession returns a context with an associated user, if the email and password are valid.
// Otherwise, it returns nil.
func (sess Session) WithUserCredentials(ctx context.Context, email, password string) (*Session, UserFacingError) {
	if users, err := sess.app.store.GetUsersByEmailAddress(ctx, email); len(users) == 0 || err != nil {
		return nil, sess.SanitizedError(err)
	} else if len(users) != 1 {
		return nil, sess.SanitizedError(fmt.Errorf("expected 1 user, got %d for email: %s", len(users), email))
	} else if !users[0].VerifyPassword(password, sess.app.config.PasswordEncryptionKey) {
		return nil, nil
	} else {
		sess.user = users[0]
		return &sess, nil
	}
}

// NewUserSession returns a context with an associated user, if the access token is valid.
// Otherwise, it returns nil.
func (sess Session) WithUserAccessToken(ctx context.Context, token []byte) (*Session, UserFacingError) {
	hash := model.TokenHash(token)
	if accessToken, err := sess.app.store.GetUserAccessTokenByHash(ctx, hash); accessToken == nil || accessToken.ExpirationTime.Before(time.Now()) || err != nil {
		return nil, sess.SanitizedError(err)
	} else if user, err := sess.app.store.GetUserById(ctx, accessToken.UserId, store.ConsistencyEventual); user == nil || err != nil {
		return nil, sess.SanitizedError(err)
	} else {
		sess.user = user
		sess.userAccessToken = accessToken
		return &sess, nil
	}
}

// NewUserSession returns a context with an associated user, if passkey authentication can be
// completed using the given data. Otherwise, it returns nil.
func (sess Session) WithUserPasskey(ctx context.Context, sessionId model.Id, assertionResponse *protocol.CredentialAssertionResponse) (*Session, UserFacingError) {
	if user, err := sess.CompleteUserPasskeyAuthentication(ctx, CompleteUserPasskeyAuthenticationInput{
		SessionId:                   sessionId,
		CredentialAssertionResponse: assertionResponse,
	}); user == nil {
		return nil, sess.SanitizedError(err)
	} else {
		sess.user = user
		return &sess, nil
	}
}

// NewUserSession returns a context with an associated user, if passkey authentication can be
// completed using the given data. Otherwise, it returns nil.
func (sess Session) WithUserEmailAuthentication(ctx context.Context, token []byte) (*Session, UserFacingError) {
	hash := model.TokenHash(token)
	if accessToken, err := sess.app.store.GetUserEmailAuthenticationTokenByHash(ctx, hash); accessToken == nil || accessToken.ExpirationTime.Before(time.Now()) || err != nil {
		return nil, sess.SanitizedError(err)
	} else if user, err := sess.app.store.GetUserById(ctx, accessToken.UserId, store.ConsistencyEventual); user == nil || err != nil {
		return nil, sess.SanitizedError(err)
	} else if err := sess.app.store.DeleteUserEmailAuthenticationTokenByHash(ctx, hash); err != nil {
		return nil, sess.SanitizedError(err)
	} else {
		sess.user = user
		return &sess, nil
	}
}
