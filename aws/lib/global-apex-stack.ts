import { aws_cloudwatch as cw, Duration, Stack, StackProps, custom_resources as cr } from 'aws-cdk-lib';
import { Construct } from 'constructs';

interface Props extends StackProps {
    allRegions: string[];
    cloudfrontDistributionId: string;
    envSlug: string;
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

        new Dashboard(this, 'Dashboard', {
            allRegions: props.allRegions,
            envSlug: props.envSlug,
        });
    }
}

interface DashboardProps {
    allRegions: string[];
    envSlug: string;
}

class Dashboard extends Construct {
    constructor(scope: Construct, id: string, props: DashboardProps) {
        super(scope, id);

        const dashboard = new cw.Dashboard(this, 'Dashboard', {
            dashboardName: `cloud-snitch-${props.envSlug}`,
            defaultInterval: Duration.days(7),
        });

        const apiGatewayRequestsGraph = new cw.GraphWidget({
            height: 6,
            width: 12,
            title: 'API Requests',
            leftYAxis: {
                showUnits: false,
            },
            rightYAxis: {
                label: 'Errors',
                showUnits: false,
            },
        });
        for (const region of props.allRegions) {
            apiGatewayRequestsGraph.addLeftMetric(
                new cw.MathExpression({
                    expression: `SEARCH('{AWS/ApiGateway, ApiId} MetricName="Count"', 'Sum', 300)`,
                    label: region,
                    period: Duration.minutes(5),
                    searchRegion: region,
                    usingMetrics: {},
                }),
            );
        }

        const apiGatewayErrorsGraph = new cw.GraphWidget({
            height: 6,
            width: 12,
            title: 'API Errors',
            leftYAxis: {
                showUnits: false,
            },
        });
        for (const region of props.allRegions) {
            for (const metric of ['4xx', '5xx']) {
                apiGatewayErrorsGraph.addLeftMetric(
                    new cw.MathExpression({
                        expression: `SEARCH('{AWS/ApiGateway, ApiId} MetricName="${metric}"', 'Sum', 300)`,
                        label: region,
                        period: Duration.minutes(5),
                        searchRegion: region,
                        usingMetrics: {},
                    }),
                );
            }
        }

        dashboard.addWidgets(apiGatewayRequestsGraph, apiGatewayErrorsGraph);

        const lambdaInvocationsGraph = new cw.GraphWidget({
            height: 6,
            width: 12,
            title: 'Lambda Invocations',
            leftYAxis: {
                showUnits: false,
            },
        });
        for (const region of props.allRegions) {
            lambdaInvocationsGraph.addLeftMetric(
                new cw.MathExpression({
                    expression: `SEARCH('{AWS/Lambda, FunctionName} MetricName="Invocations"', 'Sum', 300)`,
                    label: region,
                    period: Duration.minutes(5),
                    searchRegion: region,
                    usingMetrics: {},
                }),
            );
        }

        const maxLambdaDurationGraph = new cw.GraphWidget({
            height: 6,
            width: 12,
            title: 'Maximum Lambda Duration',
            leftYAxis: {
                label: 'Seconds',
                showUnits: false,
            },
        });
        for (const region of props.allRegions) {
            maxLambdaDurationGraph.addLeftMetric(
                new cw.MathExpression({
                    expression: `SEARCH('{AWS/Lambda, FunctionName} MetricName="Duration"', 'Maximum', 300) / 1000`,
                    label: region,
                    period: Duration.minutes(5),
                    searchRegion: region,
                    usingMetrics: {},
                }),
            );
        }

        dashboard.addWidgets(lambdaInvocationsGraph, maxLambdaDurationGraph);

        const sqsOldestMessageAgeGraph = new cw.GraphWidget({
            height: 6,
            width: 12,
            title: 'SQS Oldest Message Age',
            leftYAxis: {
                label: 'Seconds',
                showUnits: false,
            },
        });
        for (const region of props.allRegions) {
            sqsOldestMessageAgeGraph.addLeftMetric(
                new cw.MathExpression({
                    expression: `SEARCH('{AWS/SQS, QueueName} MetricName="ApproximateAgeOfOldestMessage"', 'Maximum', 300)`,
                    label: region,
                    period: Duration.minutes(5),
                    searchRegion: region,
                    usingMetrics: {},
                }),
            );
        }

        dashboard.addWidgets(sqsOldestMessageAgeGraph);
    }
}
