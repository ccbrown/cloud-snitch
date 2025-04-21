import { clsx } from 'clsx';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { ChevronRightIcon, SparklesIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { Transition } from '@headlessui/react';

import { awsServices } from '@/aws';
import { Button, ChipEditor, Dialog, ErrorMessage, SuccessMessage } from '@/components';
import { AWSAccount } from '@/generated/api';
import { useAwsRegions, useCurrentTeamId, useManagedAwsScp, useTeamAwsAccountsMap } from '@/hooks';
import { RuleSet } from '@/rules';
import { useDispatch } from '@/store';

interface PolicyPreviewProps {
    account: AWSAccount;
    ruleSet: RuleSet;
    onSuccess: () => void;
}

const PolicyPreview = ({ account, onSuccess, ruleSet }: PolicyPreviewProps) => {
    const teamId = useCurrentTeamId();
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const dispatch = useDispatch();

    const prettyContent = useMemo(() => JSON.stringify(ruleSet.scp(), null, 2), [ruleSet]);

    const apply = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.aws.putManagedScpByTeamAndAccountId({
                teamId,
                accountId: account.id,
                input: {
                    content: JSON.stringify(ruleSet.scp()),
                },
            });
            onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
        } finally {
            setIsBusy(false);
        }
    };

    return (
        <div className="flex flex-col gap-4">
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <p>
                ⚠️ If you proceed, the following service control policy will be applied to{' '}
                <span className="font-semibold">{account.name ? `${account.name} (${account.id})` : account.id}</span>.
                Please exercise caution as it is possible to lock yourself out or disrupt services running in the
                account.
            </p>
            <pre className="p-2 border-1 border-english-violet/60 rounded-lg text-sm max-h-[50vh] overflow-auto">
                <code>{prettyContent}</code>
            </pre>
            <Button disabled={isBusy} label="Apply Policy" onClick={() => apply()} />
        </div>
    );
};

interface AccountPageProps {
    account: AWSAccount;
    onBack: () => void;
}

const AccountPage = ({ account, onBack }: AccountPageProps) => {
    const [errorMessage, setErrorMessage] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
    const [isLoadingSuggestions, setIsLoadingSuggestions] = useState(false);
    const [isPreviewing, setIsPreviewing] = useState(false);
    const dispatch = useDispatch();
    const awsRegions = useAwsRegions();
    const teamId = useCurrentTeamId();
    const scp = useManagedAwsScp(teamId, account.id);
    const scpContent = scp && scp.content;
    const scpRuleSet = useMemo(
        () => (scpContent === undefined ? undefined : scpContent ? RuleSet.fromScpContent(scpContent) : new RuleSet()),
        [scpContent],
    );

    const [ruleSet, setRuleSet] = useState(new RuleSet());
    const updateRuleSet = (update: (prev: RuleSet) => void) => {
        const newRuleSet = ruleSet.clone();
        update(newRuleSet);
        setRuleSet(newRuleSet);
    };

    const hasChanges = useMemo(() => scpRuleSet && !ruleSet.equal(scpRuleSet), [ruleSet, scpRuleSet]);

    useEffect(() => {
        if (scpRuleSet) {
            setRuleSet(scpRuleSet.clone());
        }
    }, [scpRuleSet]);

    const suggest = useCallback(() => {
        const impl = async () => {
            if (isLoadingSuggestions) {
                return;
            }
            setIsLoadingSuggestions(true);
            setErrorMessage('');
            setSuccessMessage('');

            try {
                const [reports, awsAccessReport] = await Promise.all([
                    dispatch.reports.fetchTeamReports(teamId),
                    dispatch.aws.fetchAccessReportByTeamAndAccountId({
                        teamId,
                        accountId: account.id,
                    }),
                ]);

                const ruleSet = new RuleSet();

                for (const report of reports) {
                    if (report.size > 0 && report.scope.aws.accountId === account.id) {
                        ruleSet.addRegionToAllowlist(report.scope.aws.region);
                    }
                }

                const cutoffDate = new Date();
                cutoffDate.setMonth(cutoffDate.getMonth() - 2);

                for (const service of awsAccessReport.services) {
                    if (service.lastAuthenticationTime && service.lastAuthenticationTime > cutoffDate) {
                        ruleSet.addServiceToAllowlist(service.namespace);
                    }
                }

                setRuleSet(ruleSet);
                setSuccessMessage(
                    'Based on your recent account activity, we recommend the following rules. Please review them carefully and add or remove items as needed.',
                );
            } catch (err) {
                setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            } finally {
                setIsLoadingSuggestions(false);
            }
        };
        impl();
    }, [
        isLoadingSuggestions,
        setIsLoadingSuggestions,
        dispatch,
        teamId,
        account,
        setErrorMessage,
        setSuccessMessage,
        setRuleSet,
    ]);

    return (
        <div className="flex flex-col gap-4">
            {scp === undefined ? (
                <p>Loading...</p>
            ) : (
                <>
                    <div className="flex gap-2 items-end">
                        <div className="grow flex flex-col">
                            <div className="font-semibold">{account.name || account.id}</div>
                            {account.name && <span className="text-xs">{account.id}</span>}
                        </div>
                        <div
                            className={`text-sm flex gap-1 items-center ${isLoadingSuggestions ? 'opacity-50' : 'cursor-pointer hover:text-majorelle-blue'} transition-all duration-200 ease-in-out`}
                            onClick={suggest}
                        >
                            <SparklesIcon className="h-[1rem]" />
                            <span>{isLoadingSuggestions ? 'Please wait...' : 'Suggest Changes'}</span>
                        </div>
                    </div>

                    {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
                    {successMessage && <SuccessMessage>{successMessage}</SuccessMessage>}

                    <div className="border border-english-violet/60 bg-white/20 text-sm rounded-lg p-2 flex flex-col gap-2">
                        <div>
                            <span className="label">Region allowlist: </span>
                            <ChipEditor
                                options={awsRegions.map((region) => ({
                                    label: region.id,
                                    value: region.id,
                                    altLabel: region.name,
                                }))}
                                before={scpRuleSet?.regionAllowlist?.regions || new Set()}
                                after={ruleSet.regionAllowlist?.regions || new Set()}
                                onAdd={(region) => {
                                    updateRuleSet((ruleSet) => ruleSet.addRegionToAllowlist(region));
                                }}
                                onRemove={(region) => {
                                    updateRuleSet((ruleSet) => ruleSet.removeRegionFromAllowlist(region));
                                }}
                            />
                        </div>
                        <div className="text-xs">
                            {ruleSet.regionAllowlist ? (
                                ruleSet.hasRegionInAllowlist('us-east-1') ? (
                                    'Regions not listed above will be blocked.'
                                ) : (
                                    <span>
                                        Regions not listed above will be blocked, with exceptions for IAM,
                                        Organizations, and Account Management as{' '}
                                        <Link
                                            href="https://docs.aws.amazon.com/whitepapers/latest/aws-fault-isolation-boundaries/global-services.html"
                                            className="external-link"
                                            rel="noopener noreferrer"
                                            target="_blank"
                                        >
                                            these are global services depending on us-east-1
                                        </Link>
                                        .
                                    </span>
                                )
                            ) : (
                                'All regions will be allowed.'
                            )}
                        </div>
                    </div>

                    <div className="border border-english-violet/60 bg-white/20 text-sm rounded-lg p-2 flex flex-col gap-2">
                        <div>
                            <span className="label">Service allowlist: </span>
                            <ChipEditor
                                options={awsServices.map((service) => ({
                                    label: service.namespace,
                                    value: service.namespace,
                                    altLabel: service.name,
                                }))}
                                before={scpRuleSet?.serviceAllowlist?.services || new Set()}
                                after={ruleSet.serviceAllowlist?.services || new Set()}
                                onAdd={(service) => {
                                    updateRuleSet((ruleSet) => ruleSet.addServiceToAllowlist(service));
                                }}
                                onRemove={(service) => {
                                    updateRuleSet((ruleSet) => ruleSet.removeServiceFromAllowlist(service));
                                }}
                            />
                        </div>
                        <div className="text-xs">
                            {ruleSet.serviceAllowlist
                                ? 'Services not listed above will be blocked.'
                                : 'All services will be allowed.'}
                        </div>
                    </div>

                    <Dialog isOpen={isPreviewing} onClose={() => setIsPreviewing(false)} title="Policy Preview">
                        <PolicyPreview account={account} ruleSet={ruleSet} onSuccess={() => setIsPreviewing(false)} />
                    </Dialog>
                    <Button disabled={!hasChanges} label="Preview Policy" onClick={() => setIsPreviewing(true)} />
                </>
            )}
            <div className="flex justify-center">
                <span onClick={onBack} className="cursor-pointer text-sm link">
                    Back to Accounts
                </span>
            </div>
        </div>
    );
};

export const Rules = () => {
    const teamId = useCurrentTeamId();
    const teamAwsAccountsMap = useTeamAwsAccountsMap(teamId);
    const [accountId, setAccountId] = useState<string | null>(null);
    const [transitionDirection, setTransitionDirection] = useState<'forward' | 'backward'>('forward');

    const account = accountId && teamAwsAccountsMap?.get(accountId);

    const sortedAccounts =
        teamAwsAccountsMap &&
        Array.from(teamAwsAccountsMap.values()).sort((a, b) => {
            const aLabel = a.name?.toLowerCase() || a.id;
            const bLabel = b.name?.toLowerCase() || b.id;
            return aLabel.localeCompare(bLabel);
        });

    const pageClassName = clsx([
        'flex flex-col gap-2',
        'data-[closed]:opacity-0 data-[closed]:absolute',
        'data-[enter]:duration-100 data-[leave]:duration-300',
        transitionDirection === 'forward' &&
            'data-[enter]:data-[closed]:translate-x-full data-[leave]:data-[closed]:-translate-x-full',
        transitionDirection === 'backward' &&
            'data-[enter]:data-[closed]:-translate-x-full data-[leave]:data-[closed]:translate-x-full',
    ]);

    return (
        <div className="relative overflow-hidden">
            <Transition show={!account}>
                <div className={pageClassName}>
                    <p>
                        You can use Cloud Snitch to enforce rules for the following AWS accounts. Cloud Snitch does this
                        by attaching{' '}
                        <Link
                            href="https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_scps.html"
                            className="external-link"
                            rel="noopener noreferrer"
                            target="_blank"
                        >
                            Service Control Policies
                        </Link>{' '}
                        to accounts.
                    </p>
                    <div className="flex flex-col max-h-[50vh] overflow-auto">
                        {sortedAccounts &&
                            sortedAccounts.map((account) => (
                                <div
                                    key={account.id}
                                    className="flex items-center gap-2 hover:bg-white/80 rounded-md cursor-pointer p-2"
                                    onClick={() => {
                                        setTransitionDirection('forward');
                                        setAccountId(account.id);
                                    }}
                                >
                                    {account.name ? (
                                        <div className="flex flex-col">
                                            <span className="font-semibold">{account.name}</span>
                                            <span className="text-xs">{account.id}</span>
                                        </div>
                                    ) : (
                                        <span className="font-semibold">{account.id}</span>
                                    )}
                                    <div className="grow flex justify-end">
                                        <ChevronRightIcon className="h-4 w-4 text-gray-400" />
                                    </div>
                                </div>
                            ))}
                    </div>
                </div>
            </Transition>
            <Transition show={!!account}>
                <div className={pageClassName}>
                    {account && (
                        <AccountPage
                            account={account}
                            onBack={() => {
                                setTransitionDirection('backward');
                                setAccountId(null);
                            }}
                        />
                    )}
                </div>
            </Transition>
        </div>
    );
};
