import { ChevronRightIcon } from '@heroicons/react/24/outline';
import type { Metadata } from 'next';
import Link from 'next/link';

import { articles } from '.';

export const metadata: Metadata = {
    title: 'Articles',
    description: 'Articles and resources related to AWS, InfoSec, and Cloud Snitch.',
};

const Page = () => {
    const sortedArticles = Object.entries(articles).sort(([, a], [, b]) => b.date.getTime() - a.date.getTime());

    return (
        <div className="flex flex-col gap-4 [&_p]:my-4">
            <div className="translucent-snow p-4 rounded-lg">
                <h1>Articles</h1>
                <div className="flex flex-col gap-4 mt-4">
                    {sortedArticles.map(([slug, article]) => (
                        <div className="border-1 border-platinum rounded-lg p-4" key={slug}>
                            <h2>{article.title}</h2>
                            <div className="text-sm text-english-violet">
                                {article.date.toLocaleDateString('en-US', {
                                    weekday: 'long',
                                    year: 'numeric',
                                    month: 'long',
                                    day: 'numeric',
                                })}
                            </div>
                            <p>{article.description}</p>
                            <div className="flex">
                                <Link
                                    href={`/articles/${slug}`}
                                    className="button flex grow-0 items-center whitespace-nowrap"
                                >
                                    Read More
                                    <ChevronRightIcon className="h-4 w-4 ml-2" />
                                </Link>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Page;
