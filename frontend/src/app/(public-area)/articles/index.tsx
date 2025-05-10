import CapitalOneDataBreach from './capital-one-data-breach';
import AwsModifyOnlyCreatedResource from './how-to-allow-aws-principals-to-modify-only-resources-they-create';
import ListAwsRegionsAndServices from './how-to-programmatically-get-a-list-of-all-aws-regions-and-services';

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
    'how-to-programmatically-get-a-list-of-all-aws-regions-and-services': ListAwsRegionsAndServices,
};
