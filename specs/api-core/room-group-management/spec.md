---
id: api-core-room-group-management
title: "api-core 객실 그룹 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-core 객실 그룹 관리

> 객실 그룹(카테고리) CRUD

---

## 1. 개요

### 1.1 기능 범위

- 객실 그룹 목록 조회
- 객실 그룹 상세 조회 (포함 객실 + 마지막 예약 정보)
- 객실 그룹 생성/수정/삭제

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/room_group_handler.go` | HTTP 핸들러 |
| `services/room_group_service.go` | 비즈니스 로직 |
| `repositories/room_group_repository.go` | 데이터 액세스 |
| `models/room_group.go` | 객실 그룹 모델 |
| `dto/room_group.go` | DTO 정의 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/room-groups` | 그룹 목록 | USER |
| GET | `/api/v1/room-groups/{id}` | 그룹 상세 | USER |
| POST | `/api/v1/room-groups` | 그룹 생성 | ADMIN |
| PATCH | `/api/v1/room-groups/{id}` | 그룹 수정 | ADMIN |
| DELETE | `/api/v1/room-groups/{id}` | 그룹 삭제 | ADMIN |

---

## 3. 객실 그룹 모델

### 3.1 RoomGroup 엔티티

```go
type RoomGroup struct {
    ID           uint        `gorm:"primaryKey"`
    Name         string      `gorm:"size:20;not null;uniqueIndex:idx_room_group_name_deleted"`
    PeekPrice    int         `gorm:"column:peek_price;not null"`
    OffPeekPrice int         `gorm:"column:off_peek_price;not null"`
    Description  string      `gorm:"size:200;not null;default:''"`
    Rooms        []Room      `gorm:"foreignKey:RoomGroupID"`
    CreatedBy    uint        `gorm:"not null"`
    UpdatedBy    uint        `gorm:"not null"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    LegacyDeletedAt
}
```

### 3.2 필드 설명

| 필드 | 설명 |
|------|------|
| Name | 그룹명 (예: 스탠다드, 디럭스) |
| PeekPrice | 성수기 가격 |
| OffPeekPrice | 비수기 가격 |
| Description | 그룹 설명 |
| Rooms | 포함된 객실 목록 |

---

## 4. 객실 그룹 목록

### 4.1 요청

```http
GET /api/v1/room-groups?page=0&size=20&sort=name,asc
Authorization: Bearer {token}
```

### 4.2 응답

```json
{
  "data": [
    {
      "id": 1,
      "name": "스탠다드",
      "peekPrice": 150000,
      "offPeekPrice": 100000,
      "description": "기본 객실",
      "rooms": [],
      "createdAt": "2026-01-01T00:00:00",
      "createdBy": { "id": 1, "name": "관리자" },
      "updatedAt": "2026-01-07T12:00:00",
      "updatedBy": { "id": 1, "name": "관리자" }
    }
  ],
  "page": { ... }
}
```

---

## 5. 객실 그룹 상세

### 5.1 요청

```http
GET /api/v1/room-groups/1
Authorization: Bearer {token}
```

### 5.2 응답

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
        "room": {
          "id": 1,
          "number": "101",
          "status": "NORMAL",
          "note": ""
        },
        "lastReservation": {
          "id": 10,
          "name": "홍길동",
          "stayStartAt": "2026-01-05",
          "stayEndAt": "2026-01-07",
          "status": "NORMAL"
        }
      },
      {
        "room": { "id": 2, "number": "102", ... },
        "lastReservation": null
      }
    ],
    "createdAt": "2026-01-01T00:00:00",
    "createdBy": { ... },
    "updatedAt": "2026-01-07T12:00:00",
    "updatedBy": { ... }
  }
}
```

### 5.3 특이사항

- 각 객실의 마지막 예약 정보 포함
- 객실 상태별 필터링 가능 (쿼리 파라미터)

---

## 6. 객실 그룹 생성

### 6.1 요청

```http
POST /api/v1/room-groups
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "디럭스",
  "peekPrice": 250000,
  "offPeekPrice": 180000,
  "description": "넓은 객실"
}
```

### 6.2 검증 규칙

| 필드 | 규칙 |
|------|------|
| name | 필수, 1-20자, 중복 불가 |
| peekPrice | 선택, 0 이상 |
| offPeekPrice | 선택, 0 이상 |
| description | 선택, 최대 200자 |

### 6.3 비즈니스 로직

1. 그룹명 중복 확인
2. 감사 정보 설정
3. 그룹 생성

---

## 7. 객실 그룹 수정

### 7.1 요청

```http
PATCH /api/v1/room-groups/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "스탠다드 플러스",
  "peekPrice": 180000
}
```

### 7.2 수정 가능 필드

- name (중복 체크)
- peekPrice
- offPeekPrice
- description

---

## 8. 객실 그룹 삭제

### 8.1 요청

```http
DELETE /api/v1/room-groups/1
Authorization: Bearer {token}
```

### 8.2 응답

- 성공: 204 No Content
- 실패: 404 Not Found

### 8.3 제약 조건

- 포함된 객실이 있는 경우 삭제 불가 (409 Conflict)

---

## 9. 리포지토리 메서드

| 메서드 | 설명 |
|--------|------|
| `Create` | 그룹 생성 |
| `Update` | 그룹 수정 |
| `Delete` | 소프트 삭제 |
| `FindByID` | ID로 조회 |
| `FindByIDWithRooms` | ID로 조회 (객실 포함) |
| `FindAll` | 페이지네이션 목록 |
| `FindByIDWithUsers` | 생성자/수정자 포함 조회 |
| `FindAllWithUsers` | 목록 (사용자 정보 포함) |
| `ExistsByName` | 이름 중복 확인 |
| `FindByName` | 이름으로 조회 |

---

## 10. 테스트

- `room_group_handler_test.go`
- `room_group_service_test.go`

### 테스트 케이스

- 그룹 목록 조회
- 그룹 상세 조회 (객실 + 마지막 예약)
- 그룹 생성 (중복 체크)
- 그룹 수정
- 그룹 삭제 (객실 있을 때 거부)
