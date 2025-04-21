import { createModel } from '@rematch/core';

import { RootModel } from '.';
import { apiConfiguration } from './api';

import { QueueTeamReportGenerationInput, Report, ReportApi, TeamApi } from '@/generated/api';
import { Report as ReportContent } from '@/report';

interface ContentRequest {
    time: Date;
    pending: boolean;
}

interface ReportsState {
    reports: Record<string, Report>;
    teamReportIds: Record<string, string[]>;
    contents: Record<string, ReportContent>;
    content_requests: Record<string, ContentRequest>;
}

export const reports = createModel<RootModel>()({
    state: {
        reports: {},
        teamReportIds: {},
        contents: {},
        content_requests: {},
    } as ReportsState,
    reducers: {
        put(state, report: Report) {
            state.reports[report.id] = report;
        },
        setTeamReportIds(state, teamId: string, reports: string[]) {
            state.teamReportIds[teamId] = reports;
        },
        removeReport(state, reportId: string) {
            delete state.reports[reportId];
            Object.values(state.teamReportIds).forEach((ids) => {
                const idx = ids.indexOf(reportId);
                if (idx !== -1) {
                    ids.splice(idx, 1);
                }
            });
        },
        putContent(state, reportId: string, content: ReportContent) {
            state.contents[reportId] = content;
        },
        setContentRequestPending(state, reportId: string) {
            state.content_requests[reportId] = {
                time: new Date(),
                pending: true,
            };
        },
        setContentRequestNotPending(state, reportId: string) {
            if (state.content_requests[reportId]) {
                state.content_requests[reportId].pending = false;
            }
        },
        removeReportsByAwsIntegrationId(state, integrationId: string) {
            Object.entries(state.reports).forEach(([reportId, report]) => {
                if (report.awsIntegrationId === integrationId) {
                    delete state.reports[reportId];
                    const teamReportIds = state.teamReportIds[report.teamId];
                    const idx = teamReportIds.indexOf(reportId);
                    if (idx !== -1) {
                        teamReportIds.splice(idx, 1);
                    }
                }
            });
        },
    },
    effects: (dispatch) => ({
        async fetchTeamReports(teamId: string, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            const resp = await api.getReportsByTeamId({
                teamId,
            });
            resp.forEach((r) => {
                dispatch.reports.put(r);
            });
            dispatch.reports.setTeamReportIds(
                teamId,
                resp.map((r) => r.id),
            );
            return resp;
        },
        async queueTeamReportGeneration(input: { teamId: string; input: QueueTeamReportGenerationInput }, state) {
            const api = new TeamApi(apiConfiguration(state.api));
            await api.queueTeamReportGeneration({
                teamId: input.teamId,
                queueTeamReportGenerationInput: input.input,
            });
        },
        async deleteReport(id: string, state) {
            const api = new ReportApi(apiConfiguration(state.api));
            await api.deleteReportById({
                reportId: id,
            });
            dispatch.reports.removeReport(id);
        },
        async fetchContent(id: string, state) {
            const report = state.reports.reports[id];
            if (!report) {
                return;
            }
            dispatch.reports.setContentRequestPending(id);
            try {
                const resp = await fetch(report.downloadUrl);
                if (!resp.ok) {
                    throw new Error(`Failed to fetch report content: ${resp.statusText}`);
                }
                const content = await resp.json();
                dispatch.reports.putContent(id, content);
            } finally {
                dispatch.reports.setContentRequestNotPending(id);
            }
        },
    }),
});
