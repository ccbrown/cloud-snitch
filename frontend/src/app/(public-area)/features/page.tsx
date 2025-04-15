import type { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';

export const metadata: Metadata = {
    title: 'Features',
    description:
        'Learn how Cloud Snitch can help you understand your cloud activity, collaborate with teammates, and expose bugs and suspicious behavior.',
};

const Page = () => {
    return (
        <div className="flex flex-col gap-4 [&_p]:my-4">
            <div className="translucent-snow p-4 rounded-lg lg:grid lg:grid-flow-col lg:auto-cols-fr">
                <div className="pr-8">
                    <h1>
                        <span className="text-amethyst-gradient">Explore Activity</span> 🗺️
                    </h1>
                    <p>
                        Cloud Snitch gives you a comprehensive look at your AWS account activity in a sleek and
                        intuitive map view.
                    </p>
                    <p>
                        With Cloud Snitch, there&apos;s no excuse for not knowing everything your AWS accounts are up
                        to.
                    </p>
                </div>
                <div>
                    <Image
                        src="/images/explore.png"
                        alt="Explore Activity"
                        width={1399}
                        height={808}
                        className="rounded-lg border-1 border-platinum"
                    />
                </div>
            </div>
            <div className="translucent-snow p-4 rounded-lg lg:grid lg:grid-flow-col lg:auto-cols-fr">
                <div className="pr-8">
                    <h1>
                        <span className="text-light-coral-gradient">Collaborate With Teammates</span> 🤝
                    </h1>
                    <p>Invite your teammates to let them explore your AWS activity with you.</p>
                    <p>Share links to IP address, CIDR network, and AWS principal activity.</p>
                    <p>Document AWS principals with Markdown notes for your teammates.</p>
                </div>
                <div>
                    <Image
                        src="/images/collaborate.png"
                        alt="Collaborate With Teammates"
                        width={1208}
                        height={800}
                        className="rounded-lg border-1 border-platinum"
                    />
                </div>
            </div>
            <div className="translucent-snow p-4 rounded-lg lg:grid lg:grid-flow-col lg:auto-cols-fr">
                <div className="pr-8">
                    <h1>
                        <span className="text-amethyst-gradient">Expose Bugs and Suspicious Behavior</span> 👾
                    </h1>
                    <p>
                        Cloud Snitch provides summaries of activity by AWS region, principal, IP address, and CIDR
                        network.
                    </p>
                    <p>Errors are highlighted, so you can quickly spot suspicious behavior or bugs in your code.</p>
                    <p>Take the investigation further with links into to your CloudTrail event history.</p>
                </div>
                <div>
                    <Image
                        src="/images/expose.png"
                        alt="Expose Bugs and Suspicious Behavior"
                        width={1208}
                        height={800}
                        className="rounded-lg border-1 border-platinum"
                    />
                </div>
            </div>
            <div className="bg-english-violet text-snow p-4 rounded-lg">
                <h1>Want more?</h1>
                <p>
                    Looking for something specific? We&apos;d love to hear from you! Feel free to{' '}
                    <Link href="/contact" className="underline cursor-pointer">
                        contact us
                    </Link>{' '}
                    and give us your thoughts or feature requests.
                </p>
            </div>
        </div>
    );
};

export default Page;
