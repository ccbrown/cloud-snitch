'use client';

import { useParams } from 'next/navigation';
import { useEffect } from 'react';

import { Tooltip } from '@/components';
import { useDispatch, useSelector } from '@/store';

const Page = () => {
    const { teamId } = useParams<{ teamId: string }>();

    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.teams.fetch(teamId);
    }, [dispatch, teamId]);

    const team = useSelector((state) => state.teams.teams[teamId]);

    useEffect(() => {
        if (team) {
            dispatch.aws.fetchIntegrationsByTeamId(team.id);
        }
    }, [dispatch, team]);

    const awsIntegrationIds = useSelector((state) => state.aws.teamIntegrationIds[teamId]);
    const allAwsIntegrations = useSelector((state) => state.aws.integrations);

    if (!team) {
        return <div>Loading...</div>;
    }

    return (
        <div className="flex flex-col gap-2">
            <h2>{team.name}</h2>
            <div>
                <span className="label">Entitlements:</span>{' '}
                {team.entitlements.teamFeatures
                    ? 'Team Features'
                    : team.entitlements.individualFeatures
                      ? 'Individual Features'
                      : 'None'}
            </div>
            <h3>AWS Integrations</h3>
            <table className="w-full text-left">
                <thead className="uppercase text-sm text-english-violet">
                    <tr>
                        <th>Name</th>
                        <th>Orgs</th>
                        <th>SCPs</th>
                        <th>CloudTrail Trail</th>
                    </tr>
                </thead>
                <tbody>
                    {awsIntegrationIds
                        ?.map((integrationId) => allAwsIntegrations[integrationId])
                        .filter((integration) => integration)
                        .map((integration) => (
                            <tr key={integration.id}>
                                <td>{integration.name}</td>
                                <td>{integration.getAccountNamesFromOrganizations ? 'Yes' : 'No'}</td>
                                <td>{integration.manageScps ? 'Yes' : 'No'}</td>
                                <td>
                                    {integration.cloudtrailTrail ? (
                                        <Tooltip
                                            content={`s3://${integration.cloudtrailTrail.s3BucketName}${integration.cloudtrailTrail.s3KeyPrefix}`}
                                        >
                                            <span className="hoverable">Yes</span>
                                        </Tooltip>
                                    ) : (
                                        'No'
                                    )}
                                </td>
                            </tr>
                        ))}
                </tbody>
            </table>
        </div>
    );
};

export default Page;
