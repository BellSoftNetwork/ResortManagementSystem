# JWT 인증 구조

> api-core의 JWT 인증 구조 (Spring Boot 호환)

---

## 토큰 형식

api-core는 api-legacy(Spring Boot)와 호환되는 JWT 토큰 사용:

```json
{
  "username": "testadmin",
  "authorities": "SUPER_ADMIN",
  "sub": "1",
  "exp": 1749264386,
  "iat": 1749263486
}
```

---

## 토큰 설정

| 설정 | 값 |
|------|-----|
| 액세스 토큰 만료 | 15분 (900초) |
| 리프레시 토큰 만료 | 7일 (604800초) |
| 토큰 타입 | Bearer |

---

## 인증 플로우

### 1. 로그인

```
POST /api/v1/auth/login
```

**Request:**
```json
{
  "username": "testadmin",
  "password": "testadmin123"
}
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 900
}
```

### 2. 토큰 갱신

```
POST /api/v1/auth/refresh
```

**Request:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "accessToken": "...",
  "refreshToken": "...",
  "expiresIn": 900
}
```

### 3. 보호된 엔드포인트 접근

```
Authorization: Bearer <accessToken>
```

---

## 비밀번호 저장

비밀번호는 Spring Security와 호환되는 BCrypt 형식으로 저장:

```
{bcrypt}$2a$10$...
```

`{bcrypt}` 접두사에 주의 - api-legacy와의 호환성을 위해 필수.

---

## 환경 변수

| 변수 | 설명 | 기본값 |
|------|------|--------|
| `JWT_SECRET` | JWT 서명 키 | your-secret-key |
| `JWT_ACCESS_TOKEN_EXPIRY` | 액세스 토큰 만료 (초) | 900 |
| `JWT_REFRESH_TOKEN_EXPIRY` | 리프레시 토큰 만료 (초) | 604800 |
