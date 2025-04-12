import { LinkIcon } from '@heroicons/react/24/outline';
import type { Metadata } from 'next';
import Link from 'next/link';

import { INTEGRATION_TEMPLATE_URL } from '@/integration';

export const metadata: Metadata = {
    title: 'Frequently Asked Questions',
    description: "Learn all about Cloud Snitch, including how it works, who it's by, and how to get started.",
};

const FAQ = [
    {
        question: 'What is Cloud Snitch?',
        answer: (
            <div>
                <p>
                    Cloud Snitch is a web-based tool with one simple goal: to make it easy for you to understand what is
                    happening in your AWS accounts. Inspired by the incredible{' '}
                    <Link
                        href="https://www.obdev.at/products/littlesnitch/index.html"
                        className="external-link"
                        target="_blank"
                    >
                        Little Snitch
                    </Link>{' '}
                    app for macOS, Cloud Snitch gives you an easy to digest, visual representation of the activities of
                    your team, services, and anyone or anything else that may be accessing your accounts.
                </p>
                <p>
                    Whether your goal is diagnostics, intrusion detection, or just plain curiosity, you&apos;re
                    guaranteed to learn something new about your cloud with Cloud Snitch.
                </p>
            </div>
        ),
    },
    {
        question: 'Who is Cloud Snitch made by?',
        answer: (
            <div>
                <p>
                    Cloud Snitch is developed by Paragon Cybersecurity, an LLC owned by tech and security enthusiast{' '}
                    <Link href="https://www.linkedin.com/in/ccbrown1" target="_blank" className="external-link">
                        Chris Brown
                    </Link>
                    . Feel free to connect with me on LinkedIn, or for inquiries about Cloud Snitch,{' '}
                    <Link href="/contact" className="link">
                        contact us
                    </Link>{' '}
                    directly.
                </p>
            </div>
        ),
    },
    {
        question: 'Is it safe to use Cloud Snitch?',
        answer: (
            <div>
                <p>
                    Absolutely! Cloud Snitch requires minimal, read-only access to your AWS account. It will securely
                    ingest CloudTrail entries and use them only for the purposes of providing service to you. You will
                    not need to share any sensitive credentials with us, and we&apos;ll never share your data with
                    anyone else. In fact, we won&apos;t even look at your data unless it is strictly necessary to
                    support you.
                </p>
                <p>
                    If you&apos;d like to take a look at the exact permissions we require, you can read through{' '}
                    <Link href={INTEGRATION_TEMPLATE_URL} target="_blank" className="link">
                        the CloudFormation template
                    </Link>{' '}
                    we provide for setting up integrations. Once you get started, you can also use Cloud Snitch to
                    monitor Cloud Snitch itself!
                </p>
            </div>
        ),
    },
    {
        question: 'Does Cloud Snitch only work with AWS?',
        answer: (
            <div>
                <p>
                    Yes, Cloud Snitch is currently only available for AWS. If you are interested in seeing support for
                    other clouds, please{' '}
                    <Link href="/contact" className="link">
                        contact us
                    </Link>{' '}
                    and let us know!
                </p>
            </div>
        ),
    },
    {
        question: 'How often does Cloud Snitch ingest data from my AWS account?',
        answer: (
            <div>
                <p>
                    Cloud Snitch ingests data from your AWS account daily. Real-time monitoring may be provided in the
                    future, but Cloud Snitch is generally expected to complement other real-time monitoring and alerting
                    solutions by giving you additional diagnostic capabilities and insights that you can use to respond
                    to incidents or proactively identify blind spots in your monitoring.
                </p>
            </div>
        ),
    },
    {
        question: 'Can Cloud Snitch provide SOC 2 Type II or other compliance reports?',
        answer: (
            <div>
                <p>
                    At this time, Cloud Snitch cannot provide SOC 2 Type II, ISO 27001, or other compliance reports.
                    However, we&apos;re happy to talk about our security measures or fill out any questionnaires you may
                    require.{' '}
                    <Link href="/contact" className="link">
                        Contact us
                    </Link>{' '}
                    for more information.
                </p>
            </div>
        ),
    },
    {
        question: "I've found suspicious activity in my AWS account. What should I do?",
        answer: (
            <div>
                <p>
                    If you have reason to believe that your AWS account has been compromised, you should take immediate
                    action by engaging in a cycle of &quot;containment&quot;, &quot;eradication&quot;,
                    &quot;recovery&quot;, and &quot;analysis&quot; activities as outlined by the standard NIST incident
                    response cycle.
                </p>
                <p>This should always begin with contacting AWS support immediately.</p>
                <p>
                    AWS provides a number of resources to help you understand what activities should be included in your
                    incident response cycle. For example, they publish the{' '}
                    <Link
                        href="https://docs.aws.amazon.com/security-ir/latest/userguide/security-incident-response-guide.html"
                        className="link"
                        target="_blank"
                    >
                        AWS Security Incident Response Guide
                    </Link>{' '}
                    and{' '}
                    <Link
                        href="https://github.com/aws-samples/aws-customer-playbook-framework"
                        className="link"
                        target="_blank"
                    >
                        templates for security playbooks in GitHub
                    </Link>
                    . We recommend reviewing these resources before you encounter an incident so you can develop your
                    own playbooks for incident response ahead of time.
                </p>
            </div>
        ),
    },
];

const slug = (question: string) => {
    return question
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-|-$/g, '');
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4 [&_p]:my-4">
            <div className="flex gap-4">
                <div className="translucent-snow p-4 rounded-lg">
                    <h1>Frequently Asked Questions</h1>
                    <p>Here you can find answers to all of the most common questions about Cloud Snitch.</p>
                    <p>
                        Have a question that isn&apos;t answered here?{' '}
                        <Link href="/contact" className="link">
                            Contact us!
                        </Link>
                    </p>
                </div>
                <div className="grow translucent-snow p-4 rounded-lg">
                    <h2>Table of Contents</h2>
                    <ul className="list-disc pl-6 mt-4">
                        {FAQ.map((faq) => (
                            <li key={faq.question}>
                                <a href={`#${slug(faq.question)}`} className="link">
                                    {faq.question}
                                </a>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
            {FAQ.map((faq, index) => (
                <div key={index} className="translucent-snow p-4 rounded-lg">
                    <h2 className="flex gap-2 items-center">
                        <div>{faq.question}</div>
                        <a id={slug(faq.question)} href={`#${slug(faq.question)}`} className="focus:outline-none">
                            <LinkIcon className="h-[1rem]" />
                        </a>
                    </h2>
                    {faq.answer}
                </div>
            ))}
        </div>
    );
};

export default Page;
