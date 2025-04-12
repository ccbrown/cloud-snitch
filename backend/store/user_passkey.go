package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedUserPasskey struct {
	*model.UserPasskey

	PrimaryIndex
	ByteByteIndex1
}

type UserPasskeyPatch struct {
	Name       *string
	Credential *webauthn.Credential
}

func (p *UserPasskeyPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.Name != nil {
		update = update.Set(expression.Name("Name"), expression.Value(*p.Name))
	}
	if p.Credential != nil {
		update = update.Set(expression.Name("Credential"), expression.Value(*p.Credential))
	}
	return update
}

func (s *Store) PatchUserPasskeyById(ctx context.Context, id model.Id, patch *UserPasskeyPatch) (*model.UserPasskey, error) {
	update := patch.Apply(expression.UpdateBuilder{})
	return updateByPrimaryKey[model.UserPasskey](ctx, s, []byte("user_passkey:"+id), update)
}

func (s *Store) PutUserPasskey(ctx context.Context, passkey *model.UserPasskey) error {
	return s.put(ctx, &IndexedUserPasskey{
		UserPasskey: passkey,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("user_passkey:" + passkey.Id),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte([]byte("user_passkeys:" + passkey.UserId)),
			RangeKey: []byte(passkey.Id),
		},
	})
}

func (s *Store) GetUserPasskeyById(ctx context.Context, id model.Id) (*model.UserPasskey, error) {
	return getByPrimaryKey[model.UserPasskey](ctx, s, []byte("user_passkey:"+id), ConsistencyStrongInRegion)
}

func (s *Store) GetUserPasskeysByUserId(ctx context.Context, userId model.Id) ([]*model.UserPasskey, error) {
	return getAllByHashKey[model.UserPasskey](ctx, s, "_bb1", "_bb1h", []byte("user_passkeys:"+userId))
}

func (s *Store) DeleteUserPasskeyById(ctx context.Context, id model.Id) error {
	return deleteByPrimaryKey(ctx, s, []byte("user_passkey:"+id))
}

type IndexedUserPasskeySession struct {
	*model.UserPasskeySession

	PrimaryIndex

	TTL
}

func (s *Store) PutUserPasskeySession(ctx context.Context, session *model.UserPasskeySession) error {
	return s.put(ctx, &IndexedUserPasskeySession{
		UserPasskeySession: session,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("user_passkey_session:" + session.Id),
			RangeKey: []byte("_"),
		},
		TTL: NewTTL(session.Data.Expires),
	})
}

func (s *Store) GetUserPasskeySessionById(ctx context.Context, id model.Id) (*model.UserPasskeySession, error) {
	return getByPrimaryKey[model.UserPasskeySession](ctx, s, []byte("user_passkey_session:"+id), ConsistencyStrongInRegion)
}

func (s *Store) DeleteUserPasskeySessionById(ctx context.Context, id model.Id) error {
	return deleteByPrimaryKey(ctx, s, []byte("user_passkey_session:"+id))
}
