'use client';

import { PlusCircleIcon, TrashIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

import { Button, Dialog, ErrorMessage } from '@/components';
import { ReportRetention, UserRole } from '@/generated/api';
import { useCurrentTeamId, useCurrentUser, useTeamReports } from '@/hooks';
import { useDispatch } from '@/store';

interface QueueGenerationFormProps {
    teamId: string;
    onSuccess: () => void;
}

const QueueGenerationForm = (props: QueueGenerationFormProps) => {
    const dispatch = useDispatch();

    const [isBusy, setIsBusy] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const startTime = new Date();
    startTime.setDate(startTime.getDate() - 1);
    startTime.setUTCHours(0);
    startTime.setUTCMinutes(0);
    startTime.setUTCSeconds(0);

    const durationSeconds = 24 * 60 * 60;
    const retention = ReportRetention.OneWeek;

    const doCreate = async () => {
        if (isBusy) {
            return;
        }
        setIsBusy(true);

        try {
            await dispatch.reports.queueTeamReportGeneration({
                teamId: props.teamId,
                input: {
                    startTime,
                    durationSeconds,
                    retention,
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
            <p>Start Time: {startTime.toLocaleString()}</p>
            <p>Duration: {durationSeconds} seconds</p>
            <p>Retention: {retention}</p>
            <Button
                disabled={isBusy}
                label="Queue Report Generation"
                onClick={doCreate}
                type="submit"
                className="mt-4"
            />
        </form>
    );
};

const Page = () => {
    const dispatch = useDispatch();
    const teamId = useCurrentTeamId();
    const reports = useTeamReports(teamId);
    const isAdmin = useCurrentUser()?.role === UserRole.Administrator;
    const [isQueueing, setIsQueueing] = useState(false);

    return (
        <div>
            <Dialog isOpen={isQueueing} onClose={() => setIsQueueing(false)} title="Queue Report Generation">
                <QueueGenerationForm onSuccess={() => setIsQueueing(false)} teamId={teamId} />
            </Dialog>
            <h2 className="mb-4 flex items-center gap-2">
                Reports {reports && `(${reports.length})`}{' '}
                {isAdmin && (
                    <PlusCircleIcon
                        className="h-[1.5rem] cursor-pointer hover:text-amethyst"
                        onClick={() => setIsQueueing(true)}
                    />
                )}
            </h2>
            {!reports ? (
                <p>Loading...</p>
            ) : reports.length === 0 ? (
                <p>
                    No data has been collected for this team so far. Make sure you&apos;ve configured an AWS
                    integration, then check again later.
                </p>
            ) : (
                <div>
                    <p className="mb-4">
                        This is all the data we&apos;ve collected for this team. If for any reason you want your AWS
                        data to disappear, you can delete it here.
                    </p>
                    <table className="w-full text-left text-sm">
                        <thead className="uppercase text-english-violet">
                            <tr>
                                <th>Account</th>
                                <th>Region</th>
                                <th>Time</th>
                                <th />
                            </tr>
                        </thead>
                        <tbody>
                            {reports.map((report) => (
                                <tr key={report.id}>
                                    <td>{report.scope.aws.accountId}</td>
                                    <td>{report.scope.aws.region}</td>
                                    <td>{report.scope.startTime.toLocaleString()}</td>
                                    <td align="right" className="px-2">
                                        <TrashIcon
                                            className="h-[1rem] cursor-pointer hover:text-amethyst"
                                            onClick={() => {
                                                if (confirm('Are you sure you want to delete this report?')) {
                                                    dispatch.reports.deleteReport(report.id);
                                                }
                                            }}
                                        />
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    {isAdmin && (
                        <Button
                            label="Delete All"
                            onClick={() => {
                                reports?.forEach((r) => {
                                    dispatch.reports.deleteReport(r.id);
                                });
                            }}
                        />
                    )}
                </div>
            )}
        </div>
    );
};

export default Page;
