import { ChevronDownIcon } from '@heroicons/react/24/outline';
import { CloseButton, Popover, PopoverButton, PopoverPanel } from '@headlessui/react';

import { formatDurationSeconds, formatTimeRange, SECONDS_PER_DAY, SECONDS_PER_WEEK } from '@/time';

interface Props {
    availableStartTime?: Date;
    availableEndTime?: Date;

    durationSeconds?: number;

    onChange: (durationSeconds?: number) => void;
}

interface Option {
    label: string;
    description?: string;
    time?: Date;
    durationSeconds?: number;
}

export const DurationDropdown = (props: Props) => {
    const options: Option[] = [];

    if (props.availableStartTime && props.availableEndTime) {
        const availableDurationSeconds = (props.availableEndTime.getTime() - props.availableStartTime.getTime()) / 1000;

        for (const { label, durationSeconds } of [
            { label: '1 Day', durationSeconds: SECONDS_PER_DAY },
            { label: '3 Days', durationSeconds: 3 * SECONDS_PER_DAY },
            { label: '1 Week', durationSeconds: SECONDS_PER_WEEK },
        ]) {
            if (availableDurationSeconds > durationSeconds) {
                const time = new Date(props.availableEndTime.getTime() - durationSeconds * 1000);
                options.push({
                    label,
                    time,
                    durationSeconds,
                    description: formatTimeRange(time, (durationSeconds * 1000) / 1000),
                });
            }
        }
    }

    options.push({
        label: 'All Time',
        description:
            props.availableStartTime &&
            props.availableEndTime &&
            formatTimeRange(
                props.availableStartTime,
                (props.availableEndTime.getTime() - props.availableStartTime.getTime()) / 1000,
            ),
    });

    return (
        <Popover>
            <PopoverButton className="inline-flex items-center gap-2 rounded-md bg-english-violet py-1.5 px-3 text-sm/6 font-semibold text-snow cursor-pointer focus:outline-none data-[focus]:outline-1 data-[focus]:outline-majorelle-blue">
                {options.find((opt) => opt.durationSeconds === props.durationSeconds)?.label ||
                    formatDurationSeconds(props.durationSeconds || 0)}
                <ChevronDownIcon className="h-[1rem]" />
            </PopoverButton>

            <PopoverPanel
                anchor="bottom"
                className="flex flex-col translucent-snow rounded-lg border-1 border-platinum mt-1"
            >
                <div className="flex flex-col p-2 max-h-[80vh] overflow-auto">
                    {options.map((option) => (
                        <CloseButton
                            key={option.label}
                            className="flex flex-col justify-center p-2 pr-6 hover:bg-white/80 cursor-pointer rounded-md"
                            onClick={() => {
                                props.onChange(option.durationSeconds);
                            }}
                        >
                            <span>{option.label}</span>
                            {option.description && (
                                <span className="text-xs text-english-violet">{option.description}</span>
                            )}
                        </CloseButton>
                    ))}
                </div>
            </PopoverPanel>
        </Popover>
    );
};
