'use client';

import { PlusCircleIcon, PlusIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { useMemo, useState } from 'react';
import { Combobox, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/react';

import { Tooltip } from '@/components';

interface Option {
    label: string;
    altLabel?: string;
    value: string;
}

interface Props {
    options: Option[];
    before: Set<string>;
    after: Set<string>;
    onAdd: (value: string) => void;
    onRemove: (value: string) => void;
}

export const ChipEditor = (props: Props) => {
    const [isAdding, setIsAdding] = useState(false);

    const visibleOptions = props.options.filter(
        (option) => props.before.has(option.value) || props.after.has(option.value),
    );

    // Add options that are in before or after but not in the current options.
    {
        const optionValues = new Set(props.options.map((option) => option.value));
        for (const value of props.before) {
            if (!optionValues.has(value)) {
                visibleOptions.push({ label: value, value });
            }
        }
        for (const value of props.after) {
            if (!optionValues.has(value)) {
                visibleOptions.push({ label: value, value });
            }
        }
    }

    visibleOptions.sort((a, b) => a.label.localeCompare(b.label));

    const neutralChip = 'chip px-1';
    const addedChip = 'chip px-1 bg-mint';
    const removedChip = 'chip px-1 line-through bg-indian-red';

    return (
        <span>
            {visibleOptions.map((option) => {
                const isAdded = props.after.has(option.value) && !props.before.has(option.value);
                const isRemoved = props.before.has(option.value) && !props.after.has(option.value);

                return (
                    <Tooltip
                        key={option.value}
                        disabled={!option.altLabel && !isAdded && !isRemoved}
                        content={
                            <div className="text-xs">
                                {option.altLabel && <div>{option.altLabel}</div>}
                                {isAdded && <div className="text-mint">Added</div>}
                                {isRemoved && <div className="text-indian-red">Removed</div>}
                            </div>
                        }
                    >
                        <span className={`${isAdded ? addedChip : isRemoved ? removedChip : neutralChip}`}>
                            {option.label}
                            {props.after.has(option.value) ? (
                                <XMarkIcon
                                    className="h-[0.8rem] inline-block ml-1 cursor-pointer"
                                    onClick={() => props.onRemove(option.value)}
                                />
                            ) : (
                                <PlusIcon
                                    className="h-[0.8rem] inline-block ml-1 cursor-pointer"
                                    onClick={() => props.onAdd(option.value)}
                                />
                            )}
                        </span>
                    </Tooltip>
                );
            })}
            {isAdding ? (
                <InlineCombobox
                    options={props.options.filter(
                        (option) => !props.before.has(option.value) && !props.after.has(option.value),
                    )}
                    onChange={(value) => {
                        if (value) {
                            props.onAdd(value);
                            setIsAdding(false);
                        }
                    }}
                    onClose={() => setIsAdding(false)}
                />
            ) : (
                <PlusCircleIcon
                    className="h-[1rem] inline-block ml-1 cursor-pointer"
                    onClick={() => setIsAdding(true)}
                />
            )}
        </span>
    );
};

interface InlineComboboxProps {
    options: Option[];
    onChange?: (value: string) => void;
    onClose?: () => void;
}

const InlineCombobox = ({ options, onChange, onClose }: InlineComboboxProps) => {
    const [query, setQuery] = useState('');

    const filteredOptions = useMemo(
        () =>
            options
                .filter(
                    (option) =>
                        option.label.toLowerCase().includes(query.toLowerCase()) ||
                        option.altLabel?.toLowerCase().includes(query.toLowerCase()),
                )
                .sort((a, b) => a.label.localeCompare(b.label)),
        [options, query],
    );

    return (
        <Combobox
            virtual={{ options: filteredOptions }}
            onChange={(option: Option) => onChange?.(option.value)}
            onClose={() => {
                setQuery('');
                onClose?.();
            }}
            immediate
        >
            <ComboboxInput
                autoFocus
                onChange={(event) => setQuery(event.target.value)}
                onBlur={() => {
                    setQuery('');
                    onClose?.();
                }}
                className="bg-english-violet px-2 py-0.5 mx-0.5 leading-none rounded-md text-xs font-semibold text-snow focus:outline-none"
            />
            <ComboboxOptions
                anchor="bottom"
                className="empty:invisible rounded-md text-snow bg-english-violet/80 backdrop-blur-md p-1 w-3xs"
                static
            >
                {({ option }) => (
                    <ComboboxOption
                        value={option}
                        className="cursor-pointer text-xs px-2 py-1 hover:bg-white/20 rounded-md w-full flex flex-col"
                    >
                        <div className="font-semibold">{option.label}</div>
                        {option.altLabel && <div>{option.altLabel}</div>}
                    </ComboboxOption>
                )}
            </ComboboxOptions>
        </Combobox>
    );
};
