'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

import { useSelector } from '@/store';

interface Props {
    children: React.ReactNode;
}

export const RequireAuth = (props: Props) => {
    const hasAuth = useSelector((state) => !!state.api.auth);
    const router = useRouter();

    useEffect(() => {
        if (!hasAuth) {
            router.push('/sign-in');
        }
    }, [hasAuth, router]);

    return props.children;
};
