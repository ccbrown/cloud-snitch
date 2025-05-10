import { SyntaxHighlighter } from '@/components/SyntaxHighlighter';

const regionListExample = [
    'aws ssm get-parameters-by-path --path /aws/service/global-infrastructure/regions',
    '{',
    '    "Parameters": [',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/ap-northeast-1",',
    '            "Type": "String",',
    '            "Value": "ap-northeast-1",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-04-08T17:37:38.637000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/ap-northeast-1",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/ap-southeast-5",',
    '            "Type": "String",',
    '            "Value": "ap-southeast-5",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2024-08-21T13:08:58.686000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/ap-southeast-5",',
    '            "DataType": "text"',
    '        },',
    '        ...',
    '    ]',
    '}',
];

const regionExample = [
    'aws ssm get-parameters-by-path --path /aws/service/global-infrastructure/regions/us-east-1',
    '{',
    '    "Parameters": [',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/us-east-1/domain",',
    '            "Type": "String",',
    '            "Value": "amazonaws.com",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-06-21T08:15:34.835000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/us-east-1/domain",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/us-east-1/geolocationCountry",',
    '            "Type": "String",',
    '            "Value": "US",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-04-08T17:37:50.856000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/us-east-1/geolocationCountry",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/us-east-1/geolocationRegion",',
    '            "Type": "String",',
    '            "Value": "US-VA",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-04-08T17:37:51.403000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/us-east-1/geolocationRegion",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/us-east-1/longName",',
    '            "Type": "String",',
    '            "Value": "US East (N. Virginia)",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-04-08T17:37:51.477000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/us-east-1/longName",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/regions/us-east-1/partition",',
    '            "Type": "String",',
    '            "Value": "aws",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2019-04-08T17:37:51.541000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/regions/us-east-1/partition",',
    '            "DataType": "text"',
    '        }',
    '    ]',
    '}',
];

const serviceListExample = [
    'aws ssm get-parameters-by-path --path /aws/service/global-infrastructure/services',
    '{',
    '    "Parameters": [',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/services/arc-zonal-shift",',
    '            "Type": "String",',
    '            "Value": "arc-zonal-shift",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2022-12-13T14:41:37.029000-05:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/services/arc-zonal-shift",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/services/codestar-connections",',
    '            "Type": "String",',
    '            "Value": "codestar-connections",',
    '            "Version": 1003,',
    '            "LastModifiedDate": "2020-03-17T10:38:01.897000-04:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/services/codestar-connections",',
    '            "DataType": "text"',
    '        },',
    '        ...',
    '    ]',
    '}',
];

const serviceExample = [
    'aws ssm get-parameters-by-path --path /aws/service/global-infrastructure/services/ec2',
    '{',
    '    "Parameters": [',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/services/ec2/longName",',
    '            "Type": "String",',
    '            "Value": "Amazon Elastic Compute Cloud (EC2)",',
    '            "Version": 2,',
    '            "LastModifiedDate": "2020-02-05T05:41:45.323000-05:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/services/ec2/longName",',
    '            "DataType": "text"',
    '        },',
    '        {',
    '            "Name": "/aws/service/global-infrastructure/services/ec2/marketingHomeURL",',
    '            "Type": "String",',
    '            "Value": "https://aws.amazon.com/ec2/",',
    '            "Version": 1,',
    '            "LastModifiedDate": "2020-02-07T15:32:06.856000-05:00",',
    '            "ARN": "arn:aws:ssm:us-east-1::parameter/aws/service/global-infrastructure/services/ec2/marketingHomeURL",',
    '            "DataType": "text"',
    '        }',
    '    ]',
    '}',
];

const accessReportExample = [
    'aws iam generate-organizations-access-report --entity-path o-myorgid/r-myrootid/123412341234',
    '{',
    '    "JobId": "5acb5bc4-ae86-eef1-4553-9ae3564ab6d2"',
    '}',
    '',
    'aws iam get-organizations-access-report --job-id 5acb5bc4-ae86-eef1-4553-9ae3564ab6d2',
    '{',
    '    "JobStatus": "COMPLETED",',
    '    "JobCreationDate": "2025-05-10T00:06:20.450000+00:00",',
    '    "JobCompletionDate": "2025-05-10T00:06:25.689000+00:00",',
    '    "NumberOfServicesAccessible": 11,',
    '    "NumberOfServicesNotAccessed": 0,',
    '    "AccessDetails": [',
    '        {',
    '            "ServiceName": "Amazon EC2 Auto Scaling",',
    '            "ServiceNamespace": "autoscaling",',
    '            "Region": "us-east-1",',
    '            "EntityPath": "o-myorgid/r-myrootid/123412341234",',
    '            "LastAuthenticatedTime": "2025-04-20T05:57:49+00:00",',
    '            "TotalAuthenticatedEntities": 1',
    '        },',
    '        {',
    '            "ServiceName": "AWS CloudTrail",',
    '            "ServiceNamespace": "cloudtrail",',
    '            "Region": "us-east-1",',
    '            "EntityPath": "o-myorgid/r-myrootid/123412341234",',
    '            "LastAuthenticatedTime": "2025-03-05T20:19:31+00:00",',
    '            "TotalAuthenticatedEntities": 1',
    '        },',
    '        ...',
    '    ],',
    '}',
];

const article = {
    title: 'How To Programmatically Get a List of All AWS Regions and Services',
    author: {
        name: 'Chris',
        image: '/images/chris.jpg',
    },
    description:
        'Here are some reliable, but little-known ways to programmatically get a list of all AWS regions and services.',
    date: new Date(Date.parse('2025-05-10T12:05:00-04:00')),
    content: (
        <>
            <p>
                There may be times when you need to get a list of all available AWS regions or services. Cloud Snitch
                for example needs an up-to-date list of regions to be able to plot them on a map. There are repositories
                on GitHub where developers have attempted to maintain such lists, but inevitably they all become
                out-dated.
            </p>
            <p>
                Fortunately, there&apos;s a way to get what we need directly from AWS. This makes it far more practical
                to maintain up-to-date databases of regions and services.
            </p>
            <h2>SSM Parameters for AWS Global Infrastructure</h2>
            <p>
                AWS publishes information about their global infrastructure via public SSM parameters under the{' '}
                <code className="codeinline">/aws/service/global-infrastructure</code> path. These parameters include
                all available regions, services, and more.
            </p>
            <h2>Listing AWS Regions</h2>
            <p>To get a list of all available AWS regions, you can use a command like...</p>
            <SyntaxHighlighter language="bash">{regionListExample.join('\n')}</SyntaxHighlighter>
            <p>
                These parameters are hierarchical, so you can drill down to get more information about a specific region
                like so:
            </p>
            <SyntaxHighlighter language="bash">{regionExample.join('\n')}</SyntaxHighlighter>
            <h2>Listing AWS Services</h2>
            <p>Similarly, to get a list of all available AWS services, you can use a command like...</p>
            <SyntaxHighlighter language="bash">{serviceListExample.join('\n')}</SyntaxHighlighter>
            <p>And you can drill down further like so:</p>
            <SyntaxHighlighter language="bash">{serviceExample.join('\n')}</SyntaxHighlighter>
            <h2>Listing AWS Services Via Organizations Access Report</h2>
            <p>
                Based on your security configuration or private beta enrollments, your account may have more or fewer
                services available to it. As an alternative to using the global infrastructure SSM parameters, you can
                enumerate services available to your account by generating an organizations access report:
            </p>
            <SyntaxHighlighter language="bash">{accessReportExample.join('\n')}</SyntaxHighlighter>
            <p>
                If your account is locked down using service control policies, your report will only contain the
                services you have enabled for it. For example, the report shown above only has 11 services listed,
                because it&apos;s been locked down using Cloud Snitch. For an unrestricted account, your report would
                contain several hundred services.
            </p>
        </>
    ),
    relatedLinks: [
        {
            title: 'Calling public parameters for AWS services, Regions, endpoints, Availability Zones, local zones, and Wavelength Zones in Parameter Store',
            url: 'https://docs.aws.amazon.com/systems-manager/latest/userguide/parameter-store-public-parameters-global-infrastructure.html',
        },
        {
            title: 'GenerateOrganizationsAccessReport',
            url: 'https://docs.aws.amazon.com/IAM/latest/APIReference/API_GenerateOrganizationsAccessReport.html',
        },
    ],
};

export default article;
