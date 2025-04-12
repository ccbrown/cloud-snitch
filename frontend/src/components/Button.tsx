interface Props {
    children?: React.ReactNode;
    className?: string;
    disabled?: boolean;
    label?: string;
    onClick?: () => void;
    style?: 'snow' | 'subtle';
    type?: 'submit';
}

export const Button = (props: Props) => {
    return (
        <button
            className={`${props.style === 'snow' ? 'snow-button' : props.style === 'subtle' ? 'subtle-button' : 'button'} ${props.disabled ? 'bg-gray-600' : ''} ${props.className || ''}`}
            disabled={props.disabled}
            onClick={(e) => {
                e.preventDefault();
                if (props.onClick) {
                    props.onClick();
                }
            }}
            type={props.type}
        >
            {props.label || props.children}
        </button>
    );
};
