# frontend-web 개발 가이드

Vue.js 3 + Quasar Framework 기반 SPA

## 기술 스택

- **Framework**: Vue.js 3
- **UI Library**: Quasar Framework
- **Language**: TypeScript (strict 모드)
- **State Management**: Pinia
- **Build Tool**: Vite

## 개발 명령어

```bash
# Docker 컨테이너 접속
docker compose exec frontend sh

# 컨테이너 내부에서 실행
yarn                 # 의존성 설치
yarn dev            # 개발 서버 (http://localhost:9000)
yarn build          # 프로덕션 빌드
yarn lint           # ESLint 실행
yarn format         # Prettier 실행
```

## 프로젝트 구조

```
src/
├── api/            # API 호출 함수
├── components/     # 재사용 컴포넌트
├── composables/    # Vue Composables (useTable 등)
├── layouts/        # 레이아웃 컴포넌트
├── pages/          # 페이지 컴포넌트
├── router/         # Vue Router 설정
├── stores/         # Pinia 스토어 (도메인별)
├── schema/         # TypeScript 타입/인터페이스
└── util/           # 유틸리티 함수
```

## 아키텍처 패턴

### Composables

재사용 가능한 Vue Composition API 로직을 제공합니다.

#### `useTable` - QTable 공통 로직

페이지네이션, 정렬, 필터링, URL 동기화를 포함한 QTable 공통 로직을 제공합니다.

**주요 기능:**

- 페이지네이션 상태 관리
- 정렬 및 필터링 처리
- URL 쿼리 파라미터 동기화
- 로딩 상태 관리
- 에러 처리

**사용 예시:**

```typescript
import { useTable } from "src/composables/useTable";
import { fetchRooms } from "src/api/room";
import type { Room } from "src/schema/room";

const { pagination, loading, rows, onRequest } = useTable<Room>({
  fetchFn: fetchRooms,
  defaultPagination: { sortBy: "number", page: 1, rowsPerPage: 15 },
  onError: (error) =>
    $q.notify({
      message: getErrorMessage(error),
      type: "negative",
    }),
});
```

### Domain Stores

도메인별 Pinia Store 패턴을 사용하여 상태 관리를 구조화합니다.

**현재 구현된 Store:**

- `room.ts` - 객실 도메인
- `reservation.ts` - 예약 도메인
- `roomGroup.ts` - 객실 그룹 도메인

**각 Store의 공통 구조:**

- **State**: `list`, `loading`, `error`
- **Actions**:
  - `fetchList()` - 목록 조회
  - `getById(id)` - 단일 항목 조회
  - `create(data)` - 생성
  - `update(id, data)` - 수정
  - `delete(id)` - 삭제

**사용 예시:**

```typescript
import { useRoomStore } from "src/stores/room";

const roomStore = useRoomStore();

// 목록 조회
await roomStore.fetchList({ page: 1, size: 15 });

// 단일 항목 조회
const room = await roomStore.getById(1);

// 생성
await roomStore.create({ number: "101", type: "STANDARD" });

// 수정
await roomStore.update(1, { number: "102" });

// 삭제
await roomStore.delete(1);
```

### Utilities

공통 유틸리티 함수를 제공합니다.

#### `date-preset-util.ts` - 날짜 범위 계산

예약 필터링에 사용되는 날짜 범위를 계산하는 유틸리티입니다.

**제공 타입 및 함수:**

- `DueOption` - 날짜 프리셋 타입: `"ALL" | "1M" | "2M" | "3M" | "6M" | "CUSTOM"`
- `DateRange` - 날짜 범위 인터페이스: `{ startAt: string; endAt: string }`
- `calculateDateRange(dueOption, startDate?)` - 프리셋에 따른 날짜 범위 반환

**사용 예시:**

```typescript
import { calculateDateRange, DueOption } from "src/util/date-preset-util";

// 6개월 범위 계산
const range = calculateDateRange("6M");
console.log(range); // { startAt: '2026-01-26', endAt: '2026-07-26' }

// 전체 (빈 범위)
const allRange = calculateDateRange("ALL");
console.log(allRange); // { startAt: '', endAt: '' }
```

## API 프록시

개발 서버에서 백엔드 API 호출 시 프록시 사용:

```javascript
// quasar.config.js
devServer: {
  proxy: {
    '/api': {
      target: 'http://api-core:8080',
      changeOrigin: true
    }
  }
}
```

## 모바일 앱 빌드 (Capacitor)

```bash
# 호스트에서 실행 (Android Studio 필요)
yarn cap:add-android      # Android 플랫폼 추가 (처음만)
yarn android:build        # 빌드 및 동기화
yarn cap:open-android     # Android Studio에서 열기
yarn android:run          # 로컬 에뮬레이터
yarn android:run:dev      # 개발 API
yarn android:run:prod     # 프로덕션 API
```

## 코드 품질

### ESLint + Prettier

```bash
yarn lint           # 검사
yarn lint --fix     # 자동 수정
yarn format         # Prettier 실행
```

### TypeScript

- strict 모드 활성화
- any 타입 최소화
- interface/type 명확히 정의

## 테스트 (계획)

현재 유닛 테스트 없음. [frontend/unit-testing 스펙](../../specs/frontend/unit-testing/spec.md) 참조.

계획된 도구:

- **테스트 프레임워크**: Vitest
- **컴포넌트 테스트**: @vue/test-utils
- **API 모킹**: msw

## API 스키마

백엔드 API 응답 타입:

```typescript
// src/schema/revision.ts
type Revision<T> = {
  entity: T;
  historyType: "CREATED" | "UPDATED" | "DELETED";
  historyCreatedAt: string;
  updatedFields: string[];
};
```

## 환경 변수

```bash
# API URL 설정
VITE_API_URL=http://localhost:8080

# 개발/프로덕션 분리
VITE_API_URL_DEV=http://dev-api.example.com
VITE_API_URL_PROD=https://api.example.com
```

## 참조 문서

- [API 비교](../../docs/contracts/api-comparison.md)
- [프론트엔드 유닛 테스트 스펙](../../specs/frontend/unit-testing/spec.md)
