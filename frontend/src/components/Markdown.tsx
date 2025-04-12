import Link from 'next/link';
import ReactMarkdown from 'react-markdown';

interface Props {
    children: React.ReactNode;
}

export const Markdown = ({ children }: Props) => {
    return (
        <div className="flex flex-col gap-2">
            <ReactMarkdown
                components={{
                    a(props) {
                        const { href, title, children } = props;
                        const isAbsolute = href && (href.startsWith('http://') || href.startsWith('https://'));
                        return href ? (
                            <Link
                                href={href}
                                className={isAbsolute ? 'external-link' : 'link'}
                                target="_blank"
                                rel="noopener noreferrer"
                                title={title}
                            >
                                {children}
                            </Link>
                        ) : (
                            children
                        );
                    },
                    blockquote(props) {
                        const { children } = props;
                        return <blockquote className="border-l-2 border-platinum pl-1">{children}</blockquote>;
                    },
                    code(props) {
                        const { children, node } = props;
                        const isInline = node?.position && node.position.start.line === node.position.end.line;
                        return <code className={isInline ? 'codeinline' : ''}>{children}</code>;
                    },
                    ol(props) {
                        const { children } = props;
                        return <ul className="list-decimal ml-1 list-inside">{children}</ul>;
                    },
                    ul(props) {
                        const { children } = props;
                        return <ul className="list-disc ml-1 list-inside">{children}</ul>;
                    },
                    pre(props) {
                        const { children } = props;
                        return <pre className="codeblock">{children}</pre>;
                    },
                }}
            >
                {children as string}
            </ReactMarkdown>
        </div>
    );
};
