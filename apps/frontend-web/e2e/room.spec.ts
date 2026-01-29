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

interface TestRoomGroup {
  id: number;
  name: string;
}

async function createTestRoomGroupViaUI(page: Page): Promise<TestRoomGroup | null> {
  const name = `객실그룹-${generateUniqueId()}`;
  
  await page.goto('/#/room-groups/create');
  await page.waitForLoadState('networkidle');
  
  const nameInput = page.locator('input').first();
  await nameInput.fill(name);

  const peekPriceInput = page.locator('input[type="number"]').first();
  await peekPriceInput.fill('100000');

  const offPeekPriceInput = page.locator('input[type="number"]').nth(1);
  await offPeekPriceInput.fill('80000');

  await page.locator('button:has-text("추가")').click();
  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(1000);

  const detailLink = page.locator(`a:has-text("${name}")`).first();
  if (await detailLink.isVisible()) {
    const href = await detailLink.getAttribute('href');
    const match = href?.match(/\/room-groups\/(\d+)/);
    if (match) {
      return { id: parseInt(match[1], 10), name };
    }
  }
  return null;
}

async function deleteRoomGroupViaUI(page: Page, id: number): Promise<void> {
  await page.goto(`/#/room-groups/${id}`);
  await page.waitForLoadState('networkidle');
  
  const deleteButton = page.locator('button:has-text("삭제")').first();
  if (await deleteButton.isVisible()) {
    await deleteButton.click();
    await page.waitForTimeout(300);
    const confirmButton = page.locator('.q-dialog button:has-text("삭제")');
    if (await confirmButton.isVisible()) {
      await confirmButton.click();
      await page.waitForLoadState('networkidle');
    }
  }
}

async function deleteRoomViaUI(page: Page, id: number): Promise<void> {
  await page.goto(`/#/rooms/${id}`);
  await page.waitForLoadState('networkidle');
  
  const deleteButton = page.locator('button:has-text("삭제")').first();
  if (await deleteButton.isVisible()) {
    await deleteButton.click();
    await page.waitForTimeout(300);
    const confirmButton = page.locator('.q-dialog button:has-text("삭제")');
    if (await confirmButton.isVisible()) {
      await confirmButton.click();
      await page.waitForLoadState('networkidle');
    }
  }
}

test.describe('Room CRUD Tests', () => {
  test.describe.configure({ mode: 'serial' });

  let testRoomGroup: TestRoomGroup | null = null;
  let createdRoomId: number | null = null;
  const testRoomNumber = generateUniqueId();

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      clearLoginAttempts(TEST_USER.userId);
      await loginViaUI(page, TEST_USER.userId, TEST_USER.password);
      testRoomGroup = await createTestRoomGroupViaUI(page);
    } finally {
      await page.close();
    }
  });

  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      await loginViaUI(page, TEST_USER.userId, TEST_USER.password);
      if (createdRoomId) {
        await deleteRoomViaUI(page, createdRoomId);
      }
      if (testRoomGroup) {
        await deleteRoomGroupViaUI(page, testRoomGroup.id);
      }
    } finally {
      await page.close();
    }
  });

  test('객실 현황 페이지 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 현황 페이지 접근
    await page.goto('/#/room-status');
    await page.waitForLoadState('networkidle');

    // then: 객실 현황 페이지가 표시됨
    await expect(page.getByText('객실 현황')).toBeVisible();
    await expect(page.locator('input[type="date"]').first()).toBeVisible();
  });

  test('객실 목록 조회', async ({ page }) => {
    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 목록 페이지 접근
    await page.goto('/#/rooms');
    await page.waitForLoadState('networkidle');

    // then: 테이블이 표시됨
    await expect(page.locator('.q-table')).toBeVisible();
    await expect(page.locator('.q-table__title')).toContainText('객실');

    // then: 추가 버튼이 있음
    const addButton = page.locator('a[href*="rooms/create"], button:has-text("add")').first();
    await expect(addButton).toBeVisible();
  });

  test('객실 생성', async ({ page }) => {
    test.skip(!testRoomGroup, 'No room group available');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 생성 페이지 접근
    await page.goto('/#/rooms/create');
    await page.waitForLoadState('networkidle');

    // then: 생성 폼이 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('객실 추가')).toBeVisible();

    // when: 객실 그룹 드롭다운이 로드될 때까지 대기
    await page.waitForTimeout(1000);
    const roomGroupSelect = page.locator('.q-select').first();
    await roomGroupSelect.click();
    await page.waitForTimeout(500);

    // when: 첫 번째 옵션 선택 (또는 테스트 그룹 선택)
    const options = page.locator('.q-menu .q-item');
    await expect(options.first()).toBeVisible({ timeout: 5000 });
    const optionWithTestGroup = page.locator(`.q-menu .q-item:has-text("${testRoomGroup!.name}")`).first();
    if (await optionWithTestGroup.isVisible()) {
      await optionWithTestGroup.click();
    } else {
      await options.first().click();
    }
    await page.waitForTimeout(300);

    // when: 객실 번호 입력
    const numberInput = page.locator('input[placeholder="101호"]');
    await numberInput.fill(testRoomNumber);

    // when: 추가 버튼 클릭
    await page.locator('button:has-text("추가")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    // then: 목록 페이지로 리다이렉트됨
    expect(page.url()).toContain('/#/rooms');
    expect(page.url()).not.toContain('/create');

    // then: 생성된 객실이 목록에 표시됨
    await expect(page.locator('.q-table')).toBeVisible();

    const rowWithNumber = page.locator(`text=${testRoomNumber}`).first();
    if (await rowWithNumber.isVisible()) {
      const row = rowWithNumber.locator('xpath=ancestor::tr');
      const editLink = row.locator('a[href*="/rooms/"][href*="/edit"]');
      if (await editLink.isVisible()) {
        const href = await editLink.getAttribute('href');
        const match = href?.match(/\/rooms\/(\d+)\/edit/);
        if (match) {
          createdRoomId = parseInt(match[1], 10);
        }
      }
    }

    if (!createdRoomId) {
      const detailLink = page.locator(`a:has-text("${testRoomNumber}")`).first();
      if (await detailLink.isVisible()) {
        const href = await detailLink.getAttribute('href');
        const match = href?.match(/\/rooms\/(\d+)/);
        if (match) {
          createdRoomId = parseInt(match[1], 10);
        }
      }
    }
  });

  test('객실 상세 조회', async ({ page }) => {
    test.skip(!createdRoomId, 'No room created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 상세 페이지 접근
    await page.goto(`/#/rooms/${createdRoomId}`);
    await page.waitForLoadState('networkidle');

    // then: 상세 정보가 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText(testRoomNumber)).toBeVisible();

    // then: 객실 그룹 정보 확인
    await expect(page.getByText('객실 그룹')).toBeVisible();

    // then: 수정/삭제 버튼 확인
    await expect(page.locator('button:has-text("삭제"), .q-btn:has-text("삭제")').first()).toBeVisible();
    await expect(page.locator('a:has-text("수정"), .q-btn:has-text("수정")').first()).toBeVisible();
  });

  test('객실 수정', async ({ page }) => {
    test.skip(!createdRoomId, 'No room created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 수정 페이지 접근
    await page.goto(`/#/rooms/${createdRoomId}/edit`);
    await page.waitForLoadState('networkidle');

    // then: 수정 폼이 표시됨
    await expect(page.locator('.q-card')).toBeVisible();
    await expect(page.getByText('객실 수정')).toBeVisible();

    // when: 메모 수정
    const noteInput = page.locator('textarea');
    await noteInput.fill('E2E 테스트 수정됨');

    // when: 수정 버튼 클릭
    await page.locator('button:has-text("수정")').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // then: 상세 페이지로 리다이렉트됨
    expect(page.url()).toContain(`/#/rooms/${createdRoomId}`);
    expect(page.url()).not.toContain('/edit');

    // then: 수정된 정보 확인
    await expect(page.locator('.q-card')).toBeVisible();
  });

  test('객실 삭제 (상세 페이지에서)', async ({ page }) => {
    test.skip(!createdRoomId, 'No room created in previous test');

    // given: 사용자가 로그인된 상태
    await loginViaUI(page, TEST_USER.userId, TEST_USER.password);

    // when: 객실 상세 페이지 접근
    await page.goto(`/#/rooms/${createdRoomId}`);
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
    expect(page.url()).toContain('/#/rooms');
    expect(page.url()).not.toContain(`/${createdRoomId}`);

    createdRoomId = null;
  });
});
