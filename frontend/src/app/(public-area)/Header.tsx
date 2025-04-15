'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Bars3Icon } from '@heroicons/react/24/outline';

import { GitHubIcon, Logo } from '@/components';

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

export const Header = () => {
    const [isMenuOpen, setMenuOpen] = useState(false);
    const menuDisplay = isMenuOpen ? 'flex' : 'hidden lg:flex';

    return (
        <div className="relative">
            <div
                className={`${isMenuOpen ? 'absolute' : 'hidden'} top-0 left-0 right-0 h-screen lg:hidden bg-radial from-black/10 to-black/40`}
                onClick={() => setMenuOpen(false)}
            />
            <header className="translucent-snow min-h-[4rem] p-4">
                <div className="max-w-7xl lg:mx-auto lg:h-full flex flex-col lg:flex-row gap-2 lg:items-center lg:justify-between">
                    <div className="flex lg:grow lg:basis-0">
                        <Link href="/" className="grow font-bold tracking-wide text-2xl flex items-center">
                            <Logo className="h-[2rem] inline mr-2" />
                            Cloud Snitch
                        </Link>
                        <Bars3Icon
                            className="lg:hidden h-[2rem] hover:text-amethyst cursor-pointer"
                            onClick={() => setMenuOpen(!isMenuOpen)}
                        />
                    </div>
                    <ul
                        className={`${menuDisplay} flex-col lg:flex-row lg:h-full gap-2 lg:gap-12 lg:items-center lg:justify-center`}
                    >
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
                    <div
                        className={`${menuDisplay} lg:grow lg:basis-0 lg:text-right lg:h-full gap-2 lg:gap-4 items-center justify-between lg:justify-end`}
                    >
                        <Link href="/sign-in" className="text-center grow lg:grow-0 font-semibold hover:text-amethyst">
                            Sign In
                        </Link>
                        <Link href="/register" className="text-center grow lg:grow-0 button">
                            Get Started
                        </Link>
                    </div>
                </div>
            </header>
        </div>
    );
};
