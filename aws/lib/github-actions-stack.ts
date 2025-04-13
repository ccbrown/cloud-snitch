import { aws_iam as iam, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';

interface Props extends StackProps {
    ref: string;
    repo: string;
}

export class GithubActionsStack extends Stack {
    constructor(scope: Construct, id: string, props: Props) {
        super(scope, id, props);

        const oidcProvider = new iam.OpenIdConnectProvider(this, 'OidcProvider', {
            url: 'https://token.actions.githubusercontent.com',
            clientIds: ['sts.amazonaws.com'],
            thumbprints: ['ffffffffffffffffffffffffffffffffffffffff'],
        });

        new iam.Role(this, 'GithubActionsRole', {
            assumedBy: new iam.WebIdentityPrincipal(oidcProvider.openIdConnectProviderArn, {
                StringLike: {
                    'token.actions.githubusercontent.com:aud': 'sts.amazonaws.com',
                    'token.actions.githubusercontent.com:sub': `repo:${props.repo}:ref:${props.ref}`,
                },
            }),
            managedPolicies: [iam.ManagedPolicy.fromAwsManagedPolicyName('AdministratorAccess')],
        });
    }
}
