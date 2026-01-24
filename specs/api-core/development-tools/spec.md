---
id: api-core-development-tools
title: "api-core 개발 도구"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: low
effort: small
---

# api-core 개발 도구

> 테스트 데이터 생성 및 개발 지원 기능

---

## 1. 개요

### 1.1 기능 범위

- 테스트 데이터 생성 (결제 수단, 객실 그룹, 객실, 예약)
- 개발/스테이징 환경에서만 사용 가능

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/development_handler.go` | HTTP 핸들러 |
| `services/development_service.go` | 테스트 데이터 생성 로직 |
| `dto/development.go` | DTO 정의 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| POST | `/api/v1/dev/test-data` | 테스트 데이터 생성 | SUPER_ADMIN |

### 2.1 접근 제한

- SUPER_ADMIN 역할 필요
- 프로덕션 환경에서 접근 불가 (DevelopmentOnlyMiddleware)

---

## 3. 테스트 데이터 생성

### 3.1 요청

```http
POST /api/v1/dev/test-data
Authorization: Bearer {token}
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

### 3.2 생성 타입

| 타입 | 설명 |
|------|------|
| essential | 필수 데이터만 (결제 수단, 객실 그룹, 객실) |
| reservation | 예약 데이터만 |
| all | 모든 데이터 |

### 3.3 응답

```json
{
  "data": {
    "message": "테스트 데이터 생성 완료",
    "data": {
      "paymentMethods": 3,
      "roomGroups": 3,
      "rooms": 15,
      "reservations": 110
    }
  }
}
```

---

## 4. 생성되는 데이터

### 4.1 필수 데이터 (essential)

#### 결제 수단

| 이름 | 수수료율 | 기본 선택 |
|------|----------|:---------:|
| 현금 | 0% | ✅ |
| 카드 | 3% | ❌ |
| 계좌이체 | 1% | ❌ |

#### 객실 그룹

| 이름 | 성수기 | 비수기 |
|------|--------|--------|
| 스탠다드 | 100,000 | 80,000 |
| 디럭스 | 150,000 | 120,000 |
| 스위트 | 250,000 | 200,000 |

#### 객실

| 그룹 | 객실 번호 |
|------|-----------|
| 스탠다드 | 101, 102, 103, 104, 105 |
| 디럭스 | 201, 202, 203, 204, 205 |
| 스위트 | 301, 302, 303, 304, 305 |

### 4.2 예약 데이터 (reservation)

#### 옵션

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| startDate | 1년 전 | 예약 시작일 |
| endDate | 현재 | 예약 종료일 |
| regularReservations | 50 | 숙박 예약 수 |
| monthlyReservations | 5 | 월세 예약 수 |

#### 생성 로직

```go
// 랜덤 예약 생성
for i := 0; i < options.RegularReservations; i++ {
    reservation := Reservation{
        PaymentMethodID: randomPaymentMethod(),
        Name:           randomName(),
        Phone:          randomPhone(),
        PeopleCount:    rand.Intn(4) + 1,
        StayStartAt:    randomDate(options.StartDate, options.EndDate),
        StayEndAt:      stayStartAt.Add(time.Duration(rand.Intn(3)+1) * 24 * time.Hour),
        Price:          randomPrice(80000, 300000),
        Status:         NORMAL,
        Type:           STAY,
    }
    // 랜덤 객실 할당
    reservation.Rooms = randomRooms(1, 3)
}
```

---

## 5. 서비스 구현

### 5.1 GenerateTestData

```go
func (s *DevelopmentService) GenerateTestData(
    ctx context.Context, 
    request *dto.GenerateTestDataRequest,
) (*dto.GenerateTestDataResponse, error) {
    result := make(map[string]interface{})
    
    switch request.Type {
    case "essential":
        s.generateEssentialData(ctx, result)
    case "reservation":
        s.generateReservationData(ctx, request.ReservationOptions, result)
    case "all":
        s.generateEssentialData(ctx, result)
        s.generateReservationData(ctx, request.ReservationOptions, result)
    }
    
    return &dto.GenerateTestDataResponse{
        Message: "테스트 데이터 생성 완료",
        Data:    result,
    }, nil
}
```

### 5.2 중복 방지

- 기존 데이터가 있으면 건너뛰기
- 이름 기반 중복 체크 (결제 수단, 객실 그룹, 객실)

---

## 6. DTO

### 6.1 GenerateTestDataRequest

```go
type GenerateDataType string

const (
    GenerateDataTypeEssential   GenerateDataType = "essential"
    GenerateDataTypeReservation GenerateDataType = "reservation"
    GenerateDataTypeAll         GenerateDataType = "all"
)

type GenerateTestDataRequest struct {
    Type               GenerateDataType             `json:"type" binding:"required,oneof=essential reservation all"`
    ReservationOptions *ReservationGenerationOptions `json:"reservationOptions,omitempty"`
}

type ReservationGenerationOptions struct {
    StartDate           *time.Time `json:"startDate"`
    EndDate             *time.Time `json:"endDate"`
    RegularReservations *int       `json:"regularReservations"`
    MonthlyReservations *int       `json:"monthlyReservations"`
}
```

### 6.2 GenerateTestDataResponse

```go
type GenerateTestDataResponse struct {
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}
```

---

## 7. 테스트

- `development_handler_test.go`

### 테스트 케이스

- essential 타입으로 필수 데이터 생성
- reservation 타입으로 예약 데이터 생성
- all 타입으로 모든 데이터 생성
- 프로덕션 환경에서 접근 거부
- ADMIN 역할로 접근 거부 (SUPER_ADMIN만 허용)
