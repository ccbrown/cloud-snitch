'use client';

import { CogIcon, LinkIcon, CircleStackIcon, CreditCardIcon, UserGroupIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import React from 'react';

import { TabLayout as TabLayoutImpl } from '@/components';
import { useCurrentTeamId } from '@/hooks';

interface Props {
    children?: React.ReactNode;
}

export const TabLayout = ({ children }: Props) => {
    const teamId = useCurrentTeamId();

    return (
        <TabLayoutImpl
            header={
                <div>
                    <h1>Team</h1>
                    <Link className="link text-sm" href={`/teams/${teamId}`}>
                        Return to Dashboard
                    </Link>
                </div>
            }
            tabs={[
                {
                    title: 'General',
                    path: `/teams/${teamId}/settings`,
                    icon: CogIcon,
                },
                {
                    title: 'Members',
                    path: `/teams/${teamId}/settings/members`,
                    icon: UserGroupIcon,
                },
                {
                    title: 'Integrations',
                    path: `/teams/${teamId}/settings/integrations`,
                    icon: LinkIcon,
                },
                {
                    title: 'Data',
                    path: `/teams/${teamId}/settings/data`,
                    icon: CircleStackIcon,
                },
                {
                    title: 'Billing',
                    path: `/teams/${teamId}/settings/billing`,
                    icon: CreditCardIcon,
                },
            ]}
        >
            {children}
        </TabLayoutImpl>
    );
};
