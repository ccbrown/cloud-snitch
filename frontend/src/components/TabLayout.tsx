'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React from 'react';

interface Tab {
    path: string;
    title: string;
    icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
}

interface Props {
    children?: React.ReactNode;
    header?: React.ReactNode;
    tabs: Tab[];
}

export const TabLayout = (props: Props) => {
    const pathname = usePathname();

    return (
        <div className="translucent-snow w-4xl mx-auto my-8 rounded-xl p-4">
            {props.header}
            <div className="mt-4 flex">
                <div className="flex flex-col gap-2 pr-8 w-1/4">
                    {props.tabs.map((tab) => {
                        const isCurrent = pathname === tab.path;
                        const extraClasses = isCurrent ? 'bg-majorelle-blue/80 text-snow' : 'hover:bg-white/80';
                        return (
                            <Link
                                href={tab.path}
                                key={tab.path}
                                className={`whitespace-nowrap flex items-center gap-2 cursor-pointer p-2 w-full rounded-md ${extraClasses}`}
                            >
                                <tab.icon className="h-[1.5rem]" />
                                {tab.title}
                            </Link>
                        );
                    })}
                </div>
                <main className="w-3/4">{props.children}</main>
            </div>
        </div>
    );
};
