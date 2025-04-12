package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
	"go.uber.org/zap"
)

type BeginUserRegistrationInput struct {
	EmailAddress                    string
	TermsOfServiceAgreementRevision model.UserAgreementRevision
	PrivacyPolicyAgreementRevision  model.UserAgreementRevision
	CookiePolicyAgreementRevision   model.UserAgreementRevision
}

func ValidateEmailAddress(emailAddress string) UserFacingError {
	if parsed, err := mail.ParseAddress(emailAddress); err != nil || parsed.Address != emailAddress {
		return NewUserError("Invalid email address.")
	}
	if len(emailAddress) > 1000 {
		return NewUserError("Please provide a shorter email address.")
	}
	return nil
}

func (a *App) IsUserRegistrationAllowed(emailAddress string) bool {
	if len(a.config.UserRegistrationAllowlist) == 0 {
		return true
	}
	for _, allowlistEmail := range a.config.UserRegistrationAllowlist {
		if match, _ := regexp.MatchString("^"+allowlistEmail+"$", emailAddress); match {
			return true
		}
	}
	return false
}

// Registers the user and sends them a registration email.
func (s *Session) BeginUserRegistration(ctx context.Context, input BeginUserRegistrationInput) UserFacingError {
	if !s.app.IsUserRegistrationAllowed(input.EmailAddress) {
		return NewUserError("User registration is disabled at this time.")
	} else if err := ValidateEmailAddress(input.EmailAddress); err != nil {
		return err
	} else if !input.TermsOfServiceAgreementRevision.IsValid() || !input.PrivacyPolicyAgreementRevision.IsValid() || !input.CookiePolicyAgreementRevision.IsValid() {
		return NewUserError("You must agree to all of the relevant terms and policies.")
	}

	if user, err := s.app.getUserByEmailAddress(ctx, input.EmailAddress); err != nil {
		return s.SanitizedError(err)
	} else if user != nil {
		if err := s.app.Email(ctx, input.EmailAddress, "Registration", "user_registration_for_existing_user_email.html.tmpl", nil); err != nil {
			return s.SanitizedError(fmt.Errorf("unable to send registration email for existing user: %w", err))
		}
		return nil
	}

	token := model.NewToken()
	tokenHash := model.TokenHash(token)

	now := time.Now()
	if err := s.app.store.PutUserRegistrationToken(ctx, &model.UserRegistrationToken{
		EmailAddress:   input.EmailAddress,
		Hash:           tokenHash,
		ExpirationTime: now.Add(48 * time.Hour),
		TermsOfServiceAgreement: model.UserAgreement{
			Revision: input.TermsOfServiceAgreementRevision,
			Time:     now,
		},
		PrivacyPolicyAgreement: model.UserAgreement{
			Revision: input.PrivacyPolicyAgreementRevision,
			Time:     now,
		},
		CookiePolicyAgreement: model.UserAgreement{
			Revision: input.CookiePolicyAgreementRevision,
			Time:     now,
		},
	}); err != nil {
		return s.SanitizedError(fmt.Errorf("unable to put user registration token: %w", err))
	}

	tokenBase64 := base64.RawURLEncoding.EncodeToString(token)

	if err := s.app.Email(ctx, input.EmailAddress, "Registration", "user_registration_email.html.tmpl", map[string]any{
		"TokenBase64": tokenBase64,
	}); err != nil {
		return s.SanitizedError(fmt.Errorf("unable to send registration email: %w", err))
	}

	return nil
}

type CompleteUserRegistrationInput struct {
	Token []byte
}

func (s *Session) CompleteUserRegistration(ctx context.Context, input CompleteUserRegistrationInput) (retUser *model.User, retAccessToken []byte, retErr UserFacingError) {
	hash := model.TokenHash(input.Token)

	token, err := s.app.store.GetUserRegistrationTokenByHash(ctx, hash)
	if err != nil {
		return nil, nil, s.SanitizedError(fmt.Errorf("unable to get user registration token: %w", err))
	} else if token == nil || token.ExpirationTime.Before(time.Now()) {
		return nil, nil, NewUserError("The provided token is invalid or expired. Please register again.")
	}

	if !s.app.IsUserRegistrationAllowed(token.EmailAddress) {
		return nil, nil, NewUserError("User registration is disabled at this time.")
	}

	if user, err := s.app.getUserByEmailAddress(ctx, token.EmailAddress); err != nil {
		return nil, nil, s.SanitizedError(err)
	} else if user != nil {
		return nil, nil, NewUserError("Registration has already been completed for this email address.")
	}

	now := time.Now()
	user := &model.User{
		Id:                      model.NewUserId(),
		CreationTime:            now,
		Role:                    model.UserRoleCustomer,
		EmailAddress:            token.EmailAddress,
		TermsOfServiceAgreement: token.TermsOfServiceAgreement,
		PrivacyPolicyAgreement:  token.PrivacyPolicyAgreement,
		CookiePolicyAgreement:   token.CookiePolicyAgreement,
	}

	accessToken, err := s.app.createUserAccessToken(ctx, user.Id)
	if err != nil {
		return nil, nil, s.SanitizedError(fmt.Errorf("unable to create user access token: %w", err))
	} else if err := s.app.store.PutUser(ctx, user); err != nil {
		return nil, nil, s.SanitizedError(fmt.Errorf("unable to put user: %w", err))
	}

	if err := s.app.store.DeleteUserRegistrationTokenByHash(ctx, hash); err != nil {
		s.Logger().Error("unable to delete user registration token", zap.Error(err))
	}

	return user, accessToken, nil
}

func (s *Session) GetUserByEmailAddress(ctx context.Context, emailAddress string) (*model.User, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		return nil, AuthorizationError{}
	}
	ret, err := s.app.getUserByEmailAddress(ctx, emailAddress)
	return ret, s.SanitizedError(err)
}

func (a *App) getUserByEmailAddress(ctx context.Context, email string) (*model.User, error) {
	if _, err := mail.ParseAddress(email); err != nil || email == "" {
		return nil, NewUserError("Invalid email address.")
	} else if users, err := a.store.GetUsersByEmailAddress(ctx, email); err != nil {
		return nil, fmt.Errorf("unable to get user with email %s: %w", email, err)
	} else if len(users) == 0 {
		return nil, nil
	} else if len(users) != 1 {
		return nil, fmt.Errorf("expected 1 user, got %d for email %s", len(users), email)
	} else {
		return users[0], nil
	}
}

func (s *Session) GetUsers(ctx context.Context) ([]*model.User, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		return nil, AuthorizationError{}
	}

	users, err := s.app.store.GetUsers(ctx)
	return users, s.SanitizedError(err)
}

func (s *Session) GetUserById(ctx context.Context, id model.Id) (*model.User, UserFacingError) {
	if s.HasUserId(id) {
		return s.SanitizeUser(s.User()), nil
	} else if !s.HasUserRole(model.UserRoleAdministrator) {
		return nil, AuthorizationError{}
	}

	user, err := s.app.store.GetUserById(ctx, id, store.ConsistencyEventual)
	return user, s.SanitizedError(err)
}

func (s *Session) CreateUserAccessToken(ctx context.Context) ([]byte, UserFacingError) {
	if s.User() == nil {
		return nil, AuthorizationError{}
	}
	token, err := s.app.createUserAccessToken(ctx, s.User().Id)
	return token, s.SanitizedError(err)
}

func (a *App) createUserAccessToken(ctx context.Context, userId model.Id) ([]byte, error) {
	token := model.NewToken()
	tokenHash := model.TokenHash(token)
	if err := a.store.PutUserAccessToken(ctx, &model.UserAccessToken{
		UserId:         userId,
		CreationTime:   time.Now(),
		Hash:           tokenHash,
		ExpirationTime: time.Now().Add(90 * 24 * time.Hour),
	}); err != nil {
		return nil, fmt.Errorf("unable to put user access token: %w", err)
	}
	return token, nil
}

type UserPatch struct {
	Role                            *model.UserRole
	EmailAddress                    *string
	Password                        *string
	TermsOfServiceAgreementRevision *model.UserAgreementRevision
	PrivacyPolicyAgreementRevision  *model.UserAgreementRevision
	CookiePolicyAgreementRevision   *model.UserAgreementRevision
}

func ValidatePassword(password string) UserFacingError {
	if len(password) < 14 {
		return NewUserError("Password must be at least 14 characters.")
	}
	return nil
}

func (s *Session) PatchUserById(ctx context.Context, userId model.Id, patch UserPatch) (*model.User, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		allowedSelfService := UserPatch{
			Password:                        patch.Password,
			TermsOfServiceAgreementRevision: patch.TermsOfServiceAgreementRevision,
			PrivacyPolicyAgreementRevision:  patch.PrivacyPolicyAgreementRevision,
			CookiePolicyAgreementRevision:   patch.CookiePolicyAgreementRevision,
		}
		if !s.HasUserId(userId) || patch != allowedSelfService {
			return nil, AuthorizationError{}
		}
	}

	now := time.Now()
	storePatch := &store.UserPatch{
		Role:         patch.Role,
		EmailAddress: patch.EmailAddress,
	}

	if patch.Password != nil {
		var hash []byte
		if *patch.Password != "" {
			if err := ValidatePassword(*patch.Password); err != nil {
				return nil, err
			}
			hash = model.EncryptedPasswordHash(*patch.Password, s.app.config.PasswordEncryptionKey)
		}
		storePatch.EncryptedPasswordHash = &hash
	}

	if patch.TermsOfServiceAgreementRevision != nil {
		if !s.HasUserId(userId) {
			// Not even admins can agree on behalf of a user.
			return nil, AuthorizationError{}
		} else if !patch.TermsOfServiceAgreementRevision.IsValid() {
			return nil, NewUserError("Invalid terms of service agreement revision.")
		} else if *patch.TermsOfServiceAgreementRevision > s.User().TermsOfServiceAgreement.Revision {
			storePatch.TermsOfServiceAgreement = &model.UserAgreement{
				Revision: *patch.TermsOfServiceAgreementRevision,
				Time:     now,
			}
		}
	}

	if patch.PrivacyPolicyAgreementRevision != nil {
		if !s.HasUserId(userId) {
			// Not even admins can agree on behalf of a user.
			return nil, AuthorizationError{}
		} else if !patch.PrivacyPolicyAgreementRevision.IsValid() || *patch.PrivacyPolicyAgreementRevision < s.User().PrivacyPolicyAgreement.Revision {
			return nil, NewUserError("Invalid privacy policy agreement revision.")
		} else if *patch.PrivacyPolicyAgreementRevision > s.User().PrivacyPolicyAgreement.Revision {
			storePatch.PrivacyPolicyAgreement = &model.UserAgreement{
				Revision: *patch.PrivacyPolicyAgreementRevision,
				Time:     now,
			}
		}
	}

	if patch.CookiePolicyAgreementRevision != nil {
		if !s.HasUserId(userId) {
			// Not even admins can agree on behalf of a user.
			return nil, AuthorizationError{}
		} else if !patch.CookiePolicyAgreementRevision.IsValid() || *patch.CookiePolicyAgreementRevision < s.User().CookiePolicyAgreement.Revision {
			return nil, NewUserError("Invalid cookie policy agreement revision.")
		} else if *patch.CookiePolicyAgreementRevision > s.User().CookiePolicyAgreement.Revision {
			storePatch.CookiePolicyAgreement = &model.UserAgreement{
				Revision: *patch.CookiePolicyAgreementRevision,
				Time:     now,
			}
		}
	}

	user, err := s.app.store.PatchUserById(ctx, userId, storePatch)
	return user, s.SanitizedError(err)
}

func (a *App) SetUserRole(ctx context.Context, userId model.Id, role model.UserRole) (*model.User, error) {
	user, err := a.store.PatchUserById(ctx, userId, &store.UserPatch{
		Role: &role,
	})
	return user, err
}

func (s *Session) SignOut(ctx context.Context) UserFacingError {
	if s.userAccessToken == nil {
		return NewUserError("An access token was not used to authenticate.")
	}
	return s.SanitizedError(s.app.store.DeleteUserAccessTokenByHash(ctx, s.userAccessToken.Hash))
}

func (s *Session) SanitizeUser(user *model.User) *model.User {
	if s.HasUserId(user.Id) || s.HasUserRole(model.UserRoleAdministrator) {
		return user
	}
	return &model.User{
		Id:           user.Id,
		EmailAddress: user.EmailAddress,
	}
}

type BeginUserEmailAuthenticationInput struct {
	EmailAddress string
}

func (s *Session) BeginUserEmailAuthentication(ctx context.Context, input BeginUserEmailAuthenticationInput) UserFacingError {
	if err := ValidateEmailAddress(input.EmailAddress); err != nil {
		return err
	}

	user, err := s.app.getUserByEmailAddress(ctx, input.EmailAddress)
	if err != nil {
		return s.SanitizedError(err)
	}

	if user == nil {
		// Pretend we're doing something.
		time.Sleep(100 * time.Millisecond)
		// TODO: send registration email?
	} else {
		token := model.NewToken()
		tokenHash := model.TokenHash(token)

		now := time.Now()
		if err := s.app.store.PutUserEmailAuthenticationToken(ctx, &model.UserEmailAuthenticationToken{
			UserId:         user.Id,
			Hash:           tokenHash,
			ExpirationTime: now.Add(1 * time.Hour),
		}); err != nil {
			return s.SanitizedError(fmt.Errorf("unable to put user email authentication token: %w", err))
		}

		tokenBase64 := base64.RawURLEncoding.EncodeToString(token)

		if err := s.app.Email(ctx, input.EmailAddress, "Sign-In Link", "user_email_authentication_email.html.tmpl", map[string]any{
			"TokenBase64": tokenBase64,
		}); err != nil {
			return s.SanitizedError(fmt.Errorf("unable to send email authentication email: %w", err))
		}
	}

	return nil
}
