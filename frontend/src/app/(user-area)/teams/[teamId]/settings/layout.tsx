'use client';

import { Header } from '../Header';
import { TabLayout } from './TabLayout';

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
