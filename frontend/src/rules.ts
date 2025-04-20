import { AwsPolicy } from '@/aws';

export interface ServiceAllowlistRule {
    type: 'service_allowlist';
    services: Set<string>;
}

export interface RegionAllowlistRule {
    type: 'region_allowlist';
    regions: Set<string>;
}

export type Rule = ServiceAllowlistRule | RegionAllowlistRule;

const toStrings = (thing: string | string[] | undefined): string[] => {
    if (!thing) {
        return [];
    }
    if (typeof thing === 'string') {
        return [thing];
    }
    return thing;
};

export class RuleSet {
    regionAllowlist?: RegionAllowlistRule;
    serviceAllowlist?: ServiceAllowlistRule;

    static fromScpContent(content: string): RuleSet {
        const policy = JSON.parse(content) as AwsPolicy;
        const ret = new RuleSet();
        for (const statement of policy.Statement) {
            switch (statement.Sid) {
                case 'ServiceAllowlist':
                    ret.serviceAllowlist = {
                        type: 'service_allowlist',
                        services: new Set(toStrings(statement.NotAction).map((action) => action.split(':')[0])),
                    };
                    break;
                case 'RegionAllowlist':
                    ret.regionAllowlist = {
                        type: 'region_allowlist',
                        regions: new Set(statement.Condition?.StringNotEquals['aws:RequestedRegion'] || []),
                    };
                    break;
            }
        }
        return ret;
    }

    // Creates a deep clone of the RuleSet.
    clone(): RuleSet {
        const ret = new RuleSet();
        if (this.serviceAllowlist) {
            ret.serviceAllowlist = {
                type: 'service_allowlist',
                services: new Set(this.serviceAllowlist.services),
            };
        }
        if (this.regionAllowlist) {
            ret.regionAllowlist = {
                type: 'region_allowlist',
                regions: new Set(this.regionAllowlist.regions),
            };
        }
        return ret;
    }

    equal(other: RuleSet): boolean {
        if (this.serviceAllowlist && other.serviceAllowlist) {
            if (this.serviceAllowlist.services.size !== other.serviceAllowlist.services.size) {
                return false;
            }
            for (const service of this.serviceAllowlist.services) {
                if (!other.serviceAllowlist.services.has(service)) {
                    return false;
                }
            }
        } else if (this.serviceAllowlist || other.serviceAllowlist) {
            return false;
        }

        if (this.regionAllowlist && other.regionAllowlist) {
            if (this.regionAllowlist.regions.size !== other.regionAllowlist.regions.size) {
                return false;
            }
            for (const region of this.regionAllowlist.regions) {
                if (!other.regionAllowlist.regions.has(region)) {
                    return false;
                }
            }
        } else if (this.regionAllowlist || other.regionAllowlist) {
            return false;
        }

        return true;
    }

    addRegionToAllowlist(region: string) {
        if (!this.regionAllowlist) {
            this.regionAllowlist = { type: 'region_allowlist', regions: new Set() };
        }
        this.regionAllowlist.regions.add(region);
    }

    removeRegionFromAllowlist(region: string) {
        if (this.regionAllowlist) {
            this.regionAllowlist.regions.delete(region);
            if (this.regionAllowlist.regions.size === 0) {
                this.regionAllowlist = undefined;
            }
        }
    }

    hasRegionInAllowlist(region: string): boolean {
        return this.regionAllowlist ? this.regionAllowlist.regions.has(region) : false;
    }

    addServiceToAllowlist(service: string) {
        if (!this.serviceAllowlist) {
            this.serviceAllowlist = { type: 'service_allowlist', services: new Set() };
        }
        this.serviceAllowlist.services.add(service);
    }

    removeServiceFromAllowlist(service: string) {
        if (this.serviceAllowlist) {
            this.serviceAllowlist.services.delete(service);
            if (this.serviceAllowlist.services.size === 0) {
                this.serviceAllowlist = undefined;
            }
        }
    }

    hasServiceInAllowlist(service: string): boolean {
        return this.serviceAllowlist ? this.serviceAllowlist.services.has(service) : false;
    }

    scp(): AwsPolicy {
        const policy: AwsPolicy = {
            Version: '2012-10-17',
            Statement: [],
        };

        if (this.serviceAllowlist) {
            policy.Statement.push({
                Sid: 'ServiceAllowlist',
                Effect: 'Deny',
                NotAction: Array.from(this.serviceAllowlist.services).map((service) => `${service}:*`),
                Resource: '*',
            });
        }

        if (this.regionAllowlist) {
            policy.Statement.push({
                Sid: 'RegionAllowlist',
                Effect: 'Deny',
                // Take special care with these core global services. They either work in all regions or not at all.
                Action: this.hasRegionInAllowlist('us-east-1') ? '*' : undefined,
                NotAction: this.hasRegionInAllowlist('us-east-1')
                    ? undefined
                    : ['iam:*', 'organizations:*', 'account:*'],
                Resource: '*',
                Condition: {
                    StringNotEquals: {
                        'aws:RequestedRegion': Array.from(this.regionAllowlist.regions),
                    },
                },
            });
        }

        return policy;
    }
}
