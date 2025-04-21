import { ChevronRightIcon } from '@heroicons/react/24/outline';
import type { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';

import { MascotBox } from '@/components';

export const metadata: Metadata = {
    title: {
        absolute: 'Cloud Snitch: Take your relationship with your cloud to the next level',
    },
    description:
        "How well do you really know your AWS account? Uncover blind spots with Cloud Snitch's sleek and intuitive activity explorer.",
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4">
            <MascotBox>
                <div>
                    <h1>
                        <span className="text-spectrum-gradient text-4xl">
                            How well do you really know your AWS account?
                        </span>
                    </h1>
                    <div className="mt-4 lg:grid lg:grid-flow-col lg:auto-cols-fr">
                        <div className="pr-8 mb-4 lg:mb-0 flex flex-col gap-4">
                            <p>
                                Whether you&apos;re a developer, security engineer, or just a curious person, Cloud
                                Snitch is guaranteed to teach you something and take your relationship with your cloud
                                to the next level.
                            </p>
                            <p>
                                Cloud Snitch provides a sleek and intuitive way of exploring your AWS account activity.
                                It&apos;s a great addition to any toolbox, regardless of if you&apos;re a hobbyist
                                that&apos;s just getting started with the cloud or a large enterprise with complex and
                                mature cloud infrastructure.
                            </p>
                        </div>
                        <Image
                            src="/images/explore.png"
                            alt="Cloud Snitch"
                            width={1399}
                            height={808}
                            className="rounded-lg border-1 border-platinum"
                        />
                    </div>
                </div>
            </MascotBox>
            <div className="flex flex-col lg:flex-row gap-4">
                <div className="flex-1 flex flex-col bg-english-violet text-snow rounded-lg p-4">
                    <h2>
                        <Link href="/features" className="text-snow cursor-pointer">
                            ✨ Features
                        </Link>
                    </h2>
                    <p className="text-sm mt-4">
                        Cloud Snitch aims to do just one thing well: Enlighten you as to what&apos;s happening in your
                        cloud.
                    </p>
                    <p className="text-sm mt-4">Check out our features to learn how we do it.</p>
                    <div className="grow" />
                    <div className="flex mt-8">
                        <Link href="/features" className="snow-button flex grow-0 items-center whitespace-nowrap">
                            Learn More
                            <ChevronRightIcon className="h-4 w-4 ml-2" />
                        </Link>
                    </div>
                </div>
                <div className="flex-1 flex flex-col bg-dark-purple text-snow rounded-lg p-4">
                    <h2>
                        <Link href="/articles/capital-one-data-breach" className="link">
                            😱{' '}
                            <span className="text-amethyst-gradient">Capital One Data Breach – A Cautionary Tale</span>
                        </Link>
                    </h2>
                    <div className="text-sm text-amethyst">
                        What&apos;s in your <s>wallet</s> cloud?
                    </div>
                    <p className="text-sm mt-4">
                        Capital One lost hundreds of millions after being notified by a third party of an intruder that
                        had been lurking in their AWS account for four months.
                    </p>
                    <div className="grow" />
                    <div className="flex mt-8 items-center gap-2">
                        <Link
                            href="/articles/capital-one-data-breach"
                            className="snow-button flex grow-0 items-center whitespace-nowrap"
                        >
                            Read More
                            <ChevronRightIcon className="h-4 w-4 ml-2" />
                        </Link>
                        <Link href="/articles" className="text-sm underline">
                            More Articles
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Page;
