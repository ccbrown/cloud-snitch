import type { Metadata } from 'next';
import Link from 'next/link';

import { IndividualSubscriptionBox, TeamSubscriptionBox } from '@/components';

export const metadata: Metadata = {
    title: 'Pricing',
    description:
        "Whether you're a hobbyist on a budget or a large enterprise with complex and mature cloud infrastructure, Cloud Snitch has a plan for you.",
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4 [&_p]:my-4">
            <div className="translucent-snow p-4 rounded-lg">
                <h1>Simple and Predictable Pricing</h1>
                <p>
                    Cloud Snitch pricing is based solely on your subscription tier and the number of AWS accounts you
                    have, so your costs are predictable and easy to understand.
                </p>
            </div>
            <div className="flex gap-4">
                <IndividualSubscriptionBox />
                <TeamSubscriptionBox />
                <div className="flex-1 flex flex-col bg-dark-purple outline-1 outline-white/20 text-snow p-4 rounded-lg">
                    <h2>Enterprise</h2>
                    <p>
                        Need something off the menu? Longer retention? Feature prioritization? Custom invoicing? A Cloud
                        Snitch deployment in your cloud?
                    </p>
                    <p>Whatever it is, we&apos;ll work with you!</p>
                    <div className="grow" />
                    <Link href="/contact" className="snow-button text-center text-lg">
                        Contact Us
                    </Link>
                </div>
            </div>
            <div className="translucent-snow p-4 rounded-lg">
                <h2>AWS Costs</h2>
                <p>
                    Cloud Snitch requires you to log CloudTrail events to an S3 bucket in your AWS accounts. Cloud
                    Snitch will guide you through the setup, but the resulting CloudTrail and S3 costs, if any, will be
                    billed to your account. For small accounts, these costs typically fall within the AWS Free Tier,
                    costing nothing. For larger accounts, the costs may be non-zero, but are virtually always a
                    negligible portion of your overall AWS bill.
                </p>
                <p>For more information on AWS costs, see the following:</p>
                <ul className="list-disc pl-6">
                    <li>
                        <Link
                            href="https://aws.amazon.com/cloudtrail/pricing/"
                            target="_blank"
                            className="external-link"
                        >
                            AWS CloudTrail Pricing
                        </Link>
                    </li>
                    <li>
                        <Link href="https://aws.amazon.com/s3/pricing/" target="_blank" className="external-link">
                            Amazon S3 Pricing
                        </Link>
                    </li>
                </ul>
            </div>
            <div className="translucent-snow p-4 rounded-lg">
                <h2>Proration and Timing</h2>
                <p>
                    When data from a new AWS account is first ingested by Cloud Snitch, your payment method will be
                    charged for the remainder of the current billing cycle. Afterwards, you will be charged for the
                    account one month at a time, at the start of each billing cycle.
                </p>
                <p>
                    AWS accounts will be considered inactive and will be removed from your subscription within one week
                    if no new data is ingested for them. You will be granted credit for the remainder of the billing
                    cycle, which will apply to any future payments.
                </p>
                <p>
                    Switching between subscription tiers is prorated in the same way. If you upgrade, you will be
                    charged immediately for the remainder of the billing cycle. If you downgrade, you will be granted
                    credits.
                </p>
                <p>
                    <strong className="uppercase text-english-violet">Example:</strong> You sign up at the begining of
                    the month and create your one-person team using an Individual plan. You then set up an AWS
                    integration and begin ingesting CloudTrail data for 3 AWS accounts. This results in a charge of 3 x
                    $0.99 = $2.97. Halfway through the month, you upgrade to the Team plan, resulting in a charge of 3 x
                    ($9.99 - $0.99) x 50% = $13.50. You also decide to add another account, which results in an
                    additional $9.99 x 50% = $4.98 charge. At the beginning of future billing cycles, you simply pay 4 x
                    $9.99 = $39.96.
                </p>
                <p className="text-xs">*This example excludes any applicable taxes.</p>
            </div>
        </div>
    );
};

export default Page;
