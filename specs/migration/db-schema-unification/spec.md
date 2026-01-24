---
id: db-schema-unification
title: "DB 스키마 버전 관리 통합"
status: approved
type: migration
version: 2.0.0
created: 2026-01-07
updated: 2026-01-12

depends_on: [api-response-compat]
related: [production-cutover]
replaces: null
replaced_by: null

impact: backend
risk: low
effort: small

changelog:
  - date: 2026-01-12
    description: "블루-그린 마이그레이션 전략으로 업데이트, 동일 DB 사용 방식 확정"
  - date: 2026-01-07
    description: "TODO.md Phase 3에서 마이그레이션"
---

# DB 스키마 버전 관리 통합

> api-core가 api-legacy의 프로덕션 DB를 직접 사용하는 블루-그린 마이그레이션 전략

---

## 1. 개요

### 1.1 배경

두 API가 서로 다른 스키마 관리 방식을 사용하지만, 테이블 구조는 동일:
- **api-legacy**: Liquibase (약 40개 changeset)
- **api-core**: Go 마이그레이션 (001_initial_schema)

### 1.2 검증 결과 (2026-01-12)

```bash
# 스키마 비교 결과: 100% 동일
rms-legacy.user == rms-core.user  ✅
rms-legacy.room == rms-core.room  ✅
# ... 모든 테이블 동일
```

### 1.3 결정된 전략

**동일 DB 직접 사용** (블루-그린 마이그레이션)
- api-core를 프로덕션 DB(`rms-legacy`)에 직접 연결
- 별도 데이터 마이그레이션 불필요
- 스키마 버전 관리 테이블만 다름 (무시)

---

## 2. 블루-그린 마이그레이션 설정

### 2.1 환경 변수 변경

api-core의 프로덕션 설정에서 DB를 api-legacy와 동일하게 지정:

```yaml
# api-core production config
database:
  mysql:
    host: ${DATABASE_MYSQL_HOST}
    port: 3306
    user: ${DATABASE_MYSQL_USER}
    password: ${DATABASE_MYSQL_PASSWORD}
    database: rms-legacy  # <-- 기존 프로덕션 DB 사용
```

### 2.2 마이그레이션 비활성화

api-core가 프로덕션 DB에 연결할 때 마이그레이션을 실행하지 않도록 설정:

```go
// cmd/server/main.go
// production 환경에서는 마이그레이션 자동 실행 안함
autoMigrate := cfg.Environment == "local" || cfg.Environment == "development"
```

현재 코드가 이미 이렇게 구현되어 있음 ✅

### 2.3 무시할 테이블

api-core가 api-legacy DB에 연결할 때 무시해야 할 테이블:
- `database_changelog` (Liquibase 메타데이터)
- `database_changelog_lock` (Liquibase 락)

api-core가 생성하지 않아도 되는 테이블:
- `schema_migrations` (Go 마이그레이션 메타데이터)
- `audit_logs` (api-core 전용, 선택적)

---

## 3. 배포 절차

### 3.1 사전 준비

```bash
# 1. 프로덕션 DB 백업
mysqldump -h <prod-host> -u <user> -p rms-legacy > backup_$(date +%Y%m%d).sql

# 2. api-core 이미지 빌드
docker build -t api-core:v1.0.0 .
docker push registry/api-core:v1.0.0
```

### 3.2 블루-그린 전환

```yaml
# 1. api-core 배포 (그린 환경)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-core  # 새 서비스
spec:
  template:
    spec:
      containers:
      - name: api-core
        image: registry/api-core:v1.0.0
        env:
        - name: DATABASE_MYSQL_DATABASE
          value: "rms-legacy"  # 동일 DB 사용
        - name: PROFILE
          value: "production"

# 2. Ingress 트래픽 전환
# /api/v1/* → api-core (점진적으로)
```

### 3.3 롤백

```bash
# Ingress 설정만 복구하면 즉시 롤백
kubectl apply -f ingress-legacy.yaml
```

---

## 4. 향후 작업 (마이그레이션 완료 후)

### Phase 1: 안정화 (1-2주)
- api-legacy와 api-core 병행 운영
- 트래픽 비율 조정 (10% → 50% → 100%)
- 모니터링 및 에러 추적

### Phase 2: 정리 (안정화 후)
- api-legacy 서비스 제거
- Liquibase 관련 테이블 삭제 (선택)
  ```sql
  DROP TABLE IF EXISTS database_changelog;
  DROP TABLE IF EXISTS database_changelog_lock;
  ```
- DB명 변경 (선택): `rms-legacy` → `rms`

### Phase 3: audit_logs 마이그레이션 (선택)
- 기존 Hibernate Envers 히스토리 → audit_logs 변환 스크립트
- 별도 스펙으로 관리

---

## 5. 완료 조건

- [x] 스키마 동등성 확인 (user, room, room_group, reservation, payment_method 등)
- [x] api-core 프로덕션 환경에서 마이그레이션 비활성화 확인
- [x] 전략 결정: 동일 DB 직접 사용
- [ ] K8s ConfigMap/Secret 설정 (production-cutover 스펙에서)
- [ ] 프로덕션 배포 및 검증

---

## 6. 참고

### 6.1 테이블 비교 결과

| 테이블 | api-legacy | api-core | 호환성 |
|--------|:----------:|:--------:|:------:|
| user | ✅ | ✅ | 동일 |
| room | ✅ | ✅ | 동일 |
| room_group | ✅ | ✅ | 동일 |
| room_history | ✅ | ✅ | 동일 |
| reservation | ✅ | ✅ | 동일 |
| reservation_room | ✅ | ✅ | 동일 |
| reservation_history | ✅ | ✅ | 동일 |
| reservation_room_history | ✅ | ✅ | 동일 |
| payment_method | ✅ | ✅ | 동일 |
| login_attempts | ✅ | ✅ | 동일 |
| revision_info | ✅ | ✅ | 동일 |
| database_changelog | ✅ | ❌ | Liquibase 전용 |
| database_changelog_lock | ✅ | ❌ | Liquibase 전용 |
| schema_migrations | ❌ | ✅ | Go 전용 |
| audit_logs | ❌ | ✅ | api-core 전용 |

### 6.2 관련 문서

- [production-cutover](../production-cutover/spec.md): K8s 배포 설정
- [history-api-compat](../history-api-compat/spec.md): 히스토리 API 전략
