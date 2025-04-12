import type { Metadata } from 'next';

import { CompleteEmailSigninPage } from './CompleteEmailSigninPage';

export const metadata: Metadata = {
    title: 'Sign In',
    robots: {
        index: false,
        follow: false,
    },
};

const Page = () => <CompleteEmailSigninPage />;

export default Page;
