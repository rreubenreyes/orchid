"use strict";

module.exports = {
    extends: ['plugin:@typescript-eslint/recommended'],
    parser: '@typescript-eslint/parser',
    parserOptions: {
        project: './tsconfig.json'
    },
    rules: {
        indent: ['error', 4, { SwitchCase: 1, VariableDeclarator: 1, MemberExpression: 'off' }],
        '@typescript-eslint/indent': ['error', 4],
    }
};
