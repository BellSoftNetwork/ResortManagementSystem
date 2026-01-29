# Frontend E2E & Unit Testing Checklist

프론트엔드 테스트 현황 및 실행 가이드

## Test Summary

| Type | Files | Tests | Coverage |
|------|-------|-------|----------|
| E2E (Playwright) | 8 | ~50 | 주요 사용자 플로우 |
| Unit (Vitest) | 13 | 238 | Stores, Composables, Utils |

## E2E Tests

### Test Files

| File | Tests | Routes Covered |
|------|-------|----------------|
| `auth.spec.ts` | 6 | /login, /register, /my |
| `room.spec.ts` | 7 | /room-status, /rooms/* |
| `room-group.spec.ts` | 6 | /room-groups/* |
| `reservation.spec.ts` | 7 | /, /reservations/* |
| `monthly-rent.spec.ts` | 6 | /monthly-rents/* |
| `admin.spec.ts` | 10 | /payment-methods, /admin/*, /debug |
| `misc.spec.ts` | 7 | /stats, /error/403, /error/404 |
| `smoke.spec.ts` | 1 | / (basic connectivity) |

### Running E2E Tests

```bash
# Docker 컨테이너 접속
docker compose exec frontend sh

# 전체 E2E 테스트 실행
yarn test:e2e

# 특정 파일만 실행
yarn test:e2e e2e/auth.spec.ts

# UI 모드로 실행 (디버깅용)
yarn test:e2e --ui

# 특정 테스트만 실행
yarn test:e2e --grep "로그인"
```

### E2E Test Prerequisites

1. Docker 환경 실행 중 (`docker compose up -d`)
2. Frontend 서버 실행 중 (localhost:9000)
3. API 서버 실행 중 (localhost:8080)
4. 테스트 계정: `testadmin` / `password123`

### E2E Test Patterns

```typescript
// BDD 스타일 주석 사용
it("사용자가 로그인할 수 있다", async ({ page }) => {
  // given: 로그인 페이지 접속
  await page.goto("/#/login");
  
  // when: 자격 증명 입력 및 제출
  await page.locator('input').first().fill('testadmin');
  await page.locator('input[type="password"]').fill('password123');
  await page.locator('button[type="submit"]').click();
  
  // then: 대시보드로 이동
  await expect(page).toHaveURL(/\/#\/$/);
});
```

## Unit Tests

### Test Files by Category

#### Stores (`src/stores/__tests__/`)

| File | Tests | Coverage |
|------|-------|----------|
| `auth.test.ts` | 24 | login, logout, token refresh, roles |
| `network.test.ts` | 19 | error states, retry logic |
| `room.test.ts` | 27 | CRUD, filtering, history |
| `roomGroup.test.ts` | 23 | CRUD, validation |
| `reservation.test.ts` | 32 | CRUD, statistics |
| `simple.test.ts` | 5 | basic vitest validation |

#### Composables (`src/composables/__tests__/`)

| File | Tests | Coverage |
|------|-------|----------|
| `useTable.test.ts` | 20 | pagination, sorting, URL sync |
| `useReservationCalendar.test.ts` | 21 | date arrays, event formatting |
| `useStatsData.test.ts` | 10 | yearly data loading |

#### Utilities (`src/util/__tests__/`)

| File | Tests | Coverage |
|------|-------|----------|
| `date-preset-util.test.ts` | 22 | date range calculation |
| `format-util.test.ts` | 22 | price, date, stay formatting |
| `errorHandler.test.ts` | 10 | error message extraction |
| `query-string-util.test.ts` | 3 | sort param formatting |

### Running Unit Tests

```bash
# Docker 컨테이너 접속
docker compose exec frontend sh

# 전체 유닛 테스트 실행
yarn test:unit

# 특정 파일만 실행
yarn test:unit src/stores/__tests__/auth.test.ts

# Watch 모드
yarn test:unit --watch

# 커버리지 리포트
yarn test:unit --coverage
```

### Unit Test Patterns

```typescript
import { describe, it, expect, vi, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";

describe("useAuthStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("로그인 성공 시 토큰을 저장한다", async () => {
    // given: mock 응답 설정
    vi.mocked(authApi.postLogin).mockResolvedValue(mockResponse);
    
    // when: 로그인 실행
    const store = useAuthStore();
    await store.login({ username: "test", password: "test" });
    
    // then: 토큰 저장 확인
    expect(store.accessToken).toBe("mock-token");
  });
});
```

## Test Helpers

### Location

`apps/frontend-web/test/vitest/helpers.ts`

### Available Helpers

```typescript
// 단일 값 응답 생성
createMockApiResponse<T>(value: T): { value: T }

// 목록 응답 생성
createMockApiResponse<T>(value: T[], totalCount: number): { value: T[], totalCount: number }

// 에러 생성
createMockApiError(message: string): Error
```

## CI/CD Integration

### GitHub Actions / GitLab CI

```yaml
test:
  script:
    - cd apps/frontend-web
    - yarn install
    - yarn test:unit --run
    - yarn test:e2e
```

## Known Issues & Workarounds

### 1. Login Attempts Table

E2E 테스트 전 `login_attempts` 테이블 클리어 필요:

```typescript
function clearLoginAttempts(username: string): void {
  execSync(
    `docker compose exec -T mysql mysql -urms -prms123 rms-core -e "DELETE FROM login_attempts WHERE username='${username}'"`,
    { cwd: PROJECT_ROOT, stdio: 'pipe' }
  );
}
```

### 2. Hash-based Routing

모든 라우트는 `/#/path` 형식 사용:

```typescript
await page.goto('/#/reservations');  // 올바름
await page.goto('/reservations');     // 잘못됨
```

### 3. Multiple Elements with Same Selector

`.q-card` 등 여러 요소가 있을 때 `.first()` 사용:

```typescript
await page.locator('.q-card').first().click();
```

### 4. Axios Mock in Vitest

`errorHandler.test.ts`에서 실제 axios 사용 필요:

```typescript
vi.unmock("axios");
import axios, { AxiosError } from "axios";
```

## Maintenance

### Adding New Tests

1. E2E: `apps/frontend-web/e2e/` 디렉토리에 `*.spec.ts` 파일 추가
2. Unit: 해당 소스 파일 옆 `__tests__/` 디렉토리에 `*.test.ts` 파일 추가

### Test Naming Convention

- E2E: `{feature}.spec.ts` (예: `reservation.spec.ts`)
- Unit: `{source-file}.test.ts` (예: `auth.test.ts`)

### BDD Comments

테스트 내 `// given:`, `// when:`, `// then:` 형식 주석 사용 권장
