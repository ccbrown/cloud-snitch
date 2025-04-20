import { useParams, usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useMemo, useState } from 'react';

import {
    AWSAccount,
    AWSIntegration,
    AWSRegion,
    AWSSCP,
    Report,
    TeamBillingProfile,
    TeamPaymentMethod,
    TeamPrincipalSettings,
    TeamSubscription,
    TeamTeamMembership,
    UserPasskey,
    UserTeamInvite,
    UserTeamMembership,
} from '@/generated/api';
import { CombinedReport } from '@/report';
import { useDispatch, useSelector } from '@/store';
import { formatDurationSeconds, parseDuration } from '@/time';

// Returns the current user, fetching it if it's not already loaded.
export const useCurrentUser = () => {
    const [didFetch, setDidFetch] = useState(false);
    const dispatch = useDispatch();

    const currentUser = useSelector((state) =>
        state.users.currentUserId ? state.users.users[state.users.currentUserId] : null,
    );
    const isLoading = useSelector((state) => state.loading.effects.users.fetchCurrent);

    useEffect(() => {
        if (!currentUser && !isLoading && !didFetch) {
            setDidFetch(true);
            dispatch.users.fetchCurrent();
        }
    }, [currentUser, isLoading, didFetch, dispatch]);

    return currentUser;
};

// Returns the current user's team memberships, fetching them once.
export const useCurrentUserTeamMemberships = (): UserTeamMembership[] | undefined => {
    const dispatch = useDispatch();
    const currentUserId = useCurrentUser()?.id;

    useEffect(() => {
        if (currentUserId) {
            dispatch.teams.fetchUserTeamMemberships(currentUserId);
        }
    }, [currentUserId, dispatch]);

    const memberships = useSelector((state) =>
        state.users.currentUserId ? state.teams.userTeamMemberships[state.users.currentUserId] : undefined,
    );
    const teams = useSelector((state) => state.teams.teams);

    if (memberships === undefined) {
        return undefined;
    }

    const ret: Array<UserTeamMembership> = [];
    Object.values(memberships).forEach((m) => {
        if (m.teamId in teams) {
            ret.push({
                team: teams[m.teamId],
                membership: m,
            });
        }
    });
    return ret;
};

// Returns the current team's team memberships, fetching them once.
export const useCurrentTeamTeamMemberships = (): TeamTeamMembership[] | undefined => {
    const dispatch = useDispatch();
    const currentTeamId = useCurrentTeamId();

    useEffect(() => {
        if (currentTeamId) {
            dispatch.teams.fetchMemberships(currentTeamId);
        }
    }, [currentTeamId, dispatch]);

    const memberships = useSelector((state) => currentTeamId && state.teams.memberships[currentTeamId]);
    const users = useSelector((state) => state.users.users);

    if (memberships === undefined) {
        return undefined;
    }

    const ret: Array<TeamTeamMembership> = [];
    Object.values(memberships).forEach((m) => {
        if (m.userId in users) {
            ret.push({
                user: users[m.userId],
                membership: m,
            });
        }
    });
    return ret;
};

// Returns the current user's team invites, fetching them once.
export const useCurrentUserTeamInvites = (): UserTeamInvite[] | undefined => {
    const dispatch = useDispatch();
    const currentUserId = useCurrentUser()?.id;

    useEffect(() => {
        if (currentUserId) {
            dispatch.teams.fetchUserTeamInvites(currentUserId);
        }
    }, [currentUserId, dispatch]);

    const invites = useSelector((state) =>
        state.users.currentUserId ? state.teams.userTeamInvites[state.users.currentUserId] : undefined,
    );
    const teams = useSelector((state) => state.teams.teams);
    const users = useSelector((state) => state.users.users);

    if (invites === undefined) {
        return undefined;
    }

    const ret: Array<UserTeamInvite> = [];
    Object.values(invites).forEach((m) => {
        if (m.teamId in teams && m.senderId in users) {
            ret.push({
                team: teams[m.teamId],
                sender: users[m.senderId],
                invite: m,
            });
        }
    });
    return ret;
};

// Returns the current user's passkeys, fetching them once.
export const useCurrentUserPasskeys = (): UserPasskey[] | undefined => {
    const dispatch = useDispatch();
    const currentUserId = useCurrentUser()?.id;

    useEffect(() => {
        if (currentUserId) {
            dispatch.users.fetchPasskeys(currentUserId);
        }
    }, [currentUserId, dispatch]);

    const passkeys = useSelector((state) => currentUserId && state.users.passkeys[currentUserId]);
    return passkeys ? Object.values(passkeys) : undefined;
};

export const useCurrentTeamId = () => {
    const { teamId } = useParams<{ teamId: string }>();

    useEffect(() => {
        window.localStorage.setItem('team', teamId);
    }, [teamId]);

    return teamId;
};

export const useMostRecentTeamId = () => {
    const [mostRecentTeamId, setMostRecentTeamId] = useState<string | null>(null);

    useEffect(() => {
        const storedTeamId = window.localStorage.getItem('team');
        if (storedTeamId) {
            setMostRecentTeamId(storedTeamId);
        }
    }, []);

    return mostRecentTeamId;
};

export const useCurrentTeam = () => {
    const dispatch = useDispatch();
    const teamId = useCurrentTeamId();

    useEffect(() => {
        if (teamId) {
            dispatch.teams.fetch(teamId);
        }
    }, [teamId, dispatch]);

    return useSelector((state) => state.teams.teams[teamId]);
};

export const useTeamAwsIntegrations = (teamId?: string): AWSIntegration[] | undefined => {
    const dispatch = useDispatch();

    useEffect(() => {
        if (teamId) {
            dispatch.aws.fetchIntegrationsByTeamId(teamId);
        }
    }, [teamId, dispatch]);

    const ids = useSelector((state) => teamId !== undefined && state.aws.teamIntegrationIds[teamId]);
    const integrations = useSelector((state) => state.aws.integrations);

    return useMemo(() => {
        if (!ids) {
            return undefined;
        }

        const ret: AWSIntegration[] = [];
        ids.forEach((id) => {
            if (id in integrations) {
                ret.push(integrations[id]);
            }
        });
        return ret;
    }, [ids, integrations]);
};

export const useTeamAwsAccountsMap = (teamId: string): Map<string, AWSAccount> | undefined => {
    const dispatch = useDispatch();

    const ids = useSelector((state) => state.aws.teamAccountIds[teamId]);

    useEffect(() => {
        if (!ids) {
            dispatch.aws.fetchAccountsByTeamId(teamId);
        }
    }, [ids, teamId, dispatch]);

    const accounts = useSelector((state) => state.aws.accounts);

    return useMemo(() => {
        if (!ids) {
            return undefined;
        }

        const ret: Map<string, AWSAccount> = new Map();
        ids.forEach((id) => {
            if (id in accounts) {
                ret.set(id, accounts[id]);
            }
        });
        return ret;
    }, [ids, accounts]);
};

export const useTeamReports = (teamId: string): Report[] | undefined => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.reports.fetchTeamReports(teamId);
    }, [teamId, dispatch]);

    const ids = useSelector((state) => state.reports.teamReportIds[teamId]);
    const reports = useSelector((state) => state.reports.reports);

    return useMemo(() => {
        if (!ids) {
            return undefined;
        }

        const ret: Report[] = [];
        ids.forEach((id) => {
            if (id in reports) {
                ret.push(reports[id]);
            }
        });
        return ret;
    }, [ids, reports]);
};

// Returns loaded reports in combined form, and a loading progress metric between 0 and 1.
export const useCombinedReport = (reports?: Report[]): [CombinedReport, number] => {
    const dispatch = useDispatch();
    const reportContents = useSelector((state) => state.reports.contents);
    const contentRequests = useSelector((state) => state.reports.content_requests);

    useEffect(() => {
        reports?.forEach((report) => {
            if (!contentRequests[report.id]) {
                dispatch.reports.fetchContent(report.id);
            }
        });
    }, [reports, dispatch, reportContents, contentRequests]);

    const loaded = useMemo(() => {
        const loaded = [];
        if (reports) {
            for (const report of reports) {
                const content = reportContents[report.id];
                if (content) {
                    loaded.push({
                        report: content,
                        apiReport: report,
                    });
                }
            }
        }
        return loaded;
    }, [reports, reportContents]);

    const combined = useMemo(() => new CombinedReport(loaded), [loaded]);

    const progress = reports ? (reports.length ? loaded.length / reports.length : 1) : 0;

    return [combined, progress];
};

export const useAwsRegions = (): AWSRegion[] => {
    const m = useAwsRegionsMap();
    return useMemo(() => Array.from(m.values()), [m]);
};

export const useAwsRegionsMap = (): Map<string, AWSRegion> => {
    const dispatch = useDispatch();
    const regions = useSelector((state) => state.aws.regions);
    const needsRegions = Object.keys(regions).length === 0;
    const [didFetch, setDidFetch] = useState(false);
    const isLoading = useSelector((state) => state.loading.effects.aws.fetchRegions);

    useEffect(() => {
        if (needsRegions && !isLoading && !didFetch) {
            setDidFetch(true);
            dispatch.aws.fetchRegions();
        }
    }, [needsRegions, dispatch, isLoading, didFetch]);

    return useMemo(() => new Map(Object.entries(regions)), [regions]);
};

export const useModifySearchParams = (): ((updates: Record<string, string | null | undefined>) => void) => {
    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();

    return useCallback(
        (updates: Record<string, string | null | undefined>) => {
            const params = new URLSearchParams(searchParams.toString());
            let didChange = false;
            for (const [key, value] of Object.entries(updates)) {
                if (!value) {
                    if (params.has(key)) {
                        didChange = true;
                    }
                    params.delete(key);
                } else {
                    if (params.get(key) !== value) {
                        didChange = true;
                    }
                    params.set(key, value);
                }
            }
            if (didChange) {
                router.push(`${pathname}?${params.toString()}`);
            }
        },
        [searchParams, pathname, router],
    );
};

export const useSearchParamState = (key: string): [string, (newValue?: string | null) => void] => {
    const searchParams = useSearchParams();
    const modifySearchParams = useModifySearchParams();

    const value = searchParams.get(key);

    const setter = useCallback(
        (newValue?: string | null) => modifySearchParams({ [key]: newValue }),
        [key, modifySearchParams],
    );

    return [value || '', setter];
};

interface Filter {
    mode: 'include' | 'exclude';
    values: Set<string>;
}

const DEFAULT_FILTER: Filter = {
    mode: 'exclude',
    values: new Set(),
};

export const useSearchParamFilterState = (key: string): [Filter, (newValue?: Filter | null) => void] => {
    const [stringValue, setStringValue] = useSearchParamState(key);

    const setFilter = useCallback(
        (newValue?: Filter | null) => {
            if (!newValue || (newValue.mode === 'exclude' && newValue.values.size === 0)) {
                setStringValue(null);
            } else {
                setStringValue(`${newValue.mode};${Array.from(newValue.values).join(',')}`);
            }
        },
        [setStringValue],
    );

    const parts = stringValue.split(';');
    if (parts.length !== 2 || (parts[0] !== 'include' && parts[0] !== 'exclude')) {
        return [DEFAULT_FILTER, setFilter];
    }

    const values = new Set(parts[1].split(',').filter((v) => v.length > 0));
    return [
        {
            mode: parts[0],
            values,
        },
        setFilter,
    ];
};

export const useSearchParamDurationState = (
    key: string,
): [number | undefined, (newDurationSeconds?: number | null) => void] => {
    const [stringValue, setStringValue] = useSearchParamState(key);

    const setDuration = useCallback(
        (newDurationSeconds?: number | null) => {
            setStringValue(newDurationSeconds ? formatDurationSeconds(newDurationSeconds) : null);
        },
        [setStringValue],
    );

    const durationSeconds = stringValue ? parseDuration(stringValue) : undefined;

    return [!durationSeconds || isNaN(durationSeconds) ? undefined : durationSeconds, setDuration];
};

export const useTeamBillingProfile = (teamId: string): TeamBillingProfile | undefined | null => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.teams.fetchBillingProfile(teamId);
    }, [teamId, dispatch]);

    return useSelector((state) => state.teams.billingProfiles[teamId]);
};

export const useTeamPaymentMethod = (teamId: string): TeamPaymentMethod | undefined | null => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.teams.fetchPaymentMethod(teamId);
    }, [teamId, dispatch]);

    return useSelector((state) => state.teams.paymentMethods[teamId]);
};

export const useTeamSubscription = (teamId: string): TeamSubscription | undefined | null => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.teams.fetchSubscription(teamId);
    }, [teamId, dispatch]);

    return useSelector((state) => state.teams.subscriptions[teamId]);
};

// Returns the current team's principal settings for the given key, fetching them once.
export const useCurrentTeamPrincipalSettings = (principalKey: string): TeamPrincipalSettings | undefined => {
    const dispatch = useDispatch();
    const currentTeamId = useCurrentTeamId();

    useEffect(() => {
        if (currentTeamId && principalKey) {
            dispatch.teams.fetchPrincipalSettings({
                teamId: currentTeamId,
                principalKey,
            });
        }
    }, [currentTeamId, principalKey, dispatch]);

    return useSelector((state) => {
        if (currentTeamId && principalKey && state.teams.principalSettings[currentTeamId]) {
            return state.teams.principalSettings[currentTeamId][principalKey];
        }
        return undefined;
    });
};

export const useManagedAwsScp = (teamId: string, accountId: string): AWSSCP | undefined | null => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.aws.fetchManagedScpByTeamAndAccountId({ teamId, accountId });
    }, [teamId, accountId, dispatch]);

    return useSelector((state) => state.aws.managedScps[accountId]);
};
