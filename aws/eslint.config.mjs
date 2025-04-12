import { dirname } from 'path';
import { fileURLToPath } from 'url';
import { FlatCompat } from '@eslint/eslintrc';
import prettier from 'eslint-plugin-prettier';
import unusedImports from 'eslint-plugin-unused-imports';
import tseslint from 'typescript-eslint';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
    baseDirectory: __dirname,
});

const eslintConfig = [
    ...compat.extends('prettier', 'plugin:@typescript-eslint/recommended'),
    {
        ignores: ['cdk.out', '**/*.js'],
    },
    {
        plugins: {
            prettier: prettier,
            'unused-imports': unusedImports,
            '@typescript-eslint': tseslint.plugin,
        },
    },
    {
        rules: {
            'unused-imports/no-unused-imports': 'error',
            '@typescript-eslint/no-unused-vars': [
                'error',
                {
                    args: 'all',
                    argsIgnorePattern: '^_',
                    caughtErrors: 'all',
                    caughtErrorsIgnorePattern: '^_',
                    destructuredArrayIgnorePattern: '^_',
                    varsIgnorePattern: '^_',
                    ignoreRestSiblings: true,
                },
            ],
            'prettier/prettier': 1,
        },
    },
];

export default eslintConfig;
