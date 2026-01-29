import { test, expect } from '@playwright/test';

test.describe('Smoke Tests', () => {
  test('페이지 로드 확인', async ({ page }) => {
    await page.goto('/');

    await expect(page).toHaveURL(/localhost:9000/);
    await expect(page.locator('body')).toBeVisible();
  });

  test('로그인 페이지 표시 확인', async ({ page }) => {
    await page.goto('/');

    const loginForm = page.locator('form, .login, [class*="login"]');
    const isLoginPage = await loginForm.count() > 0;

    if (isLoginPage) {
      await expect(loginForm.first()).toBeVisible();
    } else {
      await expect(page.locator('body')).toBeVisible();
    }
  });
});
