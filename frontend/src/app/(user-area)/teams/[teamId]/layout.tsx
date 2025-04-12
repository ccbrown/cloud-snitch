import type { Metadata } from 'next';

import { PageTitle } from './PageTitle';

export const metadata: Metadata = {
    title: {
        absolute: '',
    },
};

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <PageTitle>{children}</PageTitle>;
}
