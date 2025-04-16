package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestReport(t *testing.T) {
	s := NewTestStore(t)

	r := &model.Report{
		Id:     model.NewReportId(),
		TeamId: model.NewTeamId(),
	}

	require.NoError(t, s.PutReport(context.Background(), r))

	t.Run("Get", func(t *testing.T) {
		got, err := s.GetReportById(context.Background(), r.Id)
		require.NoError(t, err)
		assert.Equal(t, r, got)

		reports, err := s.GetReportsByTeamId(context.Background(), r.TeamId)
		require.NoError(t, err)
		assert.Len(t, reports, 1)
	})

	t.Run("Delete", func(t *testing.T) {
		require.NoError(t, s.DeleteReportById(context.Background(), r.Id))

		got, err := s.GetReportById(context.Background(), r.Id)
		require.NoError(t, err)
		assert.Nil(t, got)

		reports, err := s.GetReportsByTeamId(context.Background(), r.TeamId)
		require.NoError(t, err)
		assert.Len(t, reports, 0)
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		var ids []model.Id
		for i := 0; i < 120; i++ {
			id := model.NewReportId()
			ids = append(ids, id)
			require.NoError(t, s.PutReport(context.Background(), &model.Report{
				Id:     id,
				TeamId: r.TeamId,
			}))
		}

		reports, err := s.GetReportsByTeamId(context.Background(), r.TeamId)
		require.NoError(t, err)
		assert.Len(t, reports, 120)

		require.NoError(t, s.DeleteReportsByIds(context.Background(), ids...))

		reports, err = s.GetReportsByTeamId(context.Background(), r.TeamId)
		require.NoError(t, err)
		assert.Empty(t, reports)
	})
}
