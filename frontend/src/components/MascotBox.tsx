import Image from 'next/image';

interface Props {
    children?: React.ReactNode;
}

export const MascotBox = (props: Props) => {
    return (
        <div className="translucent-snow rounded-xl flex flex-col lg:flex-row items-stretch">
            <div className="grow p-8">{props.children}</div>
            <div className="flex items-end justify-end min-w-1/3">
                <Image src="/images/mascot-512.png" alt="Watcher" width={512} height={512} />
            </div>
        </div>
    );
};
