---
id: current-state-data-model
title: "데이터 모델 현황"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: []
related: [current-state-overview, current-state-api-core, current-state-api-legacy]
replaces: null
replaced_by: null

impact: both
risk: low
effort: small

changelog:
  - date: 2026-01-07
    description: "현재 구현 상태 기반 초기 스펙 작성"
---

# 데이터 모델 현황

> 리조트 관리 시스템의 데이터베이스 스키마 및 모델 문서화

---

## 1. 개요

### 1.1 데이터베이스

| 항목 | 값 |
|------|-----|
| DBMS | MySQL 8.0 |
| 문자셋 | utf8mb4 |
| 콜레이션 | utf8mb4_unicode_ci |

### 1.2 공유 데이터베이스

api-core와 api-legacy는 동일한 MySQL 인스턴스의 다른 데이터베이스를 사용:
- `rms-core`: api-core 전용
- `rms-legacy`: api-legacy 전용

---

## 2. ERD (Entity Relationship Diagram)

```
┌─────────────────┐
│      User       │
├─────────────────┤
│ id              │◄──────────────────────────────────────────┐
│ user_id         │                                           │
│ email           │                                           │
│ name            │                                           │
│ password        │                                           │
│ status          │                                           │
│ role            │                                           │
│ created_at      │                                           │
│ updated_at      │                                           │
│ deleted_at      │                                           │
└─────────────────┘                                           │
         │                                                    │
         │ created_by / updated_by                            │
         ▼                                                    │
┌─────────────────┐       ┌─────────────────┐                │
│   RoomGroup     │       │      Room       │                │
├─────────────────┤       ├─────────────────┤                │
│ id              │◄──────│ id              │                │
│ name            │       │ number          │                │
│ peek_price      │       │ room_group_id   │────────────────┤
│ off_peek_price  │       │ note            │                │
│ description     │       │ status          │                │
│ created_by      │───────│ created_by      │────────────────┤
│ updated_by      │───────│ updated_by      │────────────────┤
│ created_at      │       │ created_at      │                │
│ updated_at      │       │ updated_at      │                │
│ deleted_at      │       │ deleted_at      │                │
└─────────────────┘       └─────────────────┘                │
                                   │                          │
                                   │                          │
┌─────────────────┐       ┌─────────────────┐                │
│ PaymentMethod   │       │ReservationRoom  │                │
├─────────────────┤       ├─────────────────┤                │
│ id              │◄──┐   │ id              │                │
│ name            │   │   │ reservation_id  │────────────────┤
│ commission_rate │   │   │ room_id         │────────────────┤
│ require_unpaid  │   │   │ created_by      │────────────────┤
│ is_default      │   │   │ updated_by      │────────────────┤
│ status          │   │   │ created_at      │                │
│ created_at      │   │   │ updated_at      │                │
│ updated_at      │   │   │ deleted_at      │                │
│ deleted_at      │   │   └─────────────────┘                │
└─────────────────┘   │            │                          │
         │            │            │                          │
         │            │            │                          │
         │            │   ┌─────────────────┐                │
         │            │   │   Reservation   │                │
         │            │   ├─────────────────┤                │
         │            └───│ id              │                │
         │                │ payment_method  │                │
         └────────────────│ name            │                │
                          │ phone           │                │
                          │ people_count    │                │
                          │ stay_start_at   │                │
                          │ stay_end_at     │                │
                          │ check_in_at     │                │
                          │ check_out_at    │                │
                          │ price           │                │
                          │ deposit         │                │
                          │ payment_amount  │                │
                          │ refund_amount   │                │
                          │ broker_fee      │                │
                          │ note            │                │
                          │ canceled_at     │                │
                          │ status          │                │
                          │ type            │                │
                          │ created_by      │────────────────┤
                          │ updated_by      │────────────────┘
                          │ created_at      │
                          │ updated_at      │
                          │ deleted_at      │
                          └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│  LoginAttempt   │       │  RevisionInfo   │
├─────────────────┤       ├─────────────────┤
│ id              │       │ id              │
│ username        │       │ created_at      │
│ ip_address      │       └─────────────────┘
│ successful      │
│ attempt_at      │       ┌─────────────────┐
│ os_info         │       │   AuditLog      │
│ language_info   │       ├─────────────────┤
│ user_agent      │       │ id              │
│ device_finger   │       │ entity_type     │
└─────────────────┘       │ entity_id       │
                          │ action          │
                          │ old_values      │
                          │ new_values      │
                          │ changed_fields  │
                          │ user_id         │
                          │ username        │
                          │ created_at      │
                          └─────────────────┘
```

---

## 3. 테이블 상세

### 3.1 user (사용자)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| user_id | VARCHAR(30) | ❌ | | 로그인 아이디 |
| email | VARCHAR(100) | ⭕ | | 이메일 |
| name | VARCHAR(20) | ❌ | | 이름 |
| password | VARCHAR(100) | ❌ | | 비밀번호 (bcrypt) |
| status | TINYINT | ❌ | -1 | 상태 (-1:INACTIVE, 1:ACTIVE) |
| role | TINYINT | ❌ | 0 | 역할 (0:NORMAL, 100:ADMIN, 127:SUPER_ADMIN) |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 (소프트 삭제) |

**인덱스**:
- UNIQUE: (user_id, deleted_at)
- UNIQUE: (email, deleted_at)

---

### 3.2 room_group (객실 그룹)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| name | VARCHAR(20) | ❌ | | 그룹명 |
| peek_price | INT | ❌ | | 성수기 가격 |
| off_peek_price | INT | ❌ | | 비수기 가격 |
| description | VARCHAR(200) | ❌ | '' | 설명 |
| created_by | BIGINT | ❌ | | 생성자 FK |
| updated_by | BIGINT | ❌ | | 수정자 FK |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 |

**인덱스**:
- UNIQUE: (name, deleted_at)

**외래키**:
- created_by → user(id)
- updated_by → user(id)

---

### 3.3 room (객실)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| number | VARCHAR(10) | ❌ | | 객실 번호 |
| room_group_id | BIGINT | ❌ | | 객실 그룹 FK |
| note | VARCHAR(200) | ❌ | '' | 메모 |
| status | TINYINT | ❌ | 1 | 상태 |
| created_by | BIGINT | ❌ | | 생성자 FK |
| updated_by | BIGINT | ❌ | | 수정자 FK |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 |

**인덱스**:
- UNIQUE: (number, deleted_at)

**외래키**:
- room_group_id → room_group(id)
- created_by → user(id)
- updated_by → user(id)

**상태값**:
| 값 | 코드 | 설명 |
|----|------|------|
| -10 | DAMAGED | 파손 |
| -1 | CONSTRUCTION | 공사중 |
| 0 | INACTIVE | 비활성 |
| 1 | NORMAL | 정상 |

---

### 3.4 reservation (예약)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| payment_method_id | BIGINT | ❌ | | 결제 수단 FK |
| name | VARCHAR(30) | ❌ | | 예약자명 |
| phone | VARCHAR(20) | ❌ | | 전화번호 |
| people_count | INT | ❌ | 0 | 인원수 |
| stay_start_at | DATE | ❌ | | 입실일 |
| stay_end_at | DATE | ❌ | | 퇴실일 |
| check_in_at | DATETIME | ⭕ | | 체크인 시간 |
| check_out_at | DATETIME | ⭕ | | 체크아웃 시간 |
| price | INT | ❌ | | 총 가격 |
| deposit | INT | ❌ | 0 | 보증금 |
| payment_amount | INT | ❌ | 0 | 결제액 |
| refund_amount | INT | ❌ | 0 | 환불액 |
| broker_fee | INT | ❌ | 0 | 중개료 |
| note | VARCHAR(200) | ❌ | '' | 메모 |
| canceled_at | DATETIME | ⭕ | | 취소일시 |
| status | TINYINT | ❌ | | 상태 |
| type | TINYINT | ❌ | | 유형 |
| created_by | BIGINT | ❌ | | 생성자 FK |
| updated_by | BIGINT | ❌ | | 수정자 FK |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 |

**외래키**:
- payment_method_id → payment_method(id)
- created_by → user(id)
- updated_by → user(id)

**상태값**:
| 값 | 코드 | 설명 |
|----|------|------|
| -10 | REFUND | 환불 |
| -1 | CANCEL | 취소 |
| 0 | PENDING | 대기 |
| 1 | NORMAL | 정상 |

**유형값**:
| 값 | 코드 | 설명 |
|----|------|------|
| 0 | STAY | 숙박 |
| 10 | MONTHLY_RENT | 월세 |

---

### 3.5 reservation_room (예약-객실 연결)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| reservation_id | BIGINT | ❌ | | 예약 FK |
| room_id | BIGINT | ❌ | | 객실 FK |
| created_by | BIGINT | ❌ | | 생성자 FK |
| updated_by | BIGINT | ❌ | | 수정자 FK |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 |

**인덱스**:
- UNIQUE: (reservation_id, room_id, deleted_at)

**외래키**:
- reservation_id → reservation(id)
- room_id → room(id)
- created_by → user(id)
- updated_by → user(id)

---

### 3.6 payment_method (결제 수단)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| name | VARCHAR(20) | ❌ | | 결제 수단명 |
| commission_rate | DOUBLE | ❌ | | 수수료율 (0~1) |
| require_unpaid_amount_check | BIT(1) | ❌ | 0 | 미결제 확인 필요 |
| is_default_select | BIT(1) | ❌ | 0 | 기본 선택 여부 |
| status | TINYINT | ❌ | | 상태 |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성일시 |
| updated_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 수정일시 |
| deleted_at | DATETIME | ❌ | '1970-01-01 00:00:00' | 삭제일시 |

**인덱스**:
- UNIQUE: (name, deleted_at)

**상태값**:
| 값 | 코드 | 설명 |
|----|------|------|
| -1 | INACTIVE | 비활성 |
| 1 | ACTIVE | 활성 |

---

### 3.7 login_attempts (로그인 시도)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| username | VARCHAR(50) | ❌ | | 사용자명 |
| ip_address | VARCHAR(50) | ❌ | | IP 주소 |
| successful | TINYINT(1) | ❌ | | 성공 여부 |
| attempt_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 시도 시간 |
| os_info | VARCHAR(50) | ⭕ | | OS 정보 |
| language_info | VARCHAR(50) | ⭕ | | 언어 정보 |
| user_agent | VARCHAR(500) | ⭕ | | User-Agent |
| device_fingerprint | VARCHAR(50) | ⭕ | | 디바이스 지문 |

**인덱스**:
- (username, attempt_at)
- (ip_address, attempt_at)
- (username, ip_address, attempt_at)

---

### 3.8 revision_info (리비전 정보 - api-legacy)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| created_at | DATETIME | ❌ | | 리비전 생성 시간 |

---

### 3.9 audit_logs (감사 로그 - api-core)

| 컬럼 | 타입 | Null | 기본값 | 설명 |
|------|------|:----:|--------|------|
| id | BIGINT | ❌ | AUTO_INCREMENT | PK |
| entity_type | VARCHAR(50) | ❌ | | 엔티티 타입 |
| entity_id | BIGINT | ❌ | | 엔티티 ID |
| action | VARCHAR(20) | ❌ | | 액션 (CREATE/UPDATE/DELETE) |
| old_values | JSON | ⭕ | | 이전 값 |
| new_values | JSON | ⭕ | | 새 값 |
| changed_fields | JSON | ⭕ | | 변경된 필드 |
| user_id | BIGINT | ⭕ | | 사용자 ID |
| username | VARCHAR(50) | ❌ | '' | 사용자명 |
| created_at | DATETIME | ❌ | CURRENT_TIMESTAMP | 생성 시간 |

**인덱스**:
- (entity_type, entity_id)

---

## 4. 히스토리 테이블 (api-legacy Envers)

Hibernate Envers가 자동 생성하는 히스토리 테이블:

| 테이블 | 원본 테이블 | 설명 |
|--------|-------------|------|
| room_history | room | 객실 변경 이력 |
| room_group_history | room_group | 객실 그룹 변경 이력 |
| reservation_history | reservation | 예약 변경 이력 |
| reservation_room_history | reservation_room | 예약-객실 연결 변경 이력 |

히스토리 테이블 추가 컬럼:
- `rev`: revision_info.id 참조
- `revtype`: 변경 유형 (0=INSERT, 1=UPDATE, 2=DELETE)
- `*_mod`: 변경 플래그 (withModifiedFlag=true일 때)

---

## 5. 소프트 삭제 패턴

### 5.1 구현 방식

```sql
-- 활성 레코드 조건
WHERE deleted_at = '1970-01-01 00:00:00'

-- 삭제 시
UPDATE table SET deleted_at = NOW() WHERE id = ?
```

### 5.2 유니크 제약조건

소프트 삭제와 유니크 제약조건을 함께 사용하기 위해 `deleted_at`을 유니크 인덱스에 포함:

```sql
UNIQUE KEY (number, deleted_at)  -- 활성 레코드 중 유니크
```

---

## 6. 관계 요약

| 관계 | 타입 | 설명 |
|------|------|------|
| RoomGroup → Room | 1:N | 객실 그룹은 여러 객실 포함 |
| Reservation → PaymentMethod | N:1 | 예약은 하나의 결제 수단 |
| Reservation ↔ Room | N:M | reservation_room 통해 연결 |
| User → Room | 1:N | 생성자/수정자 |
| User → RoomGroup | 1:N | 생성자/수정자 |
| User → Reservation | 1:N | 생성자/수정자 |
| User → ReservationRoom | 1:N | 생성자/수정자 |

---

## 7. 마이그레이션 관리

### 7.1 api-legacy (Liquibase)

```yaml
# db.changelog-master.yaml
databaseChangeLog:
  - include:
      file: 001-initial-schema.yaml
  - include:
      file: 002-add-payment-method.yaml
  # ... ~40개 체인지셋
```

### 7.2 api-core (Go Migrations)

```go
// internal/migrations/migrations.go
var Migrations = []Migration{
    {Version: 1, Name: "initial_schema", Up: migration001Up, Down: migration001Down},
    {Version: 2, Name: "add_audit_logs", Up: migration002Up, Down: migration002Down},
}
```

---

## 8. 참고 자료

- [docs/guides/database-migration.md](../../../docs/guides/database-migration.md)
- [docs/references/hibernate-envers.md](../../../docs/references/hibernate-envers.md)
