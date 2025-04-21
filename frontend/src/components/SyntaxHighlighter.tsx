import { Light as SyntaxHighlighterImpl } from 'react-syntax-highlighter';
import bash from 'react-syntax-highlighter/dist/esm/languages/hljs/bash';
import json from 'react-syntax-highlighter/dist/esm/languages/hljs/json';
import style from 'react-syntax-highlighter/dist/esm/styles/hljs/github';

SyntaxHighlighterImpl.registerLanguage('bash', bash);
SyntaxHighlighterImpl.registerLanguage('json', json);

interface Props {
    language: 'bash' | 'json';
    children: string | string[];
    className?: string;
}

export const SyntaxHighlighter = ({ className, language, children }: Props) => {
    return (
        <SyntaxHighlighterImpl
            language={language}
            style={style}
            className={`codeblock ${className || ''}`}
            customStyle={{ backgroundColor: 'transparent' }}
        >
            {children}
        </SyntaxHighlighterImpl>
    );
};
