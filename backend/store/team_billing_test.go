package store_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestTeamBillableAccount(t *testing.T) {
	s := NewTestStore(t)

	const n = 123
	teamId := model.NewTeamId()

	for i := 0; i < n; i++ {
		require.NoError(t, s.PutTeamBillableAccount(context.Background(), &model.TeamBillableAccount{
			Id:             strconv.Itoa(i),
			TeamId:         teamId,
			ExpirationTime: time.Now().Add(24 * time.Hour),
		}))
	}

	count, err := s.GetTeamBillableAccountCountByTeamId(context.Background(), teamId)
	require.NoError(t, err)
	assert.Equal(t, n, count)
}
