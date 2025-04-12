import type { Metadata } from 'next';

import { CompleteRegistrationPage } from './CompleteRegistrationPage';

export const metadata: Metadata = {
    title: 'Register',
    robots: {
        index: false,
        follow: false,
    },
};

const Page = () => <CompleteRegistrationPage />;

export default Page;
