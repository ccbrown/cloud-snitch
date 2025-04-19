import React, { useId } from 'react';
import { CheckIcon } from '@heroicons/react/24/outline';

interface Props {
    disabled?: boolean;
    checked: boolean;
    className?: string;
    children?: React.ReactNode;
    label?: string;
    subLabel?: string | React.ReactNode;
    onChange: (checked: boolean) => void;
}

export const Checkbox = (props: Props) => {
    const id = useId();

    return (
        <div className={props.className}>
            <div className="flex items-center my-2 gap-2 cursor-pointer">
                <div className="flex items-center relative">
                    <input
                        type="checkbox"
                        id={id}
                        disabled={props.disabled}
                        checked={props.checked}
                        onChange={(e) => props.onChange(e.target.checked)}
                        className={`w-5 h-5 appearance-none rounded-md ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-majorelle-blue ${props.disabled ? 'checked:bg-gray-300' : 'checked:bg-majorelle-blue'} cursor-pointer`}
                    />
                    {props.checked && <CheckIcon className="pointer-events-none text-white absolute inset-0" />}
                </div>
                {props.children && (
                    <label htmlFor={id} className="cursor-pointer">
                        {props.children}
                    </label>
                )}
                {props.label && (
                    <div className="flex flex-col">
                        <label htmlFor={id} className="text-sm label cursor-pointer">
                            {props.label}
                        </label>
                        {props.subLabel && (
                            <label htmlFor={id} className="text-xs text-english-violet cursor-pointer">
                                {props.subLabel}
                            </label>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};
