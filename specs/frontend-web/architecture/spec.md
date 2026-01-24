---
id: frontend-web-architecture
title: "frontend-web 아키텍처"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: low
effort: small
---

# frontend-web 아키텍처

> Vue.js 3 + Quasar Framework 기반 프론트엔드 아키텍처

---

## 1. 기술 스택

| 항목 | 기술 |
|------|------|
| 언어 | TypeScript |
| 프레임워크 | Vue.js 3 (Composition API) |
| UI 프레임워크 | Quasar Framework 2 |
| 상태관리 | Pinia |
| 라우팅 | Vue Router 4 |
| HTTP 클라이언트 | Axios |
| 차트 | ApexCharts |
| 모바일 | Capacitor (Android) |
| 테스트 | Vitest |
| 린터 | ESLint + Prettier |

---

## 2. 디렉토리 구조

```
apps/frontend-web/src/
├── api/                    # API 통신 레이어
│   ├── v1/                 # 엔드포인트별 함수
│   ├── services/           # HTTP 서비스
│   └── constants.ts
├── boot/                   # Quasar 부트 파일
├── components/             # 재사용 컴포넌트 (70+개)
├── layouts/                # 레이아웃 컴포넌트
├── pages/                  # 페이지 컴포넌트 (28개)
├── router/                 # Vue Router 설정
│   └── guards/             # 라우트 가드
├── schema/                 # TypeScript 타입 정의
├── stores/                 # Pinia 스토어 (3개)
├── util/                   # 유틸리티
└── App.vue
```

---

## 3. 페이지 요약

| 카테고리 | 페이지 수 | 주요 페이지 |
|----------|:---------:|-------------|
| 인증 | 2 | 로그인, 회원가입 |
| 대시보드 | 1 | 홈 |
| 객실 | 5 | 목록, 상세, 수정, 생성, 현황 |
| 객실 그룹 | 4 | 목록, 상세, 수정, 생성 |
| 예약 | 4 | 목록, 상세, 수정, 생성 |
| 월세 | 4 | 목록, 상세, 수정, 생성 |
| 결제 수단 | 1 | 목록 |
| 통계 | 1 | 대시보드 |
| 관리자 | 3 | 계정, 개발테스트, 디버그 |
| 에러 | 2 | 403, 404 |
| 프로필 | 1 | 내 정보 |

---

## 4. 컴포넌트 카테고리

| 카테고리 | 파일 수 | 설명 |
|----------|:-------:|------|
| auth | 3 | 로그인/회원가입 폼 |
| room | 7 | 객실 테이블, 카드 |
| room-group | 7 | 객실 그룹 컴포넌트 |
| reservation | 3 | 예약 테이블, 에디터 |
| payment-method | 2 | 결제 수단 |
| account | 3 | 계정 관리 |
| stats | 11 | 통계 차트, 카드 |
| dashboard | 1 | 대시보드 요약 |
| my | 1 | 프로필 |

---

## 5. 라우트 구조

```typescript
// 공개 라우트
/login
/register

// 인증 필요
/                   # 대시보드
/my                 # 프로필
/room-status        # 객실 현황

// ADMIN 필요
/rooms/*            # 객실 관리
/room-groups/*      # 객실 그룹
/reservations/*     # 예약 관리
/monthly-rents/*    # 월세
/payment-methods    # 결제 수단
/stats              # 통계
/admin/accounts     # 계정 관리

// SUPER_ADMIN 필요
/admin/dev-test     # 개발 테스트
/debug              # 디버그
```

---

## 6. 라우트 가드

| 가드 | 역할 |
|------|------|
| authenticate-guard | 인증 확인, 미인증 시 로그인 리다이렉트 |
| role-guard | 역할 기반 접근 제어 |
| account-info-loader | 사용자 정보 로드 |

---

## 7. 개발 명령어

```bash
docker compose exec frontend sh

yarn dev          # 개발 서버
yarn build        # 빌드
yarn lint         # 린트
yarn test:unit    # 테스트
```

---

## 8. 환경 설정

### 8.1 API 프록시 (개발)

```javascript
// quasar.config.js
devServer: {
  port: 9000,
  proxy: {
    '/api': {
      target: 'http://api-core:8080',
      changeOrigin: true
    }
  }
}
```

### 8.2 모바일 빌드

```bash
yarn android:build      # Android 빌드
yarn cap:open-android   # Android Studio 열기
```
