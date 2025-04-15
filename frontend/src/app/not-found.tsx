import Link from 'next/link';

import { MascotBox } from '@/components';

const Page = () => {
    return (
        <div className="bg-dark-purple min-h-screen flex items-center justify-center">
            <div className="lg:min-w-4xl max-w-7xl p-4 mx-auto">
                <MascotBox>
                    <h1>404</h1>
                    <p className="mb-6">You seem to be lost.</p>
                    <Link href="/" className="button">
                        Return Home
                    </Link>
                </MascotBox>
            </div>
        </div>
    );
};

export default Page;
