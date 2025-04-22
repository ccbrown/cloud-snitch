import { ArrowUturnLeftIcon } from '@heroicons/react/24/outline';
import type { Metadata, ResolvingMetadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';

import { articles } from '..';

export const dynamicParams = false;

export async function generateStaticParams() {
    return Object.keys(articles).map((slug) => ({
        slug,
    }));
}

type Props = { params: Promise<{ slug: string }> };

export async function generateMetadata({ params }: Props, parent: ResolvingMetadata): Promise<Metadata> {
    const { slug } = await params;
    const article = articles[slug];
    const previousImages = (await parent).openGraph?.images || [];

    return {
        title: article.title,
        description: article.description,
        openGraph: {
            title: article.title,
            description: article.description,
            type: 'article',
            publishedTime: article.date.toISOString(),
            images: previousImages,
        },
    };
}

export default async function Page({ params }: Props) {
    const { slug } = await params;
    const article = articles[slug];

    return (
        <div className="flex flex-col gap-4">
            <div className="translucent-snow p-4 rounded-lg flex flex-col gap-4">
                <div className="text-sm">
                    <Link href="/articles" className="link flex gap-1 items-center">
                        More Articles
                        <ArrowUturnLeftIcon className="h-[1rem]" />
                    </Link>
                </div>
                <div className="flex gap-4 items-center">
                    <Image
                        src={article.author.image}
                        alt={`By ${article.author.name}`}
                        height={48}
                        width={48}
                        className="rounded-full h-[48px] w-[48px]"
                    />
                    <div>
                        <h1>{article.title}</h1>
                        <div className="text-sm text-english-violet">
                            Published{' '}
                            {article.date.toLocaleDateString('en-US', {
                                weekday: 'long',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                            })}
                        </div>
                    </div>
                </div>
                {article.content}
            </div>

            {article.relatedLinks && (
                <div className="translucent-snow p-4 rounded-lg">
                    <h2>Further Reading</h2>
                    <ul className="list-disc mt-4 pl-6">
                        {article.relatedLinks.map((item, index) => (
                            <li key={index}>
                                <Link
                                    href={item.url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="external-link"
                                >
                                    {item.title}
                                </Link>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}
