'use client';

import { ShieldCheckIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Map as WorldMap, MapRef } from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';

import { Dialog, DurationDropdown, FilterDropdown } from '@/components';
import {
    useAwsRegionsMap,
    useCombinedReport,
    useCurrentTeam,
    useCurrentTeamId,
    useSearchParamDurationState,
    useSearchParamFilterState,
    useTeamReports,
    useTeamAwsAccountsMap,
    useTeamAwsIntegrations,
} from '@/hooks';
import { Report, TeamMembershipRole } from '@/generated/api';
import { useSelector } from '@/store';
import { Header } from './Header';
import { ContextPanel } from './ContextPanel';
import { maxZoomForClusterRect, minZoomForClusterRect, MapOverlays } from './MapOverlays';
import { NavigationPanel } from './NavigationPanel';
import { Rules } from './Rules';
import { isEqualSelection, useSelection, Selection } from './selection';

interface IdAndName {
    id: string;
    name?: string;
}

const Page = () => {
    const teamId = useCurrentTeamId();
    const team = useCurrentTeam();
    const teamReports = useTeamReports(teamId);

    const isTeamAdmin = useSelector((state) => {
        const memberships = state.users.currentUserId
            ? state.teams.userTeamMemberships[state.users.currentUserId]
            : undefined;
        return memberships && memberships[teamId]?.role === TeamMembershipRole.Administrator;
    });

    const [isRulesOpen, setIsRulesOpen] = useState(false);

    const [map, setMap] = useState<MapRef | null>(null);
    const [isMapReady, setIsMapReady] = useState(false);

    const teamAwsAccountsMap = useTeamAwsAccountsMap(teamId);

    const [awsAccountFilter, setAwsAccountFilter] = useSearchParamFilterState('awsAccount');
    const [awsRegionFilter, setAwsRegionFilter] = useSearchParamFilterState('awsRegion');

    const awsRegionsMap = useAwsRegionsMap();

    const [minUnfilteredReportTime, maxUnfilteredReportEndTime] = useMemo(() => {
        let minReportTime: Date | undefined = undefined;
        let maxReportEndTime: Date | undefined = undefined;
        teamReports?.forEach((report) => {
            const reportTime = report.scope.startTime;
            const reportEndTime = new Date(report.scope.startTime);
            reportEndTime.setSeconds(reportEndTime.getSeconds() + report.scope.durationSeconds);
            if (!minReportTime || reportTime < minReportTime) {
                minReportTime = reportTime;
            }
            if (!maxReportEndTime || reportEndTime > maxReportEndTime) {
                maxReportEndTime = reportEndTime;
            }
        });
        return [minReportTime, maxReportEndTime];
    }, [teamReports]);

    const [durationSeconds, setDurationSeconds] = useSearchParamDurationState('d');
    const [startTime, endTime] = useMemo(() => {
        if (!durationSeconds || !maxUnfilteredReportEndTime) {
            return [undefined, undefined];
        }
        const endTime = new Date(maxUnfilteredReportEndTime);
        const startTime = new Date(endTime.getTime() - durationSeconds * 1000);
        return [startTime, endTime];
    }, [maxUnfilteredReportEndTime, durationSeconds]);

    const reportFilter = useCallback(
        (report: Report) => {
            if (startTime) {
                if (report.scope.startTime < startTime) {
                    return false;
                }
            }
            if (endTime) {
                const reportEndTime = new Date(report.scope.startTime);
                reportEndTime.setSeconds(reportEndTime.getSeconds() + report.scope.durationSeconds);
                if (reportEndTime > endTime) {
                    return false;
                }
            }
            if (awsAccountFilter.mode === 'include' && !awsAccountFilter.values.has(report.scope.aws.accountId)) {
                return false;
            }
            if (awsAccountFilter.mode === 'exclude' && awsAccountFilter.values.has(report.scope.aws.accountId)) {
                return false;
            }
            if (awsRegionFilter.mode === 'include' && !awsRegionFilter.values.has(report.scope.aws.region)) {
                return false;
            }
            if (awsRegionFilter.mode === 'exclude' && awsRegionFilter.values.has(report.scope.aws.region)) {
                return false;
            }
            return true;
        },
        [awsAccountFilter, awsRegionFilter, startTime, endTime],
    );
    const filteredReports = useMemo(() => {
        if (!teamReports) {
            return undefined;
        }
        const filtered = teamReports.filter(reportFilter);

        const ret: Report[] = [];

        // Do some deduplication for robustness against duplicate backend processing, integration setups, etc.
        const seen = new Set<string>();
        filtered.forEach((report) => {
            const key = `${report.scope.aws.accountId}-${report.scope.aws.region}-${report.scope.startTime.getTime()}`;
            if (!seen.has(key)) {
                seen.add(key);
                ret.push(report);
            }
        });

        return ret;
    }, [teamReports, reportFilter]);

    const [allPrincipalsCombinedReport, loadingProgress] = useCombinedReport(filteredReports);

    const [selection, setSelection] = useSelection();

    const [zoom, setZoom] = useState(0);
    const [lastCenteredSelection, setLastCenteredSelection] = useState<Selection | null>(null);

    const [principalFilter, setPrincipalFilter] = useState<string>('');
    const [hoveredPrincipalId, setHoveredPrincipalId] = useState<string | null>(null);

    const sortedAvailableAwsAccounts: IdAndName[] = useMemo(() => {
        const accounts: Map<string, IdAndName> = new Map();
        teamReports?.forEach((report) => {
            accounts.set(report.scope.aws.accountId, {
                id: report.scope.aws.accountId,
                name: teamAwsAccountsMap?.get(report.scope.aws.accountId)?.name,
            });
        });
        return Array.from(accounts.values()).sort((a, b) => {
            const aLabel = a.name?.toLowerCase() || a.id;
            const bLabel = b.name?.toLowerCase() || b.id;
            return aLabel.localeCompare(bLabel);
        });
    }, [teamAwsAccountsMap, teamReports]);

    const sortedAvailableAwsRegions: IdAndName[] = useMemo(() => {
        const regions: Map<string, IdAndName> = new Map();
        teamReports?.forEach((report) => {
            regions.set(report.scope.aws.region, {
                id: report.scope.aws.region,
                name: awsRegionsMap.get(report.scope.aws.region)?.name,
            });
        });
        return Array.from(regions.values()).sort((a, b) => {
            const aLabel = a.name?.toLowerCase() || a.id;
            const bLabel = b.name?.toLowerCase() || b.id;
            return aLabel.localeCompare(bLabel);
        });
    }, [awsRegionsMap, teamReports]);

    const combinedReport = useMemo(() => {
        if (!principalFilter) {
            return allPrincipalsCombinedReport;
        }
        return allPrincipalsCombinedReport.withFilteredPrincipals(principalFilter);
    }, [allPrincipalsCombinedReport, principalFilter]);

    // This effect is used to validate the selection string and clear if it is invalid or center the
    // map on it otherwise.
    useEffect(() => {
        if (!selection && lastCenteredSelection) {
            setLastCenteredSelection(null);
        }

        if (!selection || isEqualSelection(lastCenteredSelection, selection) || !map || !isMapReady) {
            return;
        }

        if (!selection) {
            setSelection(null);
            return;
        }

        switch (selection.type) {
            case 'aws-region':
                const region = awsRegionsMap.get(selection.id);
                if (region) {
                    map.flyTo({
                        center: [region.longitude, region.latitude],
                        zoom: Math.max(zoom, 5),
                    });
                    setLastCenteredSelection(selection);
                } else if (awsRegionsMap.size > 0) {
                    setSelection(null);
                }
                break;
            case 'network':
                const loc = combinedReport.networkLocations[selection.cidr];
                if (loc) {
                    map.flyTo({
                        center: [loc.longitude, loc.latitude],
                        zoom: Math.max(zoom, 5),
                    });
                    setLastCenteredSelection(selection);
                } else if (loadingProgress >= 1) {
                    setSelection(null);
                }
                break;
            case 'principal':
                const principal = combinedReport.principals[selection.id];
                // TODO: adjust map based on principal activity?
                if (!principal && loadingProgress >= 1) {
                    setSelection(null);
                }
                break;
            case 'cluster':
                const minZoom = minZoomForClusterRect(selection.rect);
                const maxZoom = maxZoomForClusterRect(selection.rect);
                map.flyTo({
                    center: [selection.location.longitude, selection.location.latitude],
                    zoom: zoom < minZoom || zoom > maxZoom ? (minZoom + maxZoom) / 2 : zoom,
                });
                setLastCenteredSelection(selection);
                break;
        }
    }, [
        combinedReport,
        awsRegionsMap,
        selection,
        setSelection,
        lastCenteredSelection,
        setLastCenteredSelection,
        loadingProgress,
        map,
        isMapReady,
        zoom,
    ]);

    const needsData = teamReports && teamReports.length === 0;
    const teamAwsIntegrationsIfAdminAndNeedsData = useTeamAwsIntegrations(
        isTeamAdmin && needsData ? teamId : undefined,
    );
    const needsSubscriptionSetup = team && !team.entitlements.individualFeatures;
    const needsAwsIntegrationSetup =
        teamAwsIntegrationsIfAdminAndNeedsData && teamAwsIntegrationsIfAdminAndNeedsData.length === 0;
    const needsSetup = needsSubscriptionSetup || needsAwsIntegrationSetup || needsData;

    const canManageScps =
        teamAwsAccountsMap && [...teamAwsAccountsMap.values()].some((account) => account.canManageScps);

    return (
        <div className="w-full h-screen relative">
            <div className="absolute top-0 left-0 w-full h-full isolate">
                <WorldMap
                    mapStyle="https://tiles.openfreemap.org/styles/liberty"
                    onClick={() => {
                        setSelection(null);
                    }}
                    onStyleData={() => {
                        if (map) {
                            setZoom(map.getZoom());
                        }
                        setIsMapReady(true);
                    }}
                    onZoom={(e) => {
                        setZoom(e.viewState.zoom);
                    }}
                    ref={setMap}
                >
                    <MapOverlays
                        combinedReport={combinedReport}
                        highlightPrincipalId={hoveredPrincipalId}
                        zoom={zoom}
                    />
                </WorldMap>
            </div>
            <div className="absolute top-0 left-0 w-full h-full flex flex-col pointer-events-none">
                <Header>
                    <div className="flex gap-2">
                        <FilterDropdown
                            label="AWS Account"
                            filter={awsAccountFilter}
                            onChange={setAwsAccountFilter}
                            options={sortedAvailableAwsAccounts.map((account) => ({
                                label: account.name || account.id,
                                subLabel: account.name ? account.id : undefined,
                                value: account.id,
                            }))}
                        />
                        <FilterDropdown
                            label="AWS Region"
                            filter={awsRegionFilter}
                            onChange={setAwsRegionFilter}
                            options={sortedAvailableAwsRegions.map((region) => ({
                                label: region.name || region.id,
                                subLabel: region.name ? region.id : undefined,
                                value: region.id,
                            }))}
                        />
                        <DurationDropdown
                            durationSeconds={durationSeconds}
                            availableStartTime={minUnfilteredReportTime}
                            availableEndTime={maxUnfilteredReportEndTime}
                            onChange={setDurationSeconds}
                        />
                        <Dialog isOpen={isRulesOpen} onClose={() => setIsRulesOpen(false)} title="Rules" size="xl">
                            <Rules />
                        </Dialog>
                        {canManageScps && (
                            <div
                                className="inline-flex items-center gap-2 rounded-md bg-mint hover:bg-emerald transition-all duration-200 ease-in-out ml-8 py-1.5 px-3 text-sm/6 font-semibold text-snow cursor-pointer"
                                onClick={() => setIsRulesOpen(true)}
                            >
                                <ShieldCheckIcon className="h-[1.2rem]" />
                                Rules
                            </div>
                        )}
                    </div>
                </Header>
                <main className="w-full grow p-4 flex flex-row min-h-0 pointer-events-none relative">
                    <div className="w-1/4 flex flex-col [&>*]:pointer-events-auto">
                        <NavigationPanel
                            combinedReport={combinedReport}
                            onPrincipalHover={setHoveredPrincipalId}
                            principalFilter={principalFilter}
                            onPrincipalFilterChange={setPrincipalFilter}
                        />
                        {loadingProgress < 1 && (
                            <p className="mt-2 translucent-snow p-2 rounded-md text-sm w-full">
                                Loading... {Math.round(loadingProgress * 100)}%
                            </p>
                        )}
                    </div>
                    <div className="grow" />
                    <div className="w-1/4 flex flex-col [&>*]:pointer-events-auto">
                        {selection && <ContextPanel combinedReport={combinedReport} selection={selection} />}
                    </div>
                    {needsSetup && isTeamAdmin && (
                        <div className="absolute inset-0 flex w-screen items-center justify-center p-4 pointer-events-auto bg-radial from-black/10 to-black/40">
                            <div className="flex flex-col gap-4 translucent-snow border border-platinum rounded-xl p-8 max-w-2xl">
                                <h1>Welcome!</h1>
                                <p>This is the dashboard for your team.</p>
                                <p>There&apos;s currently nothing here, but we can fix that!</p>
                                {needsSubscriptionSetup && (
                                    <>
                                        <p className="uppercase label">Activate a subscription</p>
                                        <p>
                                            Head over to your team&apos;s{' '}
                                            <Link href={`/teams/${teamId}/settings/billing`} className="link">
                                                billing settings
                                            </Link>
                                            , where you can provide required billing information and pick your
                                            subscription tier.
                                        </p>
                                    </>
                                )}
                                {needsAwsIntegrationSetup ? (
                                    <>
                                        <p className="uppercase label">Integrate with your AWS account</p>
                                        <p>
                                            Go to your team&apos;s{' '}
                                            <Link href={`/teams/${teamId}/settings/integrations`} className="link">
                                                integration settings
                                            </Link>{' '}
                                            to configure ingestion from your AWS account.
                                        </p>
                                    </>
                                ) : needsData ? (
                                    <>
                                        <p className="uppercase label">Wait for data</p>
                                        <p>
                                            If you&apos;ve recently configured your AWS integration and chose to
                                            backfill data, it will be available here shortly.
                                        </p>
                                        <p>
                                            If you did not choose to backfill data, it will take up to 24 hours for data
                                            to appear.
                                        </p>
                                        <p>
                                            If data does not appear as expected, please double check your team&apos;s{' '}
                                            <Link href={`/teams/${teamId}/settings/integrations`} className="link">
                                                integration settings
                                            </Link>{' '}
                                            or{' '}
                                            <Link href="/contact" className="link">
                                                contact us
                                            </Link>{' '}
                                            for help.
                                        </p>
                                    </>
                                ) : null}
                            </div>
                        </div>
                    )}
                </main>
            </div>
        </div>
    );
};

export default Page;
