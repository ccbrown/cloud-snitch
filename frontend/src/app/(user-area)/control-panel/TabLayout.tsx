'use client';

import { UserGroupIcon } from '@heroicons/react/24/outline';

import { TabLayout as TabLayoutImpl } from '@/components';

const TABS = [
    {
        title: 'Teams',
        path: '/control-panel/teams',
        icon: UserGroupIcon,
    },
];

interface Props {
    children?: React.ReactNode;
}

export const TabLayout = (props: Props) => (
    <TabLayoutImpl header={<h1>Control Panel</h1>} tabs={TABS}>
        {props.children}
    </TabLayoutImpl>
);
