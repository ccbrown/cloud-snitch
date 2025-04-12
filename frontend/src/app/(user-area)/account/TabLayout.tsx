'use client';

import { ShieldCheckIcon } from '@heroicons/react/24/outline';

import { TabLayout as TabLayoutImpl } from '@/components';

const TABS = [
    {
        title: 'Security',
        path: '/account/security',
        icon: ShieldCheckIcon,
    },
];

interface Props {
    children?: React.ReactNode;
}

export const TabLayout = (props: Props) => (
    <TabLayoutImpl header={<h1>Account</h1>} tabs={TABS}>
        {props.children}
    </TabLayoutImpl>
);
