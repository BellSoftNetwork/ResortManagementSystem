---
id: api-core-user-management
title: "api-core 사용자 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-core 사용자 관리

> 사용자 프로필 및 관리자 계정 관리

---

## 1. 개요

### 1.1 기능 범위

- 내 프로필 조회/수정
- 비밀번호 변경
- 관리자 계정 목록/생성/수정 (ADMIN 이상)

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/user_handler.go` | HTTP 핸들러 |
| `services/user_service.go` | 비즈니스 로직 |
| `repositories/user_repository.go` | 데이터 액세스 |
| `models/user.go` | 사용자 모델 |
| `dto/user.go` | DTO 정의 |

---

## 2. 엔드포인트

### 2.1 프로필

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/my` | 내 정보 조회 | USER |
| POST | `/api/v1/my` | 내 정보 조회 (호환용) | USER |
| PATCH | `/api/v1/my` | 내 정보 수정 | USER |

### 2.2 관리자 계정

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/admin/accounts` | 사용자 목록 | ADMIN |
| POST | `/api/v1/admin/accounts` | 사용자 생성 | ADMIN |
| PATCH | `/api/v1/admin/accounts/{id}` | 사용자 수정 | ADMIN |

---

## 3. 사용자 모델

### 3.1 User 엔티티

```go
type User struct {
    ID        uint            `gorm:"primaryKey"`
    UserID    string          `gorm:"column:user_id;size:30;not null"`
    Email     string          `gorm:"size:100"`
    Name      string          `gorm:"size:20;not null"`
    Password  string          `gorm:"size:100;not null" json:"-"`
    Status    UserStatus      `gorm:"not null"`
    Role      UserRole        `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt LegacyDeletedAt `gorm:"not null"`
}
```

### 3.2 역할 (UserRole)

| 값 | 코드 | 설명 |
|----|------|------|
| 0 | NORMAL | 일반 사용자 |
| 100 | ADMIN | 관리자 |
| 127 | SUPER_ADMIN | 슈퍼 관리자 |

### 3.3 상태 (UserStatus)

| 값 | 코드 | 설명 |
|----|------|------|
| -1 | INACTIVE | 비활성 |
| 1 | ACTIVE | 활성 |

---

## 4. 프로필 관리

### 4.1 내 정보 조회

```http
GET /api/v1/my
Authorization: Bearer {token}
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
    "profileImageUrl": "https://gravatar.com/...",
    "createdAt": "2026-01-01T00:00:00",
    "updatedAt": "2026-01-07T12:00:00"
  }
}
```

### 4.2 내 정보 수정

```http
PATCH /api/v1/my
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "새 이름",
  "email": "new@example.com"
}
```

### 4.3 비밀번호 변경

```http
PATCH /api/v1/my
Authorization: Bearer {token}
Content-Type: application/json

{
  "currentPassword": "oldPassword",
  "newPassword": "newPassword123"
}
```

**검증**:
- currentPassword 필수
- newPassword 최소 8자
- 현재 비밀번호 일치 확인

---

## 5. 관리자 계정 관리

### 5.1 사용자 목록

```http
GET /api/v1/admin/accounts?page=0&size=20&sort=createdAt,desc
Authorization: Bearer {token}
```

**쿼리 파라미터**:
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| page | int | 페이지 번호 (0-based) |
| size | int | 페이지 크기 (기본 20) |
| sort | string | 정렬 (예: createdAt,desc) |

### 5.2 사용자 생성

```http
POST /api/v1/admin/accounts
Authorization: Bearer {token}
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

**검증**:
| 필드 | 규칙 |
|------|------|
| userId | 필수, 3-30자, 중복 불가 |
| email | 선택, 이메일 형식, 중복 불가 |
| name | 필수, 2-20자 |
| password | 필수, 최소 8자 |
| role | 선택, NORMAL/ADMIN/SUPER_ADMIN |
| status | 선택, ACTIVE/INACTIVE |

### 5.3 사용자 수정

```http
PATCH /api/v1/admin/accounts/2
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "수정된 이름",
  "role": "ADMIN",
  "status": "INACTIVE"
}
```

**수정 가능 필드**:
- email
- name
- status
- role

---

## 6. 프로필 이미지

### 6.1 Gravatar 연동

- 이메일 기반 Gravatar URL 생성
- `pkg/utils/gravatar.go`에서 처리

```go
func GetGravatarURL(email string) string {
    hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
    return fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon", hash)
}
```

---

## 7. 리포지토리 메서드

| 메서드 | 설명 |
|--------|------|
| `Create` | 사용자 생성 |
| `Update` | 사용자 수정 |
| `Delete` | 소프트 삭제 |
| `FindByID` | ID로 조회 |
| `FindByUserID` | userId로 조회 |
| `FindByEmail` | 이메일로 조회 |
| `FindAll` | 페이지네이션 목록 |
| `ExistsByUserID` | userId 존재 확인 |
| `ExistsByEmail` | 이메일 존재 확인 |
| `HasAnyUsers` | 사용자 존재 여부 |

---

## 8. 테스트

- `user_handler_test.go`
- `user_service_test.go`

### 테스트 케이스

- 프로필 조회 성공
- 프로필 수정 성공
- 비밀번호 변경 (현재 비밀번호 검증)
- 관리자 목록 페이지네이션
- 관리자 생성 (중복 체크)
- 권한 없는 사용자 접근 거부
