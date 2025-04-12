import React, { useState } from 'react';

import { PrincipalIcon } from '@/components';
import { useSearchParamState } from '@/hooks';
import { CombinedReport } from '@/report';

import { parseSelection } from './selection';

interface Props {
    combinedReport: CombinedReport;
    onPrincipalHover?: (id: string | null) => void;
    principalFilter?: string;
    onPrincipalFilterChange?: (value: string) => void;
}

const lastComponent = (name: string) => {
    const parts = name.split('/');
    const last = parts[parts.length - 1];
    const parts2 = last.split(':');
    return parts2[parts2.length - 1];
};

export const NavigationPanel = (props: Props) => {
    const [selectionString, setSelectionString] = useSearchParamState('selection');
    const selection = parseSelection(selectionString);

    const [hoveredPrincipalId, setHoveredPrincipalId] = useState<string | null>(null);

    const sortedPrincipals = Object.entries(props.combinedReport.principals).sort((a, b) => {
        const aLabel = lastComponent(a[1].name || a[0]).toLowerCase();
        const bLabel = lastComponent(b[1].name || b[0]).toLowerCase();
        return aLabel.localeCompare(bLabel);
    });

    return (
        <div className="translucent-snow p-2 rounded-lg w-full flex flex-1 flex-col overflow-y-auto">
            <input
                className={`block w-full rounded-lg border-0 outline-none text-dark-purple mb-2 shadow-none ring-1 p-2 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-majorelle-blue sm:text-sm sm:leading-6`}
                type="text"
                onChange={(e) => props.onPrincipalFilterChange && props.onPrincipalFilterChange(e.target.value)}
                value={props.principalFilter || ''}
                placeholder="Filter"
            />

            {sortedPrincipals.map(([id, principal]) => {
                const label = lastComponent(principal.name || id);
                const isSelected = selection?.type === 'principal' && selection.id === id;
                const extraClasses = isSelected ? 'bg-majorelle-blue/80 text-snow' : 'hover:bg-white/80';

                return (
                    <div
                        key={id}
                        className={`cursor-pointer flex gap-2 items-center p-2 rounded-md ${extraClasses}`}
                        onClick={() => {
                            setSelectionString(isSelected ? null : `principal:${id}`);
                        }}
                        onMouseEnter={() => {
                            setHoveredPrincipalId(id);
                            props.onPrincipalHover?.(id);
                        }}
                        onMouseLeave={() => {
                            if (hoveredPrincipalId === id) {
                                setHoveredPrincipalId(null);
                                props.onPrincipalHover?.(null);
                            }
                        }}
                    >
                        <PrincipalIcon
                            className="flex-none h-[2rem] w-[2rem] border-1 border-platinum rounded-sm"
                            type={principal.type}
                        />
                        <div className="flex flex-col justify-center min-w-0">
                            <div className="text-sm font-bold truncate">{label}</div>
                            {label !== principal.name && <div className="text-xs truncate">{principal.name}</div>}
                        </div>
                    </div>
                );
            })}
        </div>
    );
};
