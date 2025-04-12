'use client';

import { get, parseRequestOptionsFromJSON } from '@github/webauthn-json/browser-ponyfill';
import { KeyIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

import { Button, Dialog, ErrorMessage, SuccessMessage, TextField } from '@/components';
import { useDispatch, useSelector } from '@/store';

const AccountRecoveryForm = () => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [success, setSuccess] = useState(false);

    const [emailAddress, setEmailAddress] = useState('');

    const doRecovery = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);
        setErrorMessage('');

        try {
            await dispatch.users.beginEmailAuthentication({
                emailAddress,
            });
            setSuccess(true);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return success ? (
        <SuccessMessage>Please check your email for your temporary one-time sign-in link.</SuccessMessage>
    ) : (
        <form className="flex flex-col gap-4">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <p>
                Enter your email address below. If you have an account with us, we&apos;ll send you a temporary one-time
                link that you can use to sign back in.
            </p>
            <TextField
                disabled={isBusy}
                label="Email Address"
                required
                value={emailAddress}
                onChange={setEmailAddress}
            />
            <Button disabled={isBusy || !emailAddress} label="Recover Account" onClick={doRecovery} type="submit" />
        </form>
    );
};

export const SignInPage = () => {
    const isSignedIn = useSelector((state) => !!state.api.auth);
    const router = useRouter();
    const dispatch = useDispatch();

    const [isRedirecting, setIsRedirecting] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [emailAddress, setEmailAddress] = useState('');
    const [password, setPassword] = useState('');
    const [isBusy, setIsBusy] = useState(false);
    const [isRecovering, setIsRecovering] = useState(false);

    useEffect(() => {
        if (isSignedIn) {
            setIsRedirecting(true);
            router.push('/dashboard');
        }
    }, [isSignedIn, router]);

    const signInWithPasskey = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            const { sessionId, credentialAssertionOptions } = await dispatch.users.beginPasskeyAuthentication();
            const credential = await get(parseRequestOptionsFromJSON(credentialAssertionOptions));
            await dispatch.api.signIn({
                sessionId,
                credentialAssertionResponse: credential,
            });
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        } finally {
            setIsBusy(false);
        }
    };

    const signInWithPassword = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.api.signIn({
                emailAddress,
                password,
            });
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        } finally {
            setIsBusy(false);
        }
    };

    return (
        <div className="translucent-snow max-w-4xl mx-auto rounded-xl p-4">
            {isRedirecting ? (
                <p>Signing in...</p>
            ) : (
                <>
                    <Dialog isOpen={isRecovering} onClose={() => setIsRecovering(false)} title="Account Recovery">
                        <AccountRecoveryForm />
                    </Dialog>
                    <h1>Sign In</h1>
                    {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
                    <div className="flex mt-4">
                        <div className="grow flex flex-col gap-2 pr-8">
                            <p>If you already have an account, you can sign in to the right.</p>
                            <p>
                                Otherwise, please{' '}
                                <Link href="/register" className="link">
                                    register for an account
                                </Link>
                                .
                            </p>
                            <p>
                                <span className="link" onClick={() => setIsRecovering(true)}>
                                    Lost access to your account?
                                </span>
                            </p>
                        </div>
                        <div className="flex flex-col gap-4 items-center mx-8 mb-8 basis-2/5">
                            <Button
                                className="flex justify-center items-center w-full gap-2"
                                disabled={isBusy}
                                onClick={() => signInWithPasskey()}
                            >
                                <KeyIcon className="h-[1rem]" />
                                Sign in with Passkey
                            </Button>
                            <p className="text-sm">or...</p>
                            <form className="w-full flex flex-col gap-4">
                                <TextField
                                    disabled={isBusy}
                                    label="Email Address"
                                    onChange={setEmailAddress}
                                    type="email"
                                    required
                                    value={emailAddress}
                                />
                                <TextField
                                    disabled={isBusy}
                                    label="Password"
                                    onChange={setPassword}
                                    type="password"
                                    required
                                    value={password}
                                />
                                <Button
                                    className="w-full"
                                    disabled={isBusy}
                                    label="Sign in with Password"
                                    onClick={() => signInWithPassword()}
                                    type="submit"
                                />
                            </form>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};
