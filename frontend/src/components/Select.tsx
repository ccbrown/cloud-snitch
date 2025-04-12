import { Field, Label, Select as HeadlessSelect } from '@headlessui/react';

interface Option {
    label: string;
    value: string;
}

interface Props {
    disabled?: boolean;
    label?: string;
    options: Option[];
    value?: string;
    onChange?: (value: string) => void;
}

export const Select = (props: Props) => {
    return (
        <Field className="my-2 grow">
            <div className="flex text-sm">
                {props.label && (
                    <Label className="block leading-6 text-english-violet font-semibold">{props.label}</Label>
                )}
            </div>
            <HeadlessSelect
                disabled={props.disabled}
                className={`block w-full rounded-md appearance-none border-0 outline-none text-dark-purple shadow-none ring-1 p-2 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-majorelle-blue sm:text-sm sm:leading-6 ${
                    props.disabled && 'bg-gray-200'
                }`}
                value={props.value}
                onChange={(e) => {
                    if (props.onChange) {
                        props.onChange(e.target.value);
                    }
                }}
            >
                {props.options.map((option) => (
                    <option key={option.value} value={option.value}>
                        {option.label}
                    </option>
                ))}
            </HeadlessSelect>
        </Field>
    );
};
