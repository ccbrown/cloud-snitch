package api

import (
	"context"
	"fmt"
	"time"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TeamFromModel(team *model.Team) apispec.Team {
	return apispec.Team{
		Id:   team.Id.String(),
		Name: team.Name,
		Entitlements: apispec.TeamEntitlements{
			IndividualFeatures: team.Entitlements.IndividualFeatures,
			TeamFeatures:       team.Entitlements.TeamFeatures,
		},
	}
}

func UserTeamMembershipFromApp(membership *app.UserTeamMembership) apispec.UserTeamMembership {
	return apispec.UserTeamMembership{
		Team:       TeamFromModel(membership.Team),
		Membership: TeamMembershipFromModel(membership.Membership),
	}
}

func TeamTeamMembershipFromApp(membership *app.TeamTeamMembership) apispec.TeamTeamMembership {
	return apispec.TeamTeamMembership{
		User:       UserFromModel(membership.User),
		Membership: TeamMembershipFromModel(membership.Membership),
	}
}

func UserTeamInviteFromApp(invite *app.UserTeamInvite) apispec.UserTeamInvite {
	return apispec.UserTeamInvite{
		Team:   TeamFromModel(invite.Team),
		Sender: UserFromModel(invite.Sender),
		Invite: TeamInviteFromModel(invite.Invite),
	}
}

func TeamMembershipRoleFromSpec(role apispec.TeamMembershipRole) model.TeamMembershipRole {
	switch role {
	case apispec.TeamMembershipRoleADMINISTRATOR:
		return model.TeamMembershipRoleAdministrator
	case apispec.TeamMembershipRoleMEMBER:
		return model.TeamMembershipRoleMember
	default:
		panic(fmt.Sprintf("unknown team membership role: %v", role))
	}
}

func TeamMembershipRoleFromModel(role model.TeamMembershipRole) apispec.TeamMembershipRole {
	switch role {
	case model.TeamMembershipRoleAdministrator:
		return apispec.TeamMembershipRoleADMINISTRATOR
	case model.TeamMembershipRoleMember:
		return apispec.TeamMembershipRoleMEMBER
	default:
		panic(fmt.Sprintf("unexpected team membership role: %v", string(role)))
	}
}

func TeamMembershipFromModel(membership *model.TeamMembership) apispec.TeamMembership {
	return apispec.TeamMembership{
		TeamId: membership.TeamId.String(),
		UserId: membership.UserId.String(),
		Role:   TeamMembershipRoleFromModel(membership.Role),
	}
}

func TeamInviteFromModel(invite *model.TeamInvite) apispec.TeamInvite {
	return apispec.TeamInvite{
		TeamId:       invite.TeamId.String(),
		SenderId:     invite.SenderId.String(),
		EmailAddress: invite.EmailAddress,
		Role:         TeamMembershipRoleFromModel(invite.Role),
	}
}

func (api *API) GetTeams(ctx context.Context, request apispec.GetTeamsRequestObject) (apispec.GetTeamsResponseObject, error) {
	sess := ctxSession(ctx)

	if teams, err := sess.GetTeams(ctx); err != nil {
		return nil, err
	} else {
		return apispec.GetTeams200JSONResponse(mapSlice(teams, TeamFromModel)), nil
	}
}

func (api *API) CreateTeam(ctx context.Context, request apispec.CreateTeamRequestObject) (apispec.CreateTeamResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.CreateTeamInput{
		Name: request.Body.Name,
	}

	if team, err := sess.CreateTeam(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.CreateTeam200JSONResponse(TeamFromModel(team)), nil
	}
}

func (api *API) GetTeam(ctx context.Context, request apispec.GetTeamRequestObject) (apispec.GetTeamResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if team, err := sess.GetTeamById(ctx, teamId); err != nil {
		return nil, err
	} else if team == nil {
		return nil, app.NotFoundError("No such team.")
	} else {
		return apispec.GetTeam200JSONResponse(TeamFromModel(team)), nil
	}
}

func (api *API) UpdateTeam(ctx context.Context, request apispec.UpdateTeamRequestObject) (apispec.UpdateTeamResponseObject, error) {
	sess := ctxSession(ctx)

	patch := app.TeamPatch{
		Name: request.Body.Name,
	}

	if team, err := sess.PatchTeamById(ctx, model.Id(request.TeamId), patch); err != nil {
		return nil, err
	} else if team == nil {
		return nil, app.NotFoundError("No such team.")
	} else {
		return apispec.UpdateTeam200JSONResponse(TeamFromModel(team)), nil
	}
}

func (api *API) GetTeamInvitesByTeamId(ctx context.Context, request apispec.GetTeamInvitesByTeamIdRequestObject) (apispec.GetTeamInvitesByTeamIdResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if invites, err := sess.GetTeamInvitesByTeamId(ctx, teamId); err != nil {
		return nil, err
	} else {
		return apispec.GetTeamInvitesByTeamId200JSONResponse(mapSlice(invites, TeamInviteFromModel)), nil
	}
}

func (api *API) CreateTeamInvite(ctx context.Context, request apispec.CreateTeamInviteRequestObject) (apispec.CreateTeamInviteResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	input := app.InviteToTeamInput{
		TeamId:       teamId,
		EmailAddress: request.Body.EmailAddress,
		Role:         TeamMembershipRoleFromSpec(request.Body.Role),
	}

	if err := sess.InviteToTeam(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.CreateTeamInvite200JSONResponse{}, nil
	}
}

func (api *API) DeleteTeamInvite(ctx context.Context, request apispec.DeleteTeamInviteRequestObject) (apispec.DeleteTeamInviteResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if err := sess.DeleteTeamInviteByTeamIdAndEmailAddress(ctx, teamId, request.EmailAddress); err != nil {
		return nil, err
	} else {
		return apispec.DeleteTeamInvite200JSONResponse{}, nil
	}
}

func (api *API) JoinTeam(ctx context.Context, request apispec.JoinTeamRequestObject) (apispec.JoinTeamResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if membership, err := sess.JoinTeam(ctx, teamId); err != nil {
		return nil, err
	} else {
		return apispec.JoinTeam200JSONResponse(TeamMembershipFromModel(membership)), nil
	}
}

func (api *API) GetTeamMembershipsByTeamId(ctx context.Context, request apispec.GetTeamMembershipsByTeamIdRequestObject) (apispec.GetTeamMembershipsByTeamIdResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if memberships, err := sess.GetTeamMembershipsByTeamId(ctx, teamId); err != nil {
		return nil, err
	} else {
		return apispec.GetTeamMembershipsByTeamId200JSONResponse(mapSlice(memberships, TeamTeamMembershipFromApp)), nil
	}
}

func (api *API) DeleteTeamMembership(ctx context.Context, request apispec.DeleteTeamMembershipRequestObject) (apispec.DeleteTeamMembershipResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)
	userId := UserIdFromRequest(sess, request.UserId)

	if err := sess.DeleteTeamMembershipByTeamAndUserId(ctx, teamId, userId); err != nil {
		return nil, err
	} else {
		return apispec.DeleteTeamMembership200JSONResponse{}, nil
	}
}

func (api *API) UpdateTeamMembership(ctx context.Context, request apispec.UpdateTeamMembershipRequestObject) (apispec.UpdateTeamMembershipResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)
	userId := UserIdFromRequest(sess, request.UserId)

	patch := app.TeamMembershipPatch{}
	if request.Body.Role != nil {
		patch.Role = pointer(TeamMembershipRoleFromSpec(*request.Body.Role))
	}

	if membership, err := sess.PatchTeamMembershipByTeamAndUserId(ctx, teamId, userId, patch); err != nil {
		return nil, err
	} else if membership == nil {
		return nil, app.NotFoundError("No such membership.")
	} else {
		return apispec.UpdateTeamMembership200JSONResponse(TeamMembershipFromModel(membership)), nil
	}
}

func (api *API) GetTeamInvitesByUserId(ctx context.Context, request apispec.GetTeamInvitesByUserIdRequestObject) (apispec.GetTeamInvitesByUserIdResponseObject, error) {
	sess := ctxSession(ctx)
	userId := UserIdFromRequest(sess, request.UserId)

	if invites, err := sess.GetTeamInvitesByUserId(ctx, userId); err != nil {
		return nil, err
	} else {
		return apispec.GetTeamInvitesByUserId200JSONResponse(mapSlice(invites, UserTeamInviteFromApp)), nil
	}
}

func (api *API) GetTeamMembershipsByUserId(ctx context.Context, request apispec.GetTeamMembershipsByUserIdRequestObject) (apispec.GetTeamMembershipsByUserIdResponseObject, error) {
	sess := ctxSession(ctx)
	userId := UserIdFromRequest(sess, request.UserId)

	if memberships, err := sess.GetTeamMembershipsByUserId(ctx, userId); err != nil {
		return nil, err
	} else {
		return apispec.GetTeamMembershipsByUserId200JSONResponse(mapSlice(memberships, UserTeamMembershipFromApp)), nil
	}
}

func (api *API) QueueTeamReportGeneration(ctx context.Context, request apispec.QueueTeamReportGenerationRequestObject) (apispec.QueueTeamReportGenerationResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if err := sess.QueueTeamReportGeneration(ctx, app.QueueTeamReportGenerationInput{
		TeamId:    teamId,
		StartTime: request.Body.StartTime,
		Duration:  time.Second * time.Duration(request.Body.DurationSeconds),
		Retention: ReportRetentionFromSpec(request.Body.Retention),
	}); err != nil {
		return nil, err
	} else {
		return apispec.QueueTeamReportGeneration200Response{}, nil
	}
}

func TeamPrincipalSettingsFromModel(settings *model.TeamPrincipalSettings) apispec.TeamPrincipalSettings {
	return apispec.TeamPrincipalSettings{
		Description: nilIfEmpty(settings.Description),
	}
}

func (api *API) GetTeamPrincipalSettings(ctx context.Context, request apispec.GetTeamPrincipalSettingsRequestObject) (apispec.GetTeamPrincipalSettingsResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if settings, err := sess.GetTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx, teamId, request.Params.PrincipalKey); err != nil {
		return nil, err
	} else if settings == nil {
		return apispec.GetTeamPrincipalSettings200JSONResponse(apispec.TeamPrincipalSettings{}), nil
	} else {
		return apispec.GetTeamPrincipalSettings200JSONResponse(TeamPrincipalSettingsFromModel(settings)), nil
	}
}

func (api *API) UpdateTeamPrincipalSettings(ctx context.Context, request apispec.UpdateTeamPrincipalSettingsRequestObject) (apispec.UpdateTeamPrincipalSettingsResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if settings, err := sess.CreateOrPatchTeamPrincipalSettingsByTeamIdAndPrincipalKey(ctx, teamId, request.Params.PrincipalKey, app.TeamPrincipalSettingsPatch{
		Description: request.Body.Description,
	}); err != nil {
		return nil, err
	} else {
		return apispec.UpdateTeamPrincipalSettings200JSONResponse(TeamPrincipalSettingsFromModel(settings)), nil
	}
}
