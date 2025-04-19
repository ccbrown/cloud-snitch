import CapitalOneDataBreach from './capital-one-data-breach';

interface Article {
    title: string;
    description: string;
    date: Date;
    content: React.ReactNode;
    relatedLinks?: Array<{ title: string; url: string }>;
}

export const articles: Record<string, Article> = {
    'capital-one-data-breach': CapitalOneDataBreach,
};
