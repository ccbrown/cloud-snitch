'use client';

import { PencilIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

import { Button, Dialog, ErrorMessage, TextField } from '@/components';
import { useCurrentTeam } from '@/hooks';
import { useDispatch } from '@/store';

interface EditNameFormProps {
    teamId: string;
    currentName: string;
    onSuccess: () => void;
}

const EditNameForm = ({ onSuccess, teamId, currentName }: EditNameFormProps) => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const [name, setName] = useState(currentName);

    const doUpdate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.teams.update({
                teamId,
                input: {
                    name,
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
            <TextField disabled={isBusy} label="Name" required value={name} onChange={setName} />
            <Button disabled={isBusy} label="Rename Team" onClick={doUpdate} type="submit" className="mt-4" />
        </form>
    );
};

const Page = () => {
    const team = useCurrentTeam();
    const [isRenaming, setIsRenaming] = useState(false);

    return (
        <div>
            <Dialog isOpen={isRenaming} onClose={() => setIsRenaming(false)} title="Rename Team">
                <EditNameForm
                    onSuccess={() => {
                        setIsRenaming(false);
                    }}
                    currentName={team?.name || ''}
                    teamId={team?.id || ''}
                />
            </Dialog>
            <h2 className="mb-4 flex items-center gap-2">General Settings</h2>
            <h3 className="label mb-2">Team Name</h3>
            <div className="flex items-center gap-2">
                <span>{team?.name}</span>
                <PencilIcon
                    className="h-[1rem] cursor-pointer hover:text-amethyst"
                    onClick={() => setIsRenaming(true)}
                />
            </div>
        </div>
    );
};

export default Page;
