package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

func TestTeam(t *testing.T) {
	s := NewTestStore(t)

	team := &model.Team{
		Id:   model.NewTeamId(),
		Name: "foo",
	}

	require.NoError(t, s.PutTeam(context.Background(), team))

	t.Run("Get", func(t *testing.T) {
		got, err := s.GetTeamById(context.Background(), team.Id, store.ConsistencyStrongInRegion)
		require.NoError(t, err)
		assert.Equal(t, team, got)

		teams, err := s.GetTeams(context.Background())
		require.NoError(t, err)
		assert.Len(t, teams, 1)
		assert.Equal(t, team, teams[0])
	})

	t.Run("Patch", func(t *testing.T) {
		notATeam, err := s.PatchTeamById(context.Background(), model.NewTeamId(), &store.TeamPatch{
			Name: &team.Name,
		})
		require.NoError(t, err)
		assert.Nil(t, notATeam)

		team.Name = "Bar"
		newTeam, err := s.PatchTeamById(context.Background(), team.Id, &store.TeamPatch{
			Name: &team.Name,
		})
		require.NoError(t, err)
		assert.Equal(t, team, newTeam)
	})
}

func TestTeamMembership(t *testing.T) {
	s := NewTestStore(t)

	membership := &model.TeamMembership{
		TeamId: model.NewTeamId(),
		UserId: model.NewUserId(),
		Role:   model.TeamMembershipRoleAdministrator,
	}

	require.NoError(t, s.PutTeamMembership(context.Background(), membership))

	t.Run("GetByTeamAndUserId", func(t *testing.T) {
		got, err := s.GetTeamMembershipByTeamAndUserId(context.Background(), membership.TeamId, membership.UserId)
		require.NoError(t, err)
		assert.Equal(t, membership, got)
	})

	t.Run("GetByTeamId", func(t *testing.T) {
		memberships, err := s.GetTeamMembershipsByTeamId(context.Background(), membership.TeamId)
		require.NoError(t, err)
		assert.Len(t, memberships, 1)
		assert.Equal(t, membership, memberships[0])
	})

	t.Run("GetByUserId", func(t *testing.T) {
		memberships, err := s.GetTeamMembershipsByUserId(context.Background(), membership.UserId)
		require.NoError(t, err)
		assert.Len(t, memberships, 1)
		assert.Equal(t, membership, memberships[0])
	})
}

func TestTeamInvite(t *testing.T) {
	s := NewTestStore(t)

	invite := &model.TeamInvite{
		TeamId:       model.NewTeamId(),
		EmailAddress: "Foo@Example.COM",
	}

	require.NoError(t, s.PutTeamInvite(context.Background(), invite))

	t.Run("GetByTeamIdAndEmailAddress", func(t *testing.T) {
		got, err := s.GetTeamInviteByTeamIdAndEmailAddress(context.Background(), invite.TeamId, invite.EmailAddress)
		require.NoError(t, err)
		assert.Equal(t, invite, got)
	})

	t.Run("GetByTeamId", func(t *testing.T) {
		invites, err := s.GetTeamInvitesByTeamId(context.Background(), invite.TeamId)
		require.NoError(t, err)
		assert.Len(t, invites, 1)
		assert.Equal(t, invite, invites[0])
	})

	t.Run("GetByEmailAddress", func(t *testing.T) {
		invites, err := s.GetTeamInvitesByEmailAddress(context.Background(), "foo@EXAMPLE.com")
		require.NoError(t, err)
		assert.Len(t, invites, 1)
		assert.Equal(t, invite, invites[0])
	})
}
