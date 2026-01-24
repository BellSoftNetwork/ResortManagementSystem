# 데이터베이스 마이그레이션 가이드

> api-core는 Go 기반 커스텀 마이그레이션, api-legacy는 Liquibase 사용

---

## api-core (Go) 마이그레이션

### 마이그레이션 위치

```
apps/api-core/internal/migrations/
├── 001_initial_schema.go
└── ...
```

### 실행 방법

```bash
# Docker 컨테이너에서 실행
docker compose exec api-core bash

# 마이그레이션만 실행하고 종료
./main --migrate-only

# 마이그레이션 실행 후 서버 시작
./main --migrate

# 별도 마이그레이션 도구 사용
go run cmd/migrate/main.go -action=migrate    # 적용
go run cmd/migrate/main.go -action=status     # 상태 확인
go run cmd/migrate/main.go -action=rollback -steps=2  # 롤백
```

### 스키마 호환성

api-core는 api-legacy의 Liquibase 스키마와 완전히 동일한 테이블 구조를 생성합니다.
두 API는 서로 다른 데이터베이스 스키마를 사용하지만 테이블 구조는 동일합니다:
- api-legacy: `rms-legacy` 스키마
- api-core: `rms-core` 스키마

---

## api-legacy (Liquibase) 마이그레이션

### 변경셋 위치

```
apps/api-legacy/src/main/resources/db/changelog/
```

### 실행

Spring Boot 앱 시작 시 자동으로 Liquibase가 실행됩니다.

---

## 통합 전략 (마이그레이션 계획)

현재 두 방식이 공존하며, 완전 전환 시:

### 옵션 A: api-core 마이그레이션으로 완전 대체 (권장)

- Liquibase `databasechangelog` 테이블 무시
- api-core 마이그레이션만 사용
- 장점: 단순함
- 단점: Liquibase 이력 단절

### 옵션 B: 공존

- api-legacy 배포 시 Liquibase 실행
- api-core 배포 시 Go 마이그레이션 실행

자세한 내용은 [db-schema-unification 스펙](../specs/migration/db-schema-unification/spec.md) 참조.
