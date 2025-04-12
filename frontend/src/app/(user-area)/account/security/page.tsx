'use client';

import { create, parseCreationOptionsFromJSON } from '@github/webauthn-json/browser-ponyfill';
import { PlusCircleIcon, TrashIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

import { Button, Dialog, ErrorMessage, TextField } from '@/components';
import { useCurrentUser, useCurrentUserPasskeys } from '@/hooks';
import { useDispatch } from '@/store';

interface NewPasskeyFormProps {
    onSuccess: () => void;
}

const NewPasskeyForm = (props: NewPasskeyFormProps) => {
    const currentUser = useCurrentUser();
    const dispatch = useDispatch();

    const [name, setName] = useState('My Passkey');
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doCreate = async () => {
        if (isBusy || !currentUser) {
            return;
        }
        setIsBusy(true);

        try {
            const { sessionId, credentialCreationOptions } = await dispatch.users.beginPasskeyRegistration(
                currentUser.id,
            );
            const credential = await create(parseCreationOptionsFromJSON(credentialCreationOptions));
            await dispatch.users.completePasskeyRegistration({
                userId: currentUser.id,
                registration: {
                    sessionId,
                    passkeyName: name,
                    credentialCreationResponse: credential,
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
            <TextField disabled={isBusy} label="Name" required value={name} onChange={setName} />
            <Button
                disabled={isBusy || !name}
                label="Create Passkey"
                onClick={doCreate}
                type="submit"
                className="mt-4"
            />
        </form>
    );
};

interface NewPasswordFormProps {
    onSuccess: () => void;
}

const NewPasswordForm = (props: NewPasswordFormProps) => {
    const currentUser = useCurrentUser();
    const dispatch = useDispatch();

    const [password, setPassword] = useState('');
    const [passwordConfirmation, setPasswordConfirmation] = useState('');
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const passwordIsValid = password.length >= 14 && password === passwordConfirmation;

    const doUpdate = async () => {
        if (isBusy || !currentUser || !passwordIsValid) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.users.updateUser({
                id: currentUser.id,
                update: {
                    password,
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
            <TextField
                disabled={isBusy}
                label="Password"
                type="password"
                autocomplete="new-password"
                required
                value={password}
                onChange={setPassword}
            />
            <TextField
                disabled={isBusy}
                label="Again"
                type="password"
                autocomplete="new-password"
                required
                value={passwordConfirmation}
                onChange={setPasswordConfirmation}
            />
            <p className="text-sm text-gray-500 mb-2">
                Passwords must be at least 14 characters long. We highly recommend generating a random password using a
                password manager.
            </p>
            <Button
                disabled={isBusy || !passwordIsValid}
                label="Set Password"
                onClick={doUpdate}
                type="submit"
                className="mt-4"
            />
        </form>
    );
};

const Page = () => {
    const dispatch = useDispatch();
    const currentUser = useCurrentUser();
    const passkeys = useCurrentUserPasskeys();
    const [isCreatingPasskey, setIsCreatingPasskey] = useState(false);
    const [isSettingPassword, setIsSettingPassword] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const canDeletePasskey = (passkeys && passkeys.length > 1) || currentUser?.hasPassword;

    const deletePasskey = async (passkeyId: string) => {
        if (!confirm('Are you sure you want to delete this passkey?')) {
            return;
        }

        try {
            await dispatch.users.deletePasskey(passkeyId);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        }
    };

    const deletePassword = async () => {
        try {
            if (currentUser) {
                await dispatch.users.updateUser({
                    id: currentUser.id,
                    update: {
                        password: '',
                    },
                });
            }
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        }
    };

    return (
        <>
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <h2 className="mb-4 flex items-center gap-2">
                Passkeys{' '}
                <PlusCircleIcon
                    className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                    onClick={() => setIsCreatingPasskey(true)}
                />
            </h2>
            <Dialog isOpen={isCreatingPasskey} onClose={() => setIsCreatingPasskey(false)} title="New Passkey">
                <NewPasskeyForm onSuccess={() => setIsCreatingPasskey(false)} />
            </Dialog>
            {!passkeys ? (
                <p>Loading...</p>
            ) : passkeys.length === 0 ? (
                <p>You currently have no passkeys.</p>
            ) : (
                <table className="w-full text-left">
                    <thead className="uppercase text-sm text-english-violet">
                        <tr>
                            <th>Name</th>
                            <th>Created</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {passkeys.map((key) => (
                            <tr key={key.id}>
                                <td>{key.name}</td>
                                <td>{key.creationTime.toLocaleDateString()}</td>
                                <td align="right" className="p-2">
                                    {canDeletePasskey && (
                                        <TrashIcon
                                            className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                                            onClick={() => deletePasskey(key.id)}
                                        />
                                    )}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
            <h2 className="mt-8 mb-4">Password</h2>
            <Dialog isOpen={isSettingPassword} onClose={() => setIsSettingPassword(false)} title="New Passkey">
                <NewPasswordForm onSuccess={() => setIsSettingPassword(false)} />
            </Dialog>
            {!currentUser ? (
                <p>Loading...</p>
            ) : currentUser.hasPassword ? (
                <div className="flex flex-col gap-2">
                    <p>
                        Password authentication is currently <strong>enabled</strong>.
                    </p>
                    <p>
                        We recommend using passkeys exclusively and disabling your password. You can{' '}
                        <span className="link" onClick={() => deletePassword()}>
                            click here
                        </span>{' '}
                        to delete your password, or{' '}
                        <span className="link" onClick={() => setIsSettingPassword(true)}>
                            click here
                        </span>{' '}
                        to change it.
                    </p>
                </div>
            ) : (
                <div className="flex flex-col gap-2">
                    <p>Password authentication is currently disabled. 🎉</p>
                    <p>
                        We recommend you keep it this way and use passkeys exclusively, but if you want to set a
                        password,{' '}
                        <span className="link" onClick={() => setIsSettingPassword(true)}>
                            click here
                        </span>
                        .
                    </p>
                </div>
            )}
        </>
    );
};

export default Page;
