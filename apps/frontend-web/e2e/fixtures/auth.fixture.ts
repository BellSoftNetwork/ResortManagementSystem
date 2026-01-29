import { test as base } from '@playwright/test';

type AuthFixtures = {
  authenticatedPage: void;
};

export const test = base.extend<AuthFixtures>({
  authenticatedPage: async ({ page }, use) => {
    // TODO: Implement auth state injection when auth tests are ready
    await use();
  },
});

export { expect } from '@playwright/test';
