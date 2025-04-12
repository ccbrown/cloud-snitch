'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

import { useCurrentUser, useCurrentUserPasskeys, useCurrentUserTeamMemberships, useMostRecentTeamId } from '@/hooks';

import { Header } from '../Header';
import { InitialCredentialsSetup } from './InitialCredentialsSetup';
import { InitialTeamSetup } from './InitialTeamSetup';

const Page = () => {
    const router = useRouter();

    const currentUser = useCurrentUser();
    const currentUserTeamMemberships = useCurrentUserTeamMemberships();
    const currentUserPasskeys = useCurrentUserPasskeys();

    const mostRecentTeamId = useMostRecentTeamId();

    const isLoadingCredentialsCheck = !currentUser || !currentUserPasskeys;
    const showInitialCredentialsSetup =
        !isLoadingCredentialsCheck && !currentUser.hasPassword && Object.keys(currentUserPasskeys).length === 0;

    const isLoadingTeamsCheck = !currentUser || !currentUserTeamMemberships || isLoadingCredentialsCheck;
    const showInitialTeamSetup = !isLoadingTeamsCheck && currentUserTeamMemberships.length === 0;

    const isLoading = isLoadingCredentialsCheck || isLoadingTeamsCheck;

    useEffect(() => {
        if (currentUserTeamMemberships && currentUserTeamMemberships.length > 0) {
            const membership =
                currentUserTeamMemberships.find((membership) => membership.team.id === mostRecentTeamId) ||
                currentUserTeamMemberships[0];
            router.push(`/teams/${membership.team.id}`);
        }
    }, [currentUserTeamMemberships, mostRecentTeamId, router]);

    return (
        <div className="bg-platinum flex flex-col h-screen">
            <Header />
            <main className="grow p-8">
                <div className="translucent-snow max-w-4xl mx-auto rounded-xl p-4">
                    {showInitialCredentialsSetup ? (
                        <InitialCredentialsSetup />
                    ) : showInitialTeamSetup ? (
                        <InitialTeamSetup />
                    ) : isLoading ? (
                        <p>Loading...</p>
                    ) : (
                        <p>Please wait while we redirect you to your team...</p>
                    )}
                </div>
            </main>
        </div>
    );
};

export default Page;
