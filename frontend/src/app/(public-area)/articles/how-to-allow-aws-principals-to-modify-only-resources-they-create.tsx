import { SyntaxHighlighter } from '@/components/SyntaxHighlighter';
import Link from 'next/link';

const arnBasedExample = `{
    "Effect": "Allow",
    "Action": [
        "ecs:CreateCluster",
        "ecs:DeleteCluster"
    ],
    "Resource": [
        "arn:aws:ec2:us-east-1:123412341234:cluster/MyCoolService-*"
    ]
}`;

const ec2Example = `[{
    "Effect": "Allow",
    "Action": "ec2:RunInstances",
    "Resource": "*"
}, {
    "Effect": "Allow",
    "Action": "ec2:CreateTags",
    "Resource": "arn:aws:ec2:us-east-1:123412341234:instance/*",
    "Condition": {
        "StringEquals": {
            "ec2:CreateAction" : "RunInstances"
        }
    }
}, {
    "Effect": "Allow",
    "Action": "ec2:StopInstances",
    "Resource": "arn:aws:ec2:us-east-1:123412341234:instance/*",
    "Condition": {
        "StringEquals": {
            "aws:ResourceTag/ManagedByMyCoolService" : "true"
        }
    }
}]`;

const hiddenTechniqueExample = `{
    "Effect": "Allow",
    "Action": [
        "ec2:RunInstances",
        "ec2:StopInstances",
        "ec2:CreateTags"
    ],
    "Resource": "arn:aws:ec2:us-east-1:123412341234:instance/*",
    "Condition": {
        "StringEquals": {
            "aws:ResourceTag/ManagedByMyCoolService" : "true"
        }
    }
}`;

const hiddenTechniqueCombinedExample = `{
    "Effect": "Allow",
    "Action": [
        "ec2:RunInstances",
        "ec2:StopInstances",
        "ec2:CreateTags",
        "ecs:CreateCluster",
        "ecs:DeleteCluster",
        "ecs:TagResource"
    ],
    "Resource": "*",
    "Condition": {
        "StringEquals": {
            "aws:ResourceTag/ManagedByMyCoolService" : "true"
        }
    }
}`;

const tryIt = [
    '# Create two clusters, one owned by our cool service and one that is not.',
    'aws --profile admin ecs create-cluster --cluster-name uncool-service-cluster',
    `aws --profile my-cool-service ecs create-cluster --cluster-name my-cool-service-cluster --tags 'key=ManagedByMyCoolService,value=true'`,
    '',
    "# Try to delete the uncool service cluster (It won't work)",
    'aws --profile my-cool-service ecs delete-cluster --cluster uncool-service-cluster',
    'An error occurred (AccessDeniedException) when calling the DeleteCluster operation: User: arn:aws:iam::********:user/my-cool-service is not authorized to perform: ecs:DeleteCluster on resource: arn:aws:ecs:us-east-1:********:cluster/uncool-service-cluster because no identity-based policy allows the ecs:DeleteCluster action',
    '',
    '# But we can delete our own cluster',
    'aws --profile my-cool-service ecs delete-cluster --cluster my-cool-service-cluster',
    '',
    '# And we cannot add our tag to the uncool service cluster',
    `aws --profile my-cool-service ecs tag-resource --resource-arn arn:aws:ecs:us-east-1:********:cluster/uncool-service-cluster --tags 'key=ManagedByMyCoolService,value=true'`,
    'An error occurred (AccessDeniedException) when calling the TagResource operation: User: arn:aws:iam::********:user/my-cool-service is not authorized to perform: ecs:TagResource on resource: arn:aws:ecs:us-east-1:********:cluster/uncool-service-cluster because no identity-based policy allows the ecs:TagResource action',
];

const article = {
    title: 'How To Allow AWS Principals To Modify Only Resources They Create',
    author: {
        name: 'Chris',
        image: '/images/chris.jpg',
    },
    description:
        'Use this hidden technique if you want to allow users or services to create and modify resources without giving them access to any pre-existing resources.',
    date: new Date(Date.parse('2025-04-22T12:05:00-04:00')),
    content: (
        <>
            <h2>Why would you want to do that anyway?</h2>
            <p>
                If you&apos;re serious about enforcing least privilege, you&apos;ve probably run into a situation like
                this before:
            </p>
            <p>
                You&apos;ve written an amazing cloud-native service that you want to deploy to AWS, so you begin working
                on your infrastructure-as-code and get to the point where you need to write the IAM policy for your
                service. Your service is going to be deployed to an account that has other services running in it, so
                you and your InfoSec people want to make absolutely sure that it can&apos;t access or interfere with
                resources that don&apos;t belong to it.
            </p>
            <p>
                In particular, your service needs to be able to create and delete ECS clusters, but you want to make
                sure that it can&apos;t delete clusters that don&apos;t belong to it. So you decide on a prefix for your
                cluster names and grant your service permissions like so:
            </p>
            <SyntaxHighlighter language="json">{arnBasedExample}</SyntaxHighlighter>
            <p>Mission accomplished! Everyone is happy.</p>
            <h2>Trouble in Paradise</h2>
            <p>
                You continue writing out your IAM policy and realize that your service also needs to be able to start
                and stop EC2 instances. Again, you don&apos;t want your service to be able to stop any EC2 instances
                that don&apos;t belong to it.
            </p>
            <p>
                Unfortunately, you have no control over the ARN for EC2 instances, so you can&apos;t use the same trick.
                But after some digging, you come up with a clever solution that combines{' '}
                <Link
                    href="https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction_attribute-based-access-control.html"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="external-link"
                >
                    attributed-based access control
                </Link>{' '}
                (ABAC) and{' '}
                <Link
                    href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/supported-iam-actions-tagging.html"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="external-link"
                >
                    <code className="codeinline">ec2:CreateAction</code>
                </Link>
                :
            </p>
            <SyntaxHighlighter language="json">{ec2Example}</SyntaxHighlighter>
            <p>
                You&apos;ve done it again! Your service can only stop instances with your special tag, which it can add
                to new instances, but not to pre-existing instances.
            </p>
            <h2>The Hidden Technique 🥷</h2>
            <p>
                Cloud Snitch recently gained the ability to restrict account activity via{' '}
                <Link
                    href="https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_scps.html"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="external-link"
                >
                    service control policies
                </Link>
                . When implementing this functionality, we wanted to be absolutely sure that when you grant the required
                permissions to Cloud Snitch, you can rest assured knowing that it is literally impossible for Cloud
                Snitch to do anything that would reduce your security posture. That means Cloud Snitch should be able to
                create service control policies, but not modify, detach, or delete policies that it didn&apos;t create
                for you.
            </p>
            <p>
                Unfortunately, there&apos;s no equivalent to <code className="codeinline">ec2:CreateAction</code> for
                service control policies. In fact, the vast majority of resources in AWS don&apos;t have an equivalent
                condition that can be used.
            </p>
            <p>
                However, if we dig deep, we can unlock a secret jutsu that works for any AWS resource that supports
                tagging. The key is in{' '}
                <Link
                    href="https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-keys.html#condition-keys-resourcetag"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="external-link"
                >
                    the documentation for <code className="codeinline">aws:ResourceTag</code>
                </Link>
                :
            </p>
            <blockquote className="quote italic">
                This key is included in the request context when the requested resource already has attached tags or in
                requests that create a resource with an attached tag. This key is returned only for resources that
                support authorization based on tags. There is one context key for each tag key-value pair.
            </blockquote>
            <p>
                Note that the when the resource already exists, <code className="codeinline">aws:ResourceTag</code>{' '}
                refers to the tags already on the resource. However, when the resource is being created,{' '}
                <code className="codeinline">aws:ResourceTag</code> refers to the <em>desired</em> tags for the
                resource.
            </p>
            <p>This means the above example could also be written like this:</p>
            <SyntaxHighlighter language="json">{hiddenTechniqueExample}</SyntaxHighlighter>
            <p>
                This allows your service to manipulate its own instances, while enforcing least privilege, and without
                relying on service-specific conditions. It also has the benefit of being concise!
            </p>
            <p>
                You could even combine multiple services together into a single statement without sacrificing security:
            </p>
            <SyntaxHighlighter language="json">{hiddenTechniqueCombinedExample}</SyntaxHighlighter>
            <h2>Try it Yourself</h2>
            <p>
                Don&apos;t take our word for it. Try it out yourself. Configure your AWS CLI with two sets of
                credentials: one with admin powers and one with the above statement, then...
            </p>
            <SyntaxHighlighter language="bash">{tryIt.join('\n')}</SyntaxHighlighter>
        </>
    ),
    relatedLinks: [
        {
            title: 'Control access using attribute-based access',
            url: 'https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-policies-for-amazon-ec2.html#control-access-with-tags',
        },
        {
            title: 'Grant permission to tag Amazon EC2 resources during creation',
            url: 'https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/supported-iam-actions-tagging.html',
        },
        {
            title: 'Controlling access to AWS resources using tags',
            url: 'https://docs.aws.amazon.com/IAM/latest/UserGuide/access_tags.html',
        },
        {
            title: 'Service control policies (SCPs)',
            url: 'https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_scps.html',
        },
    ],
};

export default article;
