import { aws_connect as connect, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';

export class PhoneServiceStack extends Stack {
    constructor(scope: Construct, id: string, props: StackProps) {
        super(scope, id, props);

        const instance = new connect.CfnInstance(this, 'Instance', {
            identityManagementType: 'CONNECT_MANAGED',
            instanceAlias: 'cloudsnitch',
            attributes: {
                inboundCalls: true,
                outboundCalls: false,
            },
        });

        new connect.CfnPhoneNumber(this, 'PhoneNumber', {
            countryCode: 'US',
            targetArn: instance.attrArn,
            type: 'DID',
        });

        // The advertised support method is through our website, and we just need a phone number to sign up for Stripe.
        // TODO: Accept calls? Currently callers just hear a busy signal and get disconnected.
    }
}
