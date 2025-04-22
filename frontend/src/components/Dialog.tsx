import { XMarkIcon } from '@heroicons/react/24/outline';
import { CloseButton, Dialog as DialogImpl, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react';
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
            <DialogBackdrop
                className="fixed inset-0 bg-radial from-black/10 to-black/40 ease-out duration-200 data-[closed]:opacity-0"
                transition
            />
            <div className="fixed inset-0 flex w-screen items-center justify-center p-4">
                <DialogPanel
                    className={`relative min-w-lg ${width} space-y-4 translucent-snow border border-platinum rounded-xl p-8 duration-200 ease-out data-[closed]:scale-95 data-[closed]:opacity-0`}
                    transition
                >
                    <CloseButton className="absolute right-2 top-2 cursor-pointer opacity-20 hover:opacity-100 hover:text-majorelle-blue transition-opacity duration-200 ease-in-out">
                        <XMarkIcon className="h-[1.5rem]" />
                    </CloseButton>
                    {props.title && <DialogTitle className="font-bold">{props.title}</DialogTitle>}
                    {props.children}
                </DialogPanel>
            </div>
        </DialogImpl>
    );
};
