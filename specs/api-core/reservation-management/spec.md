---
id: api-core-reservation-management
title: "api-core 예약 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: high
effort: medium
---

# api-core 예약 관리

> 예약 CRUD, 통계, 변경 이력

---

## 1. 개요

### 1.1 기능 범위

- 예약 목록 조회 (필터, 검색, 페이지네이션)
- 예약 상세 조회
- 예약 생성/수정/삭제
- 예약 통계 조회
- 예약 변경 히스토리 조회

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/reservation_handler.go` | HTTP 핸들러 |
| `services/reservation_service.go` | 비즈니스 로직 |
| `services/history_service.go` | 히스토리 조회 |
| `repositories/reservation_repository.go` | 데이터 액세스 |
| `models/reservation.go` | 예약 모델 |
| `models/reservation_room.go` | 예약-객실 연결 모델 |
| `dto/reservation.go` | DTO 정의 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/reservations` | 예약 목록 | USER |
| GET | `/api/v1/reservations/{id}` | 예약 상세 | USER |
| POST | `/api/v1/reservations` | 예약 생성 | ADMIN |
| PATCH | `/api/v1/reservations/{id}` | 예약 수정 | ADMIN |
| DELETE | `/api/v1/reservations/{id}` | 예약 삭제 | ADMIN |
| GET | `/api/v1/reservations/{id}/histories` | 변경 이력 | ADMIN |
| GET | `/api/v1/reservation-statistics` | 예약 통계 | USER |

---

## 3. 예약 모델

### 3.1 Reservation 엔티티

```go
type Reservation struct {
    ID              uint              `gorm:"primaryKey"`
    PaymentMethodID uint              `gorm:"not null"`
    PaymentMethod   PaymentMethod     `gorm:"foreignKey:PaymentMethodID"`
    Rooms           []ReservationRoom `gorm:"foreignKey:ReservationID"`
    Name            string            `gorm:"size:30;not null"`
    Phone           string            `gorm:"size:20;not null"`
    PeopleCount     int               `gorm:"not null;default:0"`
    StayStartAt     time.Time         `gorm:"type:date;not null"`
    StayEndAt       time.Time         `gorm:"type:date;not null"`
    CheckInAt       *time.Time        `gorm:"type:datetime"`
    CheckOutAt      *time.Time        `gorm:"type:datetime"`
    Price           int               `gorm:"not null"`
    Deposit         int               `gorm:"not null;default:0"`
    PaymentAmount   int               `gorm:"not null;default:0"`
    RefundAmount    int               `gorm:"not null;default:0"`
    BrokerFee       int               `gorm:"not null;default:0"`
    Note            string            `gorm:"size:200;not null;default:''"`
    CanceledAt      *time.Time
    Status          ReservationStatus `gorm:"not null"`
    Type            ReservationType   `gorm:"not null"`
    CreatedBy       uint              `gorm:"not null"`
    UpdatedBy       uint              `gorm:"not null"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       LegacyDeletedAt
}
```

### 3.2 예약 상태 (ReservationStatus)

| 값 | 코드 | 설명 |
|----|------|------|
| -10 | REFUND | 환불 |
| -1 | CANCEL | 취소 |
| 0 | PENDING | 대기 |
| 1 | NORMAL | 정상 |

### 3.3 예약 유형 (ReservationType)

| 값 | 코드 | 설명 |
|----|------|------|
| 0 | STAY | 숙박 |
| 10 | MONTHLY_RENT | 월세 |

---

## 4. 예약 목록

### 4.1 요청

```http
GET /api/v1/reservations?page=0&size=20&status=NORMAL&type=STAY&stayStartAt=2026-01-01&stayEndAt=2026-01-31&roomId=1&search=홍길동
Authorization: Bearer {token}
```

### 4.2 쿼리 파라미터

| 파라미터 | 타입 | 설명 |
|----------|------|------|
| page | int | 페이지 번호 (0-based) |
| size | int | 페이지 크기 |
| sort | string | 정렬 |
| status | string | 상태 필터 |
| type | string | 유형 필터 (STAY, MONTHLY_RENT) |
| roomId | int | 객실 ID 필터 |
| stayStartAt | date | 입실일 시작 |
| stayEndAt | date | 입실일 종료 |
| search | string | 검색어 (이름, 전화번호) |

### 4.3 응답

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

---

## 5. 예약 생성

### 5.1 요청

```http
POST /api/v1/reservations
Authorization: Bearer {token}
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

### 5.2 검증 규칙

| 필드 | 규칙 |
|------|------|
| paymentMethodId | 필수, 존재하는 결제 수단 |
| roomIds | 선택, 존재하는 객실 |
| name | 필수, 1-30자 |
| phone | 필수, 1-20자 |
| peopleCount | 선택, 0 이상 |
| stayStartAt | 필수 |
| stayEndAt | 필수, stayStartAt 이후 |
| price | 선택, 0 이상 |
| deposit | 선택, 0 이상 |
| paymentAmount | 선택, 0 이상 |
| note | 선택, 최대 200자 |
| type | 선택, STAY/MONTHLY_RENT |

### 5.3 비즈니스 로직

1. 결제 수단 존재 확인
2. 각 객실 존재 및 가용성 확인
3. 중개료 자동 계산 (paymentAmount × commissionRate)
4. 예약 생성 (트랜잭션)
5. ReservationRoom 연결 생성
6. audit_logs에 CREATE 기록

---

## 6. 예약 수정

### 6.1 요청

```http
PATCH /api/v1/reservations/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "checkInAt": "2026-02-01T15:00:00",
  "paymentAmount": 400000,
  "status": "NORMAL"
}
```

### 6.2 수정 가능 필드

- paymentMethodId, roomIds
- name, phone, peopleCount
- stayStartAt, stayEndAt
- checkInAt, checkOutAt
- price, deposit, paymentAmount, refundAmount
- note, status, type

### 6.3 비즈니스 로직

1. 예약 존재 확인
2. 날짜/객실 변경 시 가용성 재확인
3. 중개료 재계산 (결제액 또는 결제 수단 변경 시)
4. 예약 수정
5. audit_logs에 UPDATE 기록

---

## 7. 예약 삭제

### 7.1 요청

```http
DELETE /api/v1/reservations/1
Authorization: Bearer {token}
```

### 7.2 응답

- 성공: 204 No Content

### 7.3 비즈니스 로직

1. 예약 존재 확인
2. 소프트 삭제
3. 연결된 ReservationRoom도 소프트 삭제
4. audit_logs에 DELETE 기록

---

## 8. 예약 통계

### 8.1 요청

```http
GET /api/v1/reservation-statistics?startDate=2026-01-01&endDate=2026-12-31&periodType=MONTHLY
Authorization: Bearer {token}
```

### 8.2 쿼리 파라미터

| 파라미터 | 타입 | 설명 |
|----------|------|------|
| startDate | date | 시작일 (필수) |
| endDate | date | 종료일 (필수) |
| periodType | string | DAILY, MONTHLY, YEARLY |

### 8.3 응답

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
      },
      {
        "period": "2026-02",
        "totalSales": 4500000,
        "totalReservations": 22,
        "totalGuests": 55
      }
    ],
    "monthlyStats": [ ... ]
  }
}
```

### 8.4 집계 로직

```sql
SELECT 
    DATE_FORMAT(stay_start_at, '%Y-%m') as period,
    SUM(price) as total_sales,
    COUNT(*) as total_reservations,
    SUM(people_count) as total_guests
FROM reservation
WHERE stay_start_at BETWEEN ? AND ?
  AND status = 1
  AND deleted_at = '1970-01-01 00:00:00'
GROUP BY DATE_FORMAT(stay_start_at, '%Y-%m')
ORDER BY period
```

---

## 9. 예약 히스토리

### 9.1 요청

```http
GET /api/v1/reservations/1/histories?page=0&size=20
Authorization: Bearer {token}
```

### 9.2 응답

```json
{
  "data": [
    {
      "entity": { ... },
      "historyType": "UPDATED",
      "historyCreatedAt": "2026-01-10T15:00:00",
      "updatedFields": ["checkInAt", "status"]
    }
  ],
  "page": { ... }
}
```

---

## 10. 테스트

- `reservation_handler_test.go`
- `reservation_service_test.go`

### 테스트 케이스

- 예약 목록 조회 (필터, 페이지네이션)
- 예약 상세 조회
- 예약 생성 (객실 가용성 체크)
- 예약 수정 (중개료 재계산)
- 예약 삭제
- 통계 조회 (기간별)
- 히스토리 조회
