'use client';

import { PencilIcon } from '@heroicons/react/24/outline';
import { loadStripe } from '@stripe/stripe-js';
import { AddressElement, Elements, PaymentElement, useElements, useStripe } from '@stripe/react-stripe-js';
import Link from 'next/link';
import { useState } from 'react';

import { Button, Dialog, ErrorMessage, SubscriptionSelector, Tooltip, WarningMessage } from '@/components';
import { TeamBillingProfile, TeamSubscription, TeamSubscriptionTier } from '@/generated/api';
import { useCurrentTeamId, useTeamBillingProfile, useTeamPaymentMethod, useTeamSubscription } from '@/hooks';
import { useDispatch } from '@/store';

interface BillingAddressFormProps {
    teamId: string;
    currentProfile?: TeamBillingProfile | null;
    onSuccess: () => void;
}

const BillingAddressForm = (props: BillingAddressFormProps) => {
    const dispatch = useDispatch();
    const elements = useElements();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doSave = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            const addressElement = elements?.getElement('address');
            if (!addressElement) {
                setErrorMessage('Address element not found. Please reload the page and try again.');
                setIsBusy(false);
                return;
            }

            const { complete, value } = await addressElement.getValue();
            if (!complete) {
                setErrorMessage('Incomplete address.');
                setIsBusy(false);
                return;
            }

            const input = {
                name: value.name,
                address: {
                    line1: value.address.line1 || undefined,
                    line2: value.address.line2 || undefined,
                    city: value.address.city || undefined,
                    state: value.address.state || undefined,
                    postalCode: value.address.postal_code,
                    country: value.address.country,
                },
            };
            if (props.currentProfile) {
                await dispatch.teams.updateBillingProfile({
                    teamId: props.teamId,
                    input,
                });
            } else {
                await dispatch.teams.createBillingProfile({
                    teamId: props.teamId,
                    input,
                });
            }
            props.onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <AddressElement
                options={{
                    defaultValues: props.currentProfile
                        ? {
                              name: props.currentProfile.name,
                              address: {
                                  line1: props.currentProfile.address.line1,
                                  line2: props.currentProfile.address.line2,
                                  city: props.currentProfile.address.city,
                                  state: props.currentProfile.address.state,
                                  postal_code: props.currentProfile.address.postalCode,
                                  country: props.currentProfile.address.country,
                              },
                          }
                        : undefined,
                    display: {
                        name: 'organization',
                    },
                    mode: 'billing',
                }}
            />
            <Button disabled={isBusy} label="Save Billing Address" onClick={doSave} type="submit" className="mt-6" />
        </form>
    );
};

const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || '', {
    betas: ['custom_checkout_beta_6'],
});

interface PaymentMethodFormProps {
    teamId: string;
    onSuccess: () => void;
}

const PaymentMethodForm = (props: PaymentMethodFormProps) => {
    const elements = useElements();
    const stripe = useStripe();
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doSave = async () => {
        if (isBusy || !elements || !stripe) {
            return;
        }
        setIsBusy(true);

        try {
            const { error: submitError } = await elements.submit();
            if (submitError) {
                throw submitError;
            }

            const { error: paymentMethodError, paymentMethod } = await stripe.createPaymentMethod({
                elements,
            });
            if (paymentMethodError) {
                throw paymentMethodError;
            }

            await dispatch.teams.updatePaymentMethod({
                teamId: props.teamId,
                input: {
                    stripePaymentMethodId: paymentMethod.id,
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
            <PaymentElement />
            <Button disabled={isBusy} label="Save Payment Method" onClick={doSave} type="submit" className="mt-6" />
        </form>
    );
};

interface SubscriptionFormProps {
    teamId: string;
    current?: TeamSubscription | null;
    onSuccess: () => void;
}

const SubscriptionForm = (props: SubscriptionFormProps) => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doUpdate = async (tier: TeamSubscriptionTier) => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            if (props.current) {
                await dispatch.teams.updateSubscription({
                    teamId: props.teamId,
                    input: {
                        tier,
                    },
                });
            } else {
                await dispatch.teams.createSubscription({
                    teamId: props.teamId,
                    input: {
                        tier,
                    },
                });
            }

            // fetch the latest billing profile as the balance may have changed. no need to wait
            // though
            dispatch.teams.fetchBillingProfile(props.teamId);

            // team entitlements update asynchronously, so kick off a small series of refreshes
            for (let i = 0; i <= 5; i++) {
                setTimeout(() => {
                    dispatch.teams.fetch(props.teamId);
                }, i * 1000);
            }

            props.onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col gap-4">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <SubscriptionSelector disabled={isBusy} onSelect={(tier) => doUpdate(tier)} />
            <p className="text-xs">
                {props.current ? (
                    <span>
                        You are currently subscribed to the <strong>{props.current.name}</strong> plan. Changes to your
                        subscription will be prorated, charging your payment method or granting credits accordingly.
                    </span>
                ) : (
                    'You are not currently subscribed a plan.'
                )}
            </p>
            <p className="text-xs">
                <Link href="/pricing" className="link" target="_blank">
                    Click here
                </Link>{' '}
                for more details on features and pricing.
            </p>
        </form>
    );
};

const Page = () => {
    const teamId = useCurrentTeamId();
    const billingProfile = useTeamBillingProfile(teamId);
    const paymentMethod = useTeamPaymentMethod(teamId);
    const subscription = useTeamSubscription(teamId);
    const [isEnteringAddress, setIsEnteringAddress] = useState(false);
    const [isEnteringPaymentMethod, setIsEnteringPaymentMethod] = useState(false);
    const [isUpdatingSubscription, setIsUpdatingSubscription] = useState(false);

    return (
        <Elements
            stripe={stripePromise}
            options={{
                currency: 'usd',
                mode: 'setup',
                setupFutureUsage: 'off_session',
                paymentMethodCreation: 'manual',
            }}
        >
            <div>
                {subscription === null && (
                    <div className="mb-4">
                        <WarningMessage>
                            This team currently has no active subscription. Please provide the following information to
                            activate a subscription.
                        </WarningMessage>
                    </div>
                )}
                <Dialog
                    isOpen={isEnteringAddress}
                    onClose={() => setIsEnteringAddress(false)}
                    title="Edit Billing Address"
                >
                    <BillingAddressForm
                        onSuccess={() => setIsEnteringAddress(false)}
                        currentProfile={billingProfile}
                        teamId={teamId}
                    />
                </Dialog>
                <h2 className="mb-4 flex items-center gap-2">
                    Billing Address{' '}
                    <PencilIcon
                        className="h-[1.2rem] cursor-pointer hover:text-amethyst"
                        onClick={() => setIsEnteringAddress(true)}
                    />
                </h2>
                {billingProfile === undefined ? (
                    <span>Loading...</span>
                ) : billingProfile === null ? (
                    <span>No billing address provided.</span>
                ) : (
                    <div>
                        {billingProfile.name}
                        <br />
                        {billingProfile.address.line1 && (
                            <span>
                                {billingProfile.address.line1}
                                <br />
                            </span>
                        )}
                        {billingProfile.address.line2 && (
                            <span>
                                {billingProfile.address.line2}
                                <br />
                            </span>
                        )}
                        {billingProfile.address.city && <span>{billingProfile.address.city}, </span>}
                        {billingProfile.address.state && <span>{billingProfile.address.state} </span>}
                        {billingProfile.address.postalCode}
                        <br />
                        {billingProfile.address.country}
                    </div>
                )}

                <Dialog
                    isOpen={isEnteringPaymentMethod}
                    onClose={() => setIsEnteringPaymentMethod(false)}
                    title="Edit Payment Method"
                >
                    <PaymentMethodForm onSuccess={() => setIsEnteringPaymentMethod(false)} teamId={teamId} />
                </Dialog>
                <h2 className="my-4 flex items-center gap-2">
                    Payment Method{' '}
                    <PencilIcon
                        className="h-[1.2rem] cursor-pointer hover:text-amethyst"
                        onClick={() =>
                            billingProfile
                                ? setIsEnteringPaymentMethod(true)
                                : alert('Please provide a billing address first.')
                        }
                    />
                </h2>
                {paymentMethod === undefined ? (
                    <span>Loading...</span>
                ) : paymentMethod === null ? (
                    <span>No payment method provided.</span>
                ) : (
                    <div>
                        {paymentMethod.type === 'CARD' &&
                            `Card ending in ${paymentMethod.last4Digits}, expiring in ${paymentMethod.expirationMonth}/${paymentMethod.expirationYear}.`}
                        {paymentMethod.type === 'US_BANK_ACCOUNT' &&
                            `US bank account ending in ${paymentMethod.last4Digits}.`}
                        {paymentMethod.type === 'OTHER' && 'Custom payment method.'}
                    </div>
                )}

                <Dialog
                    isOpen={isUpdatingSubscription}
                    onClose={() => setIsUpdatingSubscription(false)}
                    title="Update Subscription"
                >
                    <SubscriptionForm
                        onSuccess={() => setIsUpdatingSubscription(false)}
                        current={subscription}
                        teamId={teamId}
                    />
                </Dialog>
                <h2 className="my-4 flex items-center gap-2">
                    {subscription ? (
                        <span className="uppercase text-amethyst-gradient">{subscription.name}</span>
                    ) : (
                        'Subscription'
                    )}{' '}
                    <PencilIcon
                        className="h-[1.2rem] cursor-pointer hover:text-amethyst"
                        onClick={() =>
                            paymentMethod
                                ? setIsUpdatingSubscription(true)
                                : alert('Please provide a payment method first.')
                        }
                    />
                </h2>
                <div className="flex flex-col">
                    {billingProfile?.balance && (
                        <div>
                            <strong className="text-english-violet font-semibold">Balance:</strong>{' '}
                            <Tooltip
                                content={
                                    <div className="w-sm">
                                        This is your account&apos;s unpaid balance. A negative balance indicates that
                                        you have credits which will be applied to future charges.
                                    </div>
                                }
                            >
                                <span className="hoverable">{billingProfile.balance.text}</span>
                            </Tooltip>
                        </div>
                    )}
                    {subscription === undefined ? (
                        <div>Loading...</div>
                    ) : subscription === null ? (
                        <div>No active subscription.</div>
                    ) : (
                        <>
                            <div>
                                <strong className="text-english-violet font-semibold">Accounts:</strong>{' '}
                                <Tooltip
                                    content={
                                        <div className="w-sm">
                                            This is the number of AWS accounts currently being billed to your
                                            subscription. New accounts may take up to 24 hours to be reflected here.
                                        </div>
                                    }
                                >
                                    <span className="hoverable">{subscription.accounts}</span>
                                </Tooltip>
                            </div>
                            {subscription.price?.accountMonth && (
                                <div>
                                    <strong className="text-english-violet font-semibold">Monthly Price:</strong>{' '}
                                    {subscription.price.accountMonth.text} per AWS account
                                </div>
                            )}
                        </>
                    )}
                </div>
            </div>
        </Elements>
    );
};

export default Page;
