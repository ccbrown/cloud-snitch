import { aws_connect as connect, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';

// This stack sets up an Amazon Connect instance with a phone number. Stripe requires a phone
// number even though it's not actually used for anything. We don't actually provide phone support
// at the moment, but we can use this number for Stripe.
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

        // TODO: Accept calls? Currently callers just hear a busy signal and get disconnected.
    }
}
