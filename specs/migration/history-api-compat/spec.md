---
id: history-api-compat
title: "History API 완전 호환"
status: completed
type: migration
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: []
related: [api-response-compat]
replaces: null
replaced_by: null

impact: both
risk: high
effort: medium

changelog:
  - date: 2026-01-07
    description: "TODO.md Phase 1에서 마이그레이션 (완료 상태)"
---

# History API 완전 호환

> api-core의 History API 응답을 프론트엔드 기대 형식과 완전히 일치시키는 작업

---

## 1. 개요

### 1.1 배경 및 문제

프론트엔드에서 변경 이력 기능을 사용하고 있으며, api-core의 응답 형식이 api-legacy와 다름:
- api-legacy: Hibernate Envers 기반 `Revision<T>` 형식
- api-core: `AuditLogResponse` 형식

### 1.2 목표

- 프론트엔드 `Revision<T>` 스키마와 완전 호환되는 응답 형식
- HistoryService 추상화 레이어 구현
- 기존 Envers 데이터 처리 방침 결정

### 1.3 비목표 (Non-Goals)

- 기존 Hibernate Envers 히스토리 데이터 마이그레이션 (포기 결정)

---

## 2. 상세 설계

### 2.1 프론트엔드 기대 응답 형식

```typescript
// apps/frontend-web/src/schema/revision.ts
type Revision<T> = {
  entity: T;                    // 해당 시점의 엔티티 스냅샷
  historyType: HistoryType;     // "CREATED" | "UPDATED" | "DELETED"
  historyCreatedAt: string;     // 변경 시각
  updatedFields: string[];      // 변경된 필드 목록
};
```

### 2.2 구현된 DTO

```go
// dto/revision.go
type RoomRevisionResponse struct {
    Entity           RoomResponse `json:"entity"`
    HistoryType      string       `json:"historyType"`      // CREATED, UPDATED, DELETED
    HistoryCreatedAt string       `json:"historyCreatedAt"`
    UpdatedFields    []string     `json:"updatedFields"`
}

type ReservationRevisionResponse struct {
    Entity           ReservationResponse `json:"entity"`
    HistoryType      string              `json:"historyType"`
    HistoryCreatedAt string              `json:"historyCreatedAt"`
    UpdatedFields    []string            `json:"updatedFields"`
}
```

### 2.3 HistoryService 추상화

```go
// services/history_service.go
type HistoryService interface {
    GetRoomHistory(ctx context.Context, roomID int64) ([]RoomRevisionResponse, error)
    GetReservationHistory(ctx context.Context, reservationID int64) ([]ReservationRevisionResponse, error)
}
```

### 2.4 Hibernate Envers 기존 데이터 결정

**결정**: 기존 Envers 히스토리 데이터 포기, 새로운 `audit_logs` 테이블만 사용
- 기존 `*_history` 테이블 데이터는 마이그레이션하지 않음
- 마이그레이션 시점 이후부터 새로 쌓이는 데이터만 사용
- 향후 필요시 별도 마이그레이션 스크립트로 처리 가능

---

## 3. API 엔드포인트

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/rooms/{id}/histories` | 객실 변경 이력 |
| GET | `/api/v1/reservations/{id}/histories` | 예약 변경 이력 |

---

## 4. 완료 조건

- [x] 프론트엔드 `Revision<T>` 스키마 분석
- [x] api-legacy `EntityRevisionDto` 응답 구조 확인
- [x] `dto/revision.go` 생성하여 Revision DTO 정의
- [x] `services/history_service.go` 생성
- [x] `handlers/room_handler.go` - `GetRoomHistories` 응답 형식 수정
- [x] `handlers/reservation_handler.go` - `GetReservationHistories` 응답 형식 수정
- [x] `cmd/server/main.go` - `HistoryService` 의존성 주입
- [x] Room Handler 테스트 (MockHistoryService)
- [x] Reservation Handler 테스트 (MockHistoryService)
- [ ] 프론트엔드 연동 테스트 (수동)

---

## 5. 참고 자료

- Hibernate Envers RevisionType: 0=ADD, 1=MOD, 2=DEL
- 프론트엔드 HistoryType: CREATED, UPDATED, DELETED
