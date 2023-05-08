module.exports = {
    extends: [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:jest/recommended",
        "plugin:jest/style",
        "plugin:testing-library/react",
        "plugin:jest-dom/recommended",
        "plugin:react-hooks/recommended",
    ],
    plugins: ["react", "jest", "testing-library", "jest-dom"],
    env: {
        browser: true,
        es2021: true,
        node: true,
        commonjs: true,
        node: true,
        jest: true,
    },
    parserOptions: {
        ecmaFeatures: {
            jsx: true,
        },
        ecmaVersion: 12,
        sourceType: "module",
    },
    rules: {
        "no-unused-vars": "error",
        "jest/no-disabled-tests": "warn",
        "jest/no-focused-tests": "error",
        "jest/no-identical-title": "error",
        "jest/prefer-to-have-length": "warn",
        "jest/valid-expect": "error",
        "react/jsx-no-undef": [2, { allowGlobals: true }],
        "react/react-in-jsx-scope": "off",
    },
    settings: {
        react: {
            version: "detect",
        },
    },
};
