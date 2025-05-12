import { createModel } from '@rematch/core';

import { RootModel } from '.';
import { ApiError, apiConfiguration } from './api';

import {
    CreateTeamInput,
    CreateTeamBillingProfileInput,
    CreateTeamSubscriptionInput,
    CreateTeamInviteInput,
    PutTeamPaymentMethodInput,
    Team,
    TeamApi,
    TeamBillingProfile,
    TeamInvite,
    TeamMembership,
    TeamMembershipRole,
    TeamPaymentMethod,
    TeamPrincipalSettings,
    TeamSubscription,
    UpdateTeamInput,
    UpdateTeamBillingProfileInput,
    UpdateTeamMembershipInput,
    UpdateTeamSubscriptionInput,
    UpdateTeamPrincipalSettingsInput,
} from '@/generated/api';

interface TeamsState {
    teams: Record<string, Team>;
    billingProfiles: Record<string, TeamBillingProfile | null>;
    memberships: Record<string, Record<string, TeamMembership>>;
    paymentMethods: Record<string, TeamPaymentMethod | null>;
    principalSettings: Record<string, Record<string, TeamPrincipalSettings>>;
    subscriptions: Record<string, TeamSubscription | null>;
    userTeamInvites: Record<string, Record<string, TeamInvite>>;
    userTeamMemberships: Record<string, Record<string, TeamMembership>>;
}

export const teams = createModel<RootModel>()({
    state: {
        teams: {},
        billingProfiles: {},
        memberships: {},
        paymentMethods: {},
        principalSettings: {},
        subscriptions: {},
        userTeamInvites: {},
        userTeamMemberships: {},
    } as TeamsState,
    reducers: {
        put(state, team: Team) {
            state.teams[team.id] = team;
        },
        putBillingProfile(state, teamId: string, billingProfile: TeamBillingProfile | null) {
            state.billingProfiles[teamId] = billingProfile;
        },
        putPaymentMethod(state, teamId: string, paymentMethod: TeamPaymentMethod | null) {
            state.paymentMethods[teamId] = paymentMethod;
        },
        putSubscription(state, teamId: string, subscription: TeamSubscription | null) {
            state.subscriptions[teamId] = subscription;
        },
        putPrincipalSettings(
            state,
            payload: { teamId: string; principalKey: string; settings: TeamPrincipalSettings },
        ) {
            if (!state.principalSettings[payload.teamId]) {
                state.principalSettings[payload.teamId] = {};
            }
            state.principalSettings[payload.teamId][payload.principalKey] = payload.settings;
        },
        setUserTeamInvites(state, userId: string, invites: TeamInvite[]) {
            const map: Record<string, TeamInvite> = {};
            invites.forEach((inv) => {
                map[inv.teamId] = inv;
            });
            state.userTeamInvites[userId] = map;
        },
        removeInvite(state, userId: string, teamId: string) {
            if (state.userTeamInvites[userId]) {
                delete state.userTeamInvites[userId][teamId];
            }
        },
        removeInviteForEmailAddress(state, emailAddress: string, teamId: string) {
            for (const userId in state.userTeamInvites) {
                if (state.userTeamInvites[userId][teamId]?.emailAddress === emailAddress) {
                    delete state.userTeamInvites[userId][teamId];
                }
            }
        },
        putTeamMembership(state, membership: TeamMembership) {
            if (!state.userTeamMemberships[membership.userId]) {
                state.userTeamMemberships[membership.userId] = {};
            }
            state.userTeamMemberships[membership.userId][membership.teamId] = membership;
            if (!state.memberships[membership.teamId]) {
                state.memberships[membership.teamId] = {};
            }
            state.memberships[membership.teamId][membership.userId] = membership;
        },
        setUserTeamMemberships(state, userId: string, memberships: TeamMembership[]) {
            const map: Record<string, TeamMembership> = {};
            memberships.forEach((m) => {
                map[m.teamId] = m;
            });
            state.userTeamMemberships[userId] = map;
        },
        setMemberships(state, teamId: string, memberships: TeamMembership[]) {
            const map: Record<string, TeamMembership> = {};
            memberships.forEach((m) => {
                map[m.userId] = m;
            });
            state.memberships[teamId] = map;
        },
        putMembership(state, membership: TeamMembership) {
            if (!state.memberships[membership.teamId]) {
                state.memberships[membership.teamId] = {};
            }
            state.memberships[membership.teamId][membership.userId] = membership;
            if (!state.userTeamMemberships[membership.userId]) {
                state.userTeamMemberships[membership.userId] = {};
            }
            state.userTeamMemberships[membership.userId][membership.teamId] = membership;
        },
        removeMembership(state, teamId: string, userId: string) {
            if (state.memberships[teamId]) {
                delete state.memberships[teamId][userId];
            }
            if (state.userTeamMemberships[userId]) {
                delete state.userTeamMemberships[userId][teamId];
            }
        },
    },
    effects: (dispatch) => ({
        async fetch(id: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeam({
                teamId: id,
            });
            dispatch.teams.put(resp);
        },
        async fetchAll(_payload: void, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeams();
            resp.forEach((team) => {
                dispatch.teams.put(team);
            });
        },
        async fetchBillingProfile(id: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            try {
                const resp = await api.getTeamBillingProfile({
                    teamId: id,
                });
                dispatch.teams.putBillingProfile(id, resp);
            } catch (err) {
                if (err instanceof ApiError && err.status === 404) {
                    dispatch.teams.putBillingProfile(id, null);
                } else {
                    throw err;
                }
            }
        },
        async fetchPaymentMethod(id: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            try {
                const resp = await api.getTeamPaymentMethod({
                    teamId: id,
                });
                dispatch.teams.putPaymentMethod(id, resp);
            } catch (err) {
                if (err instanceof ApiError && err.status === 404) {
                    dispatch.teams.putPaymentMethod(id, null);
                } else {
                    throw err;
                }
            }
        },
        async fetchSubscription(id: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            try {
                const resp = await api.getTeamSubscription({
                    teamId: id,
                });
                dispatch.teams.putSubscription(id, resp);
            } catch (err) {
                if (err instanceof ApiError && err.status === 404) {
                    dispatch.teams.putSubscription(id, null);
                } else {
                    throw err;
                }
            }
        },
        async fetchUserTeamInvites(userId: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeamInvitesByUserId({
                userId,
            });
            resp.forEach((m) => {
                dispatch.teams.put(m.team);
                dispatch.users.put(m.sender);
            });
            dispatch.teams.setUserTeamInvites(
                userId,
                resp.map((inv) => inv.invite),
            );
        },
        async fetchUserTeamMemberships(userId: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeamMembershipsByUserId({
                userId,
            });
            resp.forEach((m) => {
                dispatch.teams.put(m.team);
            });
            dispatch.teams.setUserTeamMemberships(
                userId,
                resp.map((m) => m.membership),
            );
        },
        async fetchMemberships(teamId: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeamMembershipsByTeamId({
                teamId,
            });
            resp.forEach((m) => {
                dispatch.users.put(m.user);
            });
            dispatch.teams.setMemberships(
                teamId,
                resp.map((m) => m.membership),
            );
        },
        async fetchPrincipalSettings(payload: { teamId: string; principalKey: string }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getTeamPrincipalSettings({
                teamId: payload.teamId,
                principalKey: payload.principalKey,
            });
            dispatch.teams.putPrincipalSettings({
                teamId: payload.teamId,
                principalKey: payload.principalKey,
                settings: resp,
            });
        },
        async updatePrincipalSettings(
            payload: { teamId: string; principalKey: string; input: UpdateTeamPrincipalSettingsInput },
            state,
        ) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.updateTeamPrincipalSettings({
                teamId: payload.teamId,
                principalKey: payload.principalKey,
                updateTeamPrincipalSettingsInput: payload.input,
            });
            dispatch.teams.putPrincipalSettings({
                teamId: payload.teamId,
                principalKey: payload.principalKey,
                settings: resp,
            });
        },
        async create(input: CreateTeamInput, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.createTeam({
                createTeamInput: input,
            });
            dispatch.teams.put(resp);
            if (state.users.currentUserId) {
                dispatch.teams.putTeamMembership({
                    userId: state.users.currentUserId,
                    teamId: resp.id,
                    role: TeamMembershipRole.Administrator,
                });
            }
            return resp;
        },
        async update(payload: { teamId: string; input: UpdateTeamInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.updateTeam({
                teamId: payload.teamId,
                updateTeamInput: payload.input,
            });
            dispatch.teams.put(resp);
        },
        async createBillingProfile(payload: { teamId: string; input: CreateTeamBillingProfileInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.createTeamBillingProfile({
                teamId: payload.teamId,
                createTeamBillingProfileInput: payload.input,
            });
            dispatch.teams.putBillingProfile(payload.teamId, resp);
            return resp;
        },
        async updateBillingProfile(payload: { teamId: string; input: UpdateTeamBillingProfileInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.updateTeamBillingProfile({
                teamId: payload.teamId,
                updateTeamBillingProfileInput: payload.input,
            });
            dispatch.teams.putBillingProfile(payload.teamId, resp);
            return resp;
        },
        async updatePaymentMethod(payload: { teamId: string; input: PutTeamPaymentMethodInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.putTeamPaymentMethod({
                teamId: payload.teamId,
                putTeamPaymentMethodInput: payload.input,
            });
            dispatch.teams.putPaymentMethod(payload.teamId, resp);
            return resp;
        },
        async join(teamId: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.joinTeam({
                teamId,
                body: {},
            });
            dispatch.teams.putTeamMembership(resp);
            dispatch.teams.removeInvite(resp.userId, teamId);
        },
        async createInvite(payload: { teamId: string; input: CreateTeamInviteInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            await api.createTeamInvite({
                teamId: payload.teamId,
                createTeamInviteInput: payload.input,
            });
        },
        async deleteInvite(payload: { teamId: string; emailAddress: string }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            await api.deleteTeamInvite({
                teamId: payload.teamId,
                emailAddress: payload.emailAddress,
                body: {},
            });
            dispatch.teams.removeInviteForEmailAddress(payload.emailAddress, payload.teamId);
        },
        async createSubscription(payload: { teamId: string; input: CreateTeamSubscriptionInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.createTeamSubscription({
                teamId: payload.teamId,
                createTeamSubscriptionInput: payload.input,
            });
            dispatch.teams.putSubscription(payload.teamId, resp);
            return resp;
        },
        async updateSubscription(payload: { teamId: string; input: UpdateTeamSubscriptionInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.updateTeamSubscription({
                teamId: payload.teamId,
                updateTeamSubscriptionInput: payload.input,
            });
            dispatch.teams.putSubscription(payload.teamId, resp);
            return resp;
        },
        async updateMembership(payload: { teamId: string; userId: string; input: UpdateTeamMembershipInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.updateTeamMembership({
                teamId: payload.teamId,
                userId: payload.userId,
                updateTeamMembershipInput: payload.input,
            });
            dispatch.teams.putTeamMembership(resp);
        },
        async deleteMembership(payload: { teamId: string; userId: string }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            await api.deleteTeamMembership({
                teamId: payload.teamId,
                userId: payload.userId,
                body: {},
            });
            dispatch.teams.removeMembership(payload.teamId, payload.userId);
        },
    }),
});
