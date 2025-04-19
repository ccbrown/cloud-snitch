import { Field, Label } from '@headlessui/react';

interface Props {
    disabled?: boolean;
    label?: string;
    placeholder?: string;
    type?: 'text' | 'email' | 'password' | 'search';
    autocomplete?: 'email' | 'current-password' | 'new-password';
    required?: boolean;
    onChange?: (value: string) => void;
    value: string;
}

export const TextField = (props: Props) => {
    return (
        <Field className="grow">
            <div className="flex text-sm">
                {props.label && <Label className="block leading-6 label">{props.label}</Label>}
            </div>
            <input
                className={`block w-full rounded-md border-0 outline-none text-dark-purple mt-1 shadow-none ring-1 p-2 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-majorelle-blue sm:text-sm sm:leading-6 ${
                    props.disabled && 'bg-gray-200'
                }`}
                disabled={props.disabled}
                type={props.type || 'text'}
                required={props.required}
                autoComplete={props.autocomplete}
                onChange={(e) => props.onChange && props.onChange(e.target.value)}
                value={props.value}
                placeholder={props.placeholder}
            />
        </Field>
    );
};
