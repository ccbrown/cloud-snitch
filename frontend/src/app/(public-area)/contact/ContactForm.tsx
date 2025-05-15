'use client';

import Link from 'next/link';
import { useState } from 'react';

import { Button, ErrorMessage, SuccessMessage, TextArea, TextField } from '@/components';
import { SystemApi } from '@/generated/api';
import { apiConfiguration } from '@/models/api';

export const ContactForm = () => {
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [name, setName] = useState('');
    const [emailAddress, setEmailAddress] = useState('');
    const [subject, setSubject] = useState('');
    const [message, setMessage] = useState('');
    const [success, setSuccess] = useState(false);

    const submit = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            const systemApi = new SystemApi(apiConfiguration());
            await systemApi.contactUs({
                contactUsRequest: {
                    name,
                    emailAddress,
                    subject,
                    message,
                },
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
            Message submitted! We&apos;ll be in touch soon. If you don&apos;t see any emails from us within a few
            business days, please be sure to check your spam folder.
        </SuccessMessage>
    ) : (
        <>
            <p>
                Cloud Snitch is developed by Paragon Cybersecurity, a Limited Liability Company formed under Delaware
                jurisdiction.
            </p>
            <p>For inquiries, please fill out the form below and we will get back to you as soon as possible.</p>
            <p>
                Alternatively, you can{' '}
                <Link href="https://github.com/ccbrown/cloud-snitch/issues" target="_blank" className="external-link">
                    raise an issue on GitHub
                </Link>
                .
            </p>
            <form className="flex flex-col gap-4">
                {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
                <div className="flex gap-4 w-full">
                    <TextField disabled={isBusy} label="Name" onChange={setName} type="text" required value={name} />
                    <TextField
                        disabled={isBusy}
                        label="Email Address"
                        onChange={setEmailAddress}
                        type="email"
                        required
                        value={emailAddress}
                    />
                </div>
                <TextField
                    disabled={isBusy}
                    label="Subject"
                    onChange={setSubject}
                    type="text"
                    required
                    value={subject}
                />
                <TextArea disabled={isBusy} label="Message" onChange={setMessage} required value={message} />
                <Button className="w-full" disabled={isBusy} label="Submit" onClick={() => submit()} type="submit" />
            </form>
        </>
    );
};
