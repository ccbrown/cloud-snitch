import Link from 'next/link';

import { BackgroundAnimation, GitHubIcon, Logo } from '@/components';

const headerLinks = [
    {
        label: '✨ Features',
        href: '/features',
    },
    {
        label: '💳 Pricing',
        href: '/pricing',
    },
    {
        label: '🤔 FAQ',
        href: '/faq',
    },
    {
        label: (
            <span>
                <GitHubIcon className="h-[1.1rem] w-[1.1rem] mr-1 inline-block align-middle" />
                <span className="align-middle">Open Source</span>
            </span>
        ),
        href: 'https://github.com/ccbrown/cloud-snitch',
        target: '_blank',
    },
];

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="bg-platinum flex flex-col w-full min-w-5xl h-screen relative">
            <div className="absolute inset-0">
                <BackgroundAnimation />
            </div>
            <header className="absolute inset-0 translucent-snow h-[4rem] px-4 z-10">
                <div className="max-w-7xl mx-auto h-full flex items-center justify-between">
                    <div className="grow basis-0">
                        <Link href="/" className="font-bold tracking-wide text-2xl flex items-center">
                            <Logo className="h-[2rem] inline mr-2" />
                            Cloud Snitch
                        </Link>
                    </div>
                    <ul className="flex h-full gap-12 items-center justify-center">
                        {headerLinks.map((link, i) => (
                            <li key={i}>
                                <Link
                                    href={link.href}
                                    target={link.target}
                                    className="font-semibold hover:text-amethyst"
                                >
                                    {link.label}
                                </Link>
                            </li>
                        ))}
                    </ul>
                    <div className="grow basis-0 text-right">
                        <Link href="/sign-in" className="font-semibold hover:text-amethyst pr-6">
                            Sign In
                        </Link>
                        <Link href="/register" className="button">
                            Get Started
                        </Link>
                    </div>
                </div>
            </header>
            <div className="absolute inset-0 flex flex-col overflow-auto">
                <div className="px-4 mt-[4rem] grow">
                    <main className="max-w-7xl h-full mx-auto py-8">{children}</main>
                </div>
                <footer className="bg-dark-purple text-snow px-4">
                    <div className="max-w-7xl mx-auto py-8 flex text-sm min-w-5xl">
                        <div className="grow basis-0">&copy; Paragon Cybersecurity, LLC. All Rights Reserved.</div>
                        <div className="flex items-start justify-center font-semibold">
                            <Link href="/contact">Contact Us</Link>
                        </div>
                        <div className="flex grow basis-0 flex-col items-end">
                            <Link href="/privacy-policy">Privacy Policy</Link>
                            <Link href="/cookie-policy">Cookie Policy</Link>
                            <Link href="/terms-of-service">Terms of Service</Link>
                        </div>
                    </div>
                </footer>
            </div>
        </div>
    );
}
