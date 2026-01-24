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
├── components/     # 재사용 컴포넌트
├── layouts/        # 레이아웃 컴포넌트
├── pages/          # 페이지 컴포넌트
├── router/         # Vue Router 설정
├── stores/         # Pinia 스토어
├── schema/         # TypeScript 타입/인터페이스
├── services/       # API 호출 서비스
└── utils/          # 유틸리티 함수
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
