import type { Metadata } from 'next';

import { Header } from '../Header';
import { TabLayout } from './TabLayout';

export const metadata: Metadata = {
    title: 'Control Panel',
};

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="bg-platinum flex flex-col min-h-screen">
            <Header />
            <TabLayout>{children}</TabLayout>
        </div>
    );
}
