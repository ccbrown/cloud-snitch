'use client';

import { clsx } from 'clsx';
import Link from 'next/link';
import { PencilIcon, PlusCircleIcon, TrashIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';
import { Transition } from '@headlessui/react';

import { Button, Checkbox, Dialog, ErrorMessage, TextField, Tooltip } from '@/components';
import { AWSIntegration } from '@/generated/api';
import { useCurrentTeam, useCurrentTeamId, useTeamAwsIntegrations } from '@/hooks';
import { INTEGRATION_TEMPLATE_S3_URL } from '@/integration';
import { useDispatch } from '@/store';

interface CreateIntegrationFormProps {
    teamId: string;
    onSuccess: () => void;
}

const CreateIntegrationForm = (props: CreateIntegrationFormProps) => {
    const dispatch = useDispatch();

    const [transitionDirection, setTransitionDirection] = useState<'forward' | 'backward'>('forward');
    const [step, setStep] = useState(0);
    const [name, setName] = useState('My Integration');
    const [roleArn, setRoleArn] = useState('');
    const [s3BucketName, setS3BucketName] = useState('');
    const [s3KeyPrefix, setS3KeyPrefix] = useState('');
    const [getAccountNamesFromOrganizations, setGetAccountNamesFromOrganizations] = useState(false);
    const [manageScps, setManageScps] = useState(false);
    const [isBusy, setIsBusy] = useState(false);
    const [queueReportGeneration, setQueueReportGeneration] = useState(true);
    const [errorMessage, setErrorMessage] = useState('');

    const doCreate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.aws.createIntegration({
                teamId: props.teamId,
                input: {
                    name,
                    roleArn,
                    cloudtrailTrail: s3BucketName
                        ? {
                              s3BucketName,
                              s3KeyPrefix,
                          }
                        : undefined,
                    getAccountNamesFromOrganizations,
                    manageScps,
                    queueReportGeneration,
                },
            });
            props.onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    const cfnParameters = [
        ['CloudSnitchAWSAccountId', process.env.NEXT_PUBLIC_AWS_ACCOUNT_ID],
        ['TeamId', props.teamId],
        ['AllowOrganizationsAccess', getAccountNamesFromOrganizations ? 'Yes' : 'No'],
        ['AllowSCPManagement', manageScps ? 'Yes' : 'No'],
        ['S3BucketName', s3BucketName],
        ['S3KeyPrefix', s3KeyPrefix],
    ];

    const quickLinkParams = new URLSearchParams();
    quickLinkParams.append('templateURL', INTEGRATION_TEMPLATE_S3_URL);
    quickLinkParams.append('stackName', 'CloudSnitchIntegration');
    cfnParameters.forEach(([key, value]) => {
        if (value) {
            quickLinkParams.append(`param_${key}`, value);
        }
    });
    const quickLink = `https://console.aws.amazon.com/cloudformation/home#/stacks/create/review?${quickLinkParams.toString()}`;

    const formClassName = clsx([
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
            <Transition show={step === 0}>
                <form className={formClassName}>
                    <p>
                        Before configuring an integration, you&apos;ll need to make sure you have a CloudTrail trail set
                        up in your AWS account. If you don&apos;t already have one, we recommend creating an
                        organization trail in the management account or a delegated administrator account as described
                        in{' '}
                        <Link
                            href="https://docs.aws.amazon.com/awscloudtrail/latest/userguide/creating-trail-organization.html"
                            target="_blank"
                            className="external-link"
                        >
                            Creating a trail for an organization
                        </Link>
                        . If you need help setting this up, please{' '}
                        <Link href="/contact" target="_blank" className="link">
                            contact us
                        </Link>
                        .
                    </p>
                    <p>Once you have a trail, provide the bucket name and prefix (if any) below.</p>
                    <TextField
                        disabled={isBusy}
                        label="S3 Bucket Name"
                        required
                        value={s3BucketName}
                        onChange={setS3BucketName}
                    />
                    <TextField disabled={isBusy} label="S3 Key Prefix" value={s3KeyPrefix} onChange={setS3KeyPrefix} />
                    <Checkbox
                        disabled={isBusy}
                        checked={getAccountNamesFromOrganizations}
                        onChange={setGetAccountNamesFromOrganizations}
                        label="Enable AWS Organizations integration"
                        subLabel="If you're using an organization trail, we recommend checking this box as this will allow us to show your account names in the UI."
                    />
                    <Tooltip
                        disabled={getAccountNamesFromOrganizations}
                        content="⚠️ Requires AWS Organizations integration"
                    >
                        <Checkbox
                            disabled={isBusy || !getAccountNamesFromOrganizations}
                            checked={manageScps}
                            onChange={setManageScps}
                            label="Enable SCP management"
                            subLabel="If you're deploying to an organization management account, checking this box will allow Cloud Snitch to enforce access controls through service control policies. For example, Cloud Snitch can be configured to block activity for services and regions that you don't use. Cloud Snitch will only be able to block actions and will not be able to read or modify existing policies or grant additional access."
                        />
                    </Tooltip>
                    <Button
                        disabled={!s3BucketName && !getAccountNamesFromOrganizations}
                        label="Continue"
                        onClick={() => {
                            setTransitionDirection('forward');
                            setStep(1);
                        }}
                        type="submit"
                        className="mt-4"
                    />
                </form>
            </Transition>
            <Transition show={step === 1}>
                <form className={formClassName}>
                    <p>
                        Now you&apos;ll need to create an IAM role in your AWS account for Cloud Snitch to use. This
                        will be done via Infrastructure as Code using CloudFormation.
                    </p>
                    <div className="text-sm uppercase font-bold text-english-violet mt-2">Template URL</div>
                    <Link href={INTEGRATION_TEMPLATE_S3_URL} target="_blank" className="external-link text-sm">
                        {INTEGRATION_TEMPLATE_S3_URL}
                    </Link>
                    <table className="w-full text-sm my-2">
                        <thead className="uppercase text-english-violet">
                            <tr>
                                <th className="text-left">Parameter</th>
                                <th className="text-left">Value</th>
                            </tr>
                        </thead>
                        <tbody className="text-xs">
                            {cfnParameters
                                .filter(([_key, value]) => !!value)
                                .map(([key, value]) => (
                                    <tr key={key} className="border-t border-platinum">
                                        <td className="py-1">{key}</td>
                                        <td>{value}</td>
                                    </tr>
                                ))}
                        </tbody>
                    </table>
                    <p>
                        To open the CloudFormation console and deploy this template using the parameters above,{' '}
                        <Link href={quickLink} target="_blank" className="external-link">
                            click here
                        </Link>
                        .
                    </p>
                    <Button
                        label="Continue"
                        onClick={() => {
                            setTransitionDirection('forward');
                            setStep(2);
                        }}
                        type="submit"
                        className="mt-4"
                    />
                    <p className="mt-4 text-center">
                        <span
                            className="link"
                            onClick={() => {
                                setTransitionDirection('backward');
                                setStep(0);
                            }}
                        >
                            Back
                        </span>
                    </p>
                </form>
            </Transition>
            <Transition show={step === 2}>
                <form className={formClassName}>
                    {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
                    <p>Lastly, provide a name for your integration and the ARN of the role you created.</p>
                    <p>
                        The ARN can be found in &quot;Outputs&quot; tab of the CloudFormation stack you created and will
                        begin with &quot;arn:aws:iam&quot;.
                    </p>
                    <TextField disabled={isBusy} label="Name" required value={name} onChange={setName} />
                    <TextField disabled={isBusy} label="Role ARN" required value={roleArn} onChange={setRoleArn} />
                    <Checkbox
                        disabled={isBusy}
                        checked={queueReportGeneration}
                        onChange={setQueueReportGeneration}
                        label="Backfill data"
                        subLabel="If checked, we'll go ahead and ingest up to a week's worth of data, which will start to become available within minutes. Otherwise, data will be ingested starting tomorrow."
                    />
                    <Button
                        disabled={isBusy || !roleArn || !name}
                        label="Test and Save Integration"
                        onClick={doCreate}
                        type="submit"
                        className="mt-4"
                    />
                    <p className="mt-4 text-center">
                        <span
                            className="link"
                            onClick={() => {
                                setTransitionDirection('backward');
                                setStep(1);
                            }}
                        >
                            Back
                        </span>
                    </p>
                </form>
            </Transition>
        </div>
    );
};

interface EditIntegrationFormProps {
    current: AWSIntegration;
    onSuccess: () => void;
}

const EditIntegrationForm = ({ onSuccess, current }: EditIntegrationFormProps) => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const [name, setName] = useState(current.name);

    const doUpdate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.aws.updateIntegration({
                integrationId: current.id,
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
            <Button disabled={isBusy} label="Save Changes" onClick={doUpdate} type="submit" className="mt-4" />
        </form>
    );
};

interface DeleteIntegrationFormProps {
    integrationId: string;
    onSuccess: () => void;
}

const DeleteIntegrationForm = (props: DeleteIntegrationFormProps) => {
    const dispatch = useDispatch();

    const [deleteAssociatedData, setDeleteAssociatedData] = useState(false);
    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const doDelete = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.aws.deleteIntegration({
                id: props.integrationId,
                deleteAssociatedData,
            });
            props.onSuccess();
        } catch (err) {
            setErrorMessage(err instanceof Error ? err.message : 'An unknown error occurred.');
            setIsBusy(false);
        }
    };

    return (
        <form className="flex flex-col">
            <p className="mb-2">Are you sure you want to delete this integration?</p>
            {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
            <Checkbox
                disabled={isBusy}
                checked={deleteAssociatedData}
                onChange={setDeleteAssociatedData}
                label="Also delete associated data"
            />
            <Button disabled={isBusy} label="Delete Integration" onClick={doDelete} type="submit" className="mt-4" />
        </form>
    );
};

const Page = () => {
    const teamId = useCurrentTeamId();
    const team = useCurrentTeam();
    const integrations = useTeamAwsIntegrations(teamId);
    const [isCreating, setIsCreating] = useState(false);
    const [editIntegrationId, setEditIntegrationId] = useState<string>('');
    const [deleteIntegrationId, setDeleteIntegrationId] = useState<string>('');
    const hasSubscription = team?.entitlements?.individualFeatures;

    const editIntegration =
        editIntegrationId && integrations?.find((integration) => integration.id === editIntegrationId);

    return (
        <div>
            <Dialog isOpen={isCreating} onClose={() => setIsCreating(false)} title="Add AWS Integration">
                <CreateIntegrationForm onSuccess={() => setIsCreating(false)} teamId={teamId} />
            </Dialog>
            <Dialog isOpen={!!editIntegration} onClose={() => setEditIntegrationId('')} title="Edit AWS Integration">
                {editIntegration && (
                    <EditIntegrationForm onSuccess={() => setEditIntegrationId('')} current={editIntegration} />
                )}
            </Dialog>
            <Dialog
                isOpen={!!deleteIntegrationId}
                onClose={() => setDeleteIntegrationId('')}
                title="Delete AWS Integration"
            >
                <DeleteIntegrationForm
                    onSuccess={() => setDeleteIntegrationId('')}
                    integrationId={deleteIntegrationId}
                />
            </Dialog>
            <h2 className="mb-4 flex items-center gap-2">
                AWS Integrations ({integrations?.length}){' '}
                {hasSubscription && (
                    <PlusCircleIcon
                        className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                        onClick={() => setIsCreating(true)}
                    />
                )}
            </h2>
            {!integrations ? (
                <p>Loading...</p>
            ) : integrations.length === 0 ? (
                <p>You currently have no AWS integrations.</p>
            ) : (
                <div className="flex flex-col gap-2">
                    {integrations.map((integration) => (
                        <div key={integration.id} className="flex flex-col border-1 border-platinum p-4 rounded-lg">
                            <div className="flex gap-2 mb-2">
                                <h3 className="font-bold grow">{integration.name}</h3>
                                <PencilIcon
                                    className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                                    onClick={() => setEditIntegrationId(integration.id)}
                                />
                                <TrashIcon
                                    className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                                    onClick={() => setDeleteIntegrationId(integration.id)}
                                />
                            </div>
                            <div className="flex flex-col gap-1 text-sm">
                                <div>
                                    <span className="label">Creation Time:</span>{' '}
                                    {integration.creationTime.toLocaleString()}
                                </div>
                                <div>
                                    <span className="label">Role ARN:</span> {integration.roleArn}
                                </div>
                                <div>
                                    <span className="label">CloudTrail Trail:</span>{' '}
                                    {integration.cloudtrailTrail
                                        ? `s3://${integration.cloudtrailTrail.s3BucketName}${integration.cloudtrailTrail.s3KeyPrefix || ''}`
                                        : 'None'}
                                </div>
                                <div>
                                    <span className="label">AWS Organizations Integration:</span>{' '}
                                    {integration.getAccountNamesFromOrganizations ? 'Yes' : 'No'}
                                </div>
                                <div>
                                    <span className="label">Managed SCPs:</span> {integration.manageScps ? 'Yes' : 'No'}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
            {integrations && !hasSubscription && (
                <div className="mt-4">
                    <Link href={`/teams/${teamId}/settings/billing`} className="link">
                        Activate your subscription
                    </Link>{' '}
                    to add AWS integrations.
                </div>
            )}
        </div>
    );
};

export default Page;
