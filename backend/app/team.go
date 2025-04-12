package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

type CreateTeamInput struct {
	Name string
}

func (s *Session) CreateTeam(ctx context.Context, input CreateTeamInput) (*model.Team, UserFacingError) {
	if err := s.RequireUser(); err != nil {
		return nil, err
	} else if err := ValidateName(input.Name); err != nil {
		return nil, err
	}

	team := &model.Team{
		Id:           model.NewTeamId(),
		CreationTime: time.Now(),
		Name:         input.Name,
	}

	membership := &model.TeamMembership{
		TeamId:       team.Id,
		UserId:       s.user.Id,
		CreationTime: time.Now(),
		Role:         model.TeamMembershipRoleAdministrator,
	}

	if err := s.app.store.PutTeam(ctx, team); err != nil {
		return nil, s.SanitizedError(err)
	} else if err := s.app.store.PutTeamMembership(ctx, membership); err != nil {
		return nil, s.SanitizedError(err)
	}

	return team, nil
}

type TeamPatch struct {
	Name *string
}

func (s *Session) PatchTeamById(ctx context.Context, teamId model.Id, patch TeamPatch) (*model.Team, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	storePatch := &store.TeamPatch{
		Name: patch.Name,
	}

	if patch.Name != nil {
		if err := ValidateName(*patch.Name); err != nil {
			return nil, err
		}
	}

	team, err := s.app.store.PatchTeamById(ctx, teamId, storePatch)
	return team, s.SanitizedError(err)
}

func (s *Session) GetTeamById(ctx context.Context, teamId model.Id) (*model.Team, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}
	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyEventual)
	return team, s.SanitizedError(err)
}

func (s *Session) GetTeams(ctx context.Context) ([]*model.Team, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) {
		return nil, AuthorizationError{}
	}
	teams, err := s.app.store.GetTeams(ctx)
	return teams, s.SanitizedError(err)
}

type InviteToTeamInput struct {
	TeamId       model.Id
	EmailAddress string
	Role         model.TeamMembershipRole
}

func (s *Session) InviteToTeam(ctx context.Context, input InviteToTeamInput) UserFacingError {
	if err := s.RequireTeamAdministrator(ctx, input.TeamId); err != nil {
		return err
	} else if err := ValidateEmailAddress(input.EmailAddress); err != nil {
		return err
	}

	team, err := s.app.store.GetTeamById(ctx, input.TeamId, store.ConsistencyEventual)
	if err != nil {
		return s.SanitizedError(err)
	} else if team == nil {
		return NotFoundError("Team not found.")
	} else if !team.Entitlements.TeamFeatures {
		return NewUserError("Please upgrade your subscription to invite additional team members.")
	}

	users, err := s.app.store.GetUsersByEmailAddress(ctx, input.EmailAddress)
	if err != nil {
		return s.SanitizedError(err)
	}

	for _, user := range users {
		if membership, err := s.app.store.GetTeamMembershipByTeamAndUserId(ctx, team.Id, user.Id); err != nil {
			return s.SanitizedError(err)
		} else if membership != nil {
			return NewUserError("User is already a member of the team.")
		}
	}

	invite := &model.TeamInvite{
		TeamId:         team.Id,
		SenderId:       s.user.Id,
		EmailAddress:   input.EmailAddress,
		Role:           input.Role,
		CreationTime:   time.Now(),
		ExpirationTime: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.app.store.PutTeamInvite(ctx, invite); err != nil {
		return s.SanitizedError(err)
	}

	emailParams := map[string]any{
		"TeamName": team.Name,
	}

	if len(users) > 0 {
		if err := s.app.Email(ctx, input.EmailAddress, "Team Invite", "team_invite_for_existing_user_email.html.tmpl", emailParams); err != nil {
			return s.SanitizedError(fmt.Errorf("unable to send team invite email for existing user: %w", err))
		}
	} else {
		if err := s.app.Email(ctx, input.EmailAddress, "Team Invite", "team_invite_for_new_user_email.html.tmpl", emailParams); err != nil {
			return s.SanitizedError(fmt.Errorf("unable to send team invite email for new user: %w", err))
		}
	}

	return nil
}

func (s *Session) JoinTeam(ctx context.Context, teamId model.Id) (*model.TeamMembership, UserFacingError) {
	if err := s.RequireUser(); err != nil {
		return nil, err
	}

	if team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyEventual); err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if !team.Entitlements.TeamFeatures {
		return nil, NewUserError("This team's subscription does not allow additional members. Please contact team administrator about upgrading their subscription.")
	}

	invite, err := s.app.store.GetTeamInviteByTeamIdAndEmailAddress(ctx, teamId, s.user.EmailAddress)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if invite == nil || invite.ExpirationTime.Before(time.Now()) {
		return nil, NewUserError("The invitation to this team is invalid or expired. Please ask a team administrator to resend the invitation.")
	}

	membership := &model.TeamMembership{
		TeamId:       teamId,
		UserId:       s.user.Id,
		Role:         invite.Role,
		CreationTime: time.Now(),
	}

	if err := s.app.store.PutTeamMembership(ctx, membership); err != nil {
		return nil, s.SanitizedError(err)
	} else if err := s.app.store.DeleteTeamInviteByTeamIdAndEmailAddress(ctx, teamId, s.user.EmailAddress); err != nil {
		return nil, s.SanitizedError(err)
	}

	return membership, nil
}

func (s *Session) DeleteTeamInviteByTeamIdAndEmailAddress(ctx context.Context, teamId model.Id, emailAddress string) UserFacingError {
	if err := s.RequireUser(); err != nil {
		return err
	}

	if strings.ToLower(emailAddress) != strings.ToLower(s.user.EmailAddress) {
		if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
			return err
		}
	}

	return s.SanitizedError(s.app.store.DeleteTeamInviteByTeamIdAndEmailAddress(ctx, teamId, emailAddress))
}

type UserTeamMembership struct {
	Team       *model.Team
	Membership *model.TeamMembership
}

// Gets team memberships for the user with the given id.
func (s *Session) GetTeamMembershipsByUserId(ctx context.Context, userId model.Id) ([]*UserTeamMembership, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) && !s.HasUserId(userId) {
		return nil, AuthorizationError{}
	}

	memberships, err := s.app.store.GetTeamMembershipsByUserId(ctx, userId)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	teamIds := make([]model.Id, 0, len(memberships))
	for _, membership := range memberships {
		teamIds = append(teamIds, membership.TeamId)
	}

	teams, err := s.app.store.GetTeamsByIds(ctx, teamIds...)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	teamsById := make(map[model.Id]*model.Team, len(teams))
	for _, team := range teams {
		teamsById[team.Id] = team
	}

	ret := make([]*UserTeamMembership, 0, len(memberships))
	for _, membership := range memberships {
		if team, ok := teamsById[membership.TeamId]; ok {
			ret = append(ret, &UserTeamMembership{
				Team:       team,
				Membership: membership,
			})
		}
	}

	return ret, nil
}

// Gets the membership of the current user for the team with the given id.
func (s *Session) GetTeamMembershipByTeamId(ctx context.Context, teamId model.Id) (*model.TeamMembership, UserFacingError) {
	if err := s.RequireUser(); err != nil {
		return nil, err
	}

	membership, err := s.app.store.GetTeamMembershipByTeamAndUserId(ctx, teamId, s.user.Id)
	return membership, s.SanitizedError(err)
}

type TeamTeamMembership struct {
	User       *model.User
	Membership *model.TeamMembership
}

// Gets all memberships for the team with the given id.
func (s *Session) GetTeamMembershipsByTeamId(ctx context.Context, teamId model.Id) ([]*TeamTeamMembership, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	memberships, err := s.app.store.GetTeamMembershipsByTeamId(ctx, teamId)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	userIds := make([]model.Id, 0, len(memberships))
	for _, membership := range memberships {
		userIds = append(userIds, membership.UserId)
	}

	users, err := s.app.store.GetUsersByIds(ctx, userIds...)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	usersById := make(map[model.Id]*model.User, len(users))
	for _, user := range users {
		usersById[user.Id] = user
	}

	ret := make([]*TeamTeamMembership, 0, len(memberships))
	for _, membership := range memberships {
		if user, ok := usersById[membership.UserId]; ok {
			ret = append(ret, &TeamTeamMembership{
				User:       s.SanitizeUser(user),
				Membership: membership,
			})
		}
	}

	return ret, nil
}

func (s *Session) DeleteTeamMembershipByTeamAndUserId(ctx context.Context, teamId, userId model.Id) UserFacingError {
	if !s.HasUserId(userId) {
		if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
			return err
		}
	}

	return s.SanitizedError(s.app.store.DeleteTeamMembershipByTeamAndUserId(ctx, teamId, userId))
}

type TeamMembershipPatch struct {
	Role *model.TeamMembershipRole
}

func (s *Session) PatchTeamMembershipByTeamAndUserId(ctx context.Context, teamId, userId model.Id, patch TeamMembershipPatch) (*model.TeamMembership, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	storePatch := &store.TeamMembershipPatch{
		Role: patch.Role,
	}
	membership, err := s.app.store.PatchTeamMembershipByTeamAndUserId(ctx, teamId, userId, storePatch)
	return membership, s.SanitizedError(err)
}

type UserTeamInvite struct {
	Team   *model.Team
	Sender *model.User
	Invite *model.TeamInvite
}

// Gets team invites for the user with the given id.
func (s *Session) GetTeamInvitesByUserId(ctx context.Context, userId model.Id) ([]*UserTeamInvite, UserFacingError) {
	if !s.HasUserRole(model.UserRoleAdministrator) && !s.HasUserId(userId) {
		return nil, AuthorizationError{}
	}

	var emailAddress string
	if userId == s.user.Id {
		emailAddress = s.user.EmailAddress
	} else if user, err := s.app.store.GetUserById(ctx, userId, store.ConsistencyEventual); user == nil || err != nil {
		return nil, s.SanitizedError(err)
	} else {
		emailAddress = user.EmailAddress
	}

	invites, err := s.app.store.GetTeamInvitesByEmailAddress(ctx, emailAddress)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	teamIds := make([]model.Id, 0, len(invites))
	for _, invite := range invites {
		teamIds = append(teamIds, invite.TeamId)
	}

	teams, err := s.app.store.GetTeamsByIds(ctx, teamIds...)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	teamsById := make(map[model.Id]*model.Team, len(teams))
	for _, team := range teams {
		teamsById[team.Id] = team
	}

	userIds := make([]model.Id, 0, len(invites))
	for _, invite := range invites {
		userIds = append(userIds, invite.SenderId)
	}

	users, err := s.app.store.GetUsersByIds(ctx, userIds...)
	if err != nil {
		return nil, s.SanitizedError(err)
	}

	usersById := make(map[model.Id]*model.User, len(users))
	for _, user := range users {
		usersById[user.Id] = user
	}

	ret := make([]*UserTeamInvite, 0, len(invites))
	for _, invite := range invites {
		if team, ok := teamsById[invite.TeamId]; ok {
			if sender, ok := usersById[invite.SenderId]; ok {
				ret = append(ret, &UserTeamInvite{
					Team:   team,
					Sender: s.SanitizeUser(sender),
					Invite: invite,
				})
			}
		}
	}

	return ret, nil
}

func (s *Session) GetTeamInvitesByTeamId(ctx context.Context, teamId model.Id) ([]*model.TeamInvite, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	memberships, err := s.app.store.GetTeamInvitesByTeamId(ctx, teamId)
	return memberships, s.SanitizedError(err)
}

func ValidatePrincipalId(id string) UserFacingError {
	if len(id) == 0 {
		return NewUserError("A principal id is required.")
	} else if len(id) > 2048 {
		return NewUserError("Principal id is too long.")
	}
	return nil
}

func ValidateDescription(desc string) UserFacingError {
	if len(desc) > 4096 {
		return NewUserError("Descriptions are limited to 4096 characters.")
	}
	return nil
}

type TeamPrincipalSettingsPatch struct {
	Description *string
}

func (s *Session) CreateOrPatchTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx context.Context, teamId model.Id, principalKey string, patch TeamPrincipalSettingsPatch) (*model.TeamPrincipalSettings, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyEventual)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if !team.Entitlements.IndividualFeatures {
		return nil, NewUserError("An active subscription is required.")
	}

	storePatch := &store.TeamPrincipalSettingsPatch{
		Description: patch.Description,
	}

	if patch.Description != nil {
		if err := ValidateDescription(*patch.Description); err != nil {
			return nil, err
		}
	}

	settings, err := s.app.store.CreateOrPatchTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx, teamId, principalKey, storePatch)
	return settings, s.SanitizedError(err)
}

func (s *Session) GetTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx context.Context, teamId model.Id, principalKey string) (*model.TeamPrincipalSettings, UserFacingError) {
	if err := s.RequireTeamMember(ctx, teamId); err != nil {
		return nil, err
	}

	settings, err := s.app.store.GetTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx, teamId, principalKey)
	return settings, s.SanitizedError(err)
}
