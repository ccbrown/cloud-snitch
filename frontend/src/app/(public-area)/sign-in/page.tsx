import type { Metadata } from 'next';

import { SignInPage } from './SignInPage';

export const metadata: Metadata = {
    title: 'Sign In',
    description: 'Sign in to your Cloud Snitch account.',
};

const Page = () => {
    return <SignInPage />;
};

export default Page;
