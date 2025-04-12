'use client';

import Link from 'next/link';
import { useState } from 'react';

import { REVISION as COOKIE_POLICY_REVISION } from '../cookie-policy/revision';
import { REVISION as TERMS_OF_SERVICE_REVISION } from '../terms-of-service/revision';
import { REVISION as PRIVACY_POLICY_REVISION } from '../privacy-policy/revision';
import { Button, Checkbox, ErrorMessage, SuccessMessage, TextField } from '@/components';
import { useDispatch } from '@/store';

export const RegistrationForm = () => {
    const dispatch = useDispatch();
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [emailAddress, setEmailAddress] = useState('');
    const [agreeToEverything, setAgreeToEverything] = useState(false);
    const [success, setSuccess] = useState(false);

    const hasValidInput = emailAddress && agreeToEverything;

    const submit = async () => {
        if (isBusy || !hasValidInput) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.users.beginRegistration({
                emailAddress,
                cookiePolicyAgreementRevision: COOKIE_POLICY_REVISION,
                privacyPolicyAgreementRevision: PRIVACY_POLICY_REVISION,
                termsOfServiceAgreementRevision: TERMS_OF_SERVICE_REVISION,
            });
            setSuccess(true);
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        } finally {
            setIsBusy(false);
        }
    };

    return success ? (
        <SuccessMessage>
            Success! Check your email for the next steps. If you don&apos;t see an email from Cloud Snitch within a few
            minutes, check your spam folder and{' '}
            <Link href="/contact" className="link">
                contact us
            </Link>{' '}
            if you still don&apos;t see it.
        </SuccessMessage>
    ) : (
        <>
            <p>
                Provide your email address and agree to let us do stuff, and we&apos;ll send you a verification email to
                create your account and get started. 🚀
            </p>
            <p>
                Already have an account?{' '}
                <Link href="/sign-in" className="link">
                    Click here
                </Link>{' '}
                to sign in.
            </p>
            <form>
                {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
                <div className="max-w-md">
                    <TextField
                        disabled={isBusy}
                        label="Email Address"
                        onChange={setEmailAddress}
                        type="email"
                        required
                        value={emailAddress}
                    />
                    <Checkbox checked={agreeToEverything} className="my-4" onChange={setAgreeToEverything}>
                        <div className="text-sm">
                            I&apos;ve read and agree to the{' '}
                            <Link href="/terms-of-service" className="link" target="_blank">
                                Terms of Service
                            </Link>{' '}
                            and{' '}
                            <Link href="/privacy-policy" className="link" target="_blank">
                                Privacy Policy
                            </Link>
                            , and consent to the use of cookies as defined by the{' '}
                            <Link href="/cookie-policy" className="link" target="_blank">
                                Cookie Policy
                            </Link>
                            .
                        </div>
                    </Checkbox>
                    <Button
                        className="w-full"
                        disabled={!hasValidInput || isBusy}
                        label="Register"
                        onClick={() => submit()}
                        type="submit"
                    />
                </div>
            </form>
        </>
    );
};
