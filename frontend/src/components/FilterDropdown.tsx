import { ChevronDownIcon } from '@heroicons/react/24/outline';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/react';

import { Checkbox } from '@/components';

interface FilterOption {
    label: string;
    subLabel?: string;
    value: string;
}

type Mode = 'include' | 'exclude';

export interface Filter {
    mode: Mode;
    values: Set<string>;
}

const MODES: Mode[] = ['include', 'exclude'];

interface Props {
    label: string;
    options: FilterOption[];
    onChange: (filter: Filter) => void;
    filter: Filter;
}

export const FilterDropdown = ({ label, options, onChange, filter }: Props) => {
    const singularOrPluralLabel = filter.values.size === 1 ? label : label + 's';

    return (
        <Popover>
            <PopoverButton className="inline-flex items-center gap-2 rounded-md bg-english-violet py-1.5 px-3 text-sm/6 font-semibold text-snow cursor-pointer focus:outline-none data-[focus]:outline-1 data-[focus]:outline-majorelle-blue">
                {filter.mode === 'exclude'
                    ? filter.values.size === 0
                        ? `All ${singularOrPluralLabel}`
                        : `All But ${filter.values.size} ${singularOrPluralLabel}`
                    : `${filter.values.size} ${singularOrPluralLabel}`}
                <ChevronDownIcon className="h-[1rem]" />
            </PopoverButton>

            <PopoverPanel
                anchor="bottom"
                className="flex flex-col translucent-snow rounded-lg border-1 border-platinum mt-1"
            >
                <div className="flex gap-2 p-2 border-b border-platinum">
                    {MODES.map((mode) => (
                        <button
                            key={mode}
                            className={`flex-1 px-6 py-1 rounded-lg text-sm text-center font-semibold cursor-pointer ${
                                filter.mode === mode
                                    ? 'bg-amethyst text-snow'
                                    : 'hover:outline-1 hover:outline-platinum'
                            }`}
                            onClick={() => onChange({ ...filter, mode })}
                        >
                            {mode.charAt(0).toUpperCase() + mode.slice(1)}
                        </button>
                    ))}
                </div>
                <div className="flex gap-2 p-2 border-b border-platinum">
                    {filter.values.size > 0 ? (
                        <button
                            className="flex-1 px-6 py-1 rounded-lg button text-sm"
                            onClick={() => onChange({ ...filter, values: new Set() })}
                        >
                            Clear All
                        </button>
                    ) : (
                        <button
                            className="flex-1 px-6 py-1 rounded-lg button text-sm"
                            onClick={() =>
                                onChange({ ...filter, values: new Set(options.map((option) => option.value)) })
                            }
                        >
                            Select All
                        </button>
                    )}
                </div>
                <div className="flex flex-col p-2 max-h-[80vh] overflow-auto">
                    {options.map((option) => {
                        const toggle = () => {
                            const newValues = filter.values.has(option.value)
                                ? filter.values.difference(new Set([option.value]))
                                : filter.values.union(new Set([option.value]));
                            onChange({ ...filter, values: newValues });
                        };
                        return (
                            <div
                                key={option.value}
                                className="flex items-center gap-2 pl-2 pr-6 hover:bg-white/80 cursor-pointer rounded-md"
                                onClick={toggle}
                            >
                                <Checkbox
                                    checked={filter.values.has(option.value)}
                                    onChange={toggle}
                                    label={option.label}
                                    subLabel={option.subLabel}
                                />
                            </div>
                        );
                    })}
                </div>
            </PopoverPanel>
        </Popover>
    );
};
