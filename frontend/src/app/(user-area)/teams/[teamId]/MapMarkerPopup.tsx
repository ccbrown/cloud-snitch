import { useAwsRegionsMap } from '@/hooks';
import { MarkerData } from './MapOverlays';

interface AwsRegionPopupContentProps {
    id: string;
}

export const AwsRegionPopupContent = (props: AwsRegionPopupContentProps) => {
    const regions = useAwsRegionsMap();

    return (
        <div>
            <p className="font-semibold">AWS Region</p>
            <p>
                {regions.get(props.id)?.name} ({props.id})
            </p>
        </div>
    );
};

interface NetworkPopupContentProps {
    cidr: string;
}

export const NetworkPopupContent = (props: NetworkPopupContentProps) => {
    return (
        <div>
            <p className="font-semibold">Network</p>
            <p>{props.cidr}</p>
        </div>
    );
};

interface ClusterPopupContentProps {
    awsRegionIds: Set<string>;
    networkCidrs: Set<string>;
}

export const ClusterPopupContent = (props: ClusterPopupContentProps) => {
    const regions = useAwsRegionsMap();
    const { awsRegionIds, networkCidrs } = props;

    return (
        <div>
            <p className="font-semibold">Cluster</p>
            {awsRegionIds.size > 0 &&
                Array.from(awsRegionIds).map((id) => (
                    <p key={id}>
                        {regions.get(id)?.name} ({id})
                    </p>
                ))}
            {networkCidrs.size > 2 ? (
                <p>{networkCidrs.size.toLocaleString()} Networks</p>
            ) : networkCidrs.size > 0 ? (
                Array.from(networkCidrs).map((cidr) => <p key={cidr}>{cidr}</p>)
            ) : null}
        </div>
    );
};

interface Props {
    data: MarkerData;
    statusText?: string;
}

export const MapMarkerPopup = ({ data, statusText }: Props) => {
    return (
        <div className="absolute translucent-snow rounded-md max-w-sm p-2 text-english-violet text-sm bottom-11/10 whitespace-nowrap">
            {data.type === 'aws-region' ? (
                <AwsRegionPopupContent id={data.id} />
            ) : data.type === 'network' ? (
                <NetworkPopupContent cidr={data.cidr} />
            ) : data.type === 'cluster' ? (
                <ClusterPopupContent awsRegionIds={data.awsRegionIds} networkCidrs={data.networkCidrs} />
            ) : null}
            {statusText && <p className="text-xs text-majorelle-blue">{statusText}</p>}
        </div>
    );
};
