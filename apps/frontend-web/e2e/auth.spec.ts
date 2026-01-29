import { test, expect, Page } from '@playwright/test';
import { execSync } from 'child_process';

const TEST_USER = {
  userId: 'testadmin',
  email: 'testadmin@example.com',
  name: 'Test Admin',
  password: 'password123',
};

function isOnLoginPage(url: string): boolean {
  const hash = new URL(url).hash;
  return hash.includes('/login');
}

async function userExists(page: Page): Promise<boolean> {
  const response = await page.request.post('/api/v1/auth/login', {
    data: { username: TEST_USER.userId, password: TEST_USER.password },
  });
  return response.ok();
}

async function getServerConfig(page: Page): Promise<{ isAvailableRegistration: boolean }> {
  const response = await page.request.get('/api/v1/config');
  if (!response.ok()) throw new Error(`Failed to get server config: ${response.status()}`);
  const data = await response.json();
  return data.value;
}

async function registerFirstUser(page: Page): Promise<void> {
  await page.goto('/#/register');
  await expect(page.locator('.q-card')).toBeVisible();

  await page.locator('input').nth(0).fill(TEST_USER.userId);
  await page.locator('input').nth(1).fill(TEST_USER.email);
  await page.locator('input').nth(2).fill(TEST_USER.name);
  await page.locator('input').nth(3).fill(TEST_USER.password);
  await page.locator('input').nth(4).fill(TEST_USER.password);

  await page.locator('button[type="submit"]').click();
  await page.waitForLoadState('networkidle');
}

function upgradeToSuperAdmin(userId: string): void {
  try {
    execSync(
      `docker compose exec -T mysql mysql -urms -prms123 \\\`rms-core\\\` -e "UPDATE user SET role=127 WHERE user_id='${userId}'"`,
      { cwd: '/home/bell/programming/projects/git/intellij/resort-management-system', stdio: 'pipe' }
    );
  } catch {
  }
}

function createTestUserViaDatabaseIfNeeded(): void {
  try {
    const bcryptHash = '{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2';
    execSync(
      `docker compose exec -T mysql mysql -urms -prms123 \\\`rms-core\\\` -e "INSERT IGNORE INTO user (user_id, email, name, password, role, status, created_at, updated_at, deleted_at) VALUES ('${TEST_USER.userId}', '${TEST_USER.email}', '${TEST_USER.name}', '${bcryptHash}', 127, 1, NOW(), NOW(), '1970-01-01 00:00:00')"`,
      { cwd: '/home/bell/programming/projects/git/intellij/resort-management-system', stdio: 'pipe' }
    );
  } catch {
  }
}

async function loginViaUI(page: Page, username: string, password: string): Promise<void> {
  await page.goto('/#/login');
  await expect(page.locator('.q-card')).toBeVisible();

  const usernameInput = page.locator('input').first();
  const passwordInput = page.locator('input[type="password"]');

  await usernameInput.fill(username);
  await passwordInput.fill(password);
  await page.locator('button[type="submit"]').click();

  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(500);
}

function clearLoginAttempts(username: string): void {
  try {
    execSync(
      `docker compose exec -T mysql mysql -urms -prms123 \\\`rms-core\\\` -e "DELETE FROM login_attempts WHERE username='${username}'"`,
      { cwd: '/home/bell/programming/projects/git/intellij/resort-management-system', stdio: 'pipe' }
    );
  } catch {
  }
}

async function ensureTestUser(page: Page): Promise<void> {
  clearLoginAttempts(TEST_USER.userId);
  
  if (await userExists(page)) return;

  const config = await getServerConfig(page);
  if (config.isAvailableRegistration) {
    await registerFirstUser(page);
    upgradeToSuperAdmin(TEST_USER.userId);
  } else {
    createTestUserViaDatabaseIfNeeded();
  }
}

test.describe('Authentication Tests', () => {
  test.describe.configure({ mode: 'serial' });

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      await ensureTestUser(page);
    } finally {
      await page.close();
    }
  });

  test('로그인 성공 테스트', async ({ page }) => {
    // given: 사용자가 로그인 페이지에 있음
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // then: 홈으로 리다이렉트되고 토큰 저장됨
    expect(isOnLoginPage(page.url())).toBe(false);
    const refreshToken = await page.evaluate(() => localStorage.getItem('refresh_token'));
    expect(refreshToken).toBeTruthy();
  });

  test('로그인 실패 테스트 - 잘못된 비밀번호', async ({ page }) => {
    // given: 사용자가 로그인 페이지에 있음
    await page.goto('/#/login');
    await expect(page.locator('.q-card')).toBeVisible();

    // when: 잘못된 비밀번호로 로그인 시도
    const usernameInput = page.locator('input').first();
    const passwordInput = page.locator('input[type="password"]');
    await usernameInput.fill(TEST_USER.userId);
    await passwordInput.fill('WrongPassword123!');
    await page.locator('button[type="submit"]').click();

    // then: 에러 알림 표시되고 로그인 페이지에 유지됨
    await expect(page.locator('.q-notification')).toBeVisible({ timeout: 5000 });
    expect(isOnLoginPage(page.url())).toBe(true);
  });

  test('회원가입 테스트 - 등록 가능 여부 확인', async ({ page }) => {
    // given: 서버 설정 확인
    const config = await getServerConfig(page);

    // when: 회원가입 페이지 접근
    await page.goto('/#/register');
    await page.waitForLoadState('networkidle');

    // then: 설정에 따라 폼 표시 또는 로그인으로 리다이렉트
    if (config.isAvailableRegistration) {
      await expect(page.locator('.q-card')).toBeVisible();
      await expect(page.locator('button[type="submit"]')).toBeVisible();
    } else {
      expect(isOnLoginPage(page.url())).toBe(true);
    }
  });

  test('로그아웃 테스트', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);
    const refreshTokenBefore = await page.evaluate(() => localStorage.getItem('refresh_token'));
    expect(refreshTokenBefore).toBeTruthy();

    // when: 로그아웃 API 호출 및 토큰 제거
    const response = await page.request.post('/api/v1/auth/logout');
    expect(response.ok()).toBeTruthy();
    await page.evaluate(() => {
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('token_expires');
    });

    // then: 페이지 새로고침 후 로그인 페이지 표시
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const loginFormVisible = await page.locator('.q-card').first().isVisible().catch(() => false);
    expect(isOnLoginPage(page.url()) || loginFormVisible).toBe(true);
  });

  test('내 정보 페이지 접근 테스트', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 내 정보 페이지 접근
    await page.goto('/#/my');
    await page.waitForLoadState('networkidle');

    // then: 로그인 페이지로 리다이렉트되지 않고 페이지 내용 표시
    expect(isOnLoginPage(page.url())).toBe(false);
    const hasMyInfoTitle = await page.getByText('내 정보').isVisible();
    const hasNameInput = await page.locator('input').first().isVisible();
    expect(hasMyInfoTitle || hasNameInput).toBe(true);
  });

  test('토큰 만료 시 로그인 페이지로 리다이렉트', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);
    expect(isOnLoginPage(page.url())).toBe(false);

    // when: 토큰 만료 상태 시뮬레이션
    await page.evaluate(() => {
      localStorage.setItem('token_expires', String(Date.now() - 1000));
      localStorage.removeItem('refresh_token');
    });
    await page.goto('/#/my');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 로그인 페이지로 리다이렉트됨
    const loginCardVisible = await page.locator('.q-card').first().isVisible().catch(() => false);
    expect(isOnLoginPage(page.url()) || loginCardVisible).toBe(true);
  });
});
