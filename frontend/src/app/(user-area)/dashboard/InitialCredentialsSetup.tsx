'use client';

import React, { useState } from 'react';

import { create, parseCreationOptionsFromJSON } from '@github/webauthn-json/browser-ponyfill';

import { Button, ErrorMessage, TextField } from '@/components';
import { useCurrentUser } from '@/hooks';
import { useDispatch } from '@/store';

const PasskeyForm = () => {
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
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <p className="my-2">Passkeys are the recommended authentication mechanism.</p>
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

const PasswordForm = () => {
    const currentUser = useCurrentUser();
    const dispatch = useDispatch();

    const [password, setPassword] = useState('');
    const [passwordConfirmation, setPasswordConfirmation] = useState('');
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const passwordIsValid = password.length >= 14 && password === passwordConfirmation;

    const doSetPassword = async () => {
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
                onClick={doSetPassword}
                label="Set Password"
                type="submit"
                className="mt-6"
            />
        </form>
    );
};

export const InitialCredentialsSetup = () => {
    const [credentialsType, setCredentialsType] = useState('passkey');
    const otherCredentialsType = credentialsType === 'passkey' ? 'password' : 'passkey';

    return (
        <div className="max-w-md mx-auto m-8 flex flex-col gap-2">
            <h1 className="mb-4">Authentication</h1>
            <p>
                Before doing anything else, you should configure your authentication mechanism. This is how you&apos;ll
                sign back into Cloud Snitch in the future.
            </p>
            {credentialsType === 'passkey' ? <PasskeyForm /> : <PasswordForm />}
            <p
                className="cursor-pointer p-2 text-majorelle-blue text-center text-sm"
                onClick={() => setCredentialsType(otherCredentialsType)}
            >
                Use a {otherCredentialsType} instead
            </p>
        </div>
    );
};
