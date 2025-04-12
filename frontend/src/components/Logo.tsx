interface Props {
    className?: string;
}

// https://react-svgr.com/playground/
export const Logo = (props: Props) => (
    <svg
        xmlns="http://www.w3.org/2000/svg"
        xmlSpace="preserve"
        style={{
            fillRule: 'evenodd',
            clipRule: 'evenodd',
            strokeLinejoin: 'round',
            strokeMiterlimit: 2,
        }}
        viewBox="0 0 277 284"
        className={props.className}
    >
        <path
            d="M4305.68 2534.48c.38 3.45.58 6.95.58 10.5 0 51.58-41.88 93.46-93.47 93.46-51.58 0-93.46-41.88-93.46-93.46 0-3.55.19-7.05.58-10.5 5.22 46.66 44.85 82.98 92.88 82.98 48.04 0 87.66-36.32 92.89-82.98Z"
            style={{
                fill: '#5e4c63',
            }}
            transform="translate(-6088.295 -3615.62) scale(1.47798)"
        />
        <path
            d="M4180.3 2632.63c-35.59-13.21-60.97-47.49-60.97-87.65 0-51.59 41.88-93.47 93.46-93.47 51.59 0 93.47 41.88 93.47 93.47 0 40.16-25.39 74.44-60.97 87.65 12.48-9.7 20.52-24.86 20.52-41.88 0-29.26-23.76-53.02-53.02-53.02s-53.01 23.76-53.01 53.02c0 17.02 8.04 32.18 20.52 41.88Z"
            style={{
                fill: '#9a72d0',
            }}
            transform="translate(-5036.075 -3010.986) scale(1.22821)"
        />
    </svg>
);
