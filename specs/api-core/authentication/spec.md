---
id: api-core-authentication
title: "api-core 인증"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: high
effort: medium
---

# api-core 인증

> JWT 기반 인증, 회원가입, 브루트포스 방지

---

## 1. 개요

### 1.1 인증 방식

- JWT (JSON Web Token) 기반 Stateless 인증
- Access Token (15분) + Refresh Token (7일)
- Redis에 Refresh Token 저장

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `handlers/auth_handler.go` | 인증 엔드포인트 |
| `services/auth_service.go` | 인증 비즈니스 로직 |
| `middleware/auth.go` | JWT 검증 미들웨어 |
| `pkg/auth/jwt.go` | JWT 유틸리티 |
| `repositories/login_attempt_repository.go` | 로그인 시도 저장 |

---

## 2. 엔드포인트

| Method | Path | 설명 | 인증 |
|--------|------|------|:----:|
| POST | `/api/v1/auth/register` | 회원가입 | ❌ |
| POST | `/api/v1/auth/login` | 로그인 | ❌ |
| POST | `/api/v1/auth/refresh` | 토큰 갱신 | ❌ |

---

## 3. 회원가입

### 3.1 요청

```json
POST /api/v1/auth/register
{
  "userId": "newuser",
  "email": "user@example.com",
  "name": "홍길동",
  "password": "password123"
}
```

### 3.2 검증 규칙

| 필드 | 규칙 |
|------|------|
| userId | 필수, 3-30자 |
| email | 선택, 이메일 형식, 최대 100자 |
| name | 필수, 2-20자 |
| password | 필수, 최소 8자 |

### 3.3 비즈니스 로직

1. userId 중복 확인
2. email 중복 확인 (입력된 경우)
3. 비밀번호 bcrypt 해싱 (`{bcrypt}` 접두사)
4. 사용자 생성 (status: ACTIVE, role: NORMAL)

---

## 4. 로그인

### 4.1 요청

```json
POST /api/v1/auth/login
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

### 4.2 응답

```json
{
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "accessTokenExpiresIn": 900
  }
}
```

### 4.3 비즈니스 로직

1. 브루트포스 체크 (15분 내 5회 실패 시 차단)
2. username으로 사용자 조회 (userId 또는 email)
3. 비밀번호 검증 (`{bcrypt}` 접두사 제거 후 비교)
4. 사용자 상태 확인 (ACTIVE만 허용)
5. JWT 토큰 발급
6. 로그인 시도 기록 (성공/실패)
7. 디바이스 정보 저장

### 4.4 JWT 토큰 구조

```json
{
  "sub": "1",
  "username": "admin",
  "authorities": ["ROLE_ADMIN"],
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

## 5. 토큰 갱신

### 5.1 요청

```json
POST /api/v1/auth/refresh
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 5.2 비즈니스 로직

1. Refresh Token 유효성 검증
2. Redis에서 토큰 존재 확인
3. 사용자 정보 조회
4. 새 Access Token + Refresh Token 발급
5. 기존 Refresh Token 무효화

---

## 6. 브루트포스 방지

### 6.1 정책

| 항목 | 값 |
|------|-----|
| 시간 윈도우 | 15분 |
| 최대 실패 횟수 | 5회 |
| 차단 응답 | 429 Too Many Requests |

### 6.2 추적 기준

- username + IP 조합
- 또는 username만
- 또는 IP만

### 6.3 login_attempts 테이블

| 컬럼 | 설명 |
|------|------|
| username | 시도한 사용자명 |
| ip_address | 클라이언트 IP |
| successful | 성공 여부 |
| attempt_at | 시도 시간 |
| os_info | OS 정보 |
| user_agent | User-Agent |
| device_fingerprint | 디바이스 지문 |

---

## 7. 미들웨어

### 7.1 AuthMiddleware

```go
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
    // 1. Authorization 헤더에서 Bearer 토큰 추출
    // 2. JWT 토큰 검증
    // 3. Claims에서 사용자 정보 추출
    // 4. Context에 사용자 정보 설정
    // 5. 실패 시 401 응답
}
```

### 7.2 RoleMiddleware

```go
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    // 1. Context에서 사용자 역할 조회
    // 2. 허용된 역할인지 확인
    // 3. 불일치 시 403 응답
}
```

---

## 8. 테스트

### 8.1 테스트 파일

- `auth_handler_test.go` - 핸들러 테스트
- `auth_service_test.go` - 서비스 테스트
- `auth_security_test.go` - 보안 테스트
- `auth_brute_force_test.go` - 브루트포스 테스트
- `auth_integration_test.go` - 통합 테스트

### 8.2 테스트 케이스

- 올바른 자격 증명으로 로그인
- 잘못된 비밀번호로 로그인 실패
- 비활성 사용자 로그인 거부
- 브루트포스 차단 동작
- 토큰 갱신 성공/실패
- 회원가입 중복 체크
