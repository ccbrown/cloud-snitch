'use client';

import React, { useRef, useState } from 'react';
import { createPortal } from 'react-dom';

interface Props {
    children: React.ReactNode;
    content: React.ReactNode;
}

export const Tooltip = ({ children, content }: Props) => {
    const ref = useRef<HTMLSpanElement>(null);
    const [rect, setRect] = useState<DOMRect | undefined>();
    const [isHovered, setIsHovered] = useState(false);

    return (
        <span
            onMouseEnter={() => {
                const rect = ref.current?.getBoundingClientRect();
                setRect(rect);
                setIsHovered(true);
            }}
            onMouseLeave={() => setIsHovered(false)}
            ref={ref}
        >
            {isHovered &&
                createPortal(
                    <div
                        className="absolute z-1000 translucent-snow text-dark-purple p-2 text-sm rounded-lg border-1 border-platinum -translate-x-1/2 -translate-y-full"
                        style={{
                            left: (rect?.left || 0) + (rect?.width || 0) / 2,
                            top: rect?.top,
                        }}
                    >
                        {content}
                    </div>,
                    document.body,
                )}
            {children}
        </span>
    );
};
