'use client';

import Link from 'next/link';
import { PencilIcon, PlusCircleIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

import { Button, Dialog, ErrorMessage, Select, SuccessMessage, TextField } from '@/components';
import { TeamMembershipRole, TeamTeamMembership } from '@/generated/api';
import { useCurrentTeam, useCurrentTeamId, useCurrentTeamTeamMemberships, useCurrentUser } from '@/hooks';
import { useDispatch } from '@/store';

interface InviteMemberFormProps {
    teamId: string;
    onSuccess: () => void;
}

const InviteMemberForm = (props: InviteMemberFormProps) => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const [emailAddress, setEmailAddress] = useState('');
    const [role, setRole] = useState<TeamMembershipRole>(TeamMembershipRole.Member);

    const doInvite = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.createInvite({
                teamId: props.teamId,
                input: {
                    emailAddress,
                    role,
                },
            });
            props.onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <TextField disabled={isBusy} label="Email Address" value={emailAddress} onChange={setEmailAddress} />
            <Select
                disabled={isBusy}
                label="Role"
                value={role}
                onChange={(value) => setRole(value as TeamMembershipRole)}
                options={[
                    { label: 'Member', value: TeamMembershipRole.Member },
                    { label: 'Admin', value: TeamMembershipRole.Administrator },
                ]}
            />
            <Button disabled={isBusy} label="Invite Team Member" onClick={doInvite} type="submit" className="mt-4" />
        </form>
    );
};

interface EditMemberFormProps {
    teamId: string;
    userId: string;
    current: TeamTeamMembership;
    onSuccess: () => void;
}

const EditMemberForm = ({ teamId, userId, onSuccess, current }: EditMemberFormProps) => {
    const dispatch = useDispatch();
    const currentUserId = useCurrentUser()?.id;

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const [role, setRole] = useState<TeamMembershipRole>(current.membership.role);

    const doDelete = async () => {
        if (
            confirm(
                userId === currentUserId
                    ? 'Are you sure you want to remove yourself from the team?'
                    : 'Are you sure you want to remove this user?',
            )
        ) {
            await dispatch.teams.deleteMembership({
                teamId,
                userId,
            });
            onSuccess();
        }
    };

    const doUpdate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.updateMembership({
                teamId,
                userId,
                input: {
                    role,
                },
            });
            onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <p className="mb-2 text-sm">
                <strong>Member:</strong> {current.user.emailAddress}
            </p>
            <Select
                disabled={isBusy}
                label="Role"
                value={role}
                onChange={(value) => setRole(value as TeamMembershipRole)}
                options={[
                    { label: 'Member', value: TeamMembershipRole.Member },
                    { label: 'Admin', value: TeamMembershipRole.Administrator },
                ]}
            />
            <Button disabled={isBusy} label="Update Membership" onClick={doUpdate} type="submit" className="mt-4" />
            <div className="mt-4 text-sm flex items-center justify-center">
                <span className="text-red-900 cursor-pointer font-semibold" onClick={() => doDelete()}>
                    Remove Member
                </span>
            </div>
        </form>
    );
};

const formatMembershipRole = (role: TeamMembershipRole) => {
    switch (role) {
        case TeamMembershipRole.Member:
            return 'Member';
        case TeamMembershipRole.Administrator:
            return 'Administrator';
        default:
            return 'Unknown';
    }
};

const Page = () => {
    const teamId = useCurrentTeamId();
    const currentUserId = useCurrentUser()?.id;
    const team = useCurrentTeam();
    const [isInviting, setIsInviting] = useState(false);
    const [editUserId, setEditUserId] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState('');

    const memberships = useCurrentTeamTeamMemberships();
    const editTeamTeamMembership = editUserId && memberships?.find((m) => m.user.id === editUserId);
    const hasTeamSubscription = team?.entitlements?.teamFeatures;
    const teamHasOtherAdmins = memberships?.some(
        (m) => m.membership.role === TeamMembershipRole.Administrator && m.user.id !== currentUserId,
    );

    return (
        <div className="flex flex-col gap-4">
            <Dialog isOpen={isInviting} onClose={() => setIsInviting(false)} title="Invite Team Member">
                <InviteMemberForm
                    onSuccess={() => {
                        setSuccessMessage('Invitation sent!');
                        setIsInviting(false);
                    }}
                    teamId={teamId}
                />
            </Dialog>
            <Dialog isOpen={!!editUserId} onClose={() => setEditUserId(null)} title="Update Membership">
                {editUserId && editTeamTeamMembership && (
                    <EditMemberForm
                        onSuccess={() => {
                            setSuccessMessage('Membership updated!');
                            setEditUserId(null);
                        }}
                        teamId={teamId}
                        userId={editUserId}
                        current={editTeamTeamMembership}
                    />
                )}
            </Dialog>
            <h2 className="flex items-center gap-2">
                Members {memberships && `(${memberships.length})`}{' '}
                {hasTeamSubscription && (
                    <PlusCircleIcon
                        className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                        onClick={() => setIsInviting(true)}
                    />
                )}
            </h2>
            {successMessage && <SuccessMessage>{successMessage}</SuccessMessage>}
            {!hasTeamSubscription && (
                <span>
                    <Link href={`/teams/${teamId}/settings/billing`} className="link">
                        Upgrade your subscription
                    </Link>{' '}
                    to invite more team members.
                </span>
            )}
            {!memberships ? (
                <p>Loading...</p>
            ) : (
                <table className="w-full text-left">
                    <thead className="uppercase text-english-violet">
                        <tr>
                            <th>Email Address</th>
                            <th>Role</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {memberships.map((m) => (
                            <tr key={m.user.id}>
                                <td className="py-1">{m.user.emailAddress}</td>
                                <td>{formatMembershipRole(m.membership.role)}</td>
                                <td className="px-2 text-right">
                                    {(m.user.id !== currentUserId || teamHasOtherAdmins) && (
                                        <PencilIcon
                                            className="cursor-pointer hover:text-amethyst"
                                            onClick={() => setEditUserId(m.user.id)}
                                        />
                                    )}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default Page;
