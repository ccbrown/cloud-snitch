'use client';

import Link from 'next/link';
import React, { useEffect, useState } from 'react';

import { ErrorMessage } from '@/components';
import { useDispatch } from '@/store';

export const CompleteRegistrationPage = () => {
    const [success, setSuccess] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const dispatch = useDispatch();

    useEffect(() => {
        const completeRegistration = async () => {
            const params = new URLSearchParams(window.location.hash.slice(1));
            const token = params.get('token');
            if (!token) {
                setErrorMessage('Invalid registration link.');
                return;
            }

            try {
                await dispatch.users.completeRegistration({ token });
                setSuccess(true);
                window.location.hash = '';
            } catch (err) {
                setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            }
        };

        completeRegistration();
    }, [dispatch]);

    return (
        <div className="translucent-snow max-w-4xl mx-auto rounded-xl p-4">
            {success ? (
                <div>
                    <h1 className="mb-4">Registration complete! 🎉</h1>
                    <p>
                        You should now continue to your{' '}
                        <Link href="/dashboard" className="link">
                            dashboard
                        </Link>
                        .
                    </p>
                </div>
            ) : errorMessage ? (
                <ErrorMessage>{errorMessage}</ErrorMessage>
            ) : (
                <p>Completing registration...</p>
            )}
        </div>
    );
};
