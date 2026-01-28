import { defineConfig } from "vitest/config";
import { quasar, transformAssetUrls } from "@quasar/vite-plugin";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath } from "node:url";

export default defineConfig({
  test: {
    environment: "happy-dom",
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html", "cobertura"],
      reportsDirectory: "./coverage",
      exclude: [
        "node_modules/",
        "src/css/**",
        "src-*/**",
        "dist/**",
        "**/*.d.ts",
        "**/*.config.*",
        "**/coverage/**",
        "**/node_modules/**",
        "**/*.spec.*",
        "**/*.test.*",
      ],
    },
    globals: true,
    setupFiles: ["./test/vitest/setup.ts"],
    include: ["src/**/__tests__/**/*.{test,spec}.{js,jsx,ts,tsx,vue}"],
  },
  plugins: [
    vue({
      template: { transformAssetUrls },
    }),
    quasar({
      sassVariables: "src/css/quasar.variables.scss",
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      "src": fileURLToPath(new URL("./src", import.meta.url)),
      "stores": fileURLToPath(new URL("./src/stores", import.meta.url)),
      "app": fileURLToPath(new URL(".", import.meta.url)),
      "test": fileURLToPath(new URL("./test", import.meta.url)),
      "boot": fileURLToPath(new URL("./src/boot", import.meta.url)),
    },
  },
});
