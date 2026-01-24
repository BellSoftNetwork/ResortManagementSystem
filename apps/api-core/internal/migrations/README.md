# Database Migrations

이 디렉토리는 데이터베이스 마이그레이션 파일들을 포함합니다.

## 구조

- `migration.go`: 마이그레이션 시스템의 핵심 로직
- `migrations.go`: 모든 마이그레이션을 등록하는 레지스트리
- `001_initial_schema.go`: 초기 데이터베이스 스키마 (api-legacy와 동일한 구조)

## 마이그레이션 실행

### 애플리케이션 시작 시 자동 실행
```bash
docker compose up api-core
```

### 수동 실행
```bash
# Docker 컨테이너 내부에서
go run cmd/migrate/main.go -action=migrate

# 또는 빌드된 바이너리로
./migrate -action=migrate
```

### 마이그레이션 상태 확인
```bash
go run cmd/migrate/main.go -action=status
```

### 롤백
```bash
# 마지막 1개 마이그레이션 롤백
go run cmd/migrate/main.go -action=rollback -steps=1

# 마지막 3개 마이그레이션 롤백
go run cmd/migrate/main.go -action=rollback -steps=3
```

## 새 마이그레이션 추가

1. 새 마이그레이션 파일 생성:
```go
// 002_add_new_feature.go
package migrations

import "gorm.io/gorm"

var Migration002AddNewFeature = Migration{
    ID:          "002_add_new_feature",
    Description: "Add new feature tables",
    Up: func(db *gorm.DB) error {
        // 테이블 생성 또는 변경 로직
        return nil
    },
    Down: func(db *gorm.DB) error {
        // 롤백 로직
        return nil
    },
}
```

2. `migrations.go`에 등록:
```go
func AllMigrations() []Migration {
    return []Migration{
        Migration001InitialSchema,
        Migration002AddNewFeature, // 추가
    }
}
```

## 주의사항

- 마이그레이션 ID는 순차적으로 증가하도록 명명 (001, 002, 003...)
- 한 번 실행된 마이그레이션은 수정하지 않음
- 새로운 변경사항은 항상 새 마이그레이션으로 추가
- Down 함수는 선택사항이지만, 롤백이 필요한 경우 반드시 구현
