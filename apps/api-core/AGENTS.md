# api-core 개발 가이드

Go + Gin 기반 API 서버 (마이그레이션 대상, 메인 API)

## 아키텍처

```
internal/
├── config/         # 설정 관리
├── context/        # 요청 컨텍스트 유틸리티
├── database/       # 데이터베이스 및 Redis 설정
├── dto/            # 데이터 전송 객체 (검증 포함)
├── handlers/       # HTTP 핸들러 (컨트롤러)
├── middleware/     # HTTP 미들웨어 (인증, 에러 처리)
├── migrations/     # 데이터베이스 마이그레이션
├── models/         # GORM 모델
├── repositories/   # 데이터 액세스 레이어
├── services/       # 비즈니스 로직 레이어
└── utils/          # 공유 유틸리티
```

## 개발 명령어

```bash
# Docker 컨테이너 접속
docker compose exec api-core bash

# 컨테이너 내부에서 실행
make dev              # 개발 모드 (Air 핫 리로드)
make test            # 테스트 (자세한 출력)
make lint            # golangci-lint 실행
make fmt             # 코드 포맷팅
make build           # 현재 플랫폼용 빌드
make build-linux     # Linux용 빌드
```

## DB 마이그레이션

```bash
# 마이그레이션만 실행하고 종료
./main --migrate-only

# 마이그레이션 실행 후 서버 시작
./main --migrate

# 별도 도구 사용
go run cmd/migrate/main.go -action=migrate
go run cmd/migrate/main.go -action=status
go run cmd/migrate/main.go -action=rollback -steps=2
```

## 테스트 작성 규칙

### BDD 스타일 (한글)

```go
func TestPagination(t *testing.T) {
    t.Run("페이지가 0부터 시작하면 Spring Boot와 호환된다", func(t *testing.T) {
        // Given
        query := PaginationQuery{Page: 0, Size: 15}
        
        // When
        result := query.Validate()
        
        // Then
        assert.NoError(t, result)
    })
}
```

### HTTP 요청 시뮬레이션

```go
router := gin.New()
router.Use(middleware.ErrorHandler())
req := httptest.NewRequest("GET", "/test?page=0&size=15", nil)
```

## Spring Boot 호환성

api-legacy와 완전 호환 필요:

### 페이지네이션

```go
type PaginationQuery struct {
    Page int `form:"page" binding:"min=0"`      // 0 기반
    Size int `form:"size" binding:"min=1,max=2000"`
}
```

### 에러 응답

```go
type ErrorResponse struct {
    Message     string   `json:"message"`
    Errors      []string `json:"errors,omitempty"`
    FieldErrors []string `json:"fieldErrors,omitempty"`
}
```

### JSON 필드

- `omitempty` 주의: 프론트엔드가 키 존재를 기대하는 경우 제거
- 날짜 형식: `2024-01-15T10:30:00` (타임존 없음)

### 비밀번호 저장

```
{bcrypt}$2a$10$...
```

`{bcrypt}` 접두사 필수 (Spring Security 호환).

## 환경 변수

```bash
# 데이터베이스
DATABASE_MYSQL_HOST=mysql
DATABASE_MYSQL_PORT=3306
DATABASE_MYSQL_USER=root
DATABASE_MYSQL_PASSWORD=root
DATABASE_MYSQL_DATABASE=rms-core

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# API
API_PORT=8080
API_PROFILE=local

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_EXPIRY=900
JWT_REFRESH_TOKEN_EXPIRY=604800
```

## 참조 문서

- [JWT 인증 구조](../../docs/references/jwt-auth.md)
- [Spring Boot 호환성](../../docs/references/spring-boot-compatibility.md)
- [Hibernate Envers](../../docs/references/hibernate-envers.md)
- [API 비교](../../docs/contracts/api-comparison.md)
