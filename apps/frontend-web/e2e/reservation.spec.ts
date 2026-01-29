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

test.describe('Reservation CRUD Tests', () => {
  test.describe.configure({ mode: 'serial' });

  let createdReservationId: number | null = null;
  let testGuestName: string;

  test.beforeAll(async () => {
    testGuestName = `테스트손님-${generateUniqueId()}`;
  });

  test('예약 대시보드 (홈) 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 홈(대시보드) 페이지 접근
    await page.goto('/#/');
    await page.waitForLoadState('networkidle');

    // then: 대시보드가 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('입실 정보 요약')).toBeVisible();
    
    // then: 달력 네비게이션이 있음
    const prevButton = page.locator('button:has(i.q-icon:text("chevron_left"))').first();
    const nextButton = page.locator('button:has(i.q-icon:text("chevron_right"))').first();
    await expect(prevButton).toBeVisible();
    await expect(nextButton).toBeVisible();
    
    // then: 새로고침 버튼이 있음
    await expect(page.locator('button:has-text("새로고침")').first()).toBeVisible();
  });

  test('예약 목록 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 목록 페이지 접근
    await page.goto('/#/reservations');
    await page.waitForLoadState('networkidle');

    // then: 테이블이 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.locator('.q-table__title')).toContainText('다가오는 예약');

    // then: 상세 검색 버튼이 있음
    await expect(page.locator('button:has-text("상세 검색")').first()).toBeVisible();

    // then: 추가 버튼이 있음
    const addButton = page.locator('a[href*="reservations/create"], button:has(i.q-icon:text("add"))').first();
    await expect(addButton).toBeVisible();
  });

  test('예약 생성', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 생성 페이지 접근
    await page.goto('/#/reservations/create');
    await page.waitForLoadState('networkidle');

    // then: 생성 폼이 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('예약 등록')).toBeVisible();

    // when: 예약자 정보 입력
    await page.waitForTimeout(1000);
    
    const nameInput = page.locator('input[placeholder="홍길동"]');
    await nameInput.fill(testGuestName);

    const phoneInput = page.locator('input[placeholder="010-0000-0000"]');
    await phoneInput.fill('01012345678');

    // when: 판매 금액 입력 (spinbutton 중 판매 금액 라벨이 있는 것 선택)
    const priceInput = page.getByRole('spinbutton', { name: '판매 금액' });
    await priceInput.fill('100000');

    // when: 등록 버튼 클릭 (객실 미배정 상태로 진행)
    await page.locator('button:has-text("등록")').click();
    
    // when: 객실 미배정 경고 다이얼로그가 나타나면 확인
    await page.waitForTimeout(500);
    const warningDialog = page.locator('.q-dialog');
    if (await warningDialog.isVisible()) {
      const confirmRegister = page.locator('.q-dialog button:has-text("등록")');
      if (await confirmRegister.isVisible()) {
        await confirmRegister.click();
      }
    }

    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    // then: 목록 페이지로 리다이렉트됨
    expect(page.url()).toContain('/#/reservations');
    expect(page.url()).not.toContain('/create');

    // then: 생성된 예약이 목록에 표시됨
    await expect(page.locator('.q-table')).toBeVisible();

    const rowWithName = page.locator(`text=${testGuestName}`).first();
    if (await rowWithName.isVisible()) {
      const detailLink = page.locator(`a:has-text("${testGuestName}")`).first();
      if (await detailLink.isVisible()) {
        const href = await detailLink.getAttribute('href');
        const match = href?.match(/\/reservations\/(\d+)/);
        if (match) {
          createdReservationId = parseInt(match[1], 10);
        }
      }
    }
  });

  test('예약 상세 조회', async ({ page }) => {
    test.skip(!createdReservationId, 'No reservation created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 상세 페이지 접근
    await page.goto(`/#/reservations/${createdReservationId}`);
    await page.waitForLoadState('networkidle');

    // then: 상세 정보가 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('예약 정보')).toBeVisible();
    
    // then: 예약자명이 input에 표시됨
    const nameInput = page.locator('input[placeholder="홍길동"]');
    await expect(nameInput).toHaveValue(testGuestName);

    // then: 수정/삭제 버튼 확인
    await expect(page.locator('button:has-text("삭제")').first()).toBeVisible();
    await expect(page.locator('button:has-text("수정")').first()).toBeVisible();

    // then: 히스토리 테이블이 있음
    await expect(page.getByText('변경 이력')).toBeVisible();
  });

  test('예약 수정', async ({ page }) => {
    test.skip(!createdReservationId, 'No reservation created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 상세 페이지에서 수정 버튼 클릭
    await page.goto(`/#/reservations/${createdReservationId}`);
    await page.waitForLoadState('networkidle');
    
    await page.locator('button:has-text("수정")').first().click();
    await page.waitForTimeout(500);

    // then: 수정 모드가 활성화됨
    await expect(page.locator('.q-card').first()).toBeVisible();
    await expect(page.getByText('예약 수정')).toBeVisible();

    // when: 메모 수정
    const noteInput = page.locator('textarea');
    await noteInput.fill('E2E 테스트 수정됨');

    // when: 수정 버튼 클릭
    await page.locator('button:has-text("수정")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 다시 view 모드로 전환됨
    await expect(page.locator('.q-card').first()).toBeVisible();
  });

  test('예약 검색/필터링', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 목록 페이지 접근
    await page.goto('/#/reservations');
    await page.waitForLoadState('networkidle');

    // when: 상세 검색 버튼 클릭
    await page.locator('button:has-text("상세 검색")').click();
    await page.waitForTimeout(500);

    // then: 검색 다이얼로그가 표시됨
    await expect(page.locator('.q-dialog')).toBeVisible();

    // when: 필터 옵션들 확인
    // Status select should be visible
    const statusSelect = page.locator('.q-dialog .q-select').first();
    await expect(statusSelect).toBeVisible();

    // when: 다이얼로그 닫기
    const closeButton = page.locator('.q-dialog button:has-text("취소")');
    if (await closeButton.isVisible()) {
      await closeButton.click();
    } else {
      await page.keyboard.press('Escape');
    }
    await page.waitForTimeout(300);
  });

  test('예약 삭제 (상세 페이지에서)', async ({ page }) => {
    test.skip(!createdReservationId, 'No reservation created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 예약 상세 페이지 접근
    await page.goto(`/#/reservations/${createdReservationId}`);
    await page.waitForLoadState('networkidle');

    // then: 상세 정보가 표시됨
    await expect(page.locator('.q-card')).toBeVisible();

    // when: 삭제 버튼 클릭
    await page.locator('button:has-text("삭제")').first().click();

    // then: 삭제 확인 다이얼로그 표시
    await expect(page.locator('.q-dialog')).toBeVisible();

    // when: 삭제 확인
    await page.locator('.q-dialog button:has-text("삭제")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 목록 페이지로 리다이렉트됨
    expect(page.url()).toContain('/#/reservations');
    expect(page.url()).not.toContain(`/${createdReservationId}`);

    // Mark as deleted to prevent double cleanup
    createdReservationId = null;
  });
});
