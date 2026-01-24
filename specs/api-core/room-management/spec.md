---
id: api-core-room-management
title: "api-core 객실 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-core 객실 관리

> 객실 CRUD 및 변경 이력 관리

---

## 1. 개요

### 1.1 기능 범위

- 객실 목록 조회 (필터, 검색, 페이지네이션)
- 객실 상세 조회
- 객실 생성/수정/삭제
- 객실 변경 히스토리 조회

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/room_handler.go` | HTTP 핸들러 |
| `services/room_service.go` | 비즈니스 로직 |
| `services/history_service.go` | 히스토리 조회 |
| `repositories/room_repository.go` | 데이터 액세스 |
| `models/room.go` | 객실 모델 |
| `dto/room.go` | DTO 정의 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/rooms` | 객실 목록 | USER |
| GET | `/api/v1/rooms/{id}` | 객실 상세 | USER |
| POST | `/api/v1/rooms` | 객실 생성 | ADMIN |
| PATCH | `/api/v1/rooms/{id}` | 객실 수정 | ADMIN |
| DELETE | `/api/v1/rooms/{id}` | 객실 삭제 | ADMIN |
| GET | `/api/v1/rooms/{id}/histories` | 변경 이력 | ADMIN |

---

## 3. 객실 모델

### 3.1 Room 엔티티

```go
type Room struct {
    ID          uint        `gorm:"primaryKey"`
    Number      string      `gorm:"size:10;not null;uniqueIndex:idx_room_number_deleted"`
    RoomGroupID uint        `gorm:"not null"`
    RoomGroup   RoomGroup   `gorm:"foreignKey:RoomGroupID"`
    Note        string      `gorm:"size:200;not null;default:''"`
    Status      RoomStatus  `gorm:"not null;default:1"`
    CreatedBy   uint        `gorm:"not null"`
    UpdatedBy   uint        `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   LegacyDeletedAt
}
```

### 3.2 객실 상태 (RoomStatus)

| 값 | 코드 | 설명 |
|----|------|------|
| -10 | DAMAGED | 파손 |
| -1 | CONSTRUCTION | 공사중 |
| 0 | INACTIVE | 비활성 |
| 1 | NORMAL | 정상 |

---

## 4. 객실 목록

### 4.1 요청

```http
GET /api/v1/rooms?page=0&size=20&roomGroupId=1&status=NORMAL&search=101&sort=number,asc
Authorization: Bearer {token}
```

### 4.2 쿼리 파라미터

| 파라미터 | 타입 | 설명 |
|----------|------|------|
| page | int | 페이지 번호 (0-based) |
| size | int | 페이지 크기 |
| sort | string | 정렬 (예: number,asc) |
| roomGroupId | int | 객실 그룹 필터 |
| status | string | 상태 필터 (DAMAGED, CONSTRUCTION, INACTIVE, NORMAL) |
| search | string | 검색어 (객실 번호) |

### 4.3 응답

```json
{
  "data": [
    {
      "id": 1,
      "number": "101",
      "roomGroupId": 1,
      "roomGroup": {
        "id": 1,
        "name": "스탠다드",
        "peekPrice": 150000,
        "offPeekPrice": 100000
      },
      "note": "",
      "status": "NORMAL",
      "createdAt": "2026-01-01T00:00:00",
      "updatedAt": "2026-01-07T12:00:00",
      "createdBy": { "id": 1, "userId": "admin", "name": "관리자" },
      "updatedBy": { "id": 1, "userId": "admin", "name": "관리자" }
    }
  ],
  "page": {
    "size": 20,
    "number": 0,
    "totalElements": 50,
    "totalPages": 3
  }
}
```

---

## 5. 객실 생성

### 5.1 요청

```http
POST /api/v1/rooms
Authorization: Bearer {token}
Content-Type: application/json

{
  "number": "201",
  "roomGroupId": 1,
  "note": "바다 전망",
  "status": "NORMAL"
}
```

### 5.2 검증 규칙

| 필드 | 규칙 |
|------|------|
| number | 필수, 1-10자, 중복 불가 |
| roomGroupId | 필수, 존재하는 그룹 |
| note | 선택, 최대 200자 |
| status | 선택, DAMAGED/CONSTRUCTION/INACTIVE/NORMAL |

### 5.3 비즈니스 로직

1. 객실 번호 중복 확인
2. 객실 그룹 존재 확인
3. 감사 정보 설정 (createdBy, updatedBy)
4. 객실 생성
5. audit_logs에 CREATE 기록

---

## 6. 객실 수정

### 6.1 요청

```http
PATCH /api/v1/rooms/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "number": "101A",
  "note": "리모델링 완료",
  "status": "NORMAL"
}
```

### 6.2 수정 가능 필드

- number (중복 체크)
- roomGroupId
- note
- status

### 6.3 비즈니스 로직

1. 객실 존재 확인
2. 번호 변경 시 중복 확인
3. 감사 정보 업데이트 (updatedBy)
4. 객실 수정
5. audit_logs에 UPDATE 기록 (변경 필드 포함)

---

## 7. 객실 삭제

### 7.1 요청

```http
DELETE /api/v1/rooms/1
Authorization: Bearer {token}
```

### 7.2 응답

- 성공: 204 No Content
- 실패: 404 Not Found

### 7.3 비즈니스 로직

1. 객실 존재 확인
2. 소프트 삭제 (deleted_at = NOW())
3. audit_logs에 DELETE 기록

---

## 8. 객실 히스토리

### 8.1 요청

```http
GET /api/v1/rooms/1/histories?page=0&size=20
Authorization: Bearer {token}
```

### 8.2 응답

```json
{
  "data": [
    {
      "entity": {
        "id": 1,
        "number": "101",
        "roomGroupId": 1,
        "note": "수정된 메모",
        "status": "NORMAL"
      },
      "historyType": "UPDATED",
      "historyCreatedAt": "2026-01-07T14:00:00",
      "updatedFields": ["note"]
    },
    {
      "entity": { ... },
      "historyType": "CREATED",
      "historyCreatedAt": "2026-01-01T10:00:00",
      "updatedFields": []
    }
  ],
  "page": { ... }
}
```

### 8.3 히스토리 타입

| 타입 | 설명 |
|------|------|
| CREATED | 생성 |
| UPDATED | 수정 |
| DELETED | 삭제 |

---

## 9. 리포지토리 메서드

| 메서드 | 설명 |
|--------|------|
| `Create` | 객실 생성 |
| `Update` | 객실 수정 |
| `Delete` | 소프트 삭제 |
| `FindByID` | ID로 조회 |
| `FindByIDWithGroup` | ID로 조회 (그룹 포함) |
| `FindAll` | 필터/페이지네이션 목록 |
| `FindAvailableRooms` | 가용 객실 조회 |
| `ExistsByNumber` | 번호 중복 확인 |
| `IsRoomAvailable` | 객실 가용 여부 |
| `FindByNumber` | 번호로 조회 |

---

## 10. 테스트

- `room_handler_test.go`
- `room_service_test.go`

### 테스트 케이스

- 객실 목록 조회 (필터, 페이지네이션)
- 객실 상세 조회
- 객실 생성 (중복 체크)
- 객실 수정
- 객실 삭제
- 히스토리 조회
- 권한 없는 사용자 접근 거부
