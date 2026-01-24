---
id: api-core-architecture
title: "api-core 아키텍처"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: low
effort: small
---

# api-core 아키텍처

> Go + Gin 기반 메인 API 서버의 전체 아키텍처

---

## 1. 개요

### 1.1 역할

api-core는 Kotlin/Spring Boot 기반 api-legacy를 Go/Gin으로 마이그레이션한 메인 API 서버입니다.

### 1.2 기술 스택

| 항목 | 기술 |
|------|------|
| 언어 | Go 1.21+ |
| 웹 프레임워크 | Gin |
| ORM | GORM |
| 인증 | JWT (golang-jwt) |
| 설정 | Viper |
| 캐시 | Redis (go-redis) |
| 테스트 | testify |
| 린터 | golangci-lint |

---

## 2. 디렉토리 구조

```
apps/api-core/
├── cmd/
│   ├── server/main.go      # 애플리케이션 진입점
│   └── migrate/main.go     # DB 마이그레이션 도구
├── internal/
│   ├── config/             # 설정 관리
│   ├── handlers/           # HTTP 핸들러 (10개 파일)
│   ├── services/           # 비즈니스 로직 (9개 서비스)
│   ├── repositories/       # 데이터 액세스 (7개 리포지토리)
│   ├── models/             # GORM 모델 (12개 파일)
│   ├── dto/                # DTO (13개 파일, 49개 구조체)
│   ├── middleware/         # 미들웨어 (4개)
│   ├── migrations/         # DB 마이그레이션
│   ├── database/           # DB/Redis 연결
│   ├── audit/              # 감사 로깅
│   └── context/            # 요청 컨텍스트
├── pkg/
│   ├── auth/               # JWT 유틸리티
│   ├── response/           # 응답 헬퍼
│   └── utils/              # 공용 유틸리티
└── config/                 # 설정 파일 (YAML)
```

---

## 3. 레이어 아키텍처

```
┌─────────────────────────────────────────────────────┐
│                    Handlers                          │
│   (HTTP 요청/응답 처리, 검증, DTO 변환)              │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                    Services                          │
│   (비즈니스 로직, 트랜잭션, 규칙 적용)               │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                  Repositories                        │
│   (데이터 액세스, GORM 쿼리, 페이지네이션)           │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                    Models                            │
│   (GORM 엔티티, 관계 정의, 소프트 삭제)              │
└─────────────────────────────────────────────────────┘
```

---

## 4. 컴포넌트 요약

### 4.1 핸들러 (10개)

| 핸들러 | 역할 |
|--------|------|
| AuthHandler | 인증 (로그인, 회원가입, 토큰 갱신) |
| UserHandler | 사용자 관리 (프로필, 관리자 계정) |
| RoomHandler | 객실 CRUD + 히스토리 |
| RoomGroupHandler | 객실 그룹 CRUD |
| ReservationHandler | 예약 CRUD + 히스토리 + 통계 |
| PaymentMethodHandler | 결제 수단 CRUD |
| ConfigHandler | 설정/환경 정보 |
| HealthHandler | 헬스체크 (Actuator 호환) |
| DocsHandler | API 문서 (OpenAPI) |
| DevelopmentHandler | 테스트 데이터 생성 |

### 4.2 서비스 (9개)

| 서비스 | 역할 |
|--------|------|
| AuthService | 인증/인가 |
| UserService | 사용자 관리 |
| RoomService | 객실 관리 |
| RoomGroupService | 객실 그룹 관리 |
| ReservationService | 예약 관리 |
| PaymentMethodService | 결제 수단 관리 |
| HistoryService | 히스토리 조회 |
| ConfigService | 설정 관리 |
| DevelopmentService | 개발 도구 |

### 4.3 리포지토리 (7개)

| 리포지토리 | 메서드 수 |
|------------|:---------:|
| UserRepository | 10 |
| RoomRepository | 10 |
| RoomGroupRepository | 10 |
| ReservationRepository | 7 |
| ReservationRoomRepository | 3 |
| PaymentMethodRepository | 9 |
| LoginAttemptRepository | 4 |

---

## 5. 설정 관리

### 5.1 설정 파일

```
config/
├── application.yaml           # 기본 설정
├── application-local.yaml     # 로컬 개발
└── application-production.yaml # 운영 환경
```

### 5.2 주요 환경 변수

| 변수 | 설명 |
|------|------|
| DATABASE_MYSQL_HOST | MySQL 호스트 |
| DATABASE_MYSQL_PORT | MySQL 포트 |
| REDIS_HOST | Redis 호스트 |
| JWT_SECRET | JWT 시크릿 |
| JWT_ACCESS_TOKEN_EXPIRY | 액세스 토큰 만료(초) |
| JWT_REFRESH_TOKEN_EXPIRY | 리프레시 토큰 만료(초) |

---

## 6. Spring Boot 호환성

| 항목 | 설명 |
|------|------|
| 페이지네이션 | 0-기반 페이지 인덱스 |
| 정렬 | "field,direction" 형식 |
| 날짜 형식 | `2006-01-02T15:04:05` (타임존 없음) |
| 비밀번호 | `{bcrypt}$2a$10$...` 접두사 |
| 에러 응답 | 동일 JSON 구조 |

---

## 7. 개발 명령어

```bash
docker compose exec api-core bash

make dev      # 개발 모드 (핫 리로드)
make test     # 테스트
make lint     # 린트
make build    # 빌드
```
