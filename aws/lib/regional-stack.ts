import {
    Aws,
    aws_apigatewayv2 as apigwv2,
    aws_apigatewayv2_integrations as apigwv2_integrations,
    aws_certificatemanager as acm,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
    aws_ecr_assets as ecr_assets,
    aws_events as events,
    aws_events_targets as events_targets,
    aws_iam as iam,
    aws_lambda as lambda,
    aws_lambda_event_sources as eventsources,
    aws_logs as logs,
    aws_route53 as route53,
    aws_s3 as s3,
    aws_s3_deployment as s3deploy,
    aws_ses as ses,
    aws_sqs as sqs,
    aws_route53_targets as route53_targets,
    Duration,
    Environment,
    Stack,
    StackProps,
} from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as crypto from 'crypto';
import * as fs from 'fs';
import * as path from 'path';

import { CopyImageFiles } from './copy-image-files';

interface PricingProps {
    individualSubscriptionStripePriceId: string;
    teamSubscriptionStripePriceId: string;
}

interface Props extends StackProps {
    allRegions: string[];
    cloudfrontKeyGroupId: string;
    cloudfrontPrivateKeySecretName: string;
    cloudfrontPublicKeyId: string;
    contactEmailAddress: string;
    dynamodbTableName: string;
    domainName: string;
    env: Environment & {
        region: string;
        account: string;
    };
    envHash: string;
    globalCertificate: acm.ICertificate;
    passwordEncryptionKeyBase64SecretName: string;
    pricing: PricingProps;
    proxySecretName: string;
    stripeEventSourceName?: string;
    stripePublishableKey: string;
    stripeSecretKeySecretName: string;
    triggerReportGeneration?: boolean;
    userRegistrationAllowlist?: string[];
    noIndex?: boolean;
}

export class RegionalStack extends Stack {
    constructor(scope: Construct, id: string, props: Props) {
        super(scope, id, props);

        // Unfortunately CloudFormation templates must be made available via public S3 bucket:
        // https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cloudformation-stack.html#cfn-cloudformation-stack-templateurl
        const publicS3BucketName =
            'cloudsnitch-public-' +
            crypto
                .createHash('sha256')
                .update(`${props.env.account}:${props.env.region}:nonce-w0RbEQhTF7`)
                .digest('hex')
                .slice(0, 12);

        const publicS3Bucket = new s3.Bucket(this, 'PublicS3Bucket', {
            bucketName: publicS3BucketName,
            blockPublicAccess: {
                blockPublicAcls: false,
                ignorePublicAcls: false,
                blockPublicPolicy: false,
                restrictPublicBuckets: false,
            },
        });
        publicS3Bucket.grantPublicAccess();

        const integrationTemplateUpload = new s3deploy.BucketDeployment(this, 'UploadIntegrationTemplate', {
            sources: [
                s3deploy.Source.data(
                    'integration-v3.cfn.yaml',
                    fs.readFileSync(path.join(__dirname, '../../frontend/public/integration-v3.cfn.yaml'), 'utf8'),
                ),
            ],
            destinationBucket: publicS3Bucket,
            contentType: 'application/yaml',
        });

        const s3ExpirationPadding = Duration.days(1);
        const s3Bucket = new s3.Bucket(this, 'S3Bucket', {
            bucketName: `${props.stackName}-${props.envHash}`,
            cors: [
                {
                    allowedMethods: [s3.HttpMethods.GET],
                    allowedOrigins: ['*'],
                },
            ],
            lifecycleRules: [
                {
                    expiration: Duration.days(7).plus(s3ExpirationPadding),
                    tagFilters: { retention: '1w' },
                },
                {
                    expiration: Duration.days(14).plus(s3ExpirationPadding),
                    tagFilters: { retention: '2w' },
                },
            ],
        });

        // Keep the queue names predictable so that we know what the URLs will be for other regions.
        const queueName = 'CloudSnitchBackend';

        const dlq = new sqs.Queue(this, 'DeadLetterQueue', {
            queueName: `${queueName}-DLQ`,
        });

        const queue = new sqs.Queue(this, 'Queue', {
            deadLetterQueue: {
                maxReceiveCount: 3,
                queue: dlq,
            },
            queueName,
            // Use a visibility timeout just a bit longer than the Lambda timeout of 15 minutes.
            visibilityTimeout: Duration.minutes(17),
        });

        if (props.triggerReportGeneration) {
            const dailyRule = new events.Rule(this, 'DailyRule', {
                schedule: events.Schedule.cron({
                    minute: '30',
                    hour: '0',
                }),
            });
            dailyRule.addTarget(
                new events_targets.SqsQueue(queue, {
                    message: events.RuleTargetInput.fromObject({
                        QueueReportGeneration: {
                            Duration: 24 * 60 * 60 * 1000000000,
                        },
                    }),
                }),
            );

            const morningRule = new events.Rule(this, 'MorningRule', {
                schedule: events.Schedule.cron({
                    minute: '0',
                    hour: '2',
                }),
            });
            morningRule.addTarget(
                new events_targets.SqsQueue(queue, {
                    message: events.RuleTargetInput.fromObject({
                        QueueTeamStripeSubscriptionUpdates: {},
                    }),
                }),
            );

            const afternoonRule = new events.Rule(this, 'AfternoonRule', {
                schedule: events.Schedule.cron({
                    minute: '0',
                    hour: '14',
                }),
            });
            afternoonRule.addTarget(
                new events_targets.SqsQueue(queue, {
                    message: events.RuleTargetInput.fromObject({
                        QueueTeamStripeSubscriptionUpdates: {},
                    }),
                }),
            );
            afternoonRule.addTarget(
                new events_targets.SqsQueue(queue, {
                    message: events.RuleTargetInput.fromObject({
                        QueueTeamEntitlementRefreshes: {},
                    }),
                }),
            );
        }

        const s3BucketDistDomainName = `cdn-${props.env.region}.${props.domainName}`;

        const backendEnvironment: { [key: string]: string } = {
            API_PROXYCOUNT: '1',
            API_PROXYSECRET: `secret:arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.proxySecretName}`,
            APP_PASSWORDENCRYPTIONKEY: `secret:arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.passwordEncryptionKeyBase64SecretName}`,
            APP_FRONTENDURL: `https://${props.domainName}`,
            APP_STORE_DYNAMODB_TABLENAME: props.dynamodbTableName,
            APP_EMAIL_SES_FROMADDRESS: `no-reply@${props.domainName}`,
            APP_CLOUDFRONTPRIVATEKEY: `secret:arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.cloudfrontPrivateKeySecretName}`,
            APP_CLOUDFRONTKEYID: props.cloudfrontPublicKeyId,
            APP_CONTACTEMAILADDRESS: props.contactEmailAddress,
            APP_S3CDNURL: `https://${s3BucketDistDomainName}`,
            APP_S3BUCKETNAME: s3Bucket.bucketName,
            APP_SQSQUEUENAME: queue.queueName,
            APP_STRIPESECRETKEY: `secret:arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.stripeSecretKeySecretName}`,
            APP_PRICING_INDIVIDUALSUBSCRIPTIONSTRIPEPRICEID: props.pricing.individualSubscriptionStripePriceId,
            APP_PRICING_TEAMSUBSCRIPTIONSTRIPEPRICEID: props.pricing.teamSubscriptionStripePriceId,
            APP_AWSACCOUNTID: Aws.ACCOUNT_ID,
            APP_AWSREGIONS: props.allRegions.join(','),
            APP_USERREGISTRATIONALLOWLIST: (props.userRegistrationAllowlist || []).join(','),
        };

        const backendDockerImageAsset = new ecr_assets.DockerImageAsset(this, 'BackendDockerImage', {
            directory: path.join(__dirname, '../../backend'),
            extraHash: props.env.region,
            platform: ecr_assets.Platform.LINUX_ARM64,
        });

        const apiHandler = new lambda.DockerImageFunction(this, 'ApiHandler', {
            architecture: lambda.Architecture.ARM_64,
            code: lambda.DockerImageCode.fromEcr(backendDockerImageAsset.repository, {
                cmd: ['lambda-api-handler'],
                tagOrDigest: backendDockerImageAsset.imageTag,
            }),
            environment: backendEnvironment,
            memorySize: 2048,
            timeout: Duration.seconds(30),
            logRetention: logs.RetentionDays.ONE_MONTH,
        });

        const queueHandler = new lambda.DockerImageFunction(this, 'QueueHandler', {
            architecture: lambda.Architecture.ARM_64,
            code: lambda.DockerImageCode.fromEcr(backendDockerImageAsset.repository, {
                cmd: ['lambda-queue-handler'],
                tagOrDigest: backendDockerImageAsset.imageTag,
            }),
            environment: backendEnvironment,
            memorySize: 2048,
            timeout: Duration.minutes(15),
            logRetention: logs.RetentionDays.ONE_MONTH,
        });
        queueHandler.addEventSource(
            new eventsources.SqsEventSource(queue, {
                batchSize: 1,
            }),
        );

        const handlers = [apiHandler, queueHandler];

        if (props.stripeEventSourceName) {
            const stripeEventBus = new events.EventBus(this, 'StripeEventBus', {
                eventSourceName: props.stripeEventSourceName,
            });

            const dlq = new sqs.Queue(this, 'StripeEventHandlerDLQ', {
                queueName: `CloudSnitchStripeEventHandler-DLQ`,
            });

            const stripeEventHandler = new lambda.DockerImageFunction(this, 'StripeEventHandler', {
                architecture: lambda.Architecture.ARM_64,
                code: lambda.DockerImageCode.fromEcr(backendDockerImageAsset.repository, {
                    cmd: ['lambda-stripe-event-handler'],
                    tagOrDigest: backendDockerImageAsset.imageTag,
                }),
                deadLetterQueue: dlq,
                environment: backendEnvironment,
                memorySize: 2048,
                timeout: Duration.minutes(15),
                logRetention: logs.RetentionDays.ONE_MONTH,
            });
            handlers.push(stripeEventHandler);

            const rule = new events.Rule(this, 'StripeEventRule', {
                eventBus: stripeEventBus,
                eventPattern: {
                    source: events.Match.prefix('aws.partner/stripe.com'),
                },
            });
            rule.addTarget(new events_targets.LambdaFunction(stripeEventHandler));
        }

        for (const handler of handlers) {
            s3Bucket.grantReadWrite(handler);

            // The integration should be able to assume any role in *other* accounts, but must not be
            // allowed to assume roles in this account.
            handler.addToRolePolicy(
                new iam.PolicyStatement({
                    actions: ['sts:AssumeRole'],
                    notResources: [`arn:${Aws.PARTITION}:iam::${Aws.ACCOUNT_ID}:role/*`],
                }),
            );

            handler.addToRolePolicy(
                new iam.PolicyStatement({
                    actions: ['secretsmanager:GetSecretValue'],
                    resources: [
                        `arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.proxySecretName}-*`,
                        `arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.passwordEncryptionKeyBase64SecretName}-*`,
                        `arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.cloudfrontPrivateKeySecretName}-*`,
                        `arn:aws:secretsmanager:${Aws.REGION}:${Aws.ACCOUNT_ID}:secret:${props.stripeSecretKeySecretName}-*`,
                    ],
                }),
            );

            handler.addToRolePolicy(
                new iam.PolicyStatement({
                    actions: ['dynamodb:*'],
                    resources: [
                        `arn:aws:dynamodb:${Aws.REGION}:${Aws.ACCOUNT_ID}:table/${props.dynamodbTableName}`,
                        `arn:aws:dynamodb:${Aws.REGION}:${Aws.ACCOUNT_ID}:table/${props.dynamodbTableName}/*`,
                    ],
                }),
            );

            handler.addToRolePolicy(
                new iam.PolicyStatement({
                    actions: ['sqs:SendMessage'],
                    resources: [
                        // Allow sending messages to the queue in any region.
                        `arn:aws:sqs:*:${Aws.ACCOUNT_ID}:${queue.queueName}`,
                    ],
                }),
            );
        }

        const frontendDockerImageAsset = new ecr_assets.DockerImageAsset(this, 'FrontendDockerImage', {
            buildArgs: {
                NEXT_PUBLIC_AWS_ACCOUNT_ID: props.env.account,
                NEXT_PUBLIC_API_URL: `https://${props.domainName}/api`,
                NEXT_PUBLIC_CDN_URL: `https://${s3BucketDistDomainName}/public/frontend`,
                NEXT_PUBLIC_PUBLIC_S3_BUCKET_NAME: publicS3BucketName,
                NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY: props.stripePublishableKey,
                NEXT_PUBLIC_NO_INDEX: props.noIndex ? 'true' : '',
                OPENAPI_YAML: fs.readFileSync(path.join(__dirname, '../../backend/api/apispec/openapi.yaml'), 'utf8'),
            },
            directory: path.join(__dirname, '../../frontend'),
            extraHash: props.env.region,
            platform: ecr_assets.Platform.LINUX_ARM64,
        });

        const nextJsHandler = new lambda.DockerImageFunction(this, 'NextJsHandler', {
            architecture: lambda.Architecture.ARM_64,
            code: lambda.DockerImageCode.fromEcr(frontendDockerImageAsset.repository, {
                tagOrDigest: frontendDockerImageAsset.imageTag,
            }),
            memorySize: 2048,
            timeout: Duration.seconds(30),
            loggingFormat: lambda.LoggingFormat.JSON,
            logRetention: logs.RetentionDays.ONE_MONTH,
            systemLogLevelV2: lambda.SystemLogLevel.WARN,
        });
        nextJsHandler.node.addDependency(integrationTemplateUpload);

        const copyNextJsPublicFiles = new CopyImageFiles(this, 'CopyNextJsPublicFiles', {
            image: frontendDockerImageAsset,
            directory: '/opt/frontend/public',
            bucket: s3Bucket,
            prefix: 'public/frontend/',
        });
        copyNextJsPublicFiles.executeBefore(nextJsHandler);

        const copyNextJsAssetFiles = new CopyImageFiles(this, 'CopyNextJsAssetFiles', {
            image: frontendDockerImageAsset,
            directory: '/opt/frontend/.next/static',
            bucket: s3Bucket,
            prefix: 'public/frontend/_next/static/',
        });
        copyNextJsAssetFiles.executeBefore(nextJsHandler);

        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.domainName,
        });

        new ses.EmailIdentity(this, 'EmailIdentity', {
            identity: ses.Identity.publicHostedZone(hostedZone),
        });

        apiHandler.addToRolePolicy(
            new iam.PolicyStatement({
                actions: ['ses:SendEmail'],
                resources: ['*'],
                conditions: {
                    StringEquals: {
                        'ses:FromAddress': `no-reply@${props.domainName}`,
                    },
                },
            }),
        );

        const certificate = new acm.Certificate(this, 'Certificate', {
            domainName: `any-region.${props.domainName}`,
            subjectAlternativeNames: [`${props.env.region}.${props.domainName}`],
            validation: acm.CertificateValidation.fromDns(hostedZone),
        });

        const apiGatewayAnyRegionDomainName = new apigwv2.DomainName(this, 'ApiGatewayAnyRegionDomainName', {
            domainName: `any-region.${props.domainName}`,
            certificate,
        });

        const apiGateway = new apigwv2.HttpApi(this, 'ApiGateway', {
            defaultDomainMapping: {
                domainName: apiGatewayAnyRegionDomainName,
            },
            disableExecuteApiEndpoint: true,
        });
        apiGateway.addRoutes({
            path: '/api/{proxy+}',
            methods: [apigwv2.HttpMethod.ANY],
            integration: new apigwv2_integrations.HttpLambdaIntegration('ApiLambdaIntegration', apiHandler, {
                parameterMapping: new apigwv2.ParameterMapping().overwritePath(
                    apigwv2.MappingValue.requestPathParam('proxy'),
                ),
            }),
        });
        apiGateway.addRoutes({
            path: '/{proxy+}',
            methods: [apigwv2.HttpMethod.ANY],
            integration: new apigwv2_integrations.HttpLambdaIntegration('NextJsLambdaIntegration', nextJsHandler),
        });

        const apiGatewayRegionDomainName = new apigwv2.DomainName(this, 'ApiGatewayRegionDomainName', {
            domainName: `${props.env.region}.${props.domainName}`,
            certificate,
        });

        new apigwv2.ApiMapping(this, 'ApiGatewayRegionDomainNameMapping', {
            api: apiGateway,
            domainName: apiGatewayRegionDomainName,
        });

        new route53.ARecord(this, 'RegionAliasRecord', {
            recordName: `${props.env.region}.${props.domainName}`,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.ApiGatewayv2DomainProperties(
                    apiGatewayRegionDomainName.regionalDomainName,
                    apiGatewayRegionDomainName.regionalHostedZoneId,
                ),
            ),
            zone: hostedZone,
        });

        const healthCheck = new route53.HealthCheck(this, 'HealthCheck', {
            type: route53.HealthCheckType.HTTPS,
            fqdn: `${props.env.region}.${props.domainName}`,
            resourcePath: '/api/health-check',
        });

        new route53.ARecord(this, 'AnyRegionAliasRecord', {
            healthCheck,
            recordName: `any-region.${props.domainName}`,
            setIdentifier: props.env.region,
            region: props.env.region,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.ApiGatewayv2DomainProperties(
                    apiGatewayAnyRegionDomainName.regionalDomainName,
                    apiGatewayAnyRegionDomainName.regionalHostedZoneId,
                ),
            ),
            zone: hostedZone,
        });

        const cloudfrontKeyGroup = cloudfront.KeyGroup.fromKeyGroupId(
            this,
            'CloudFrontKeyGroup',
            props.cloudfrontKeyGroupId,
        );
        const s3Dist = new cloudfront.Distribution(this, 'S3Distribution', {
            certificate: props.globalCertificate,
            defaultBehavior: {
                origin: origins.S3BucketOrigin.withOriginAccessControl(s3Bucket),
                trustedKeyGroups: [cloudfrontKeyGroup],
                viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            },
            additionalBehaviors: {
                '/public/*': {
                    origin: origins.S3BucketOrigin.withOriginAccessControl(s3Bucket),
                    viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
                },
            },
            domainNames: [s3BucketDistDomainName],
        });

        new route53.ARecord(this, 'S3DistributionARecord', {
            recordName: s3BucketDistDomainName,
            target: route53.RecordTarget.fromAlias(new route53_targets.CloudFrontTarget(s3Dist)),
            zone: hostedZone,
        });
    }
}
