# Resort Management System 스펙 인덱스

> **마지막 업데이트**: 2026-01-07

모든 스펙을 관리하는 중앙 인덱스입니다.

---

## 구조

```
specs/
├── api-core/           # api-core (Go/Gin) 관련 스펙
├── api-legacy/         # api-legacy (Kotlin/Spring) 관련 스펙
├── frontend-web/       # frontend-web (Vue.js/Quasar) 관련 스펙
├── system/             # 시스템 전반 (아키텍처, 데이터 모델, API)
├── migration/          # Kotlin → Go 마이그레이션 스펙
├── infra/              # 인프라 스펙
└── _templates/         # 스펙/플랜 템플릿
```

---

## api-core (Go/Gin)

| 스펙 | 제목 | 상태 | 설명 |
|------|------|:----:|------|
| [architecture](./api-core/architecture/) | 아키텍처 | completed | 프로젝트 구조, 레이어 아키텍처, 의존성 |
| [authentication](./api-core/authentication/) | 인증 | completed | JWT 인증, 토큰 관리, 로그인/로그아웃 |
| [user-management](./api-core/user-management/) | 사용자 관리 | completed | CRUD, 역할 기반 접근 제어, 프로필 |
| [room-management](./api-core/room-management/) | 객실 관리 | completed | 객실 CRUD, 상태 관리, 가용성 |
| [room-group-management](./api-core/room-group-management/) | 객실 그룹 관리 | completed | 그룹 CRUD, 객실 배정 |
| [reservation-management](./api-core/reservation-management/) | 예약 관리 | completed | 예약 CRUD, 상태 변경, 필터링 |
| [payment-method](./api-core/payment-method/) | 결제 수단 | completed | 결제 수단 CRUD |
| [middleware](./api-core/middleware/) | 미들웨어 | completed | 인증, CORS, 에러 핸들링, 로깅 |
| [development-tools](./api-core/development-tools/) | 개발 도구 | completed | DB 시딩, 마이그레이션, 빌드 |

---

## api-legacy (Kotlin/Spring Boot)

| 스펙 | 제목 | 상태 | 설명 |
|------|------|:----:|------|
| [architecture](./api-legacy/architecture/) | 아키텍처 | completed | 프로젝트 구조, 레이어 아키텍처, 의존성 |
| [authentication](./api-legacy/authentication/) | 인증 | completed | JWT 인증, Spring Security |
| [user-management](./api-legacy/user-management/) | 사용자 관리 | completed | CRUD, 역할 기반 접근 제어 |
| [room-management](./api-legacy/room-management/) | 객실 관리 | completed | 객실 CRUD, 상태 관리 |
| [room-group-management](./api-legacy/room-group-management/) | 객실 그룹 관리 | completed | 그룹 CRUD, 객실 배정 |
| [reservation-management](./api-legacy/reservation-management/) | 예약 관리 | completed | 예약 CRUD, 상태 변경 |
| [payment-method](./api-legacy/payment-method/) | 결제 수단 | completed | 결제 수단 CRUD |
| [history-envers](./api-legacy/history-envers/) | Envers 히스토리 | completed | 감사 로그, 변경 이력 추적 |

---

## frontend-web (Vue.js/Quasar)

| 스펙 | 제목 | 상태 | 설명 |
|------|------|:----:|------|
| [architecture](./frontend-web/architecture/) | 아키텍처 | completed | 프로젝트 구조, Vue 3, Quasar |
| [authentication](./frontend-web/authentication/) | 인증 | completed | 로그인/로그아웃, 토큰 관리 |
| [state-management](./frontend-web/state-management/) | 상태 관리 | completed | Pinia 스토어, 전역 상태 |
| [api-integration](./frontend-web/api-integration/) | API 연동 | completed | Axios, 인터셉터, 에러 핸들링 |
| [dashboard](./frontend-web/dashboard/) | 대시보드 | completed | 메인 화면, 위젯, 통계 요약 |
| [room-management](./frontend-web/room-management/) | 객실 관리 | completed | 객실 목록, 상세, 편집 |
| [reservation-management](./frontend-web/reservation-management/) | 예약 관리 | completed | 예약 목록, 캘린더, 상태 관리 |
| [statistics](./frontend-web/statistics/) | 통계 | completed | 매출, 점유율 차트 |
| [admin-management](./frontend-web/admin-management/) | 관리자 관리 | completed | 사용자/결제수단/시스템 설정 |
| [unit-testing](./frontend-web/unit-testing/) | 유닛 테스트 | draft | 프론트엔드 테스트 전략 |

---

## system (시스템 전반)

| 스펙 | 제목 | 상태 | 설명 |
|------|------|:----:|------|
| [overview](./system/overview/) | 시스템 개요 | completed | 아키텍처, 기술 스택, 마이그레이션 현황 |
| [data-model](./system/data-model/) | 데이터 모델 | completed | DB 스키마, ERD, 테이블 정의 |
| [api-endpoints](./system/api-endpoints/) | API 엔드포인트 | completed | 전체 API 문서 (40개 엔드포인트) |

---

## migration (마이그레이션)

Kotlin Spring Boot → Golang Gin 마이그레이션 관련 스펙

| 스펙 | 제목 | Spec | Plan | 우선순위 |
|------|------|:----:|:----:|:--------:|
| [history-api-compat](./migration/history-api-compat/) | History API 완전 호환 | completed | - | - |
| [api-response-compat](./migration/api-response-compat/) | API 응답 완전 호환성 검증 | approved | not-started | 높음 |
| [db-schema-unification](./migration/db-schema-unification/) | DB 스키마 버전 관리 통합 | draft | not-started | 중간 |
| [production-cutover](./migration/production-cutover/) | 운영 환경 전환 준비 | draft | not-started | 높음 |
| [legacy-cleanup](./migration/legacy-cleanup/) | 레거시 정리 | draft | not-started | 낮음 |

---

## infra (인프라)

| 스펙 | 제목 | Spec | Plan | 우선순위 |
|------|------|:----:|:----:|:--------:|
| [gitlab-ci-restructure](./infra/gitlab-ci-restructure/) | GitLab CI 빌드 파이프라인 재구성 | approved | in-progress | 중간 |

---

## 스펙 상태 정의

### spec.md 상태

| 상태 | 설명 |
|------|------|
| `draft` | 초안 작성 중, 미해결 질문 있음 |
| `approved` | 검토 완료, 구현 가능 |
| `completed` | 구현 완료, 코드와 일치 |
| `deprecated` | 더 이상 유효하지 않음 |

### plan.md 상태

| 상태 | 설명 |
|------|------|
| `not-started` | 아직 시작 안 함 |
| `in-progress` | 구현 진행 중 |
| `blocked` | 차단 요소로 중단 |
| `done` | 구현 완료 |

---

## 사용법

### 새 스펙 생성

```bash
# 컴포넌트별 기능 스펙
/new-spec api-core caching "Redis 캐싱 전략"
/new-spec frontend-web dashboard-redesign "대시보드 리디자인"

# 마이그레이션/인프라 스펙
/new-spec migration db-cleanup "DB 정리"
```

### 스펙 상태 확인

```bash
/spec-status
```

---

## 참고 문서

- [AGENTS.md](../AGENTS.md): 개발 가이드
- [docs/](../docs/): 기술 가이드
- [_templates/](./_templates/): 스펙/플랜 템플릿
