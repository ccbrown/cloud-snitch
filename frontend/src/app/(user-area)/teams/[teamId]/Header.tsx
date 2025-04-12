'use client';

import { Popover, PopoverButton, PopoverPanel } from '@headlessui/react';
import { ChevronDownIcon, CogIcon, PlusIcon, UserGroupIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { useState } from 'react';

import { Header as UserAreaHeader } from '../../Header';
import { CreateTeamForm } from '../../dashboard/InitialTeamSetup';
import { Dialog } from '@/components';
import { TeamMembershipRole } from '@/generated/api';
import { useCurrentTeam, useCurrentUserTeamMemberships } from '@/hooks';

interface Props {
    children?: React.ReactNode;
}

export const Header = (props: Props) => {
    const team = useCurrentTeam();
    const teamId = team?.id;
    const [isCreatingTeam, setIsCreatingTeam] = useState(false);

    const memberships = useCurrentUserTeamMemberships();

    const isTeamAdmin = memberships?.some(
        (m) => m.team.id === teamId && m.membership.role === TeamMembershipRole.Administrator,
    );

    return (
        <UserAreaHeader
            left={
                <>
                    <Dialog isOpen={isCreatingTeam} onClose={() => setIsCreatingTeam(false)} title="New Team">
                        <CreateTeamForm />
                    </Dialog>
                    <Popover className="relative">
                        <PopoverButton className="outline-none cursor-pointer flex items-center gap-2">
                            {team?.name}
                            <ChevronDownIcon className="h-[1rem]" />
                        </PopoverButton>
                        <PopoverPanel
                            anchor="bottom start"
                            className="flex flex-col translucent-snow rounded-lg border-1 border-platinum"
                        >
                            {isTeamAdmin && (
                                <div className="p-2 border-b border-platinum">
                                    <Link
                                        className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                        href={`/teams/${teamId}/settings`}
                                    >
                                        <CogIcon className="h-[1.5rem]" />
                                        <span>Team Settings</span>
                                    </Link>
                                </div>
                            )}
                            {memberships && memberships.length > 1 && (
                                <div className="p-2 border-b border-platinum">
                                    <div className="p-2 font-semibold uppercase text-sm text-english-violet">
                                        Switch Team
                                    </div>
                                    {memberships
                                        .filter((membership) => membership.team.id !== teamId)
                                        .map((membership) => (
                                            <Link
                                                key={membership.team.id}
                                                className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                                href={`/teams/${membership.team.id}`}
                                            >
                                                <UserGroupIcon className="h-[1.5rem]" />
                                                <span>{membership.team.name}</span>
                                            </Link>
                                        ))}
                                </div>
                            )}
                            <div className="p-2">
                                <div
                                    className="whitespace-nowrap flex items-center gap-2 cursor-pointer hover:bg-white/80 p-2 rounded-md"
                                    onClick={() => setIsCreatingTeam(true)}
                                >
                                    <PlusIcon className="h-[1.5rem]" />
                                    <span>Create New Team</span>
                                </div>
                            </div>
                        </PopoverPanel>
                    </Popover>
                </>
            }
            center={props.children}
        />
    );
};
