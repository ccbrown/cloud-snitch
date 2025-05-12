'use client';

import Link from 'next/link';
import { useEffect } from 'react';

import { useDispatch, useSelector } from '@/store';

const Page = () => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch.teams.fetchAll();
    }, [dispatch]);

    const teams = useSelector((state) => state.teams.teams);

    return (
        <>
            <h2>Teams</h2>
            <table className="w-full text-left">
                <thead className="uppercase text-sm text-english-violet">
                    <tr>
                        <th>Name</th>
                        <th>Entitlements</th>
                    </tr>
                </thead>
                <tbody>
                    {Object.values(teams).map((team) => (
                        <tr key={team.id}>
                            <td>
                                <Link href={`/control-panel/teams/${team.id}`} className="link">
                                    {team.name}
                                </Link>
                            </td>
                            <td>
                                {team.entitlements.teamFeatures
                                    ? 'Team Features'
                                    : team.entitlements.individualFeatures
                                      ? 'Individual Features'
                                      : 'None'}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </>
    );
};

export default Page;
