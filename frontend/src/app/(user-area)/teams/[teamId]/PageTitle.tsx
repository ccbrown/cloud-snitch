'use client';

import { useCurrentTeamId } from '@/hooks';
import { useSelector } from '@/store';

interface Props {
    children: React.ReactNode;
}

export const PageTitle = (props: Props) => {
    const teamId = useCurrentTeamId();
    const team = useSelector((state) => state.teams.teams[teamId]);
    const title = team ? `${team.name} | Cloud Snitch` : 'Cloud Snitch';

    return (
        <>
            <title>{title}</title>
            {props.children}
        </>
    );
};
