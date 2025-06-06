import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';

import { Providers } from './providers';

import './globals.css';

const geistSans = Geist({
    variable: '--font-geist-sans',
    subsets: ['latin'],
});

const geistMono = Geist_Mono({
    variable: '--font-geist-mono',
    subsets: ['latin'],
});

export const metadata: Metadata = {
    metadataBase: process.env.NEXT_PUBLIC_CDN_URL ? new URL(process.env.NEXT_PUBLIC_CDN_URL) : undefined,
    openGraph: {
        images: ['/images/opengraph.png'],
    },
    robots: process.env.NEXT_PUBLIC_NO_INDEX
        ? {
              index: false,
              follow: false,
          }
        : undefined,
    title: {
        template: '%s | Cloud Snitch',
        default: 'Cloud Snitch',
    },
};

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body className={`${geistSans.variable} ${geistMono.variable} font-sans antialiased lg:min-w-5xl`}>
                <Providers>{children}</Providers>
            </body>
        </html>
    );
}
