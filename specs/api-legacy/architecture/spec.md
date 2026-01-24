---
id: api-legacy-architecture
title: "api-legacy 아키텍처"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: low
effort: small
---

# api-legacy 아키텍처

> Kotlin + Spring Boot 기반 레거시 API 서버 아키텍처

---

## 1. 개요

### 1.1 역할

api-legacy는 원본 API 서버로, 현재 api-core로 마이그레이션 진행 중입니다.

### 1.2 기술 스택

| 항목 | 기술 |
|------|------|
| 언어 | Kotlin |
| 프레임워크 | Spring Boot 3.x |
| ORM | JPA/Hibernate |
| 히스토리 | Hibernate Envers |
| 인증 | Spring Security + JWT |
| 쿼리 | QueryDSL |
| DTO 매핑 | MapStruct |
| DB 마이그레이션 | Liquibase |
| 테스트 | Kotest |
| 린터 | Ktlint |

---

## 2. 디렉토리 구조

```
apps/api-legacy/src/main/kotlin/net/bellsoft/rms/
├── authentication/         # 인증 모듈
│   ├── controller/
│   ├── service/
│   ├── filter/
│   ├── entity/
│   └── dto/
├── common/                 # 공통 모듈
│   ├── entity/             # 베이스 엔티티
│   ├── controller/         # 예외 처리
│   ├── dto/
│   └── config/
├── main/                   # 메인 모듈
│   ├── controller/
│   ├── service/
│   └── dto/
├── payment/                # 결제 모듈
├── reservation/            # 예약 모듈
├── revision/               # Envers 모듈
├── room/                   # 객실 모듈
├── user/                   # 사용자 모듈
└── migration/              # Liquibase 커스텀
```

---

## 3. 레이어 아키텍처

```
┌─────────────────────────────────────────────────────┐
│                   Controllers                        │
│   (@RestController, @RequestMapping, @Valid)        │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                    Services                          │
│   (@Service, @Transactional, Interface/Impl)        │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                  Repositories                        │
│   (JpaRepository, QueryDSL, Custom)                 │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│                    Entities                          │
│   (@Entity, @Audited, @SQLDelete, @Where)           │
└─────────────────────────────────────────────────────┘
```

---

## 4. 컨트롤러

| 컨트롤러 | 경로 | 역할 |
|----------|------|------|
| AuthController | /api/v1/auth | 인증 |
| MainController | /api/v1 | 환경/설정 |
| MyController | /api/v1/my | 프로필 |
| AdminAccountController | /api/v1/admin/accounts | 계정 관리 |
| RoomController | /api/v1/rooms | 객실 |
| RoomGroupController | /api/v1/room-groups | 객실 그룹 |
| ReservationController | /api/v1/reservations | 예약 |
| PaymentMethodController | /api/v1/payment-methods | 결제 수단 |

---

## 5. 서비스

| 인터페이스 | 구현체 | 역할 |
|------------|--------|------|
| AuthService | AuthServiceImpl | 인증 |
| ConfigService | ConfigServiceImpl | 설정 |
| UserService | UserServiceImpl | 사용자 |
| RoomService | RoomServiceImpl | 객실 |
| RoomGroupService | RoomGroupServiceImpl | 객실 그룹 |
| ReservationService | ReservationServiceImpl | 예약 |
| PaymentMethodService | PaymentMethodServiceImpl | 결제 수단 |
| LoginAttemptService | - | 로그인 시도 추적 |

---

## 6. 베이스 엔티티

```kotlin
@MappedSuperclass
abstract class BaseEntity {
    @Id @GeneratedValue(strategy = IDENTITY)
    val id: Long = 0
}

@MappedSuperclass
@EntityListeners(AuditingEntityListener::class)
@Audited
abstract class BaseTimeEntity : BaseEntity() {
    @CreatedDate
    var createdAt: LocalDateTime
    
    @LastModifiedDate
    var updatedAt: LocalDateTime
    
    var deletedAt: LocalDateTime = LocalDateTime.of(1970, 1, 1, 0, 0, 0)
}

@MappedSuperclass
abstract class BaseMustAuditEntity : BaseTimeEntity() {
    @CreatedBy @ManyToOne(fetch = LAZY)
    var createdBy: User? = null
    
    @LastModifiedBy @ManyToOne(fetch = LAZY)
    var updatedBy: User? = null
}
```

---

## 7. 응답 래퍼

```kotlin
data class SingleResponse<T>(val data: T)

data class ListResponse<T>(
    val data: List<T>,
    val page: PageDto
)

data class PageDto(
    val size: Int,
    val number: Int,
    val totalElements: Long,
    val totalPages: Int
)

data class ErrorResponse(
    val message: String,
    val errors: List<String>? = null,
    val fieldErrors: List<String>? = null
)
```

---

## 8. 개발 명령어

```bash
docker compose exec api-legacy bash

./gradlew bootRun       # 실행
./gradlew test          # 테스트
./gradlew ktlintCheck   # 린트 검사
./gradlew ktlintFormat  # 린트 수정
./gradlew build         # 빌드
```
