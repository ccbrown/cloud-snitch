import { createModel } from '@rematch/core';

import { RootModel } from '.';
import { apiConfiguration, ApiError } from './api';

import {
    AwsApi,
    AWSAccount,
    AWSIntegration,
    AWSRegion,
    AWSSCP,
    CreateAWSIntegrationInput,
    PutAWSSCPInput,
    UpdateAWSIntegrationInput,
} from '@/generated/api';

interface AwsState {
    integrations: Record<string, AWSIntegration>;
    teamIntegrationIds: Record<string, string[]>;
    regions: Record<string, AWSRegion>;
    accounts: Record<string, AWSAccount>;
    teamAccountIds: Record<string, string[]>;
    managedScps: Record<string, AWSSCP | null>;
}

export const aws = createModel<RootModel>()({
    state: {
        integrations: {},
        teamIntegrationIds: {},
        regions: {},
        accounts: {},
        teamAccountIds: {},
        managedScps: {},
    } as AwsState,
    reducers: {
        putIntegration(state, integration: AWSIntegration) {
            state.integrations[integration.id] = integration;
        },
        putRegion(state, region: AWSRegion) {
            state.regions[region.id] = region;
        },
        putTeamIntegrationId(state, teamId: string, integrationId: string) {
            if (!state.teamIntegrationIds[teamId]) {
                state.teamIntegrationIds[teamId] = [];
            }
            state.teamIntegrationIds[teamId].push(integrationId);
        },
        putAccount(state, account: AWSAccount) {
            state.accounts[account.id] = account;
        },
        putManagedScp(state, accountId: string, scp: AWSSCP | null) {
            state.managedScps[accountId] = scp;
        },
        setTeamAccountIds(state, teamId: string, accountIds: string[]) {
            state.teamAccountIds[teamId] = accountIds;
        },
        setTeamIntegrationIds(state, teamId: string, integrationIds: string[]) {
            state.teamIntegrationIds[teamId] = integrationIds;
        },
        removeIntegration(state, integrationId: string) {
            delete state.integrations[integrationId];
            Object.values(state.teamIntegrationIds).forEach((ids) => {
                const index = ids.indexOf(integrationId);
                if (index !== -1) {
                    ids.splice(index, 1);
                }
            });
        },
    },
    effects: (dispatch) => ({
        async fetchRegions(_payload: void, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.getAWSRegions();
            resp.forEach((region) => {
                dispatch.aws.putRegion(region);
            });
        },
        async fetchIntegrationsByTeamId(teamId: string, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.getAWSIntegrationsByTeamId({
                teamId,
            });
            resp.forEach((integration) => {
                dispatch.aws.putIntegration(integration);
            });
            dispatch.aws.setTeamIntegrationIds(
                teamId,
                resp.map((i) => i.id),
            );
        },
        async fetchAccountsByTeamId(teamId: string, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.getAWSAccountsByTeamId({
                teamId,
            });
            resp.forEach((account) => {
                dispatch.aws.putAccount(account);
            });
            dispatch.aws.setTeamAccountIds(
                teamId,
                resp.map((a) => a.id),
            );
        },
        async fetchManagedScpByTeamAndAccountId(payload: { teamId: string; accountId: string }, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            try {
                const resp = await api.getManagedAWSSCP(payload);
                dispatch.aws.putManagedScp(payload.accountId, resp);
            } catch (err) {
                if (err instanceof ApiError && err.status === 404) {
                    dispatch.aws.putManagedScp(payload.accountId, null);
                } else {
                    throw err;
                }
            }
        },
        async putManagedScpByTeamAndAccountId(
            payload: { teamId: string; accountId: string; input: PutAWSSCPInput },
            state,
        ) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.putManagedAWSSCP({
                teamId: payload.teamId,
                accountId: payload.accountId,
                putAWSSCPInput: payload.input,
            });
            dispatch.aws.putManagedScp(payload.accountId, resp);
        },
        async createIntegration(payload: { teamId: string; input: CreateAWSIntegrationInput }, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.createAWSIntegration({
                teamId: payload.teamId,
                createAWSIntegrationInput: payload.input,
            });
            dispatch.aws.putIntegration(resp);
            dispatch.aws.putTeamIntegrationId(payload.teamId, resp.id);
        },
        async updateIntegration(payload: { integrationId: string; input: UpdateAWSIntegrationInput }, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            const resp = await api.updateAWSIntegration({
                integrationId: payload.integrationId,
                updateAWSIntegrationInput: payload.input,
            });
            dispatch.aws.putIntegration(resp);
        },
        async deleteIntegration(payload: { id: string; deleteAssociatedData: boolean }, state) {
            const api = new AwsApi(apiConfiguration(state.api));
            await api.deleteAWSIntegration({
                integrationId: payload.id,
                deleteAWSIntegrationInput: {
                    deleteAssociatedData: payload.deleteAssociatedData,
                },
            });
            dispatch.aws.removeIntegration(payload.id);
            if (payload.deleteAssociatedData) {
                dispatch.reports.removeReportsByAwsIntegrationId(payload.id);
            }
        },
    }),
});
