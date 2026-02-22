/* eslint-env node */
module.exports = {
  root: true,
  env: { node: true, es2020: true },
  parserOptions: { ecmaVersion: 2020, sourceType: "module" },
  extends: ["eslint:recommended"],
  ignorePatterns: ["dist/", "node_modules/"],
};
