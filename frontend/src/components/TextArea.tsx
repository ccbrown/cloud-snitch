interface Props {
    disabled?: boolean;
    label?: string;
    placeholder?: string;
    required?: boolean;
    onChange?: (value: string) => void;
    rows?: number;
    value: string;
}

export const TextArea = (props: Props) => {
    return (
        <div className="my-2 grow">
            <div className="flex text-sm">
                {props.label && (
                    <label className="block leading-6 text-english-violet font-semibold">{props.label}</label>
                )}
            </div>
            <textarea
                className={`block w-full rounded-md border-0 outline-none text-dark-purple shadow-none ring-1 p-2 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-majorelle-blue sm:text-sm sm:leading-6 ${
                    props.disabled && 'bg-gray-200'
                }`}
                disabled={props.disabled}
                required={props.required}
                onChange={(e) => props.onChange && props.onChange(e.target.value)}
                rows={props.rows || 4}
                value={props.value}
                placeholder={props.placeholder}
            />
        </div>
    );
};
