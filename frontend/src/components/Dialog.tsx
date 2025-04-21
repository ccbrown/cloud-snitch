import { Dialog as DialogImpl, DialogPanel, DialogTitle } from '@headlessui/react';
import React from 'react';

interface Props {
    children: React.ReactNode;
    isOpen?: boolean;
    onClose?: () => void;
    title?: string;
    size?: 'xl';
}

export const Dialog = (props: Props) => {
    const width = props.size === 'xl' ? 'w-xl' : 'max-w-xl';

    return (
        <DialogImpl open={!!props.isOpen} onClose={() => props.onClose && props.onClose()} className="relative z-50">
            <div className="fixed inset-0 flex w-screen items-center justify-center p-4 bg-radial from-black/10 to-black/40">
                <DialogPanel
                    className={`min-w-lg ${width} space-y-4 translucent-snow border border-platinum rounded-xl p-8`}
                >
                    {props.title && <DialogTitle className="font-bold">{props.title}</DialogTitle>}
                    {props.children}
                </DialogPanel>
            </div>
        </DialogImpl>
    );
};
