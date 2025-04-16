import { Report as ApiReport } from '@/generated/api';

export interface Report {
    networkLocations?: Record<string, Location>;
    ipAddressNetworks?: Record<string, string>;
    principals?: Record<string, Principal>;
}

export interface Location {
    latitude: number;
    longitude: number;
    countryCode: string;
    countryName: string;
    cityName: string;
    subdivisionNames?: string[];
}

export type PrincipalType =
    | 'AWSAssumedRole'
    | 'AWSService'
    | 'AWSIAMUser'
    | 'AWSAccount'
    | 'AWSRole'
    | 'WebIdentityUser';

export const formatPrincipalType = (type: PrincipalType): string => {
    switch (type) {
        case 'AWSAssumedRole':
            return 'AWS Assumed Role';
        case 'AWSService':
            return 'AWS Service';
        case 'AWSIAMUser':
            return 'AWS IAM User';
        case 'AWSAccount':
            return 'AWS Account';
        case 'AWSRole':
            return 'AWS Role';
        case 'WebIdentityUser':
            return 'Web Identity User';
        default:
            return 'Unknown';
    }
};

export interface Principal {
    name?: string;
    type?: PrincipalType;
    arn?: string;
    ipAddresses?: Record<string, number>;
    userAgents?: Record<string, number>;
    events: Record<string, EventSummary>;
}

export interface EventSummary {
    name: string;
    source: string;
    count: number;
    errorCodes?: Record<string, number>;
}

export class CombinedReport {
    networkLocations: Record<string, Location>;
    ipAddressNetworks: Record<string, string>;
    principals: Record<string, CombinedReportPrincipal>;
    awsRegionIds: Set<string>;

    constructor(reports?: ReportWithApiModel[]) {
        this.networkLocations = {};
        this.ipAddressNetworks = {};
        this.principals = {};
        this.awsRegionIds = new Set();

        reports?.forEach(({ report, apiReport }) => {
            if (report.networkLocations) {
                Object.entries(report.networkLocations).forEach(([cidr, location]) => {
                    this.addNetworkLocation(cidr, location);
                });
            }

            if (report.ipAddressNetworks) {
                Object.entries(report.ipAddressNetworks).forEach(([ip, network]) => {
                    this.addIpAddressNetwork(ip, network);
                });
            }

            if (report.principals) {
                Object.entries(report.principals).forEach(([id, principal]) => {
                    this.addPrincipalFromApiReport(id, principal, apiReport);
                });
            }
        });
    }

    addNetworkLocation(cidr: string, location: Location) {
        this.networkLocations[cidr] = location;
    }

    addIpAddressNetwork(ip: string, network: string) {
        this.ipAddressNetworks[ip] = network;
    }

    getOrAddPrincipal(
        id: string,
        principal: { name?: string; type?: PrincipalType; arn?: string },
    ): CombinedReportPrincipal {
        let p = this.principals[id];
        if (!p) {
            p = new CombinedReportPrincipal(principal);
            this.principals[id] = p;
        }
        return p;
    }

    addPrincipalFromApiReport(id: string, principal: Principal, apiReport: ApiReport) {
        const p = this.getOrAddPrincipal(id, principal);

        p.accountIds.add(sanitize(apiReport.scope.aws.accountId));
        p.awsRegionIds.add(apiReport.scope.aws.region);
        this.awsRegionIds.add(apiReport.scope.aws.region);

        if (principal.ipAddresses) {
            for (const [ip, count] of Object.entries(principal.ipAddresses)) {
                p.ipAddresses.set(ip, (p.ipAddresses.get(ip) || 0) + count);
                const cidr = this.ipAddressNetworks[ip];
                if (cidr) {
                    p.networkCidrs.add(cidr);
                }
            }
        }

        if (principal.userAgents) {
            for (const [agent, count] of Object.entries(principal.userAgents)) {
                p.userAgents.set(agent, (p.userAgents.get(agent) || 0) + count);
            }
        }

        Object.entries(principal.events).forEach(([id, summary]) => {
            p.addEventSummary(id, summary);
        });
    }

    addCombinedReportPrincipal(other: CombinedReport, id: string) {
        const otherPrincipal = other.principals[id];
        const p = this.getOrAddPrincipal(id, otherPrincipal);

        otherPrincipal.accountIds.forEach((id) => {
            p.accountIds.add(id);
        });

        otherPrincipal.awsRegionIds.forEach((id) => {
            p.awsRegionIds.add(id);
            this.awsRegionIds.add(id);
        });

        for (const [ip, count] of otherPrincipal.ipAddresses) {
            p.ipAddresses.set(ip, (p.ipAddresses.get(ip) || 0) + count);
            const cidr = other.ipAddressNetworks[ip];
            if (cidr) {
                const location = other.networkLocations[cidr];
                if (location) {
                    this.addNetworkLocation(cidr, location);
                }
                this.addIpAddressNetwork(ip, cidr);
                p.networkCidrs.add(cidr);
            }
        }

        for (const [agent, count] of otherPrincipal.userAgents) {
            p.userAgents.set(agent, (p.userAgents.get(agent) || 0) + count);
        }

        for (const [id, summary] of otherPrincipal.events) {
            p.addEventSummary(id, summary);
        }
    }

    withFilteredPrincipals(filter: string): CombinedReport {
        const ret = new CombinedReport();

        const normalizedFilter = filter.toLowerCase();
        Object.entries(this.principals).forEach(([id, principal]) => {
            if (
                principal.name?.toLowerCase().includes(normalizedFilter) ||
                id.toLowerCase().includes(normalizedFilter) ||
                principal.arn?.toLowerCase().includes(normalizedFilter)
            ) {
                ret.addCombinedReportPrincipal(this, id);
            }
        });

        return ret;
    }
}

const sanitizedAccountIds = new Map<string, string>();

const randomAccountId = () => {
    const random = Math.floor(Math.random() * 1e12);
    return random.toString().padStart(12, '0');
};

// If NEXT_PUBLIC_STREAMER_MODE is set, replace all account IDs in the string with random 12 digit
// numbers. Useful for public screenshots.
const sanitize = (s: string): string => {
    if (!process.env.NEXT_PUBLIC_STREAMER_MODE || !s) {
        return s;
    }

    // find everything that looks like a 12 digit number in the string using regex
    const accountIds = s.match(/\b\d{12}\b/g);
    if (accountIds) {
        for (const accountId of accountIds) {
            let replacement = sanitizedAccountIds.get(accountId);
            if (!replacement) {
                replacement = randomAccountId();
                sanitizedAccountIds.set(accountId, replacement);
            }
            s = s.replace(accountId, replacement);
        }
    }

    return s;
};

export class CombinedReportPrincipal {
    name?: string;
    type?: PrincipalType;
    arn?: string;
    ipAddresses: Map<string, number>;
    userAgents: Map<string, number>;
    accountIds: Set<string>;
    awsRegionIds: Set<string>;
    networkCidrs: Set<string>;
    eventCount: number;
    events: Map<string, EventSummary>;

    constructor(params: { name?: string; type?: PrincipalType; arn?: string }) {
        this.name = params.name && sanitize(params.name);
        this.type = params.type;
        this.arn = params.arn && sanitize(params.arn);
        this.ipAddresses = new Map();
        this.userAgents = new Map();
        this.accountIds = new Set();
        this.awsRegionIds = new Set();
        this.networkCidrs = new Set();
        this.eventCount = 0;
        this.events = new Map();
    }

    addEventSummary(id: string, summary: EventSummary) {
        let e = this.events.get(id);
        if (!e) {
            e = {
                name: summary.name,
                source: summary.source,
                count: 0,
            };
            this.events.set(id, e);
        }
        this.eventCount += summary.count;
        e.count += summary.count;
        for (const [code, count] of Object.entries(summary.errorCodes || {})) {
            if (!e.errorCodes) {
                e.errorCodes = {};
            }
            e.errorCodes[code] = (e.errorCodes[code] || 0) + count;
        }
    }
}

interface ReportWithApiModel {
    report: Report;
    apiReport: ApiReport;
}
