package store

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedTeam struct {
	*model.Team

	PrimaryIndex
	ByteByteIndex1
}

func (s *Store) PutTeam(ctx context.Context, team *model.Team) error {
	return s.put(ctx, &IndexedTeam{
		Team: team,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("team:" + team.Id),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("team_ids"),
			RangeKey: []byte(team.Id),
		},
	})
}

func (s *Store) GetTeamById(ctx context.Context, id model.Id, consistency Consistency) (*model.Team, error) {
	return getByPrimaryKey[model.Team](ctx, s, []byte("team:"+id), consistency)
}

func (s *Store) GetTeamsByIds(ctx context.Context, ids ...model.Id) ([]*model.Team, error) {
	return getByPrimaryKeys[model.Team](ctx, s, prefixIds("team:", ids)...)
}

func (s *Store) GetTeams(ctx context.Context) ([]*model.Team, error) {
	return getAllByHashKey[model.Team](ctx, s, "_bb1", "_bb1h", []byte("team_ids"))
}

type TeamPatch struct {
	Name             *string
	StripeCustomerId *string
	Entitlements     *model.TeamEntitlements
}

func (p *TeamPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.Name != nil {
		update = update.Set(expression.Name("Name"), expression.Value(p.Name))
	}
	if p.StripeCustomerId != nil {
		update = update.Set(expression.Name("StripeCustomerId"), expression.Value(p.StripeCustomerId))
	}
	if p.Entitlements != nil {
		update = update.Set(expression.Name("Entitlements"), expression.Value(p.Entitlements))
	}
	return update
}

func (s *Store) PatchTeamById(ctx context.Context, id model.Id, patch *TeamPatch) (*model.Team, error) {
	update := patch.Apply(expression.UpdateBuilder{})
	return updateByPrimaryKey[model.Team](ctx, s, []byte("team:"+id), update)
}

type IndexedTeamMembership struct {
	*model.TeamMembership

	PrimaryIndex
	ByteByteIndex1
	ByteByteIndex2
}

func (s *Store) PutTeamMembership(ctx context.Context, membership *model.TeamMembership) error {
	return s.put(ctx, &IndexedTeamMembership{
		TeamMembership: membership,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("team_membership:" + membership.TeamId + ":" + membership.UserId),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("team_memberships:" + membership.TeamId),
			RangeKey: []byte(membership.UserId),
		},
		ByteByteIndex2: ByteByteIndex2{
			HashKey:  []byte("team_memberships:" + membership.UserId),
			RangeKey: []byte(membership.TeamId),
		},
	})
}

func (s *Store) GetTeamMembershipByTeamAndUserId(ctx context.Context, teamId, userId model.Id) (*model.TeamMembership, error) {
	return getByPrimaryKey[model.TeamMembership](ctx, s, []byte("team_membership:"+teamId+":"+userId), ConsistencyEventual)
}

func (s *Store) GetTeamMembershipsByTeamId(ctx context.Context, teamId model.Id) ([]*model.TeamMembership, error) {
	return getAllByHashKey[model.TeamMembership](ctx, s, "_bb1", "_bb1h", []byte("team_memberships:"+teamId))
}

func (s *Store) GetTeamMembershipsByUserId(ctx context.Context, userId model.Id) ([]*model.TeamMembership, error) {
	return getAllByHashKey[model.TeamMembership](ctx, s, "_bb2", "_bb2h", []byte("team_memberships:"+userId))
}

func (s *Store) DeleteTeamMembershipByTeamAndUserId(ctx context.Context, teamId, userId model.Id) error {
	return deleteByPrimaryKey(ctx, s, []byte("team_membership:"+teamId+":"+userId))
}

type TeamMembershipPatch struct {
	Role *model.TeamMembershipRole
}

func (p *TeamMembershipPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.Role != nil {
		update = update.Set(expression.Name("Role"), expression.Value(p.Role))
	}
	return update
}

func (s *Store) PatchTeamMembershipByTeamAndUserId(ctx context.Context, teamId, userId model.Id, patch *TeamMembershipPatch) (*model.TeamMembership, error) {
	update := patch.Apply(expression.UpdateBuilder{})
	return updateByPrimaryKey[model.TeamMembership](ctx, s, []byte("team_membership:"+teamId+":"+userId), update)
}

type IndexedTeamInvite struct {
	*model.TeamInvite

	PrimaryIndex
	ByteByteIndex1
	ByteByteIndex2

	TTL
}

func (s *Store) PutTeamInvite(ctx context.Context, invite *model.TeamInvite) error {
	return s.put(ctx, &IndexedTeamInvite{
		TeamInvite: invite,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("team_invite:" + invite.TeamId.String() + ":" + strings.ToLower(invite.EmailAddress)),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("team_invites:" + invite.TeamId),
			RangeKey: []byte(strings.ToLower(invite.EmailAddress)),
		},
		ByteByteIndex2: ByteByteIndex2{
			HashKey:  []byte("team_invites:" + strings.ToLower(invite.EmailAddress)),
			RangeKey: []byte(invite.TeamId),
		},
		TTL: NewTTL(invite.ExpirationTime.Add(30 * time.Hour)),
	})
}

func (s *Store) GetTeamInviteByTeamIdAndEmailAddress(ctx context.Context, teamId model.Id, emailAddress string) (*model.TeamInvite, error) {
	return getByPrimaryKey[model.TeamInvite](ctx, s, []byte("team_invite:"+teamId.String()+":"+strings.ToLower(emailAddress)), ConsistencyEventual)
}

func (s *Store) GetTeamInvitesByTeamId(ctx context.Context, teamId model.Id) ([]*model.TeamInvite, error) {
	return getAllByHashKey[model.TeamInvite](ctx, s, "_bb1", "_bb1h", []byte("team_invites:"+teamId))
}

func (s *Store) GetTeamInvitesByEmailAddress(ctx context.Context, emailAddress string) ([]*model.TeamInvite, error) {
	return getAllByHashKey[model.TeamInvite](ctx, s, "_bb2", "_bb2h", []byte("team_invites:"+strings.ToLower(emailAddress)))
}

func (s *Store) DeleteTeamInviteByTeamIdAndEmailAddress(ctx context.Context, teamId model.Id, emailAddress string) error {
	return deleteByPrimaryKey(ctx, s, []byte("team_invite:"+teamId.String()+":"+strings.ToLower(emailAddress)))
}

type IndexedTeamPrincipalSettings struct {
	*model.TeamPrincipalSettings

	PrimaryIndex
	ByteByteIndex1
}

type TeamPrincipalSettingsPatch struct {
	Description *string
}

func (p *TeamPrincipalSettingsPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.Description != nil {
		update = update.Set(expression.Name("Description"), expression.Value(p.Description))
	}
	return update
}

func (s *Store) CreateOrPatchTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx context.Context, teamId model.Id, principalKey string, patch *TeamPrincipalSettingsPatch) (*model.TeamPrincipalSettings, error) {
	hk := []byte("team_principal_settings:" + teamId.String() + ":" + principalKey)
	update := patch.Apply(
		expression.UpdateBuilder{}.
			Set(expression.Name("_bb1h"), expression.Value([]byte("principal_settings_by_team_id:"+teamId))).
			Set(expression.Name("_bb1r"), expression.Value([]byte(principalKey))).
			Set(expression.Name("TeamId"), expression.Value(teamId)).
			Set(expression.Name("PrincipalId"), expression.Value(principalKey)),
	)
	return createOrUpdateByPrimaryKey[model.TeamPrincipalSettings](ctx, s, hk, update)
}

func (s *Store) GetTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx context.Context, teamId model.Id, principalKey string) (*model.TeamPrincipalSettings, error) {
	return getByPrimaryKey[model.TeamPrincipalSettings](ctx, s, []byte("team_principal_settings:"+teamId.String()+":"+principalKey), ConsistencyEventual)
}
