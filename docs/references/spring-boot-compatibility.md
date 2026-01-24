# Spring Boot 호환성 이슈 및 해결책

> api-core (Go)가 api-legacy (Spring Boot)와 호환되기 위한 주요 이슈들

---

## 1. 페이지네이션

### 문제

- Spring Boot: 0 기반 페이지
- 초기 Go 구현: 1 기반 페이지

### 해결

```go
// PaginationQuery를 0 기반으로 수정
type PaginationQuery struct {
    Page int `form:"page" binding:"min=0"`
    Size int `form:"size" binding:"min=1,max=2000"`
}
```

### 테스트

`common_test.go`에서 Spring Boot 호환성 검증

---

## 2. 에러 응답 형식

### 문제

Spring Boot와 다른 에러 응답 구조

### 해결

```go
type ErrorResponse struct {
    Message     string   `json:"message"`
    Errors      []string `json:"errors,omitempty"`
    FieldErrors []string `json:"fieldErrors,omitempty"`
}
```

---

## 3. 데이터베이스 컬럼명

### 문제

모델과 실제 DB 컬럼명 불일치

### 해결

GORM 태그로 실제 컬럼명 매핑:

```go
type Reservation struct {
    GuestName   string `gorm:"column:guest_name"`
    GuestPhone  string `gorm:"column:guest_phone"`
    CheckInAt   time.Time `gorm:"column:check_in_at"`
    CheckOutAt  time.Time `gorm:"column:check_out_at"`
    AdultCount  int    `gorm:"column:adult_count"`
}
```

---

## 4. 날짜 파라미터명

### 문제

- api-legacy: `stayStartAt`/`stayEndAt`
- 초기 api-core: `startDate`/`endDate`

### 해결

프론트엔드가 사용하는 이름으로 통일

---

## 5. 시간 형식

### 문제

타임존 처리 차이

### 해결

Spring Boot 기본값과 일치하는 타임존 없는 JSON 타임스탬프 사용:

```
2024-01-15T10:30:00  (타임존 없음)
```

---

## 6. JSON 키 누락

### 문제

Go의 `omitempty` 태그로 인해 null/빈 값 필드가 누락

```go
// ❌ 값이 없으면 키 자체가 사라짐
Note *string `json:"note,omitempty"`
```

### 해결

```go
// ✅ 키 항상 포함
Note *string `json:"note"`
```

---

## 검증 방법

```bash
# API 응답 비교 스크립트
python3 scripts/compare-api-responses.py
```

자세한 내용은 [api-response-compat 스펙](../../specs/migration/api-response-compat/spec.md) 참조.
