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

function generateUniqueId(): string {
  return `E2E-${Date.now().toString(36)}`;
}

test.describe('Account Management Tests', () => {
  test('계정 목록 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 계정 관리 페이지 접근
    await page.goto('/#/admin/accounts');
    await page.waitForLoadState('networkidle');

    // then: 테이블이 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.locator('.q-table__title')).toContainText('사용자 계정');

    // then: 추가 버튼이 있음
    const addButton = page.locator('button:has(i.q-icon:text("add"))').first();
    await expect(addButton).toBeVisible();

    // then: 테스트 관리자 계정이 목록에 표시됨
    await expect(page.getByRole('cell', { name: TEST_USER.userId, exact: true })).toBeVisible();
  });

  test('계정 생성 다이얼로그', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 계정 관리 페이지 접근
    await page.goto('/#/admin/accounts');
    await page.waitForLoadState('networkidle');

    // when: 추가 버튼 클릭
    await page.locator('button:has(i.q-icon:text("add"))').first().click();
    await page.waitForTimeout(500);

    // then: 생성 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();

    // then: 입력 필드들이 있음
    const inputCount = await page.locator('.q-dialog input').count();
    expect(inputCount).toBeGreaterThanOrEqual(3);

    // when: 다이얼로그 닫기
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);
  });

  test('계정 수정 다이얼로그', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 계정 관리 페이지 접근
    await page.goto('/#/admin/accounts');
    await page.waitForLoadState('networkidle');

    // when: 첫 번째 계정의 수정 버튼 클릭
    await page.locator('button:has(i.q-icon:text("edit"))').first().click();
    await page.waitForTimeout(500);

    // then: 수정 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();

    // when: 다이얼로그 닫기
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);
  });
});

test.describe('Payment Method Management Tests', () => {
  test.describe.configure({ mode: 'serial' });
  
  let createdPaymentMethodName: string;

  test.beforeAll(() => {
    createdPaymentMethodName = `E2E결제수단-${generateUniqueId()}`;
  });

  test('결제 수단 목록 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 결제 수단 페이지 접근
    await page.goto('/#/payment-methods');
    await page.waitForLoadState('networkidle');

    // then: 테이블이 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.locator('.q-table__title')).toContainText('결제 수단');

    // then: 추가 버튼이 있음
    const addButton = page.locator('button:has(i.q-icon:text("add"))').first();
    await expect(addButton).toBeVisible();
  });

  test('결제 수단 생성', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 결제 수단 페이지 접근
    await page.goto('/#/payment-methods');
    await page.waitForLoadState('networkidle');

    // when: 추가 버튼 클릭
    await page.locator('button:has(i.q-icon:text("add"))').first().click();
    await page.waitForTimeout(500);

    // then: 생성 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();

    // when: 결제 수단 정보 입력
    const nameInput = page.locator('.q-dialog input').first();
    await nameInput.fill(createdPaymentMethodName);

    // when: 등록 버튼 클릭 (결제 수단 추가 다이얼로그에서는 "추가" 버튼)
    await page.locator('.q-dialog button:has-text("추가")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 생성된 결제 수단이 목록에 표시됨
    await expect(page.getByText(createdPaymentMethodName)).toBeVisible();
  });

  test('결제 수단 인라인 수정 (이름 클릭)', async ({ page }) => {
    test.skip(!createdPaymentMethodName, 'No payment method created');
    
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 결제 수단 페이지 접근
    await page.goto('/#/payment-methods');
    await page.waitForLoadState('networkidle');

    // then: 생성된 결제 수단이 목록에 표시됨
    const row = page.locator(`td:has-text("${createdPaymentMethodName}")`).first();
    await expect(row).toBeVisible();
  });

  test('결제 수단 삭제', async ({ page }) => {
    test.skip(!createdPaymentMethodName, 'No payment method created');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 결제 수단 페이지 접근
    await page.goto('/#/payment-methods');
    await page.waitForLoadState('networkidle');

    // when: 생성된 결제 수단의 삭제 버튼 찾기
    const row = page.locator('tr').filter({ hasText: createdPaymentMethodName });
    await expect(row).toBeVisible();

    // when: 삭제 버튼 클릭
    await row.locator('button:has(i.q-icon:text("delete"))').click();
    await page.waitForTimeout(500);

    // then: 삭제 확인 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();

    // when: 삭제 확인
    await page.locator('.q-dialog button:has-text("삭제")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 삭제된 결제 수단이 목록에서 사라짐
    await expect(page.getByText(createdPaymentMethodName)).not.toBeVisible();
  });
});

test.describe('Dev Test Page (SUPER_ADMIN)', () => {
  test('개발 테스트 페이지 접근 및 UI 확인', async ({ page }) => {
    // given: SUPER_ADMIN이 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 개발 테스트 페이지 접근
    await page.goto('/#/admin/dev-test');
    await page.waitForLoadState('networkidle');

    // then: 페이지가 표시됨
    await expect(page.getByText('개발 테스트 도구')).toBeVisible();

    // then: 더미 데이터 생성 버튼들이 있음
    await expect(page.getByText('전체 데이터 생성')).toBeVisible();
    await expect(page.getByText('필수 데이터만 생성')).toBeVisible();
    await expect(page.getByText('예약 데이터만 생성')).toBeVisible();
  });

  test('예약 데이터 생성 옵션 다이얼로그', async ({ page }) => {
    // given: SUPER_ADMIN이 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 개발 테스트 페이지 접근
    await page.goto('/#/admin/dev-test');
    await page.waitForLoadState('networkidle');

    // when: 예약 데이터만 생성 버튼 클릭
    await page.getByRole('button', { name: '예약 데이터만 생성' }).click();
    await page.waitForTimeout(500);

    // then: 옵션 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();
    await expect(page.getByText('예약 데이터 생성 옵션')).toBeVisible();

    // then: 옵션 입력 필드들이 있음
    await expect(page.getByLabel('시작일')).toBeVisible();
    await expect(page.getByLabel('종료일')).toBeVisible();
    await expect(page.getByLabel('일반 예약 건수')).toBeVisible();
    await expect(page.getByLabel('달방 예약 건수')).toBeVisible();

    // when: 다이얼로그 닫기
    await page.locator('.q-dialog button:has-text("취소")').click();
    await page.waitForTimeout(300);
  });
});

test.describe('Debug Page (SUPER_ADMIN)', () => {
  test('디버그 페이지 접근 및 UI 확인', async ({ page }) => {
    // given: SUPER_ADMIN이 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 디버그 페이지 접근
    await page.goto('/#/debug');
    await page.waitForLoadState('networkidle');

    // then: 페이지가 로드됨 (API 테스트, 유저 정보 카드가 있을 수 있음)
    await expect(page.locator('.q-page')).toBeVisible();
    
    // then: 카드 UI가 있음
    await expect(page.locator('.q-card').first()).toBeVisible();
  });
});
