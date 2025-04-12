import Image from 'next/image';

import AwsIcon from '@/assets/aws/Architecture-Group-Icons_02072025/AWS-Cloud-logo_32.svg';
import AwsAccountIcon from '@/assets/aws/Architecture-Group-Icons_02072025/AWS-Account_32.svg';
import AwsRoleIcon from '@/assets/aws/Resource-Icons_02072025/Res_Security-Identity-Compliance/Res_AWS-Identity-Access-Management_Role_48.svg';
import AwsIamUserIcon from '@/assets/aws/Resource-Icons_02072025/Res_General-Icons/Res_48_Dark/Res_User_48_Dark.svg';

import { formatPrincipalType, PrincipalType } from '@/report';

interface Props {
    className?: string;
    type?: PrincipalType;
}

export const PrincipalIcon = ({ className, type }: Props) => {
    const alt = type ? formatPrincipalType(type) : 'Unknown';
    switch (type) {
        case 'AWSAccount':
            return <Image src={AwsAccountIcon} className={className} alt={alt} />;
        case 'AWSRole':
        case 'AWSAssumedRole':
            return (
                <div className={`bg-white p-1 ${className}`}>
                    <Image src={AwsRoleIcon} alt={alt} className="w-full h-full" />
                </div>
            );
        case 'AWSIAMUser':
            return (
                <div className={`bg-gray-400 p-1 ${className}`}>
                    <Image src={AwsIamUserIcon} alt={alt} className="w-full h-full" />
                </div>
            );
        case 'AWSService':
            return <Image src={AwsIcon} alt={alt} className={className} />;
        default:
            return <div className={`bg-gray-600 p-1 ${className}`} />;
    }
};
