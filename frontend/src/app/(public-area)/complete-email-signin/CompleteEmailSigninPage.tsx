'use client';

import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';

import { ErrorMessage } from '@/components';
import { useDispatch } from '@/store';

export const CompleteEmailSigninPage = () => {
    const [errorMessage, setErrorMessage] = useState('');

    const router = useRouter();
    const dispatch = useDispatch();

    useEffect(() => {
        const completeSignin = async () => {
            const params = new URLSearchParams(window.location.hash.slice(1));
            const token = params.get('token');
            if (!token) {
                setErrorMessage('Invalid sign-in link.');
                return;
            }

            try {
                await dispatch.api.signIn({ token });
                router.push('/dashboard');
            } catch (err) {
                setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            }
        };

        completeSignin();
    }, [dispatch, router]);

    return (
        <div className="translucent-snow max-w-4xl mx-auto rounded-xl p-4">
            {errorMessage ? <ErrorMessage>{errorMessage}</ErrorMessage> : <p>Signing in...</p>}
        </div>
    );
};
