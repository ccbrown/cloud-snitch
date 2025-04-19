#!/usr/bin/env node
import { App } from 'aws-cdk-lib';

import { Environment } from '../lib/environment';
import { GithubActionsStack } from '../lib/github-actions-stack';

const app = new App();

new Environment(app, {
    accountId: '774305579662',
    cloudfrontPublicKey: [
        '-----BEGIN PUBLIC KEY-----',
        'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA23gKbJDDmPBviNUU062p',
        'bnMu7R2ImeYf6ZNiuU7yBHuKZk+j2Z5iu+aQylASIBD6GNnDfpzpA84LLJQUfw6+',
        'vl5m/aDlnxSO/tBo3ztn57k0G5tAzbbuUKHp0MXrUWeyQICIUNRGryOHhYNRjQNa',
        'DWGy951QKMWLd7frCBJYoLRXAanMh7Of8ErmGH4+xFCyHEMTi7ik0QaBA85Rv4q6',
        'AfSJEcOblJ0P/vFOWqAcDjdHOQZeCNHoO5A6ZOlnvRcX2gYTdo2CFMvPF3X7Afm+',
        'lkaJYrARAs80iNexcTdZQEFmI7HtbQUhe1t08wyD9JJopN5p87KCySLAoSDyNJsp',
        'fwIDAQAB',
        '-----END PUBLIC KEY-----',
    ].join('\n'),
    contactEmailAddress: 'cloudsnitch-dev-support@paragoncybersecurity.sh',
    domainName: 'cloud-snitch.ccbrown.dev',
    pricing: {
        individualSubscriptionStripePriceId: 'price_1R9dvs2ejpbHZUu9R4HllgRg',
        teamSubscriptionStripePriceId: 'price_1R9dwk2ejpbHZUu9ByxeAAjC',
    },
    regions: ['us-east-1', 'us-west-2'],
    slug: 'dev',
    stripeEventSourceName: 'aws.partner/stripe.com/ed_test_61SLFDXYPHDgkHhXb16SIF02BUSQ4aTGPsjUwpTDs8o4',
    stripePublishableKey:
        'pk_test_51R8r2g2ejpbHZUu9RxWKmyDJTK7amXkB4vE5nRhrd0qvWCnViJsazl9oNjM144gwopnJi1zi3abUk3W4qEk7aWLy00fVUZIeTO',
    userRegistrationAllowlist: ['.*@ccb\\.sh', '.*@paragoncybersecurity\\.sh'],
});

new Environment(app, {
    accountId: '449678530274',
    cloudfrontPublicKey: [
        '-----BEGIN PUBLIC KEY-----',
        'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs3hMvXUvITfcIAh19vAt',
        'GEm7/AUW8rjcDv+nJqw6vFgKwkRvwiKmvKSwRay7sJ1R5UhZBq/AZRtaG6n8RcFs',
        'UuxRxjmfpyCZm6EzPRTEwuCVW2vQR5/urtnGzm3STu4V5ZLJ19yKYZLPle/ZBSMc',
        'ptz2iAYDUOrVb0X964W3wtJ+2GKGRa1QDzgALLWuGQ5yVKrZEhD9Simjg72QuN8E',
        'oCIU+QSVLoLQ21HUS1vAqfQrJYYB5k3P7UNzUom52g8r2fhlVIxuDVZsQojyXdSP',
        '8nDYzoynb61hw49RC3gLhpkY7pkAKgNQp6jROAbTlrN6vLCRC2D4tBt5/t+vSP9k',
        'PwIDAQAB',
        '-----END PUBLIC KEY-----',
    ].join('\n'),
    contactEmailAddress: 'support@paragoncybersecurity.sh',
    domainName: 'cloudsnitch.io',
    pricing: {
        individualSubscriptionStripePriceId: 'price_1RBAJAGsawEQFubmUMPpZQ0a',
        teamSubscriptionStripePriceId: 'price_1RBAJEGsawEQFubmBvPsqPHx',
    },
    regions: ['us-east-1', 'us-west-2'],
    slug: 'prod',
    stripeEventSourceName: 'aws.partner/stripe.com/ed_61SLGZb7XoQdDArer16SIEztBUSQ9uc1GwYb5aX7wIYK',
    stripePublishableKey:
        'pk_live_51R8r2XGsawEQFubmkaZnBAnZcaYT8lGqacq3iMTi0RMpKxbKR4LpDLBerS2WoP266mM2r0KWKFXWPnqgJlpqxUu900PH2XDT5j',
});

new GithubActionsStack(app, 'github-actions-dev', {
    ref: 'refs/heads/main',
    repo: 'ccbrown/cloud-snitch',
    env: {
        account: '774305579662',
        region: 'us-east-1',
    },
    stackName: 'cloud-snitch-github-actions',
});

new GithubActionsStack(app, 'github-actions-prod', {
    ref: 'refs/tags/*',
    repo: 'ccbrown/cloud-snitch',
    env: {
        account: '449678530274',
        region: 'us-east-1',
    },
    stackName: 'cloud-snitch-github-actions',
});
