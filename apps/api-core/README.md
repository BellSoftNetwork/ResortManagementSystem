# 리조트 관리 시스템 - Go API Core

Kotlin/Spring Boot 버전에서 마이그레이션된 리조트 관리 시스템 API의 Go 구현체입니다.

## 주요 기능

- Gin 프레임워크 기반 RESTful API
- JWT 인증 및 권한 부여
- GORM ORM을 사용한 MySQL 데이터베이스
- 세션 관리를 위한 Redis
- 역할 기반 접근 제어 (USER, ADMIN, SUPER_ADMIN)
- 소프트 삭제 지원
- 감사 추적 (사용자 추적)
- 요청 검증
- CORS 지원
- 환경별 설정 관리

## 필수 요구사항

- Go 1.21 이상
- MySQL 8.0 이상
- Redis 6.0 이상
- Docker & Docker Compose (권장)

## 설치 방법

1. 저장소 클론
2. API 디렉토리로 이동:
   ```bash
   cd apps/api-core
   ```

3. 의존성 설치:
   ```bash
   go mod download
   ```

## 설정

애플리케이션은 YAML 설정 파일을 사용합니다:
- `config/application.yaml` - 기본 설정
- `config/application-local.yaml` - 로컬 개발용
- `config/application-production.yaml` - 운영 환경 설정

특정 설정을 로드하려면 `PROFILE` 환경 변수를 설정하세요:
```bash
export PROFILE=local  # 또는 production
```

### 환경 변수

필수 환경 변수:
- `DATABASE_MYSQL_HOST` - MySQL 서버 호스트
- `DATABASE_MYSQL_USER` - MySQL 사용자명
- `DATABASE_MYSQL_PASSWORD` - MySQL 비밀번호
- `REDIS_HOST` - Redis 서버 호스트

## 애플리케이션 실행

### Docker Compose 사용 (권장)
```bash
# 전체 스택 실행 (MySQL, Redis, API)
docker compose up -d

# API 컨테이너만 재시작
docker compose restart api

# 로그 확인
docker compose logs -f api

# 종료
docker compose down
```

### 로컬 개발 모드
```bash
# 데이터베이스 마이그레이션
go run cmd/migrate/main.go

# 서버 실행
go run cmd/server/main.go
```

### 프로덕션 빌드
```bash
# 바이너리 빌드
go build -o bin/server cmd/server/main.go

# 실행
./bin/server
```

## API 엔드포인트

### 공개 엔드포인트 (인증 불필요)
- `POST /api/v1/auth/register` - 사용자 등록
- `POST /api/v1/auth/login` - 사용자 로그인
- `POST /api/v1/auth/refresh` - 액세스 토큰 갱신

### 인증 필요 엔드포인트
- `GET /api/v1/env` - 환경 정보 조회
- `GET /api/v1/config` - 앱 설정 조회
- `GET/POST /api/v1/my` - 현재 사용자 정보 조회
- `PATCH /api/v1/my` - 현재 사용자 정보 수정

### 관리자 엔드포인트 (ADMIN/SUPER_ADMIN 권한 필요)
- 사용자 관리: `/api/v1/admin/accounts`
- 객실 관리: `/api/v1/rooms`
- 객실 그룹 관리: `/api/v1/room-groups`
- 예약 관리: `/api/v1/reservations`
- 결제 수단 관리: `/api/v1/payment-methods`

## 개발

### 프로젝트 구조
```
api-core/
├── cmd/
│   ├── server/         # 애플리케이션 진입점
│   └── migrate/        # DB 마이그레이션 도구
├── internal/           # 내부 애플리케이션 코드
│   ├── config/         # 설정 관리
│   ├── database/       # 데이터베이스 설정
│   ├── dto/            # 데이터 전송 객체
│   ├── handlers/       # HTTP 핸들러
│   ├── middleware/     # HTTP 미들웨어
│   ├── models/         # 데이터베이스 모델
│   ├── repositories/   # 데이터 접근 계층
│   ├── services/       # 비즈니스 로직
│   └── migrations/     # DB 마이그레이션
├── pkg/                # 공개 패키지
│   ├── auth/          # JWT 인증
│   ├── response/      # HTTP 응답 헬퍼
│   └── utils/         # 유틸리티
├── config/            # 설정 파일
├── scripts/           # 스크립트 파일
└── docker-compose.yml # Docker Compose 설정
```

### 테스트 실행
```bash
go test ./...
```

### 코드 포맷팅
```bash
go fmt ./...
```

### 린팅
```bash
golangci-lint run
```

## 데이터베이스

GORM을 사용한 데이터베이스 작업:
- 자동 마이그레이션
- 소프트 삭제
- 감사 필드 (created_by, updated_by)
- 트랜잭션 지원

## 인증

JWT 기반 인증:
- 액세스 토큰 (15분 만료)
- 리프레시 토큰 (7일 만료)
- Redis에 토큰 저장
- 로그인 시도 추적
- 무차별 대입 공격 방지

## 에러 처리

표준화된 에러 응답:
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "에러 메시지",
    "details": {}
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "requestId": "xxx"
  }
}
```

## 배포

### K8s 배포를 위한 Docker 이미지 빌드
```bash
# Docker 이미지 빌드
docker build -t resort-api-core:latest .

# 태그 지정
docker tag resort-api-core:latest your-registry/resort-api-core:latest

# 레지스트리에 푸시
docker push your-registry/resort-api-core:latest
```

### 운영 환경 변수
```yaml
# K8s ConfigMap 또는 Secret에서 설정
- name: PROFILE
  value: "production"
- name: DATABASE_MYSQL_HOST
  valueFrom:
    secretKeyRef:
      name: mysql-secret
      key: host
- name: DATABASE_MYSQL_USER
  valueFrom:
    secretKeyRef:
      name: mysql-secret
      key: username
- name: DATABASE_MYSQL_PASSWORD
  valueFrom:
    secretKeyRef:
      name: mysql-secret
      key: password
- name: REDIS_HOST
  valueFrom:
    configMapKeyRef:
      name: redis-config
      key: host
- name: JWT_SECRET
  valueFrom:
    secretKeyRef:
      name: jwt-secret
      key: secret
```

## Kotlin/Spring Boot에서 마이그레이션

이 Go 구현체는 원본 Kotlin/Spring Boot 버전과 완전한 API 호환성을 유지합니다:
- 동일한 엔드포인트 경로 및 메소드
- 동일한 요청/응답 형식
- 호환 가능한 인증 토큰 (JWT)
- 동일한 데이터베이스 스키마
- 동등한 비즈니스 로직

### 주요 호환성 기능

1. **JWT 토큰 형식**: Spring Boot 형식과 호환되는 토큰
   - Claims 포함: `username`, `authorities`, `sub`, `exp`, `iat`
   - 액세스 토큰 만료: 15분
   - 리프레시 토큰 만료: 7일

2. **비밀번호 저장**: Spring Security 호환 BCrypt 형식 사용
   - `{bcrypt}` 접두사와 함께 저장
   - 예시: `{bcrypt}$2a$10$...`

3. **시간 형식**: 타임존 없는 JSON 타임스탬프
   - 형식: `2006-01-02T15:04:05` (타임존 접미사 없음)
   - Spring Boot의 기본 JSON 형식과 일치

4. **CORS 설정**: 설정 가능한 허용 출처
   - 개발: 모든 출처 허용을 위해 `*` 설정
   - 운영: 특정 도메인 설정

주요 차이점:
- Kotlin/Spring Boot 대신 Go/Gin
- JPA/Hibernate 대신 GORM
- Spring DI 대신 수동 의존성 주입
- 단순화된 설정 관리

## 보안 주의사항

- 항상 환경별 설정 파일 사용
- 민감한 데이터(토큰, 비밀번호)를 버전 관리에 커밋하지 말 것
- 운영 환경에서는 강력한 JWT 시크릿 사용
- 운영 환경을 위해 CORS를 적절히 설정

## 문제 해결

### API 컨테이너가 시작되지 않을 때
```bash
# 로그 확인
docker compose logs api

# 컨테이너 재시작
docker compose restart api

# 전체 스택 재시작
docker compose down && docker compose up -d
```

### 데이터베이스 연결 오류
```bash
# MySQL 컨테이너 상태 확인
docker compose ps mysql

# MySQL 로그 확인
docker compose logs mysql

# 네트워크 확인
docker network ls
docker network inspect api-core_default
```

### 포트 충돌 문제
```bash
# 사용 중인 포트 확인
lsof -i :8080
lsof -i :3306
lsof -i :6379

# 프로세스 종료
kill -9 <PID>
```