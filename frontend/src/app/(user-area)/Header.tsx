'use client';

import { Popover, PopoverButton, PopoverPanel } from '@headlessui/react';
import {
    ArrowRightStartOnRectangleIcon,
    ChevronDownIcon,
    CogIcon,
    ComputerDesktopIcon,
    UserCircleIcon,
    UserGroupIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

import { UserRole } from '@/generated/api';
import { Logo } from '@/components';
import { useCurrentUser } from '@/hooks';
import { useDispatch } from '@/store';

interface Props {
    left?: React.ReactNode;
    center?: React.ReactNode;
}

export const Header = (props: Props) => {
    const router = useRouter();
    const dispatch = useDispatch();
    const currentUser = useCurrentUser();

    const signOut = async () => {
        await dispatch.api.signOut();
        router.push('/');
    };

    return (
        <div className="translucent-snow px-2 pointer-events-auto flex items-center justify-between">
            <div className="flex grow basis-0 items-center">
                <Link href="/" className="font-bold tracking-wide text-2xl flex items-center ml-2 mr-8">
                    <Logo className="h-[2rem] inline mr-2" />
                    Cloud Snitch
                </Link>
                <div>{props.left}</div>
            </div>
            <div className="flex grow items-center justify-center">{props.center}</div>
            <div className="h-[4rem] flex grow basis-0 items-center justify-end mr-2">
                <Popover className="relative">
                    <PopoverButton className="outline-none cursor-pointer flex items-center gap-2">
                        <UserCircleIcon className="h-[2rem]" />
                        <ChevronDownIcon className="h-[1rem]" />
                    </PopoverButton>
                    <PopoverPanel
                        anchor="bottom end"
                        className="flex flex-col translucent-snow pt-2 rounded-lg border-1 border-platinum"
                    >
                        <div className="text-sm p-2 px-4 border-b border-platinum text-english-violet">
                            {currentUser?.emailAddress}
                        </div>
                        <div className="p-2">
                            <Link
                                className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                href="/account"
                            >
                                <CogIcon className="h-[1.5rem]" />
                                <span>Account Settings</span>
                            </Link>
                            <Link
                                className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                href="/teams"
                            >
                                <UserGroupIcon className="h-[1.5rem]" />
                                <span>Teams</span>
                            </Link>
                            {currentUser?.role === UserRole.Administrator && (
                                <Link
                                    className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                    href="/control-panel"
                                >
                                    <ComputerDesktopIcon className="h-[1.5rem]" />
                                    <span>Control Panel</span>
                                </Link>
                            )}
                            <div
                                className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                onClick={() => signOut()}
                            >
                                <ArrowRightStartOnRectangleIcon className="h-[1.5rem]" />
                                <span>Sign Out</span>
                            </div>
                        </div>
                    </PopoverPanel>
                </Popover>
            </div>
        </div>
    );
};
