package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

func TestAWSIntegration(t *testing.T) {
	s := NewTestStore(t)

	integration := &model.AWSIntegration{
		Id:     model.NewAWSIntegrationId(),
		TeamId: model.NewTeamId(),
		Name:   "foo",
	}

	require.NoError(t, s.PutAWSIntegration(context.Background(), integration))

	t.Run("Get", func(t *testing.T) {
		got, err := s.GetAWSIntegrationById(context.Background(), integration.Id)
		require.NoError(t, err)
		assert.Equal(t, integration, got)

		integrations, err := s.GetAWSIntegrationsByTeamId(context.Background(), integration.TeamId)
		require.NoError(t, err)
		assert.Len(t, integrations, 1)
		assert.Equal(t, integration, integrations[0])
	})

	t.Run("Patch", func(t *testing.T) {
		notAAWSIntegration, err := s.PatchAWSIntegrationById(context.Background(), model.NewAWSIntegrationId(), &store.AWSIntegrationPatch{
			Name: &integration.Name,
		})
		require.NoError(t, err)
		assert.Nil(t, notAAWSIntegration)

		integration.Name = "Bar"
		newAWSIntegration, err := s.PatchAWSIntegrationById(context.Background(), integration.Id, &store.AWSIntegrationPatch{
			Name: &integration.Name,
		})
		require.NoError(t, err)
		assert.Equal(t, integration, newAWSIntegration)
	})

	t.Run("Delete", func(t *testing.T) {
		err := s.DeleteAWSIntegrationById(context.Background(), integration.Id)
		require.NoError(t, err)

		integrations, err := s.GetAWSIntegrationsByTeamId(context.Background(), integration.TeamId)
		require.NoError(t, err)
		assert.Empty(t, integrations)
	})
}
