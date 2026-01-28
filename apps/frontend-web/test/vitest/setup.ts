import { config } from "@vue/test-utils";
import { Quasar } from "quasar";
import { createTestingPinia } from "@pinia/testing";

// Global configuration for Vue Test Utils
config.global.plugins = [Quasar, createTestingPinia()];

// Mock modules that might cause issues in test environment
vi.mock("axios", () => ({
  default: {
    create: vi.fn(() => ({
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
      patch: vi.fn(),
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() },
      },
    })),
  },
}));

// Mock Quasar plugins if needed
vi.mock("quasar", async () => {
  const actual = await vi.importActual("quasar");
  return {
    ...actual,
    Notify: {
      create: vi.fn(),
    },
    Dialog: {
      create: vi.fn(),
    },
    Loading: {
      show: vi.fn(),
      hide: vi.fn(),
    },
  };
});
