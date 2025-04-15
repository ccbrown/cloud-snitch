import Link from 'next/link';

import { Header } from './Header';
import { BackgroundAnimation } from '@/components';

export default function Layout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="bg-platinum flex flex-col w-full lg:min-w-5xl h-screen relative">
            <div className="absolute inset-0">
                <BackgroundAnimation />
            </div>
            <div className="absolute top-0 left-0 right-0 z-10">
                <Header />
            </div>
            <div className="absolute inset-0 flex flex-col overflow-auto">
                <div className="px-4 mt-[4rem] grow">
                    <main className="max-w-7xl h-full mx-auto py-8">{children}</main>
                </div>
                <footer className="bg-dark-purple text-snow px-4">
                    <div className="max-w-7xl mx-auto py-8 flex flex-col lg:flex-row gap-2 text-sm lg:min-w-5xl">
                        <div className="grow basis-0">&copy; Paragon Cybersecurity, LLC. All Rights Reserved.</div>
                        <div className="flex items-start lg:justify-center font-semibold">
                            <Link href="/contact">Contact Us</Link>
                        </div>
                        <div className="flex grow basis-0 flex-col lg:items-end">
                            <Link href="/privacy-policy">Privacy Policy</Link>
                            <Link href="/cookie-policy">Cookie Policy</Link>
                            <Link href="/terms-of-service">Terms of Service</Link>
                        </div>
                    </div>
                </footer>
            </div>
        </div>
    );
}
