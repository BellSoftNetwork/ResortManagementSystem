import { test, expect, Page } from '@playwright/test';
import { execSync } from 'child_process';

const TEST_USER = {
  userId: 'testadmin',
  password: 'password123',
};

const PROJECT_ROOT = '/home/bell/programming/projects/git/intellij/resort-management-system';

function clearLoginAttempts(username: string): void {
  try {
    execSync(
      `docker compose exec -T mysql mysql -urms -prms123 \\\`rms-core\\\` -e "DELETE FROM login_attempts WHERE username='${username}'"`,
      { cwd: PROJECT_ROOT, stdio: 'pipe' }
    );
  } catch {
  }
}

async function loginViaUI(page: Page, username: string, password: string): Promise<void> {
  clearLoginAttempts(username);
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

test.describe('Stats Page Tests', () => {
  test('통계 페이지 접근 및 UI 확인', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 통계 페이지 접근
    await page.goto('/#/stats');
    await page.waitForLoadState('networkidle');

    // then: 페이지 제목이 표시됨
    await expect(page.getByText('통계')).toBeVisible();

    // then: 월별 통계 카드들이 있음 (데이터 로딩 완료 대기)
    await page.waitForTimeout(2000);
    const cardCount = await page.locator('.q-card').count();
    expect(cardCount).toBeGreaterThanOrEqual(1);
  });

  test('통계 페이지 월 선택기 동작', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 통계 페이지 접근
    await page.goto('/#/stats');
    await page.waitForLoadState('networkidle');

    // then: 월 선택기 버튼들이 있음
    const prevButton = page.locator('button:has(i.q-icon:text("chevron_left"))').first();
    const nextButton = page.locator('button:has(i.q-icon:text("chevron_right"))').first();
    
    await expect(prevButton).toBeVisible();
    await expect(nextButton).toBeVisible();
  });

  test('통계 페이지 데이터 로딩 확인', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 통계 페이지 접근
    await page.goto('/#/stats');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    // then: 로딩 완료 후 데이터 카드가 표시됨
    await expect(page.locator('.q-page')).toBeVisible();
    
    // then: 스켈레톤 로더가 사라지고 실제 데이터가 표시됨 (또는 0 값)
    const salesCard = page.locator('.q-card').first();
    await expect(salesCard).toBeVisible();
  });
});

test.describe('Error Page Tests', () => {
  test('403 에러 페이지 확인', async ({ page }) => {
    // when: 403 에러 페이지 직접 접근
    await page.goto('/#/error/403');
    await page.waitForLoadState('networkidle');

    // then: 권한 없음 메시지가 표시됨
    await expect(page.getByText('해당 기능을 사용할 권한이 없습니다.')).toBeVisible();
  });

  test('404 에러 페이지 확인', async ({ page }) => {
    // when: 404 에러 페이지 직접 접근
    await page.goto('/#/error/404');
    await page.waitForLoadState('networkidle');

    // then: 404 코드가 표시됨
    await expect(page.getByText('404')).toBeVisible();

    // then: 에러 메시지가 표시됨
    await expect(page.getByText('이런... 여기에는 아무것도 존재하지 않아요.')).toBeVisible();

    // then: 홈으로 가기 버튼이 있음
    await expect(page.getByRole('link', { name: 'Go Home' })).toBeVisible();
  });

  test('존재하지 않는 페이지 접근 시 404 처리', async ({ page }) => {
    // when: 존재하지 않는 페이지 접근
    await page.goto('/#/nonexistent-page-12345');
    await page.waitForLoadState('networkidle');

    // then: 404 페이지로 리다이렉트되거나 404 에러가 표시됨
    const is404Page = await page.getByText('404').isVisible();
    const isNotFoundText = await page.getByText('아무것도 존재하지 않아요').isVisible();
    
    expect(is404Page || isNotFoundText).toBeTruthy();
  });
});

test.describe('Permission Access Tests', () => {
  test('권한 없는 사용자가 관리자 페이지 접근 시 리다이렉트', async ({ page }) => {
    // given: 로그인하지 않은 상태

    // when: 관리자 페이지 직접 접근 시도
    await page.goto('/#/admin/accounts');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 로그인 페이지로 리다이렉트되거나 403 에러 표시
    const currentUrl = page.url();
    const isRedirectedToLogin = currentUrl.includes('/login');
    const is403Page = await page.getByText('권한이 없습니다').isVisible().catch(() => false);
    
    expect(isRedirectedToLogin || is403Page).toBeTruthy();
  });
});
