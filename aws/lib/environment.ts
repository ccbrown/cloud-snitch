import { Construct } from 'constructs';
import * as crypto from 'crypto';

import { GlobalApexStack } from './global-apex-stack';
import { GlobalBaseStack } from './global-base-stack';
import { RegionalStack } from './regional-stack';

interface Pricing {
    individualSubscriptionStripePriceId: string;
    teamSubscriptionStripePriceId: string;
}

interface Props {
    accountId: string;
    cloudfrontPublicKey: string;
    contactEmailAddress: string;
    domainName: string;
    regions: string[];
    slug: string;
    stripePublishableKey: string;
    pricing: Pricing;
    stripeEventSourceName?: string;
    userRegistrationAllowlist?: string[];
    noIndex?: boolean;
}

export class Environment {
    constructor(scope: Construct, props: Props) {
        // Generate a hash specific to this environment for the purpose of creating unique but deterministic names.
        const envHash = crypto
            .createHash('sha256')
            .update(`${props.accountId}:${props.slug}:nonce-BqJxMPNvgs`)
            .digest('hex')
            .slice(0, 12);

        const globalBaseStack = new GlobalBaseStack(scope, `global-base-${props.slug}`, {
            cloudfrontPublicKey: props.cloudfrontPublicKey,
            crossRegionReferences: true,
            domainName: props.domainName,
            env: {
                account: props.accountId,
                region: 'us-east-1',
            },
            envHash,
            envSlug: props.slug,
            regions: props.regions,
            stackName: `cloud-snitch-global-base-${props.slug}`,
        });

        const regionalStacks: RegionalStack[] = [];
        for (const region of props.regions) {
            const regionalStack = new RegionalStack(scope, `${region}-${props.slug}`, {
                allRegions: props.regions,
                cloudfrontKeyGroupId: globalBaseStack.cloudfrontKeyGroupId,
                cloudfrontPrivateKeySecretName: globalBaseStack.cloudfrontPrivateKeySecretName,
                cloudfrontPublicKeyId: globalBaseStack.cloudfrontPublicKeyId,
                contactEmailAddress: props.contactEmailAddress,
                crossRegionReferences: true,
                domainName: props.domainName,
                dynamodbTableName: globalBaseStack.dynamodbTableName,
                env: {
                    account: props.accountId,
                    region: region,
                },
                envHash,
                globalCertificate: globalBaseStack.certificate,
                passwordEncryptionKeyBase64SecretName: globalBaseStack.passwordEncryptionKeyBase64SecretName,
                pricing: props.pricing,
                proxySecretName: globalBaseStack.proxySecretName,
                stackName: `cloud-snitch-${region}-${props.slug}`,
                stripeEventSourceName: region === props.regions[0] ? props.stripeEventSourceName : undefined,
                stripePublishableKey: props.stripePublishableKey,
                stripeSecretKeySecretName: globalBaseStack.stripeSecretKeySecretName,
                triggerReportGeneration: region === props.regions[0],
                userRegistrationAllowlist: props.userRegistrationAllowlist,
                noIndex: props.noIndex,
            });
            regionalStacks.push(regionalStack);
            regionalStack.addDependency(globalBaseStack);
        }

        const globalApexStack = new GlobalApexStack(scope, `global-apex-${props.slug}`, {
            cloudfrontDistributionId: globalBaseStack.cloudfrontDistributionId,
            crossRegionReferences: true,
            env: {
                account: props.accountId,
                region: 'us-east-1',
            },
            stackName: `cloud-snitch-global-apex-${props.slug}`,
        });
        regionalStacks.forEach((s) => globalApexStack.addDependency(s));
    }
}
