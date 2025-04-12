import { Stack, StackProps, custom_resources as cr } from 'aws-cdk-lib';
import { Construct } from 'constructs';

interface Props extends StackProps {
    cloudfrontDistributionId: string;
}

export class GlobalApexStack extends Stack {
    constructor(scope: Construct, id: string, props: Props) {
        super(scope, id, props);

        new cr.AwsCustomResource(this, `CloudFrontInvalidation-${Date.now()}`, {
            onCreate: {
                physicalResourceId: cr.PhysicalResourceId.of(`${props.cloudfrontDistributionId}-${Date.now()}`),
                service: 'CloudFront',
                action: 'createInvalidation',
                parameters: {
                    DistributionId: props.cloudfrontDistributionId,
                    InvalidationBatch: {
                        CallerReference: Date.now().toString(),
                        Paths: {
                            Quantity: 1,
                            Items: ['/*'],
                        },
                    },
                },
            },
            policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
        });
    }
}
