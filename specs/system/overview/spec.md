---
id: current-state-overview
title: "시스템 현황 개요"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: []
related: [api-core, api-legacy, frontend-web, data-model, api-endpoints]
replaces: null
replaced_by: null

impact: both
risk: low
effort: small

changelog:
  - date: 2026-01-07
    description: "현재 구현 상태 기반 초기 스펙 작성"
---

# 시스템 현황 개요

> 리조트 통합 관리 시스템(RMS)의 전체 아키텍처와 현재 구현 상태를 문서화

---

## 1. 개요

### 1.1 시스템 소개

리조트 통합 관리 시스템(Resort Management System, RMS)은 리조트의 객실 예약, 결제, 사용자 관리 등을 통합 관리하는 웹 애플리케이션입니다.

### 1.2 아키텍처

```
┌─────────────────────────────────────────────────────────────────┐
│                        frontend-web                              │
│                   Vue.js 3 + Quasar Framework                   │
│                      (Port: 9000)                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (Proxy)                         │
└─────────────────────────────────────────────────────────────────┘
                    │                       │
                    ▼                       ▼
┌───────────────────────────┐   ┌───────────────────────────────┐
│       api-core            │   │        api-legacy             │
│    Go + Gin (메인)        │   │  Kotlin + Spring Boot (레거시)│
│     (Port: 8080)          │   │      (Port: 8081)             │
└───────────────────────────┘   └───────────────────────────────┘
                    │                       │
                    └───────────┬───────────┘
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Shared Infrastructure                        │
│          MySQL 8.0 (Port: 3306)  │  Redis 7 (Port: 6379)        │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 기술 스택

| 컴포넌트 | api-core | api-legacy | frontend-web |
|----------|----------|------------|--------------|
| 언어 | Go 1.21+ | Kotlin | TypeScript |
| 프레임워크 | Gin | Spring Boot 3.x | Vue.js 3 + Quasar 2 |
| ORM | GORM | JPA/Hibernate | - |
| 인증 | JWT | JWT (Spring Security) | Pinia Store |
| 테스트 | testify | Kotest | Vitest |

---

## 2. 앱 구성

### 2.1 apps/api-core (Go + Gin)

**역할**: 메인 API 서버 (마이그레이션 대상)

**아키텍처 패턴**: Clean Architecture

```
internal/
├── config/         # 설정 관리 (Viper)
├── handlers/       # HTTP 핸들러 (10개)
├── services/       # 비즈니스 로직 (9개)
├── repositories/   # 데이터 액세스 (7개)
├── models/         # GORM 모델 (12개)
├── dto/            # 요청/응답 DTO (49개 구조체)
├── middleware/     # 미들웨어 (4개)
├── migrations/     # DB 마이그레이션
└── audit/          # 감사 로깅
```

**핵심 기능**:
- 완전한 CRUD API (35+ 엔드포인트)
- JWT 인증 및 역할 기반 접근 제어
- 브루트포스 공격 방지
- 히스토리/감사 로그 API
- Spring Boot 호환 응답 형식

### 2.2 apps/api-legacy (Kotlin + Spring Boot)

**역할**: 레거시 API 서버 (마이그레이션 원본)

**아키텍처 패턴**: DDD (Domain-Driven Design)

```
src/main/kotlin/net/bellsoft/rms/
├── authentication/     # 인증 (Spring Security)
├── common/             # 공통 (엔티티, 예외, 응답)
├── main/               # 설정 서비스
├── migration/          # Liquibase 마이그레이션
├── payment/            # 결제 수단
├── reservation/        # 예약 관리
├── revision/           # Hibernate Envers
├── room/               # 객실 관리
└── user/               # 사용자 관리
```

**핵심 기능**:
- Hibernate Envers 기반 히스토리 추적
- Liquibase DB 마이그레이션
- MapStruct DTO 매핑
- Kotest 기반 테스트

### 2.3 apps/frontend-web (Vue.js 3 + Quasar)

**역할**: 웹 프론트엔드 SPA

**구조**:

```
src/
├── pages/          # 페이지 컴포넌트 (28개)
├── components/     # 재사용 컴포넌트 (70+개)
├── stores/         # Pinia 상태관리 (3개)
├── api/            # API 통신 레이어
├── router/         # Vue Router 설정
├── schema/         # TypeScript 타입 정의
└── layouts/        # 레이아웃 컴포넌트
```

**핵심 기능**:
- 역할 기반 라우팅 (NORMAL, ADMIN, SUPER_ADMIN)
- JWT 자동 갱신
- 네트워크 상태 모니터링
- 반응형 모바일 지원 (Capacitor)

---

## 3. 핵심 도메인

### 3.1 사용자 관리

| 기능 | 설명 |
|------|------|
| 회원가입 | 이메일/아이디 기반 가입 |
| 로그인 | JWT 토큰 발급 (Access 15분, Refresh 7일) |
| 역할 | NORMAL, ADMIN, SUPER_ADMIN |
| 상태 | ACTIVE, INACTIVE |

### 3.2 객실 관리

| 기능 | 설명 |
|------|------|
| 객실 그룹 | 객실 카테고리 (성수기/비수기 가격) |
| 객실 | 개별 객실 (번호, 상태, 메모) |
| 상태 | NORMAL, INACTIVE, CONSTRUCTION, DAMAGED |
| 히스토리 | 변경 이력 추적 |

### 3.3 예약 관리

| 기능 | 설명 |
|------|------|
| 예약 유형 | STAY (숙박), MONTHLY_RENT (월세) |
| 예약 상태 | PENDING, NORMAL, CANCEL, REFUND |
| 결제 | 가격, 보증금, 결제액, 환불액, 중개료 |
| 히스토리 | 변경 이력 추적 |

### 3.4 결제 수단

| 기능 | 설명 |
|------|------|
| 수수료율 | 결제 수단별 수수료 |
| 기본 선택 | 기본 결제 수단 설정 |
| 미결제 확인 | 미결제 확인 필요 여부 |

---

## 4. 마이그레이션 현황

### 4.1 완료된 항목

| 영역 | 상태 | 비고 |
|------|:----:|------|
| JWT 인증 | ✅ | Spring Security 호환 |
| 기본 CRUD | ✅ | 모든 엔티티 |
| 페이지네이션 | ✅ | 0-기반, Spring 호환 |
| 에러 처리 | ✅ | 한글 검증 메시지 |
| History API | ✅ | audit_logs 테이블 |
| 브루트포스 방지 | ✅ | login_attempts |

### 4.2 진행 중/예정

| 영역 | 상태 | 우선순위 |
|------|:----:|:--------:|
| API 응답 호환성 검증 | 🚧 | 높음 |
| DB 스키마 통합 | 📋 | 중간 |
| 운영 환경 전환 | 📋 | 높음 |
| 레거시 정리 | 📋 | 낮음 |

---

## 5. 개발 환경

### 5.1 Docker Compose 구성

```yaml
services:
  mysql:      # MySQL 8.0 (Port: 3306)
  redis:      # Redis 7 (Port: 6379)
  api-core:   # Go API (Port: 8080)
  api-legacy: # Spring Boot API (Port: 8081)
  frontend:   # Vue.js (Port: 9000)
```

### 5.2 주요 명령어

```bash
# 전체 환경 시작
docker compose up -d

# api-core 테스트
docker compose exec api-core make test

# api-legacy 테스트
docker compose exec api-legacy ./gradlew test

# frontend 빌드
docker compose exec frontend yarn build
```

---

## 6. 관련 스펙 문서

| 스펙 | 설명 |
|------|------|
| [api-core](../api-core/spec.md) | api-core 상세 구현 |
| [api-legacy](../api-legacy/spec.md) | api-legacy 상세 구현 |
| [frontend-web](../frontend-web/spec.md) | 프론트엔드 상세 구현 |
| [data-model](../data-model/spec.md) | 데이터 모델 상세 |
| [api-endpoints](../api-endpoints/spec.md) | API 엔드포인트 상세 |

---

## 7. 참고 자료

- [AGENTS.md](../../../AGENTS.md): 개발 가이드
- [README.md](../../../README.md): 사용자/설치 가이드
- [docs/](../../../docs/): 기술 가이드
