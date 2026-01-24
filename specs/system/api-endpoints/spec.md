---
id: current-state-api-endpoints
title: "API 엔드포인트 현황"
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

# API 엔드포인트 현황

> 리조트 관리 시스템 API의 전체 엔드포인트 문서화

---

## 1. 개요

### 1.1 API 버전

- **Base Path**: `/api/v1`
- **Content-Type**: `application/json`
- **인증**: JWT Bearer Token

### 1.2 엔드포인트 통계

| 카테고리 | 엔드포인트 수 | 인증 필요 | 관리자 전용 |
|----------|:-------------:|:---------:|:-----------:|
| 헬스체크 | 3 | ❌ | ❌ |
| 문서 | 2 | ❌ | ❌ |
| 인증 | 3 | ❌ | ❌ |
| 설정 | 2 | ❌/✅ | ❌ |
| 프로필 | 3 | ✅ | ❌ |
| 사용자 관리 | 3 | ✅ | ✅ |
| 객실 | 6 | ✅ | ✅ (쓰기) |
| 객실 그룹 | 5 | ✅ | ✅ (쓰기) |
| 예약 | 7 | ✅ | ✅ (쓰기) |
| 결제 수단 | 5 | ✅ | ✅ (쓰기) |
| 개발 도구 | 1 | ✅ | ✅ |
| **총계** | **40** | | |

---

## 2. 헬스체크 (Health Check)

Spring Boot Actuator 호환 헬스체크 엔드포인트입니다.

### 2.1 엔드포인트 목록

| Method | Path | 설명 |
|--------|------|------|
| GET | `/actuator/health` | 전체 헬스 상태 |
| GET | `/actuator/health/liveness` | 컨테이너 liveness probe |
| GET | `/actuator/health/readiness` | 컨테이너 readiness probe |

### 2.2 응답 예시

```json
// GET /actuator/health
{
  "status": "UP",
  "components": {
    "mysql": { "status": "UP" },
    "redis": { "status": "UP" }
  }
}

// GET /actuator/health/liveness
{ "status": "UP" }

// GET /actuator/health/readiness
{ "status": "UP" }
```

---

## 3. 문서 (Documentation)

### 3.1 엔드포인트 목록

| Method | Path | 설명 |
|--------|------|------|
| GET | `/docs/schema` | OpenAPI 3.0 스펙 (JSON) |
| GET | `/docs/swagger-ui` | Swagger UI |

---

## 4. 인증 (Authentication)

### 4.1 엔드포인트 목록

| Method | Path | 설명 | 인증 |
|--------|------|------|:----:|
| POST | `/api/v1/auth/register` | 회원가입 | ❌ |
| POST | `/api/v1/auth/login` | 로그인 | ❌ |
| POST | `/api/v1/auth/refresh` | 토큰 갱신 | ❌ |

### 4.2 회원가입

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "userId": "newuser",
  "email": "user@example.com",
  "name": "홍길동",
  "password": "password123"
}
```

**응답**:
```json
{
  "data": {
    "id": 1,
    "userId": "newuser",
    "email": "user@example.com",
    "name": "홍길동",
    "status": "ACTIVE",
    "role": "NORMAL",
    "profileImageUrl": "https://...",
    "createdAt": "2026-01-07T12:00:00",
    "updatedAt": "2026-01-07T12:00:00"
  }
}
```

### 4.3 로그인

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123",
  "device": {
    "osInfo": "Windows 10",
    "languageInfo": "ko-KR",
    "userAgent": "Mozilla/5.0...",
    "deviceFingerprint": "abc123"
  }
}
```

**응답**:
```json
{
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "accessTokenExpiresIn": 900
  }
}
```

### 4.4 토큰 갱신

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

## 5. 설정 (Configuration)

### 5.1 엔드포인트 목록

| Method | Path | 설명 | 인증 |
|--------|------|------|:----:|
| GET | `/api/v1/env` | 환경 정보 | ❌ |
| GET | `/api/v1/config` | 앱 설정 | ❌ |

### 5.2 환경 정보

```http
GET /api/v1/env
```

**응답**:
```json
{
  "data": {
    "profile": "local",
    "hostname": "api-core-5f8b7c...",
    "version": "1.0.0",
    "uptime": "24h 30m"
  }
}
```

### 5.3 앱 설정

```http
GET /api/v1/config
```

**응답**:
```json
{
  "data": {
    "isAvailableRegistration": true
  }
}
```

---

## 6. 프로필 (My Profile)

### 6.1 엔드포인트 목록

| Method | Path | 설명 | 인증 |
|--------|------|------|:----:|
| GET | `/api/v1/my` | 내 정보 조회 | ✅ |
| POST | `/api/v1/my` | 내 정보 조회 (호환용) | ✅ |
| PATCH | `/api/v1/my` | 내 정보 수정 | ✅ |

### 6.2 내 정보 조회

```http
GET /api/v1/my
Authorization: Bearer {accessToken}
```

**응답**:
```json
{
  "data": {
    "id": 1,
    "userId": "admin",
    "email": "admin@example.com",
    "name": "관리자",
    "status": "ACTIVE",
    "role": "SUPER_ADMIN",
    "profileImageUrl": "https://...",
    "createdAt": "2026-01-01T00:00:00",
    "updatedAt": "2026-01-07T12:00:00"
  }
}
```

### 6.3 내 정보 수정

```http
PATCH /api/v1/my
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "name": "새 이름",
  "currentPassword": "oldPassword",
  "newPassword": "newPassword123"
}
```

---

## 7. 사용자 관리 (Admin Accounts)

**권한**: ADMIN, SUPER_ADMIN

### 7.1 엔드포인트 목록

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/v1/admin/accounts` | 사용자 목록 |
| POST | `/api/v1/admin/accounts` | 사용자 생성 |
| PATCH | `/api/v1/admin/accounts/{id}` | 사용자 수정 |

### 7.2 사용자 목록

```http
GET /api/v1/admin/accounts?page=0&size=20&sort=createdAt,desc
Authorization: Bearer {accessToken}
```

**응답**:
```json
{
  "data": [
    {
      "id": 1,
      "userId": "admin",
      "email": "admin@example.com",
      "name": "관리자",
      "status": "ACTIVE",
      "role": "SUPER_ADMIN",
      "profileImageUrl": "https://...",
      "createdAt": "2026-01-01T00:00:00",
      "updatedAt": "2026-01-07T12:00:00"
    }
  ],
  "page": {
    "size": 20,
    "number": 0,
    "totalElements": 5,
    "totalPages": 1
  }
}
```

### 7.3 사용자 생성

```http
POST /api/v1/admin/accounts
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "userId": "newadmin",
  "email": "newadmin@example.com",
  "name": "새 관리자",
  "password": "password123",
  "role": "ADMIN",
  "status": "ACTIVE"
}
```

### 7.4 사용자 수정

```http
PATCH /api/v1/admin/accounts/2
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "name": "수정된 이름",
  "role": "ADMIN",
  "status": "INACTIVE"
}
```

---

## 8. 객실 (Rooms)

### 8.1 엔드포인트 목록

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/rooms` | 객실 목록 | USER |
| GET | `/api/v1/rooms/{id}` | 객실 상세 | USER |
| POST | `/api/v1/rooms` | 객실 생성 | ADMIN |
| PATCH | `/api/v1/rooms/{id}` | 객실 수정 | ADMIN |
| DELETE | `/api/v1/rooms/{id}` | 객실 삭제 | ADMIN |
| GET | `/api/v1/rooms/{id}/histories` | 객실 히스토리 | ADMIN |

### 8.2 객실 목록

```http
GET /api/v1/rooms?page=0&size=20&roomGroupId=1&status=NORMAL&search=101
Authorization: Bearer {accessToken}
```

**쿼리 파라미터**:
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| page | int | 페이지 번호 (0-based) |
| size | int | 페이지 크기 |
| sort | string | 정렬 (예: createdAt,desc) |
| roomGroupId | int | 객실 그룹 필터 |
| status | string | 상태 필터 |
| search | string | 검색어 (번호) |

**응답**:
```json
{
  "data": [
    {
      "id": 1,
      "number": "101",
      "roomGroupId": 1,
      "roomGroup": {
        "id": 1,
        "name": "스탠다드"
      },
      "note": "",
      "status": "NORMAL",
      "createdAt": "2026-01-01T00:00:00",
      "updatedAt": "2026-01-07T12:00:00",
      "createdBy": { "id": 1, "userId": "admin", "name": "관리자" },
      "updatedBy": { "id": 1, "userId": "admin", "name": "관리자" }
    }
  ],
  "page": { ... }
}
```

### 8.3 객실 생성

```http
POST /api/v1/rooms
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "number": "201",
  "roomGroupId": 1,
  "note": "바다 전망",
  "status": "NORMAL"
}
```

### 8.4 객실 히스토리

```http
GET /api/v1/rooms/1/histories?page=0&size=20
Authorization: Bearer {accessToken}
```

**응답**:
```json
{
  "data": [
    {
      "entity": { ... },
      "historyType": "UPDATED",
      "historyCreatedAt": "2026-01-07T14:00:00",
      "updatedFields": ["note", "status"]
    }
  ],
  "page": { ... }
}
```

---

## 9. 객실 그룹 (Room Groups)

### 9.1 엔드포인트 목록

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/room-groups` | 객실 그룹 목록 | USER |
| GET | `/api/v1/room-groups/{id}` | 객실 그룹 상세 | USER |
| POST | `/api/v1/room-groups` | 객실 그룹 생성 | ADMIN |
| PATCH | `/api/v1/room-groups/{id}` | 객실 그룹 수정 | ADMIN |
| DELETE | `/api/v1/room-groups/{id}` | 객실 그룹 삭제 | ADMIN |

### 9.2 객실 그룹 상세

```http
GET /api/v1/room-groups/1
Authorization: Bearer {accessToken}
```

**응답**:
```json
{
  "data": {
    "id": 1,
    "name": "스탠다드",
    "peekPrice": 150000,
    "offPeekPrice": 100000,
    "description": "기본 객실",
    "rooms": [
      {
        "room": { "id": 1, "number": "101", ... },
        "lastReservation": { ... }
      }
    ],
    "createdAt": "2026-01-01T00:00:00",
    "createdBy": { ... },
    "updatedAt": "2026-01-07T12:00:00",
    "updatedBy": { ... }
  }
}
```

### 9.3 객실 그룹 생성

```http
POST /api/v1/room-groups
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "name": "디럭스",
  "peekPrice": 250000,
  "offPeekPrice": 180000,
  "description": "넓은 객실"
}
```

---

## 10. 예약 (Reservations)

### 10.1 엔드포인트 목록

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/reservations` | 예약 목록 | USER |
| GET | `/api/v1/reservations/{id}` | 예약 상세 | USER |
| POST | `/api/v1/reservations` | 예약 생성 | ADMIN |
| PATCH | `/api/v1/reservations/{id}` | 예약 수정 | ADMIN |
| DELETE | `/api/v1/reservations/{id}` | 예약 삭제 | ADMIN |
| GET | `/api/v1/reservations/{id}/histories` | 예약 히스토리 | ADMIN |
| GET | `/api/v1/reservation-statistics` | 예약 통계 | USER |

### 10.2 예약 목록

```http
GET /api/v1/reservations?page=0&size=20&status=NORMAL&type=STAY&stayStartAt=2026-01-01&stayEndAt=2026-01-31&roomId=1&search=홍길동
Authorization: Bearer {accessToken}
```

**쿼리 파라미터**:
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| status | string | 상태 (PENDING, NORMAL, CANCEL, REFUND) |
| type | string | 유형 (STAY, MONTHLY_RENT) |
| roomId | int | 객실 ID |
| stayStartAt | date | 입실일 시작 |
| stayEndAt | date | 입실일 종료 |
| search | string | 검색어 (이름, 전화번호) |

**응답**:
```json
{
  "data": [
    {
      "id": 1,
      "paymentMethodId": 1,
      "paymentMethod": { "id": 1, "name": "현금" },
      "rooms": [{ "id": 1, "number": "101" }],
      "name": "홍길동",
      "phone": "010-1234-5678",
      "peopleCount": 2,
      "stayStartAt": "2026-01-10",
      "stayEndAt": "2026-01-12",
      "checkInAt": "2026-01-10T15:00:00",
      "checkOutAt": null,
      "price": 200000,
      "deposit": 50000,
      "paymentAmount": 150000,
      "refundAmount": 0,
      "brokerFee": 0,
      "note": "",
      "canceledAt": null,
      "status": "NORMAL",
      "type": "STAY",
      "createdAt": "2026-01-05T10:00:00",
      "updatedAt": "2026-01-10T15:00:00",
      "createdBy": { ... },
      "updatedBy": { ... }
    }
  ],
  "page": { ... }
}
```

### 10.3 예약 생성

```http
POST /api/v1/reservations
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "paymentMethodId": 1,
  "roomIds": [1, 2],
  "name": "홍길동",
  "phone": "010-1234-5678",
  "peopleCount": 4,
  "stayStartAt": "2026-02-01",
  "stayEndAt": "2026-02-03",
  "price": 400000,
  "deposit": 100000,
  "paymentAmount": 300000,
  "note": "특별 요청 사항",
  "type": "STAY"
}
```

### 10.4 예약 통계

```http
GET /api/v1/reservation-statistics?startDate=2026-01-01&endDate=2026-12-31&periodType=MONTHLY
Authorization: Bearer {accessToken}
```

**쿼리 파라미터**:
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| startDate | date | 시작일 (필수) |
| endDate | date | 종료일 (필수) |
| periodType | string | 기간 유형 (DAILY, MONTHLY, YEARLY) |

**응답**:
```json
{
  "data": {
    "periodType": "MONTHLY",
    "stats": [
      {
        "period": "2026-01",
        "totalSales": 5000000,
        "totalReservations": 25,
        "totalGuests": 60
      }
    ],
    "monthlyStats": [ ... ]
  }
}
```

---

## 11. 결제 수단 (Payment Methods)

### 11.1 엔드포인트 목록

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/payment-methods` | 결제 수단 목록 | USER |
| GET | `/api/v1/payment-methods/{id}` | 결제 수단 상세 | USER |
| POST | `/api/v1/payment-methods` | 결제 수단 생성 | ADMIN |
| PATCH | `/api/v1/payment-methods/{id}` | 결제 수단 수정 | ADMIN |
| DELETE | `/api/v1/payment-methods/{id}` | 결제 수단 삭제 | ADMIN |

### 11.2 결제 수단 생성

```http
POST /api/v1/payment-methods
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "name": "카드",
  "commissionRate": 0.03,
  "requireUnpaidAmountCheck": false
}
```

**응답**:
```json
{
  "data": {
    "id": 2,
    "name": "카드",
    "commissionRate": 0.03,
    "requireUnpaidAmountCheck": false,
    "isDefaultSelect": false,
    "status": "ACTIVE",
    "createdAt": "2026-01-07T12:00:00",
    "updatedAt": "2026-01-07T12:00:00"
  }
}
```

---

## 12. 개발 도구 (Development)

**권한**: SUPER_ADMIN (비운영 환경만)

### 12.1 엔드포인트 목록

| Method | Path | 설명 |
|--------|------|------|
| POST | `/api/v1/dev/test-data` | 테스트 데이터 생성 |

### 12.2 테스트 데이터 생성

```http
POST /api/v1/dev/test-data
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "type": "all",
  "reservationOptions": {
    "startDate": "2026-01-01",
    "endDate": "2026-12-31",
    "regularReservations": 100,
    "monthlyReservations": 10
  }
}
```

---

## 13. 에러 응답

### 13.1 에러 응답 형식

```json
{
  "message": "에러 메시지",
  "errors": ["상세 에러 1", "상세 에러 2"],
  "fieldErrors": ["필드명: 검증 에러 메시지"]
}
```

### 13.2 HTTP 상태 코드

| 코드 | 설명 |
|------|------|
| 200 | 성공 |
| 201 | 생성 성공 |
| 204 | 삭제 성공 (응답 본문 없음) |
| 400 | 잘못된 요청 |
| 401 | 인증 실패 |
| 403 | 권한 없음 |
| 404 | 리소스 없음 |
| 409 | 충돌 (중복) |
| 429 | 요청 과다 (브루트포스 방지) |
| 500 | 서버 에러 |

### 13.3 검증 에러 메시지 (한글)

```json
{
  "fieldErrors": [
    "이름: 2자 이상 20자 이하로 입력해주세요",
    "이메일: 올바른 이메일 형식이 아닙니다",
    "비밀번호: 8자 이상 입력해주세요"
  ]
}
```

---

## 14. 참고 자료

- [docs/guides/api-testing.md](../../../docs/guides/api-testing.md)
- [docs/contracts/api-comparison.md](../../../docs/contracts/api-comparison.md)
