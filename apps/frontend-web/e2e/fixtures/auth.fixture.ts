import { test as base, expect, Page } from "@playwright/test";
import { execSync } from "child_process";
import path from "path";
import fs from "fs";

export const TEST_CREDENTIALS = {
  userId: "testadmin",
  email: "testadmin@example.com",
  name: "Test Admin",
  password: "password123",
};

const AUTH_FILE = path.join(__dirname, "../playwright/.auth/user.json");

async function getServerConfig(page: Page): Promise<{ isAvailableRegistration: boolean }> {
  const response = await page.request.get("/api/v1/config");
  if (!response.ok()) throw new Error(`Failed to get server config: ${response.status()}`);
  const data = await response.json();
  return data.value;
}

async function userExists(page: Page): Promise<boolean> {
  const response = await page.request.post("/api/v1/auth/login", {
    data: { username: TEST_CREDENTIALS.userId, password: TEST_CREDENTIALS.password },
  });
  return response.ok();
}

async function registerFirstUser(page: Page): Promise<void> {
  await page.goto("/register");
  await expect(page.locator(".q-card")).toBeVisible();

  await page.locator("input").nth(0).fill(TEST_CREDENTIALS.userId);
  await page.locator("input").nth(1).fill(TEST_CREDENTIALS.email);
  await page.locator("input").nth(2).fill(TEST_CREDENTIALS.name);
  await page.locator("input").nth(3).fill(TEST_CREDENTIALS.password);
  await page.locator("input").nth(4).fill(TEST_CREDENTIALS.password);

  await page.locator('button[type="submit"]').click();
  await page.waitForURL("**/");
}

function upgradeToSuperAdmin(userId: string): void {
  try {
    execSync(
      `docker compose exec -T mysql mysql -urms -prms123 rms-core -e "UPDATE users SET role='SUPER_ADMIN' WHERE user_id='${userId}'"`,
      { cwd: "/home/bell/programming/projects/git/intellij/resort-management-system", stdio: "pipe" },
    );
  } catch {}
}

async function performLogin(page: Page): Promise<void> {
  await page.goto("/login");
  await expect(page.locator(".q-card")).toBeVisible();

  await page.locator("input").first().fill(TEST_CREDENTIALS.userId);
  await page.locator('input[type="password"]').fill(TEST_CREDENTIALS.password);
  await page.locator('button[type="submit"]').click();

  await page.waitForURL("**/");
  await expect(page.locator("body")).toBeVisible();
}

type AuthFixtures = {
  authenticatedPage: Page;
};

export const test = base.extend<AuthFixtures>({
  authenticatedPage: async ({ browser }, use) => {
    const authDir = path.dirname(AUTH_FILE);
    if (!fs.existsSync(authDir)) {
      fs.mkdirSync(authDir, { recursive: true });
    }

    let context;
    if (fs.existsSync(AUTH_FILE)) {
      context = await browser.newContext({ storageState: AUTH_FILE });
      const page = await context.newPage();
      const refreshToken = await page.evaluate(() => localStorage.getItem("refresh_token"));
      if (refreshToken) {
        await use(page);
        await context.close();
        return;
      }
      await context.close();
    }

    context = await browser.newContext();
    const page = await context.newPage();

    if (!(await userExists(page))) {
      const config = await getServerConfig(page);
      if (config.isAvailableRegistration) {
        await registerFirstUser(page);
        upgradeToSuperAdmin(TEST_CREDENTIALS.userId);
      } else {
        throw new Error("Cannot create test user: registration not available");
      }
    } else {
      await performLogin(page);
    }

    await page.context().storageState({ path: AUTH_FILE });
    await use(page);
    await context.close();
  },
});

export { expect } from "@playwright/test";
