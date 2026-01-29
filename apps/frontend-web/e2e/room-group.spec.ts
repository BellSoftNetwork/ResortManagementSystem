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
    // Ignore errors
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

test.describe('Room Group CRUD Tests', () => {
  test.describe.configure({ mode: 'serial' });

  let createdRoomGroupId: number | null = null;
  const testRoomGroupName = `테스트그룹-${generateUniqueId()}`;

  test.beforeAll(async ({ browser }) => {
    // Ensure test user can login
    const page = await browser.newPage();
    try {
      clearLoginAttempts(TEST_USER.userId);
    } finally {
      await page.close();
    }
  });

  test.afterAll(async ({ browser }) => {
    // Cleanup: Delete created room group via API
    if (createdRoomGroupId) {
      const page = await browser.newPage();
      try {
        await loginViaUI(page, TEST_USER.userId, TEST_USER.password);
        const response = await page.request.delete(`/api/v1/room-groups/${createdRoomGroupId}`);
        // Ignore 404 if already deleted
        if (!response.ok() && response.status() !== 404) {
          console.warn(`Failed to cleanup room group ${createdRoomGroupId}: ${response.status()}`);
        }
      } finally {
        await page.close();
      }
    }
  });

  test('객실 그룹 목록 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 그룹 목록 페이지 접근
    await page.goto('/#/room-groups');
    await page.waitForLoadState('networkidle');

    // then: 테이블이 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.locator('.q-table__title')).toContainText('객실 그룹');

    // 추가 버튼이 있음
    const addButton = page.locator('a[href*="room-groups/create"], button:has-text("add")').first();
    await expect(addButton).toBeVisible();
  });

  test('객실 그룹 생성', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 그룹 생성 페이지 접근
    await page.goto('/#/room-groups/create');
    await page.waitForLoadState('networkidle');

    // then: 생성 폼이 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('객실 그룹 추가')).toBeVisible();

    // when: 폼 작성
    const nameInput = page.locator('input').first();
    await nameInput.fill(testRoomGroupName);

    // 성수기 예약금
    const peekPriceInput = page.locator('input[type="number"]').first();
    await peekPriceInput.fill('100000');

    // 비성수기 예약금
    const offPeekPriceInput = page.locator('input[type="number"]').nth(1);
    await offPeekPriceInput.fill('80000');

    // 설명
    const descriptionInput = page.locator('textarea');
    await descriptionInput.fill('E2E 테스트용 객실 그룹입니다.');

    // when: 추가 버튼 클릭
    await page.locator('button:has-text("추가")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 목록 페이지로 리다이렉트됨
    expect(page.url()).toContain('/#/room-groups');
    expect(page.url()).not.toContain('/create');

    // then: 생성된 그룹이 목록에 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.getByText(testRoomGroupName)).toBeVisible();

    // Save the ID for later tests
    const rowWithName = page.locator(`text=${testRoomGroupName}`).first();
    const row = rowWithName.locator('xpath=ancestor::tr');
    const editLink = row.locator('a[href*="/room-groups/"][href*="/edit"]');
    if (await editLink.isVisible()) {
      const href = await editLink.getAttribute('href');
      const match = href?.match(/\/room-groups\/(\d+)\/edit/);
      if (match) {
        createdRoomGroupId = parseInt(match[1], 10);
      }
    }
    
    // Alternative: Get ID from detail link
    if (!createdRoomGroupId) {
      const detailLink = page.locator(`a:has-text("${testRoomGroupName}")`).first();
      const href = await detailLink.getAttribute('href');
      const match = href?.match(/\/room-groups\/(\d+)/);
      if (match) {
        createdRoomGroupId = parseInt(match[1], 10);
      }
    }
  });

  test('객실 그룹 상세 조회', async ({ page }) => {
    // Skip if no room group was created
    test.skip(!createdRoomGroupId, 'No room group created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 그룹 상세 페이지 접근
    await page.goto(`/#/room-groups/${createdRoomGroupId}`);
    await page.waitForLoadState('networkidle');

    // then: 상세 정보가 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText(testRoomGroupName)).toBeVisible();

    // then: 가격 정보 확인
    await expect(page.getByText('성수기 예약금', { exact: true })).toBeVisible();
    await expect(page.getByText('비성수기 예약금', { exact: true })).toBeVisible();

    // 수정/삭제 버튼 확인
    await expect(page.locator('button:has-text("삭제"), .q-btn:has-text("삭제")').first()).toBeVisible();
    await expect(page.locator('a:has-text("수정"), .q-btn:has-text("수정")').first()).toBeVisible();
  });

  test('객실 그룹 수정', async ({ page }) => {
    // Skip if no room group was created
    test.skip(!createdRoomGroupId, 'No room group created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 그룹 수정 페이지 접근
    await page.goto(`/#/room-groups/${createdRoomGroupId}/edit`);
    await page.waitForLoadState('networkidle');

    // then: 수정 폼이 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('객실 그룹 수정')).toBeVisible();

    // when: 가격 수정
    const peekPriceInput = page.locator('input[type="number"]').first();
    await peekPriceInput.clear();
    await peekPriceInput.fill('150000');

    // when: 수정 버튼 클릭
    await page.locator('button:has-text("수정")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 상세 페이지로 리다이렉트됨
    expect(page.url()).toContain(`/#/room-groups/${createdRoomGroupId}`);
    expect(page.url()).not.toContain('/edit');

    // then: 수정된 정보 확인
    await expect(page.locator('.q-card')).toBeVisible();
  });

  test('객실 그룹 삭제 (목록에서)', async ({ page }) => {
    // Skip if no room group was created
    test.skip(!createdRoomGroupId, 'No room group created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 그룹 목록 페이지 접근
    await page.goto('/#/room-groups');
    await page.waitForLoadState('networkidle');

    // then: 테이블에 그룹이 표시됨
    await expect(page.getByText(testRoomGroupName)).toBeVisible();

    // when: 삭제 버튼 클릭
    const rowWithName = page.locator(`text=${testRoomGroupName}`).first();
    const row = rowWithName.locator('xpath=ancestor::tr');
    const deleteButton = row.locator('button:has(i.material-icons:text("delete")), button[icon="delete"]').first();
    
    // Alternative selector for delete button in actions column
    if (!(await deleteButton.isVisible())) {
      const actionsCell = row.locator('td').last();
      await actionsCell.locator('button').last().click();
    } else {
      await deleteButton.click();
    }

    // then: 삭제 확인 다이얼로그 표시
    await expect(page.locator('.q-dialog')).toBeVisible();
    await expect(page.locator('.q-dialog__title')).toBeVisible();

    // when: 삭제 확인
    await page.locator('.q-dialog button:has-text("삭제")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 그룹이 목록에서 제거됨
    await expect(page.getByText(testRoomGroupName)).not.toBeVisible();

    // Mark as deleted so afterAll doesn't try to delete again
    createdRoomGroupId = null;
  });
});
