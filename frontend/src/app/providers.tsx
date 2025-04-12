'use client';

import { useServerInsertedHTML } from 'next/navigation';
import { PropsWithChildren } from 'react';
import { Provider as ReduxProvider } from 'react-redux';

import { store } from '@/store';

export const Providers = ({ children }: PropsWithChildren) => {
    useServerInsertedHTML(() => {
        return <></>;
    });

    return <ReduxProvider store={store}>{children}</ReduxProvider>;
};
