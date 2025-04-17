import RegionIcon from '@/assets/aws/Architecture-Group-Icons_02072025/Region_32.svg';
import ServerIcon from '@/assets/aws/Architecture-Group-Icons_02072025/Server-contents_32.svg';
import { CheckIcon, PencilIcon, XMarkIcon } from '@heroicons/react/24/outline';
import IPCIDR from 'ip-cidr';
import Image from 'next/image';
import Link from 'next/link';
import { useMemo, useState } from 'react';

import { MapLocation, MapRect } from './map';
import { Selection, useSelection } from './selection';
import { Markdown, PrincipalIcon, TextArea, Tooltip } from '@/components';
import {
    useAwsRegions,
    useCurrentTeamId,
    useCurrentTeamPrincipalSettings,
    useSearchParamFilterState,
    useTeamAwsAccountsMap,
} from '@/hooks';
import { CombinedReport, formatPrincipalType } from '@/report';
import { useDispatch, useSelector } from '@/store';

const COLLAPSED_LIST_LENGTH = 5;

interface AwsRegionContextProps {
    id: string;
}

const AwsRegionContext = (props: AwsRegionContextProps) => {
    const [awsRegionFilter, setAwsRegionFilter] = useSearchParamFilterState('awsRegion');

    const region = useSelector((state) => state.aws.regions[props.id]);
    if (!region) {
        return null;
    }

    return (
        <div>
            <div className="flex gap-2 items-center mb-2">
                <Image src={RegionIcon} alt={region.id} className="rounded-lg shrink-0 h-[4rem] w-[4rem]" />
                <div className="flex flex-col">
                    <h1>{region.name}</h1>
                    <span className="text-sm text-english-violet">AWS Region</span>
                </div>
            </div>
            <p>
                <strong>Id:</strong> {props.id}
            </p>
            <p>
                <strong>Name:</strong> {region.name}
            </p>
            <p>
                <strong>Geolocation:</strong>{' '}
                {MapLocation.fromLatitudeAndLongitude(region.latitude, region.longitude).toString()}
            </p>
            <p>
                <strong>Geolocation Country:</strong> {region.geolocationCountry}
            </p>
            <p>
                <strong>Geolocation Region:</strong> {region.geolocationRegion}
            </p>
            <p>
                <strong>Partition:</strong> {region.partition}
            </p>
            {(awsRegionFilter.mode !== 'include' ||
                awsRegionFilter.values.size !== 1 ||
                !awsRegionFilter.values.has(props.id)) && (
                <div>
                    <button
                        className="button text-sm mt-2"
                        onClick={() =>
                            setAwsRegionFilter({
                                mode: 'include',
                                values: new Set([props.id]),
                            })
                        }
                    >
                        Show Only This Region
                    </button>
                </div>
            )}
        </div>
    );
};

interface NetworkContextProps {
    cidr: string;
    combinedReport: CombinedReport;
}

const NetworkContext = ({ cidr, combinedReport }: NetworkContextProps) => {
    const setSelection = useSelection()[1];
    const [showAllIpAddresses, setShowAllIpAddresses] = useState(false);
    const [showAllPrincipals, setShowAllPrincipals] = useState(false);

    const ipAddresses = useMemo(() => {
        const ipAddresses = new Set<string>();
        Object.entries(combinedReport.ipAddressNetworks).forEach(([ip, cidr]) => {
            if (cidr === cidr) {
                ipAddresses.add(ip);
            }
        });
        return Array.from(ipAddresses.values()).sort((a, b) => a.localeCompare(b));
    }, [combinedReport]);

    const principals = useMemo(
        () =>
            Object.entries(combinedReport.principals)
                .filter(([_id, principal]) => {
                    for (const addr of principal.ipAddresses.values()) {
                        if (cidr === combinedReport.ipAddressNetworks[addr]) {
                            return true;
                        }
                    }
                    return false;
                })
                .sort((a, b) => {
                    const aLabel = (a[1].name || a[0]).toLowerCase();
                    const bLabel = (b[1].name || b[0]).toLowerCase();
                    return aLabel.localeCompare(bLabel);
                }),
        [cidr, combinedReport],
    );

    const location = combinedReport.networkLocations[cidr];
    if (!location) {
        return null;
    }

    const cityAndSubdivisionNames = [location.cityName, ...(location.subdivisionNames || [])];

    const ipCidr = new IPCIDR(cidr);

    return (
        <div>
            <div className="flex gap-2 items-center mb-2">
                <Image src={ServerIcon} alt={cidr} className="rounded-lg shrink-0 h-[4rem] w-[4rem]" />
                <div className="flex flex-col">
                    <h1>{cidr}</h1>
                    <span className="text-sm text-english-violet">Network</span>
                </div>
            </div>
            <p>
                <strong>CIDR:</strong> {cidr}
            </p>
            <p>
                <strong>Range:</strong> {ipCidr.start()} - {ipCidr.end()}
            </p>
            <p>
                <strong>Address Count:</strong> {ipCidr.size.toLocaleString()}
            </p>
            <p>
                <strong>Country:</strong> {location.countryName} ({location.countryCode})
            </p>
            <p>
                <strong>City:</strong> {cityAndSubdivisionNames.join(', ')}
            </p>
            <div>
                <h2 className="my-2">Active IP Addresses ({ipAddresses.length.toLocaleString()})</h2>
                <ul className="text-sm list-disc ml-6 pr-2">
                    {ipAddresses
                        .slice(0, showAllIpAddresses ? ipAddresses.length : COLLAPSED_LIST_LENGTH)
                        .map((addr) => (
                            <li key={addr}>{addr}</li>
                        ))}
                </ul>
                {!showAllIpAddresses && ipAddresses.length > COLLAPSED_LIST_LENGTH && (
                    <div className="link text-xs mt-2" onClick={() => setShowAllIpAddresses(true)}>
                        Show All IP Addresses
                    </div>
                )}
            </div>
            {principals?.length > 0 && (
                <div>
                    <h2 className="my-2">Principals ({principals?.length.toLocaleString()})</h2>
                    <ul className="text-sm list-disc ml-6 pr-2">
                        {principals
                            .slice(0, showAllPrincipals ? principals.length : COLLAPSED_LIST_LENGTH)
                            .map(([id, p]) => (
                                <li key={id}>
                                    <div
                                        className="link w-full"
                                        onClick={() => {
                                            setSelection({
                                                type: 'principal',
                                                id,
                                            });
                                        }}
                                    >
                                        {p.name || id}
                                    </div>
                                </li>
                            ))}
                    </ul>
                    {!showAllPrincipals && principals.length > COLLAPSED_LIST_LENGTH && (
                        <div className="link text-xs mt-2" onClick={() => setShowAllPrincipals(true)}>
                            Show All Principals
                        </div>
                    )}
                </div>
            )}
            <p className="mt-4 text-xs text-right opacity-50 hover:opacity-100 transition-opacity duration-200">
                IP geolocation data by{' '}
                <Link href="https://db-ip.com" target="_blank" rel="noopener noreferrer" className="external-link">
                    DB-IP
                </Link>
                .
            </p>
        </div>
    );
};

interface PrincipalContextProps {
    id: string;
    combinedReport: CombinedReport;
}

interface IdAndName {
    id: string;
    name?: string;
}

const PrincipalContext = ({ id, combinedReport }: PrincipalContextProps) => {
    const setSelection = useSelection()[1];
    const dispatch = useDispatch();
    const [showAllIpAddresses, setShowAllIpAddresses] = useState(false);
    const [showAllUserAgents, setShowAllUserAgents] = useState(false);
    const [isEditingDescription, setIsEditingDescription] = useState(false);
    const [descriptionDraft, setDescriptionDraft] = useState<string>('');
    const [isBusy, setIsBusy] = useState(false);

    const principal = combinedReport.principals[id];
    const settings = useCurrentTeamPrincipalSettings(id);

    const allAwsRegions = useAwsRegions();

    const currentTeamId = useCurrentTeamId();
    const teamAwsAccountsMap = useTeamAwsAccountsMap(currentTeamId);

    const saveDescription = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.updatePrincipalSettings({
                teamId: currentTeamId,
                principalKey: id,
                input: {
                    description: descriptionDraft,
                },
            });
            setIsEditingDescription(false);
        } catch (err) {
            alert(err instanceof Error ? err.message : 'An unknown error occurred.');
        } finally {
            setIsBusy(false);
        }
    };

    const awsRegions = useMemo(
        () =>
            allAwsRegions
                .filter((region) => principal?.awsRegionIds.has(region.id))
                .sort((a, b) => a.id.localeCompare(b.id)),
        [allAwsRegions, principal],
    );

    const awsAccounts: IdAndName[] = useMemo(
        () =>
            Array.from(principal?.accountIds.values() || [])
                .map((id) => ({
                    id,
                    name: teamAwsAccountsMap?.get(id)?.name,
                }))
                .sort((a, b) => {
                    const aLabel = (a.name || a.id).toLowerCase();
                    const bLabel = (b.name || b.id).toLowerCase();
                    return aLabel.localeCompare(bLabel);
                }),
        [principal, teamAwsAccountsMap],
    );

    const ipAddresses = useMemo(
        () => Array.from(principal?.ipAddresses.keys() || []).sort((a, b) => a.localeCompare(b)),
        [principal],
    );

    const userAgents = useMemo(
        () => Array.from(principal?.userAgents.entries() || []).sort((a, b) => a[0].localeCompare(b[0])),
        [principal],
    );

    const events = useMemo(
        () => Array.from(principal?.events.entries() || []).sort((a, b) => a[1].name.localeCompare(b[1].name)),
        [principal],
    );

    if (!principal) {
        return null;
    }

    return (
        <div>
            <div className="flex gap-2 items-center mb-2 w-full">
                <PrincipalIcon className="rounded-lg shrink-0 h-[4rem] w-[4rem]" type={principal.type} />
                <div className="flex flex-col pr-2 overflow-hidden">
                    <h1 className="truncate">{principal.name || id}</h1>
                    {principal.type && (
                        <span className="text-sm text-english-violet">{formatPrincipalType(principal.type)}</span>
                    )}
                </div>
            </div>
            {id !== principal.arn && (
                <p className="break-words">
                    <strong>Id:</strong> {id}
                </p>
            )}
            {principal.arn && (
                <p className="break-words">
                    <strong>ARN:</strong> {principal.arn}
                </p>
            )}
            <div className="flex flex-col gap-1 my-4">
                <div className="flex items-center gap-1">
                    <strong>Description</strong>
                    {isEditingDescription ? (
                        <>
                            <CheckIcon
                                className="h-[1rem] cursor-pointer hover:text-amethyst"
                                onClick={() => saveDescription()}
                            />
                            <XMarkIcon
                                className="h-[1rem] cursor-pointer hover:text-amethyst"
                                onClick={() => setIsEditingDescription(false)}
                            />
                        </>
                    ) : (
                        <PencilIcon
                            className="h-[1rem] cursor-pointer hover:text-amethyst"
                            onClick={() => {
                                setDescriptionDraft(settings?.description || '');
                                setIsEditingDescription(true);
                            }}
                        />
                    )}
                </div>
                {isEditingDescription ? (
                    <div>
                        <TextArea disabled={isBusy} value={descriptionDraft} onChange={setDescriptionDraft} />
                    </div>
                ) : settings?.description ? (
                    <div>
                        <Markdown>{settings.description}</Markdown>
                    </div>
                ) : settings ? (
                    <div className="italic text-sm">No description set.</div>
                ) : (
                    <div className="italic text-sm">Loading...</div>
                )}
            </div>
            <div>
                <h2 className="my-2">AWS Accounts ({awsAccounts.length})</h2>
                <ul className="text-sm list-disc ml-6">
                    {awsAccounts.map((account) => (
                        <li key={account.id}>{account.name ? `${account.name} (${account.id})` : account.id}</li>
                    ))}
                </ul>
            </div>
            {awsRegions?.length > 0 && (
                <div>
                    <h2 className="my-2">AWS Regions ({awsRegions.length})</h2>
                    <ul className="text-sm list-disc ml-6">
                        {awsRegions.map((region) => (
                            <li key={region.id}>
                                <span
                                    className="link"
                                    onClick={() => {
                                        setSelection({
                                            type: 'aws-region',
                                            id: region.id,
                                        });
                                    }}
                                >
                                    {region.name} ({region.id})
                                </span>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
            {ipAddresses.length > 0 && (
                <div>
                    <h2 className="my-2">IP Addresses ({ipAddresses.length.toLocaleString()})</h2>
                    <ul className="text-sm list-disc ml-6 pr-2">
                        {ipAddresses
                            .slice(0, showAllIpAddresses ? ipAddresses.length : COLLAPSED_LIST_LENGTH)
                            .map((addr) => {
                                const cidr = combinedReport.ipAddressNetworks[addr];
                                return (
                                    <li key={addr}>
                                        {addr}
                                        {cidr && (
                                            <>
                                                {' '}
                                                (
                                                <span
                                                    className="link"
                                                    onClick={() => {
                                                        setSelection({
                                                            type: 'network',
                                                            cidr,
                                                        });
                                                    }}
                                                >
                                                    {cidr}
                                                </span>
                                                )
                                            </>
                                        )}
                                    </li>
                                );
                            })}
                    </ul>
                    {!showAllIpAddresses && ipAddresses.length > COLLAPSED_LIST_LENGTH && (
                        <div className="link text-xs mt-2" onClick={() => setShowAllIpAddresses(true)}>
                            Show All IP Addresses
                        </div>
                    )}
                </div>
            )}
            <div>
                <h2 className="my-2">Event Types ({events.length})</h2>
                <ul className="text-sm list-disc ml-6 pr-2">
                    {events.map(([key, summary]) => (
                        <li key={key}>
                            <Tooltip
                                content={
                                    <div className="whitespace-nowrap flex flex-col gap-1">
                                        <div>
                                            <strong>Source:</strong> {summary.source}
                                        </div>
                                        <div>
                                            <Link
                                                href={`https://console.aws.amazon.com/cloudtrailv2/home#/events?EventName=${summary.name}`}
                                                className="external-link text-xs"
                                                target="_blank"
                                            >
                                                View Events of This Type in CloudTrail
                                            </Link>
                                        </div>
                                    </div>
                                }
                            >
                                <span className="hoverable">{summary.name}</span>
                            </Tooltip>{' '}
                            <span className="chip">{summary.count}</span>
                            {summary.errorCodes && (
                                <Tooltip
                                    content={
                                        <div className="whitespace-nowrap">
                                            {Object.entries(summary.errorCodes).map(([code, count]) => (
                                                <div key={code}>
                                                    {code} <span className="chip">{count}</span>
                                                </div>
                                            ))}
                                        </div>
                                    }
                                >
                                    <span className="chip hoverable bg-indian-red">
                                        {Object.entries(summary.errorCodes).reduce(
                                            (acc, [_code, count]) => acc + count,
                                            0,
                                        )}{' '}
                                        Errors
                                    </span>
                                </Tooltip>
                            )}
                        </li>
                    ))}
                </ul>
            </div>
            <div>
                <h2 className="my-2">User Agents ({userAgents.length.toLocaleString()})</h2>
                <ul className="text-sm list-disc ml-6 pr-2">
                    {userAgents
                        .slice(0, showAllUserAgents ? userAgents.length : COLLAPSED_LIST_LENGTH)
                        .map(([agent, eventCount]) => {
                            return (
                                <li key={agent}>
                                    {agent} <span className="chip">{eventCount}</span>
                                </li>
                            );
                        })}
                </ul>
                {!showAllUserAgents && userAgents.length > COLLAPSED_LIST_LENGTH && (
                    <div className="link text-xs mt-2" onClick={() => setShowAllUserAgents(true)}>
                        Show All User Agents
                    </div>
                )}
            </div>
        </div>
    );
};

interface ClusterContextProps {
    combinedReport: CombinedReport;
    mapRect: MapRect;
    mapLocation: MapLocation;
}

const ClusterContext = ({ combinedReport, mapRect, mapLocation }: ClusterContextProps) => {
    const setSelection = useSelection()[1];
    const [showAllNetworks, setShowAllNetworks] = useState(false);
    const [showAllPrincipals, setShowAllPrincipals] = useState(false);

    const allAwsRegions = useAwsRegions();

    const awsRegions = useMemo(
        () =>
            allAwsRegions
                .filter((region) => {
                    const location = MapLocation.fromLatitudeAndLongitude(region.latitude, region.longitude);
                    return mapRect.containsLocation(location) && combinedReport.awsRegionIds.has(region.id);
                })
                .sort((a, b) => a.id.localeCompare(b.id)),
        [combinedReport, allAwsRegions, mapRect],
    );

    const networks = useMemo(
        () =>
            new Set(
                Object.entries(combinedReport.networkLocations)
                    .filter(([_cidr, location]) =>
                        mapRect.containsLocation(
                            MapLocation.fromLatitudeAndLongitude(location.latitude, location.longitude),
                        ),
                    )
                    .map(([cidr]) => cidr),
            ),
        [combinedReport, mapRect],
    );
    const sortedNetworks = networks && Array.from(networks.values()).sort((a, b) => a.localeCompare(b));

    const principals = useMemo(
        () =>
            Object.entries(combinedReport.principals)
                .filter(([_id, principal]) => {
                    for (const addr of principal.ipAddresses.values()) {
                        const network = combinedReport.ipAddressNetworks[addr];
                        if (networks.has(network)) {
                            return true;
                        }
                    }
                    return false;
                })
                .sort((a, b) => {
                    const aLabel = (a[1].name || a[0]).toLowerCase();
                    const bLabel = (b[1].name || b[0]).toLowerCase();
                    return aLabel.localeCompare(bLabel);
                }),
        [combinedReport, networks],
    );

    const clusterSize = (awsRegions?.length || 0) + (networks?.size || 0);

    return (
        <div>
            <div className="flex gap-2 items-center mb-2">
                <div className="rounded-lg shrink-0 h-[4rem] w-[4rem] bg-[#7D8998] text-snow flex items-center justify-center font-semibold">
                    {clusterSize.toLocaleString()}
                </div>
                <div className="flex flex-col">
                    <h1>Cluster</h1>
                    <span className="text-sm text-english-violet">{mapLocation.toString()}</span>
                </div>
            </div>
            <p>
                <strong>Geolocation:</strong> {mapLocation.toString()}
            </p>
            {awsRegions?.length > 0 && (
                <div>
                    <h2 className="my-2">AWS Regions ({awsRegions.length})</h2>
                    <ul className="text-sm list-disc ml-6">
                        {awsRegions.map((region) => (
                            <li key={region.id}>
                                <span
                                    className="link"
                                    onClick={() => {
                                        setSelection({
                                            type: 'aws-region',
                                            id: region.id,
                                        });
                                    }}
                                >
                                    {region.name} ({region.id})
                                </span>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
            {sortedNetworks?.length > 0 && (
                <div>
                    <h2 className="my-2">Networks ({sortedNetworks?.length.toLocaleString()})</h2>
                    <ul className="text-sm list-disc ml-6">
                        {sortedNetworks
                            .slice(0, showAllNetworks ? sortedNetworks.length : COLLAPSED_LIST_LENGTH)
                            .map((cidr) => (
                                <li key={cidr}>
                                    <span
                                        className="link"
                                        onClick={() => {
                                            setSelection({
                                                type: 'network',
                                                cidr,
                                            });
                                        }}
                                    >
                                        {cidr}
                                    </span>
                                </li>
                            ))}
                    </ul>
                    {!showAllNetworks && sortedNetworks.length > COLLAPSED_LIST_LENGTH && (
                        <div className="link text-xs mt-2" onClick={() => setShowAllNetworks(true)}>
                            Show All Networks
                        </div>
                    )}
                </div>
            )}
            {principals?.length > 0 && (
                <div>
                    <h2 className="my-2">Principals ({principals?.length.toLocaleString()})</h2>
                    <ul className="text-sm list-disc ml-6 pr-2">
                        {principals
                            .slice(0, showAllPrincipals ? principals.length : COLLAPSED_LIST_LENGTH)
                            .map(([id, p]) => (
                                <li key={id}>
                                    <div
                                        className="link w-full"
                                        onClick={() => {
                                            setSelection({
                                                type: 'principal',
                                                id,
                                            });
                                        }}
                                    >
                                        {p.name || id}
                                    </div>
                                </li>
                            ))}
                    </ul>
                    {!showAllPrincipals && principals.length > COLLAPSED_LIST_LENGTH && (
                        <div className="link text-xs mt-2" onClick={() => setShowAllPrincipals(true)}>
                            Show All Principals
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

interface Props {
    combinedReport: CombinedReport;
    selection: Selection;
}

export const ContextPanel = (props: Props) => {
    return (
        <div className="translucent-snow p-4 rounded-lg w-full overflow-y-auto">
            {props.selection.type === 'aws-region' && <AwsRegionContext id={props.selection.id} />}
            {props.selection.type === 'network' && (
                <NetworkContext cidr={props.selection.cidr} combinedReport={props.combinedReport} />
            )}
            {props.selection.type === 'principal' && (
                <PrincipalContext id={props.selection.id} combinedReport={props.combinedReport} />
            )}
            {props.selection.type === 'cluster' && (
                <ClusterContext
                    combinedReport={props.combinedReport}
                    mapRect={props.selection.rect}
                    mapLocation={props.selection.location}
                />
            )}
        </div>
    );
};
