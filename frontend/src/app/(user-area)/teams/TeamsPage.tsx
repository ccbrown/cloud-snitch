'use client';

import { PlusCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { useState } from 'react';

import { Header } from '../Header';
import { CreateTeamForm, formatMembershipRole } from '../dashboard/InitialTeamSetup';
import { TeamInvite } from '@/generated/api';
import { Button, Dialog, ErrorMessage } from '@/components';
import { useCurrentUserTeamInvites, useCurrentUserTeamMemberships } from '@/hooks';
import { useDispatch } from '@/store';

export const TeamsPage = () => {
    const dispatch = useDispatch();
    const [isBusy, setIsBusy] = useState(false);
    const [isCreatingTeam, setIsCreatingTeam] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const memberships = useCurrentUserTeamMemberships();
    const invites = useCurrentUserTeamInvites();

    const acceptInvite = async (teamId: string) => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.join(teamId);
            setIsBusy(false);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    const declineInvite = async (invite: TeamInvite) => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.deleteInvite(invite);
            setIsBusy(false);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <div className="bg-platinum flex flex-col h-screen">
            <Header />
            <main className="grow p-8">
                <Dialog isOpen={isCreatingTeam} onClose={() => setIsCreatingTeam(false)} title="New Team">
                    <CreateTeamForm />
                </Dialog>
                <div className="translucent-snow max-w-4xl mx-auto rounded-xl p-4">
                    <h1 className="flex items-center gap-2">
                        Teams
                        <PlusCircleIcon
                            className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                            onClick={() => setIsCreatingTeam(true)}
                        />
                    </h1>

                    {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}

                    <h2 className="my-4">Memberships</h2>
                    {!memberships ? (
                        <p>Loading...</p>
                    ) : !memberships.length ? (
                        <p>You&apos;re currently not a member of any teams.</p>
                    ) : (
                        <table className="w-full text-left">
                            <thead className="uppercase text-english-violet">
                                <tr>
                                    <th className="py-1">Name</th>
                                    <th>Role</th>
                                </tr>
                            </thead>
                            <tbody>
                                {memberships.map((m) => (
                                    <tr key={m.team.id}>
                                        <td className="py-1">
                                            <Link href={`/teams/${m.team.id}`} className="link">
                                                {m.team.name}
                                            </Link>
                                        </td>
                                        <td>{formatMembershipRole(m.membership.role)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}

                    <h2 className="my-4">Invites</h2>
                    {!invites ? (
                        <p>Loading...</p>
                    ) : !invites.length ? (
                        <p>No pending invites.</p>
                    ) : (
                        <table className="w-full text-left">
                            <thead className="uppercase text-english-violet">
                                <tr>
                                    <th className="py-1">Name</th>
                                    <th>Sender</th>
                                    <th>Role</th>
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {invites.map((inv) => (
                                    <tr key={inv.team.id}>
                                        <td className="py-1">{inv.team.name}</td>
                                        <td>{inv.sender.emailAddress}</td>
                                        <td>{formatMembershipRole(inv.invite.role)}</td>
                                        <td className="px-2 text-right">
                                            <Button
                                                className="text-sm py-1"
                                                disabled={isBusy}
                                                onClick={() => acceptInvite(inv.team.id)}
                                            >
                                                Join
                                            </Button>
                                            <Button
                                                className="text-sm ml-2 py-1"
                                                disabled={isBusy}
                                                style="subtle"
                                                onClick={() => declineInvite(inv.invite)}
                                            >
                                                Decline
                                            </Button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>
            </main>
        </div>
    );
};
