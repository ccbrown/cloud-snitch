package store

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedUser struct {
	*model.User

	PrimaryIndex
	ByteByteIndex1
}

func (s *Store) PutUser(ctx context.Context, user *model.User) error {
	return s.put(ctx, &IndexedUser{
		User: user,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("user:" + user.Id),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("user_email_addresses"),
			RangeKey: []byte(strings.ToLower(user.EmailAddress)),
		},
	})
}

type UserPatch struct {
	EmailAddress          *string
	Role                  *model.UserRole
	EncryptedPasswordHash *[]byte

	TermsOfServiceAgreement *model.UserAgreement
	PrivacyPolicyAgreement  *model.UserAgreement
	CookiePolicyAgreement   *model.UserAgreement
}

func (p *UserPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.EmailAddress != nil {
		update = update.Set(expression.Name("_bb1r"), expression.Value([]byte(strings.ToLower(*p.EmailAddress))))
		update = update.Set(expression.Name("EmailAddress"), expression.Value(p.EmailAddress))
	}
	if p.Role != nil {
		update = update.Set(expression.Name("Role"), expression.Value(p.Role))
	}
	if p.EncryptedPasswordHash != nil {
		update = update.Set(expression.Name("EncryptedPasswordHash"), expression.Value(p.EncryptedPasswordHash))
	}
	if p.TermsOfServiceAgreement != nil {
		update = update.Set(expression.Name("TermsOfServiceAgreement"), expression.Value(*p.TermsOfServiceAgreement))
	}
	if p.PrivacyPolicyAgreement != nil {
		update = update.Set(expression.Name("PrivacyPolicyAgreement"), expression.Value(*p.PrivacyPolicyAgreement))
	}
	if p.CookiePolicyAgreement != nil {
		update = update.Set(expression.Name("CookiePolicyAgreement"), expression.Value(*p.CookiePolicyAgreement))
	}
	return update
}

func (s *Store) GetUserById(ctx context.Context, id model.Id, consistency Consistency) (*model.User, error) {
	return getByPrimaryKey[model.User](ctx, s, []byte("user:"+id), consistency)
}

func (s *Store) GetUsersByIds(ctx context.Context, ids ...model.Id) ([]*model.User, error) {
	return getByPrimaryKeys[model.User](ctx, s, prefixIds("user:", ids)...)
}

func (s *Store) GetUsers(ctx context.Context) ([]*model.User, error) {
	return getAllByHashKey[model.User](ctx, s, "_bb1", "_bb1h", []byte("user_email_addresses"))
}

func (s *Store) GetUsersByEmailAddress(ctx context.Context, emailAddress string) ([]*model.User, error) {
	return getAllByHashAndRangeKey[model.User](ctx, s, "_bb1", []byte("user_email_addresses"), []byte(strings.ToLower(emailAddress)))
}

func (s *Store) PatchUserById(ctx context.Context, id model.Id, patch *UserPatch) (*model.User, error) {
	update := patch.Apply(expression.UpdateBuilder{})
	return updateByPrimaryKey[model.User](ctx, s, []byte("user:"+id), update)
}

type IndexedUserRegistrationToken struct {
	*model.UserRegistrationToken

	PrimaryIndex

	TTL
}

func (s *Store) PutUserRegistrationToken(ctx context.Context, token *model.UserRegistrationToken) error {
	return s.put(ctx, &IndexedUserRegistrationToken{
		UserRegistrationToken: token,
		PrimaryIndex: PrimaryIndex{
			HashKey:  append([]byte("user_registration_token:"), token.Hash...),
			RangeKey: []byte("_"),
		},
		TTL: NewTTL(token.ExpirationTime),
	})
}

func (s *Store) GetUserRegistrationTokenByHash(ctx context.Context, hash []byte) (*model.UserRegistrationToken, error) {
	return getByPrimaryKey[model.UserRegistrationToken](ctx, s, append([]byte("user_registration_token:"), hash...), ConsistencyStrongInRegion)
}

func (s *Store) DeleteUserRegistrationTokenByHash(ctx context.Context, hash []byte) error {
	return deleteByPrimaryKey(ctx, s, append([]byte("user_registration_token:"), hash...))
}

type IndexedUserAccessToken struct {
	*model.UserAccessToken

	PrimaryIndex
	ByteByteIndex1

	TTL
}

func (s *Store) PutUserAccessToken(ctx context.Context, token *model.UserAccessToken) error {
	return s.put(ctx, &IndexedUserAccessToken{
		UserAccessToken: token,
		PrimaryIndex: PrimaryIndex{
			HashKey:  append([]byte("user_access_token:"), token.Hash...),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("user_access_token:" + token.UserId.String()),
			RangeKey: []byte("_"),
		},
		TTL: NewTTL(token.ExpirationTime),
	})
}

func (s *Store) GetUserAccessTokenByHash(ctx context.Context, hash []byte) (*model.UserAccessToken, error) {
	return getByPrimaryKey[model.UserAccessToken](ctx, s, append([]byte("user_access_token:"), hash...), ConsistencyStrongInRegion)
}

func (s *Store) DeleteUserAccessTokenByHash(ctx context.Context, hash []byte) error {
	return deleteByPrimaryKey(ctx, s, append([]byte("user_access_token:"), hash...))
}

type IndexedUserEmailAuthenticationToken struct {
	*model.UserEmailAuthenticationToken

	PrimaryIndex
	ByteByteIndex1

	TTL
}

func (s *Store) PutUserEmailAuthenticationToken(ctx context.Context, token *model.UserEmailAuthenticationToken) error {
	return s.put(ctx, &IndexedUserEmailAuthenticationToken{
		UserEmailAuthenticationToken: token,
		PrimaryIndex: PrimaryIndex{
			HashKey:  append([]byte("user_email_authentication_token:"), token.Hash...),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("user_email_authentication_token:" + token.UserId.String()),
			RangeKey: []byte("_"),
		},
		TTL: NewTTL(token.ExpirationTime),
	})
}

func (s *Store) GetUserEmailAuthenticationTokenByHash(ctx context.Context, hash []byte) (*model.UserEmailAuthenticationToken, error) {
	return getByPrimaryKey[model.UserEmailAuthenticationToken](ctx, s, append([]byte("user_email_authentication_token:"), hash...), ConsistencyStrongInRegion)
}

func (s *Store) DeleteUserEmailAuthenticationTokenByHash(ctx context.Context, hash []byte) error {
	return deleteByPrimaryKey(ctx, s, append([]byte("user_email_authentication_token:"), hash...))
}
