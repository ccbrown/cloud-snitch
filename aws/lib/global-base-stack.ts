import {
    aws_certificatemanager as acm,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
    aws_dynamodb as dynamodb,
    aws_route53 as route53,
    aws_route53_targets as route53_targets,
    aws_secretsmanager as secretsmanager,
    Duration,
    Environment,
    Stack,
    StackProps,
} from 'aws-cdk-lib';
import { Construct } from 'constructs';

interface Props extends StackProps {
    cloudfrontPublicKey: string;
    domainName: string;
    env: Environment & {
        region: string;
    };
    envHash: string;
    envSlug: string;
    regions: string[];
    stackName: string;
}

const newSecret = (stack: Stack, regions: string[], key: string): [secretsmanager.Secret, string] => {
    const secretName = stack.stackName + '/' + key + 'Secret';
    return [
        new secretsmanager.Secret(stack, key + 'Secret', {
            replicaRegions: regions
                .filter((r) => r !== stack.region)
                .map((r) => ({
                    region: r,
                })),
            secretName,
        }),
        secretName,
    ];
};

export class GlobalBaseStack extends Stack {
    certificate: acm.ICertificate;
    cloudfrontDistributionId: string;
    cloudfrontKeyGroupId: string;
    cloudfrontPublicKeyId: string;
    cloudfrontPrivateKeySecretName: string;
    dynamodbTableName: string;
    proxySecretName: string;
    passwordEncryptionKeyBase64SecretName: string;
    stripeSecretKeySecretName: string;

    constructor(scope: Construct, id: string, props: Props) {
        super(scope, id, props);

        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.domainName,
        });

        // CloudFront certificates currently MUST be created in us-east-1:
        // https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/cnames-and-https-requirements.html#https-requirements-aws-region
        //
        // To make things simple, we'll just require that this global stack lives in us-east-1.
        if (props.env.region !== 'us-east-1') {
            throw new Error('Global stack must be deployed in us-east-1');
        }
        const certificate = new acm.Certificate(this, 'Certificate', {
            domainName: props.domainName,
            subjectAlternativeNames: [`*.${props.domainName}`],
            validation: acm.CertificateValidation.fromDns(hostedZone),
        });
        this.certificate = certificate;

        const [proxySecret, proxySecretName] = newSecret(this, props.regions, 'Proxy');
        this.proxySecretName = proxySecretName;

        this.passwordEncryptionKeyBase64SecretName = newSecret(this, props.regions, 'PasswordEncryptionKeyBase64')[1];
        this.cloudfrontPrivateKeySecretName = newSecret(this, props.regions, 'CloudFrontPrivateKey')[1];
        this.stripeSecretKeySecretName = newSecret(this, props.regions, 'StripeSecretKey')[1];

        const dist = new cloudfront.Distribution(this, 'Distribution', {
            certificate,
            defaultBehavior: {
                allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
                cachePolicy: new cloudfront.CachePolicy(this, 'CachePolicy', {
                    defaultTtl: Duration.seconds(0),
                    headerBehavior: cloudfront.CacheHeaderBehavior.allowList(
                        'Access-Control-Request-Headers',
                        'Access-Control-Request-Method',
                        'Authorization',
                        'Accept-Encoding',
                        'Content-Type',
                        'Origin',
                        'RSC',
                        'Next-Router-Prefetch',
                        'Next-Router-State-Tree',
                    ),
                    queryStringBehavior: cloudfront.CacheQueryStringBehavior.all(),
                }),
                origin: new origins.HttpOrigin(`any-region.${props.domainName}`, {
                    customHeaders: {
                        // This is what AWS recommends for authenticating CloudFront with origins:
                        // https://docs.aws.amazon.com/whitepapers/latest/secure-content-delivery-amazon-cloudfront/custom-origin-with-cloudfront.html
                        // TODO: Is it worth using a Lambda@Edge function to add the header more securely?
                        'Proxy-Secret': proxySecret.secretValue.unsafeUnwrap(),
                    },
                }),
                viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            },
            domainNames: [props.domainName],
        });
        this.cloudfrontDistributionId = dist.distributionId;

        const distPublicKey = new cloudfront.PublicKey(this, 'DistPublicKey', {
            encodedKey: props.cloudfrontPublicKey,
        });
        this.cloudfrontPublicKeyId = distPublicKey.publicKeyId;

        const distKeyGroup = new cloudfront.KeyGroup(this, 'DistKeyGroup', {
            items: [distPublicKey],
        });
        this.cloudfrontKeyGroupId = distKeyGroup.keyGroupId;

        new route53.ARecord(this, 'ARecord', {
            recordName: props.domainName,
            target: route53.RecordTarget.fromAlias(new route53_targets.CloudFrontTarget(dist)),
            zone: hostedZone,
        });

        this.dynamodbTableName = props.stackName;
        new dynamodb.TableV2(this, 'DynamoDbTable', {
            partitionKey: {
                name: '_hk',
                type: dynamodb.AttributeType.BINARY,
            },
            sortKey: {
                name: '_rk',
                type: dynamodb.AttributeType.BINARY,
            },
            globalSecondaryIndexes: [
                {
                    indexName: '_bb1',
                    partitionKey: {
                        name: '_bb1h',
                        type: dynamodb.AttributeType.BINARY,
                    },
                    sortKey: {
                        name: '_bb1r',
                        type: dynamodb.AttributeType.BINARY,
                    },
                },
                {
                    indexName: '_bb2',
                    partitionKey: {
                        name: '_bb2h',
                        type: dynamodb.AttributeType.BINARY,
                    },
                    sortKey: {
                        name: '_bb2r',
                        type: dynamodb.AttributeType.BINARY,
                    },
                },
            ],
            billing: dynamodb.Billing.onDemand({
                maxReadRequestUnits: 100,
                maxWriteRequestUnits: 100,
            }),
            replicas: props.regions
                .filter((region) => region !== props.env.region)
                .map((region) => ({
                    region,
                })),
            tableName: this.dynamodbTableName,
            timeToLiveAttribute: '_ttl',
        });
    }
}
