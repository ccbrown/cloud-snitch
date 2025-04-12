import { Tooltip } from '@/components';

interface FeatureProps {
    name: string;
    tooltip?: string;
}

interface Props {
    className?: string;
    compact?: boolean;
    name: string;
    description: string;
    features?: FeatureProps[];
    footer?: React.ReactNode;
}

export const SubscriptionBox = (props: Props) => (
    <div className={`flex-1 flex flex-col outline-1 outline-white/20 text-snow p-4 rounded-lg ${props.className}`}>
        <h2>{props.name}</h2>
        <p className={`mb-4 ${props.compact ? 'text-sm' : ''}`}>{props.description}</p>
        <div className="grow" />
        {props.features && (
            <ul className={`list-disc pl-6 ${props.compact ? 'text-sm' : ''}`}>
                {props.features.map((feature) => (
                    <li key={feature.name}>
                        {feature.tooltip ? (
                            <Tooltip content={<div className="w-sm">{feature.tooltip}</div>}>
                                <span className="hoverable">{feature.name}</span>
                            </Tooltip>
                        ) : (
                            feature.name
                        )}
                    </li>
                ))}
            </ul>
        )}
        {props.footer}
    </div>
);

interface PredefinedProps {
    footer?: React.ReactNode;
    compact?: boolean;
}

export const IndividualSubscriptionBox = (props: PredefinedProps) => (
    <SubscriptionBox
        className="bg-amethyst"
        compact={props.compact}
        name="Individual"
        description="This tier is intended for hobbyists or indepedent developers. We've brought the price down as low as it can go so Cloud Snitch is affordable for any budget!"
        features={[
            { name: '1 Member Per Team' },
            { name: '1 Week Retention' },
            {
                name: '20MB CloudTrail Events Per Day',
                tooltip:
                    "Limit is per region, per AWS account. This limit exists to ensure quality of service. It's exceedingly rare for personal accounts to approach this.",
            },
            { name: 'Best Effort Support' },
        ]}
        footer={
            <>
                <div className="flex flex-col items-center mt-4">
                    <div className="text-4xl font-extrabold">$0.99</div>
                    <div className="text-xs">Per AWS Account/Month</div>
                </div>
                {props.footer}
            </>
        }
    />
);

export const TeamSubscriptionBox = (props: PredefinedProps) => (
    <SubscriptionBox
        className="bg-english-violet"
        compact={props.compact}
        name="Team"
        description="For professionals, this tier provides the features you need to collaborate with your team while still costing less than a small EC2 instance!"
        features={[
            { name: 'Unlimited Members Per Team' },
            { name: '2 Week Retention' },
            {
                name: '200MB CloudTrail Events Per Day',
                tooltip:
                    "Limit is per region, per AWS account. This limit exists to ensure quality of service. It's rare for accounts to approach this.",
            },
            { name: 'Priority Support' },
        ]}
        footer={
            <>
                <div className="flex flex-col items-center mt-4">
                    <div className="text-4xl font-extrabold">$9.99</div>
                    <div className="text-xs">Per AWS Account/Month</div>
                </div>
                {props.footer}
            </>
        }
    />
);
