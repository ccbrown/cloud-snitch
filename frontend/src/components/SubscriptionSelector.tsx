import { IndividualSubscriptionBox, TeamSubscriptionBox } from '@/components';
import { TeamSubscriptionTier } from '@/generated/api';

interface ButtonProps {
    className?: string;
    disabled?: boolean;
    onClick?: () => void;
}

const Button = (props: ButtonProps) => {
    return (
        <button
            className={`snow-button mt-4 ${props.disabled ? 'bg-gray-200' : ''}`}
            disabled={props.disabled}
            onClick={(e) => {
                e.preventDefault();
                if (props.onClick) {
                    props.onClick();
                }
            }}
        >
            Select
        </button>
    );
};

interface Props {
    disabled?: boolean;
    onSelect: (subscription: TeamSubscriptionTier) => void;
}

export const SubscriptionSelector = (props: Props) => (
    <div className="flex gap-2">
        <IndividualSubscriptionBox
            compact
            footer={
                <Button disabled={props.disabled} onClick={() => props.onSelect(TeamSubscriptionTier.Individual)} />
            }
        />
        <TeamSubscriptionBox
            compact
            footer={<Button disabled={props.disabled} onClick={() => props.onSelect(TeamSubscriptionTier.Team)} />}
        />
    </div>
);
