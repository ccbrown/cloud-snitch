'use client';

export default function imageLoader({ src }: { src: string }) {
    if (!src.startsWith('/')) {
        return src;
    }
    return `${process.env.NEXT_PUBLIC_CDN_URL}${src}`;
}
