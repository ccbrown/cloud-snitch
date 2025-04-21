import CapitalOneDataBreach from './capital-one-data-breach';
import AwsModifyOnlyCreatedResource from './how-to-allow-aws-principals-to-modify-only-resources-they-create';

interface Article {
    title: string;
    author: {
        name: string;
        image: string;
    };
    description: string;
    date: Date;
    content: React.ReactNode;
    relatedLinks?: Array<{ title: string; url: string }>;
}

export const articles: Record<string, Article> = {
    'capital-one-data-breach': CapitalOneDataBreach,
    'how-to-allow-aws-principals-to-modify-only-resources-they-create': AwsModifyOnlyCreatedResource,
};
