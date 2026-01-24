export default [
  {
    ignores: ["node_modules/", "dist/", ".quasar/", "src-*/**", "coverage/", "*.d.ts", "*.config.*"],
  },
  {
    files: ["**/*.{js,mjs,cjs,ts,vue}"],
    rules: {
      // Basic rules only to prevent build failures
      "no-unused-vars": "warn",
      "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
      "prefer-promise-reject-errors": "off",
    },
  },
];
