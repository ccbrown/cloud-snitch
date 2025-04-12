import type { Metadata } from 'next';

import { RequireAuth } from './RequireAuth';

export const metadata: Metadata = {
    robots: {
        index: false,
        follow: false,
    },
};

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <RequireAuth>{children}</RequireAuth>;
}
