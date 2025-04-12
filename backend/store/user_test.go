package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

func TestUser(t *testing.T) {
	s := NewTestStore(t)

	u := &model.User{
		Id:           model.NewUserId(),
		EmailAddress: "foo@bar.com",
	}

	require.NoError(t, s.PutUser(context.Background(), u))

	t.Run("Get", func(t *testing.T) {
		got, err := s.GetUserById(context.Background(), u.Id, store.ConsistencyEventual)
		require.NoError(t, err)
		assert.Equal(t, u, got)

		users, err := s.GetUsersByEmailAddress(context.Background(), u.EmailAddress)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, u, users[0])
	})

	t.Run("Patch", func(t *testing.T) {
		notAUser, err := s.PatchUserById(context.Background(), model.NewUserId(), &store.UserPatch{
			EmailAddress: &u.EmailAddress,
		})
		require.NoError(t, err)
		assert.Nil(t, notAUser)

		u.EmailAddress = "foo2@bar.com"
		newUser, err := s.PatchUserById(context.Background(), u.Id, &store.UserPatch{
			EmailAddress: &u.EmailAddress,
		})
		require.NoError(t, err)
		assert.Equal(t, u, newUser)

		users, err := s.GetUsersByEmailAddress(context.Background(), u.EmailAddress)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, u, users[0])
	})
}
