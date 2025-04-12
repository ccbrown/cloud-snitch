import { useRouter } from 'next/navigation';
import React, { useState } from 'react';

import { Button, ErrorMessage, TextField } from '@/components';
import { TeamMembershipRole } from '@/generated/api';
import { useCurrentUserTeamInvites } from '@/hooks';
import { useDispatch } from '@/store';

export const formatMembershipRole = (role: TeamMembershipRole) => {
    switch (role) {
        case TeamMembershipRole.Member:
            return 'Member';
        case TeamMembershipRole.Administrator:
            return 'Administrator';
        default:
            return 'Unknown';
    }
};

const JoinTeamForm = () => {
    const dispatch = useDispatch();
    const router = useRouter();
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const invites = useCurrentUserTeamInvites();

    const doJoin = async (teamId: string) => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.join(teamId);
            router.push(`/teams/${teamId}`);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <>
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <p>
                {!invites
                    ? 'Loading...'
                    : !invites.length
                      ? "You currently have no team invites. If you'd like to join an existing team, ask a team administrator to invite you."
                      : `You have ${invites.length} team invite${invites.length === 1 ? '' : 's'}:`}
            </p>
            {invites?.length && (
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
                                        onClick={() => doJoin(inv.team.id)}
                                    >
                                        Join
                                    </Button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </>
    );
};

export const CreateTeamForm = () => {
    const dispatch = useDispatch();
    const router = useRouter();
    const [name, setName] = useState('My Team');
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doCreate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            const team = await dispatch.teams.create({
                name,
            });
            router.push(`/teams/${team.id}`);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <TextField disabled={isBusy} label="Name" required value={name} onChange={setName} />
            <Button disabled={isBusy || !name} label="Create Team" onClick={doCreate} type="submit" className="mt-4" />
        </form>
    );
};

export const InitialTeamSetup = () => {
    const [isCreating, setIsCreating] = useState(false);

    return (
        <div className="max-w-md mx-auto m-8 flex flex-col gap-4">
            <h1 className="mb-4">Create or Join a Team</h1>
            <p>To get started, you&apos;ll need to create or join a team.</p>
            {isCreating ? <CreateTeamForm /> : <JoinTeamForm />}
            <p className="cursor-pointer text-majorelle-blue text-center" onClick={() => setIsCreating(!isCreating)}>
                {isCreating ? 'Join an existing team' : 'Create a new team'}
            </p>
        </div>
    );
};
