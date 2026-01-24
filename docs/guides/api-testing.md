# API 테스트 스크립트 가이드

> `scripts/api-test.py` - 인증이 필요한 API를 쉽게 테스트하는 Python 스크립트

---

## 주요 기능

이 스크립트는 모든 인증 관련 문제를 자동으로 해결합니다:

- 429 응답 (Too Many Requests) → 자동으로 login_attempts 테이블 정리
- 401 응답 (Unauthorized) → 자동으로 패스워드 재설정
- 계정이 없을 경우 → 자동 생성
- 토큰 만료 시 → 자동 갱신

**중요**: curl 명령을 직접 사용하지 마세요. 항상 이 스크립트를 사용하세요.

---

## 기본 사용법

```bash
# 간단한 GET 요청
python3 scripts/api-test.py /api/v1/users

# 관리자 권한으로 POST 요청
python3 scripts/api-test.py /api/v1/dev/generate-essential-data -m POST -r SUPER_ADMIN

# 데이터와 함께 POST 요청
python3 scripts/api-test.py /api/v1/reservations -m POST -d '{"roomId": 1, "guestName": "테스트 손님"}'

# 쿼리 파라미터와 함께 GET 요청
python3 scripts/api-test.py /api/v1/reservations -p page=0 -p size=20

# 도움말 보기
python3 scripts/api-test.py -h
```

---

## 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-m, --method` | HTTP 메소드 | GET |
| `-r, --role` | 사용자 권한 (USER, ADMIN, SUPER_ADMIN) | USER |
| `-d, --data` | JSON 형식의 요청 본문 | - |
| `-p, --param` | 쿼리 파라미터 (여러 개 사용 가능) | - |
| `-u, --url` | API 기본 URL | http://localhost:8080 |
| `-v, --verbose` | 상세 출력 모드 | false |

---

## 권한별 기본 계정

| 권한 | Username | Password |
|------|----------|----------|
| SUPER_ADMIN | testadmin | testadmin123 |
| ADMIN | testmanager | testmanager123 |
| USER | testuser | testuser123 |

---

## 활용 예시

### 더미 데이터 생성

```bash
python3 scripts/api-test.py /api/v1/dev/generate-essential-data -m POST -r SUPER_ADMIN
```

### 예약 목록 조회

```bash
python3 scripts/api-test.py /api/v1/reservations -p page=0 -p size=20 -p sort=id,desc
```

### 새 예약 생성

```bash
python3 scripts/api-test.py /api/v1/reservations -m POST -r ADMIN -d '{
  "roomId": 1,
  "guestName": "홍길동",
  "guestPhone": "010-1234-5678",
  "checkInAt": "2024-03-20",
  "checkOutAt": "2024-03-22",
  "adultCount": 2
}'
```

### 특정 예약 조회

```bash
python3 scripts/api-test.py /api/v1/reservations/1 -r ADMIN
```

---

## 자동 복구 시나리오

### 계정 잠김 (429)

```bash
python3 scripts/api-test.py /api/v1/my -r SUPER_ADMIN
# 출력: "Account locked due to too many login attempts"
# 출력: "Attempting to unlock account..."
# 출력: "Cleared login attempts for testadmin"
# 정상적으로 API 호출 성공
```

### 패스워드 오류 (401)

```bash
python3 scripts/api-test.py /api/v1/my -r ADMIN
# 출력: "Login failed: Invalid credentials"
# 출력: "Attempting to reset password..."
# 출력: "Reset password for testmanager"
# 정상적으로 API 호출 성공
```

### 계정 없음

```bash
python3 scripts/api-test.py /api/v1/my -r USER
# 출력: "Login failed for testuser, attempting to register..."
# 출력: "Successfully created user: testuser with role: USER"
# 정상적으로 API 호출 성공
```

---

## Claude Code 사용 시 주의사항

```bash
# ❌ 잘못된 방법 (사용하지 마세요)
curl -X POST http://localhost:8080/api/v1/auth/login ...

# ✅ 올바른 방법 (항상 이 스크립트 사용)
python3 scripts/api-test.py /api/v1/auth/login -m POST -d '{"username":"test","password":"test"}'
```

스크립트는 모든 인증 관련 문제를 자동으로 해결하므로, 별도의 토큰 관리나 DB 조작이 필요하지 않습니다.
