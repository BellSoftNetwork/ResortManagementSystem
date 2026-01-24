---
id: api-core-payment-method
title: "api-core 결제 수단 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-core 결제 수단 관리

> 결제 수단 CRUD 및 수수료율 관리

---

## 1. 개요

### 1.1 기능 범위

- 결제 수단 목록 조회
- 결제 수단 상세 조회
- 결제 수단 생성/수정/삭제
- 기본 결제 수단 설정

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/payment_method_handler.go` | HTTP 핸들러 |
| `services/payment_method_service.go` | 비즈니스 로직 |
| `repositories/payment_method_repository.go` | 데이터 액세스 |
| `models/payment_method.go` | 결제 수단 모델 |
| `dto/payment_method.go` | DTO 정의 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/payment-methods` | 목록 | USER |
| GET | `/api/v1/payment-methods/{id}` | 상세 | USER |
| POST | `/api/v1/payment-methods` | 생성 | ADMIN |
| PATCH | `/api/v1/payment-methods/{id}` | 수정 | ADMIN |
| DELETE | `/api/v1/payment-methods/{id}` | 삭제 | ADMIN |

---

## 3. 결제 수단 모델

### 3.1 PaymentMethod 엔티티

```go
type PaymentMethod struct {
    ID                       uint                `gorm:"primaryKey"`
    Name                     string              `gorm:"size:20;not null;uniqueIndex:idx_pm_name_deleted"`
    CommissionRate           float64             `gorm:"column:commission_rate;not null"`
    RequireUnpaidAmountCheck BitBool             `gorm:"column:require_unpaid_amount_check;type:bit(1);not null;default:0"`
    IsDefaultSelect          BitBool             `gorm:"column:is_default_select;type:bit(1);not null;default:0"`
    Status                   PaymentMethodStatus `gorm:"not null"`
    CreatedAt                time.Time
    UpdatedAt                time.Time
    DeletedAt                LegacyDeletedAt
}
```

### 3.2 상태 (PaymentMethodStatus)

| 값 | 코드 | 설명 |
|----|------|------|
| -1 | INACTIVE | 비활성 |
| 1 | ACTIVE | 활성 |

### 3.3 필드 설명

| 필드 | 설명 |
|------|------|
| Name | 결제 수단명 (예: 현금, 카드, 계좌이체) |
| CommissionRate | 수수료율 (0~1, 예: 0.03 = 3%) |
| RequireUnpaidAmountCheck | 미결제 확인 필요 여부 |
| IsDefaultSelect | 기본 선택 여부 (하나만 true) |

---

## 4. 결제 수단 목록

### 4.1 요청

```http
GET /api/v1/payment-methods?page=0&size=20
Authorization: Bearer {token}
```

### 4.2 응답

```json
{
  "data": [
    {
      "id": 1,
      "name": "현금",
      "commissionRate": 0,
      "requireUnpaidAmountCheck": false,
      "isDefaultSelect": true,
      "status": "ACTIVE",
      "createdAt": "2026-01-01T00:00:00",
      "updatedAt": "2026-01-01T00:00:00"
    },
    {
      "id": 2,
      "name": "카드",
      "commissionRate": 0.03,
      "requireUnpaidAmountCheck": false,
      "isDefaultSelect": false,
      "status": "ACTIVE",
      "createdAt": "2026-01-01T00:00:00",
      "updatedAt": "2026-01-01T00:00:00"
    }
  ],
  "page": { ... }
}
```

---

## 5. 결제 수단 생성

### 5.1 요청

```http
POST /api/v1/payment-methods
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "계좌이체",
  "commissionRate": 0.01,
  "requireUnpaidAmountCheck": true
}
```

### 5.2 검증 규칙

| 필드 | 규칙 |
|------|------|
| name | 필수, 2-20자, 중복 불가 |
| commissionRate | 선택, 0~1 |
| requireUnpaidAmountCheck | 선택, boolean |

### 5.3 비즈니스 로직

1. 이름 중복 확인
2. 상태 기본값 ACTIVE 설정
3. 결제 수단 생성

---

## 6. 결제 수단 수정

### 6.1 요청

```http
PATCH /api/v1/payment-methods/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "현금 (수정)",
  "commissionRate": 0,
  "isDefaultSelect": true,
  "status": "ACTIVE"
}
```

### 6.2 수정 가능 필드

- name (중복 체크)
- commissionRate
- requireUnpaidAmountCheck
- isDefaultSelect
- status

### 6.3 기본 선택 로직

`isDefaultSelect: true` 설정 시:
1. 기존의 모든 기본 선택 해제
2. 해당 결제 수단만 기본 선택으로 설정

```go
func (r *PaymentMethodRepository) ResetAllDefaultSelects(ctx context.Context) error {
    return r.db.WithContext(ctx).
        Model(&models.PaymentMethod{}).
        Where("is_default_select = ?", true).
        Update("is_default_select", false).Error
}
```

---

## 7. 결제 수단 삭제

### 7.1 요청

```http
DELETE /api/v1/payment-methods/1
Authorization: Bearer {token}
```

### 7.2 응답

- 성공: 204 No Content

### 7.3 제약 조건

- 연결된 예약이 있는 경우 삭제 불가 (향후 구현)

---

## 8. 중개료 계산

### 8.1 계산 시점

예약 생성/수정 시 결제 수단의 수수료율을 사용하여 중개료 자동 계산:

```go
brokerFee = int(float64(paymentAmount) * paymentMethod.CommissionRate)
```

### 8.2 예시

| 결제액 | 수수료율 | 중개료 |
|--------|----------|--------|
| 100,000 | 0.03 (3%) | 3,000 |
| 200,000 | 0.05 (5%) | 10,000 |

---

## 9. 리포지토리 메서드

| 메서드 | 설명 |
|--------|------|
| `Create` | 결제 수단 생성 |
| `Update` | 결제 수단 수정 |
| `Delete` | 소프트 삭제 |
| `FindByID` | ID로 조회 |
| `FindAll` | 페이지네이션 목록 |
| `FindActive` | 활성 결제 수단만 조회 |
| `ExistsByName` | 이름 중복 확인 |
| `FindByName` | 이름으로 조회 |
| `ResetAllDefaultSelects` | 모든 기본 선택 해제 |

---

## 10. 테스트

- `payment_method_handler_test.go`
- `payment_method_service_test.go`

### 테스트 케이스

- 결제 수단 목록 조회
- 결제 수단 생성 (중복 체크)
- 결제 수단 수정
- 기본 선택 설정 (다른 것 해제)
- 결제 수단 삭제
