---
id: frontend-unit-testing
title: "프론트엔드 유닛 테스트 작성"
status: draft
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: []
related: []
replaces: null
replaced_by: null

impact: frontend
risk: low
effort: large

changelog:
  - date: 2026-01-07
    description: "TODO.md에서 마이그레이션"
---

# 프론트엔드 유닛 테스트 작성

> 현재 테스트가 없는 프론트엔드에 주요 기능별 유닛 테스트 추가

---

## 1. 개요

### 1.1 배경 및 문제

현재 프론트엔드(apps/frontend-web)에 유닛 테스트가 하나도 없는 상태:
- 리팩토링 시 안전성 부족
- 회귀 버그 감지 불가
- CI/CD에서 품질 검증 부재

### 1.2 목표

- 주요 컴포넌트 테스트 추가
- Pinia 스토어 테스트 추가
- 유틸리티 함수 테스트 추가
- API 호출 모킹 및 테스트

### 1.3 비목표 (Non-Goals)

- E2E 테스트 (별도 계획)
- 100% 커버리지 (주요 기능만)

---

## 2. 테스트 대상

### 2.1 Pinia 스토어

| 스토어 | 우선순위 |
|--------|:--------:|
| auth.store | 🔴 높음 |
| reservation.store | 🔴 높음 |
| room.store | 🟡 중간 |
| user.store | 🟡 중간 |

### 2.2 컴포넌트

| 컴포넌트 | 우선순위 |
|---------|:--------:|
| ReservationForm | 🔴 높음 |
| RoomSelector | 🔴 높음 |
| DateRangePicker | 🟡 중간 |
| PaymentMethodSelect | 🟡 중간 |

### 2.3 유틸리티

| 유틸리티 | 우선순위 |
|---------|:--------:|
| 날짜 포맷팅 | 🔴 높음 |
| 가격 계산 | 🔴 높음 |
| 유효성 검증 | 🟡 중간 |

---

## 3. 테스트 환경

### 3.1 도구

- **테스트 프레임워크**: Vitest
- **컴포넌트 테스트**: @vue/test-utils
- **API 모킹**: msw (Mock Service Worker)

### 3.2 설정

```typescript
// vitest.config.ts
export default defineConfig({
  test: {
    globals: true,
    environment: 'happy-dom',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
    },
  },
})
```

---

## 4. 실행 방법

```bash
# Docker 컨테이너에서 실행
docker compose exec frontend yarn test
docker compose exec frontend yarn test:coverage
```

---

## 5. 완료 조건

- [ ] Vitest 설정
- [ ] auth.store 테스트
- [ ] reservation.store 테스트
- [ ] ReservationForm 컴포넌트 테스트
- [ ] 날짜/가격 유틸리티 테스트
- [ ] CI 파이프라인에 테스트 추가
- [ ] 최소 50% 커버리지 달성
