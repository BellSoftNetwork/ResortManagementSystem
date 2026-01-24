# Resort Management System - 개발 가이드

리조트 통합 관리 시스템 모노레포

## Documents

| Document | Purpose |
|----------|---------|
| **AGENTS.md** | 개발 가이드 (현재 문서) |
| **README.md** | 사용자/설치 가이드 |
| **[specs/](specs/README.md)** | 기능 스펙 문서 |
| **[docs/](docs/README.md)** | 개발 참조 문서 |

### Detailed Guides (docs/)

| Guide | When to Read |
|-------|--------------|
| [docs/guides/docker-development.md](docs/guides/docker-development.md) | Docker 개발 환경 설정 시 |
| [docs/guides/api-testing.md](docs/guides/api-testing.md) | API 테스트 시 |
| [docs/guides/database-migration.md](docs/guides/database-migration.md) | DB 마이그레이션 작업 시 |
| [docs/references/spring-boot-compatibility.md](docs/references/spring-boot-compatibility.md) | api-core 호환성 이슈 시 |
| [docs/references/jwt-auth.md](docs/references/jwt-auth.md) | 인증 구조 확인 시 |
| [docs/references/hibernate-envers.md](docs/references/hibernate-envers.md) | History API 작업 시 |
| [docs/contracts/api-comparison.md](docs/contracts/api-comparison.md) | API 비교 시 |

---

## Project Architecture

### 앱 구조

```
apps/
├── api-core/        # Go + Gin (마이그레이션 대상, 메인)
├── api-legacy/      # Kotlin + Spring Boot (레거시)
└── frontend-web/    # Vue.js + Quasar
```

### api-core 아키텍처

```
apps/api-core/internal/
├── config/         # 설정 관리
├── context/        # 요청 컨텍스트 유틸리티
├── database/       # 데이터베이스 및 Redis 설정
├── dto/            # 데이터 전송 객체
├── handlers/       # HTTP 핸들러
├── middleware/     # HTTP 미들웨어
├── migrations/     # 데이터베이스 마이그레이션
├── models/         # GORM 모델
├── repositories/   # 데이터 액세스 레이어
├── services/       # 비즈니스 로직 레이어
└── utils/          # 공유 유틸리티
```

---

## Development Environment (Docker 필수)

**⚠️ 모든 개발 작업은 반드시 Docker 기반으로 실행. 로컬 환경 직접 실행 금지.**

```bash
# 전체 개발 환경 시작
docker compose up -d

# 로그 확인
docker compose logs -f [service-name]

# 컨테이너 내부 접속
docker compose exec api-core bash
docker compose exec api-legacy bash
docker compose exec frontend sh
```

### 서비스 포트

| 서비스 | 포트 | URL |
|--------|------|-----|
| MySQL | 3306 | localhost:3306 |
| Redis | 6379 | localhost:6379 |
| API Core (Go) | 8080 | http://localhost:8080 |
| API Legacy (Spring) | 8081 | http://localhost:8081 |
| Frontend | 9000 | http://localhost:9000 |

상세 가이드: [docs/guides/docker-development.md](docs/guides/docker-development.md)

---

## Quick Commands

### api-core (Go)

```bash
docker compose exec api-core bash
make dev              # 개발 모드
make test            # 테스트
make lint            # 린트
make build           # 빌드
```

### api-legacy (Kotlin)

```bash
docker compose exec api-legacy bash
./gradlew bootRun         # 실행
./gradlew test           # 테스트
./gradlew ktlintCheck    # 린트
```

### frontend-web (Vue.js)

```bash
docker compose exec frontend sh
yarn dev            # 개발 서버
yarn build          # 빌드
yarn lint           # 린트
```

### API 테스트

```bash
# 항상 이 스크립트 사용 (curl 직접 사용 금지)
python3 scripts/api-test.py /api/v1/users
python3 scripts/api-test.py /api/v1/reservations -m POST -r ADMIN -d '{"roomId": 1}'
```

상세 가이드: [docs/guides/api-testing.md](docs/guides/api-testing.md)

---

## Tech Stack

| Component | api-core | api-legacy | frontend-web |
|-----------|----------|------------|--------------|
| Language | Go 1.21+ | Kotlin | TypeScript |
| Framework | Gin | Spring Boot | Vue.js 3 + Quasar |
| Database | MySQL 8.0 | MySQL 8.0 | - |
| Cache | Redis 7 | Redis 7 | - |
| ORM | GORM | JPA/Hibernate | - |

---

## Migration Status (Kotlin → Go)

현재 api-legacy에서 api-core로 마이그레이션 진행 중:

| 영역 | 상태 |
|------|:----:|
| 인증 및 JWT | ✅ |
| 기본 CRUD | ✅ |
| 페이지네이션/필터링 | ✅ |
| 에러 처리 | ✅ |
| History API | ✅ |
| API 응답 호환성 검증 | 🚧 |
| DB 스키마 통합 | 📋 |
| 운영 환경 전환 | 📋 |

마이그레이션 스펙: [specs/migration/](specs/migration/)

---

## Coding Rules

### 공통

- 파일 끝은 빈 줄로 끝나야 함 (EOF newline)
- 린터 규칙 준수 필수
- 죽은 코드 제거 필수

### api-core (Go)

- BDD 스타일 테스트 (한글 설명 사용)
- golangci-lint 준수
- Spring Boot 호환 응답 형식 유지

### api-legacy (Kotlin)

- Ktlint 준수
- JaCoCo 커버리지 30% 이상

### frontend-web

- ESLint + Prettier 준수
- TypeScript strict 모드

---

## Verification Checklist

작업 완료 시 검증:

```bash
# api-core
docker compose exec api-core make test
docker compose exec api-core make lint

# api-legacy
docker compose exec api-legacy ./gradlew test
docker compose exec api-legacy ./gradlew ktlintCheck

# frontend-web
docker compose exec frontend yarn lint
docker compose exec frontend yarn build
```

---

## Work Management

### 스펙 기반 개발

모든 기능 개발은 스펙 문서 기반으로 진행:

```
specs/
├── migration/      # Kotlin → Go 마이그레이션 스펙
├── frontend/       # 프론트엔드 기능 스펙
├── infra/          # 인프라 스펙
└── _templates/     # 스펙/플랜 템플릿
```

### 스펙 상태

| spec.md 상태 | 설명 |
|-------------|------|
| `draft` | 초안 작성 중 |
| `approved` | 검토 완료, 구현 가능 |
| `completed` | 구현 완료 |
| `deprecated` | 더 이상 유효하지 않음 |

---

## Git Workflow

- 기능 브랜치: `feature/XXX-설명`
- pre-commit 훅: 코드 포맷팅
- main 브랜치: 프로덕션 릴리스
