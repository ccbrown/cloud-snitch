import type { Metadata } from 'next';

import { TeamsPage } from './TeamsPage';

export const metadata: Metadata = {
    title: 'Teams',
};

const Page = () => <TeamsPage />;

export default Page;
