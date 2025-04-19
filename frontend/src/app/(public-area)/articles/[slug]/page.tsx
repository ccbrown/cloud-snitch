import type { Metadata } from 'next';
import Link from 'next/link';

import { articles } from '..';

export const dynamicParams = false;

export async function generateStaticParams() {
    return Object.keys(articles).map((slug) => ({
        slug,
    }));
}

type Props = { params: Promise<{ slug: string }> };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
    const { slug } = await params;
    const article = articles[slug];

    return {
        title: article.title,
        description: article.description,
        openGraph: {
            title: article.title,
            description: article.description,
            type: 'article',
            publishedTime: article.date.toISOString(),
        },
    };
}

export default async function Page({ params }: Props) {
    const { slug } = await params;
    const article = articles[slug];

    return (
        <div className="flex flex-col gap-4">
            <div className="translucent-snow p-4 rounded-lg [&_p]:my-4">
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
                {article.content}
            </div>

            {article.relatedLinks && (
                <div className="translucent-snow p-4 rounded-lg">
                    <h2>Further Reading</h2>
                    <ul className="list-disc mt-4 pl-6">
                        {article.relatedLinks.map((item, index) => (
                            <li key={index}>
                                <Link href={item.url} target="_blank" className="external-link">
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
