---
id: gitlab-ci-restructure
title: "GitLab CI 빌드 파이프라인 재구성"
status: approved
type: infra
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: []
related: []
replaces: null
replaced_by: null

impact: infra
risk: medium
effort: medium

changelog:
  - date: 2026-01-07
    description: "TODO.md에서 마이그레이션 (진행 중)"
---

# GitLab CI 빌드 파이프라인 재구성

> 모노레포 구조에 맞게 앱별 독립 빌드/테스트/배포 파이프라인 구성

---

## 1. 개요

### 1.1 배경 및 문제

기존 모놀리스 방식에서 앱별 분리가 필요:
- 기존: `resort-management-system/main:latest`
- 목표: `resort-management-system/api-core/main:latest`

### 1.2 목표

- 개별 앱 단위로 빌드 잡 구성
- 변경된 앱만 빌드되도록 조건 설정
- 테스트/린트 → 빌드 순서 보장
- 병렬 실행으로 효율화

### 1.3 비목표 (Non-Goals)

- 배포 자동화 (별도 계획)

---

## 2. 파이프라인 구조

### 2.1 앱별 잡 구성

#### api-core (Go)

| Job | 설명 | 조건 |
|-----|------|------|
| `test:api-core` | Go 테스트 + 린트 | `/apps/api-core/**` 변경 시 |
| `build:api-core` | Docker 이미지 빌드 | test 성공 시 |
| `integration-test:api-core` | 통합 테스트 | 수동 실행 |

#### api-legacy (Spring Boot)

| Job | 설명 | 조건 |
|-----|------|------|
| `test:api-legacy` | Gradle 테스트 + ktlint | `/apps/api-legacy/**` 변경 시 |
| `build:api-legacy` | Docker 이미지 빌드 | test 성공 시 |

#### frontend-web (Vue.js)

| Job | 설명 | 조건 |
|-----|------|------|
| `test:frontend-web` | Yarn lint + 테스트 | `/apps/frontend-web/**` 변경 시 |
| `build:frontend-web` | Docker 이미지 빌드 (nginx) | test 성공 시 |
| `analyze:frontend-web` | 번들 크기 분석 | 선택적 |

### 2.2 이미지 태그 형식

```
registry.gitlab.com/project/resort-management-system/{app}/{branch}:{commit-hash}
registry.gitlab.com/project/resort-management-system/{app}/{branch}:latest
```

---

## 3. 조건부 실행

```yaml
# .gitlab-ci.yml 예시
test:api-core:
  rules:
    - changes:
        - apps/api-core/**/*
```

---

## 4. 완료 조건

- [x] api-core 빌드 잡 구성 (test → build)
- [x] api-legacy 빌드 잡 구성 (test → build)
- [x] frontend-web 빌드 잡 구성 (test → build)
- [x] 변경된 앱만 빌드되도록 조건 설정
- [ ] MySQL/Redis 서비스 포함 테스트 환경
- [ ] 통합 테스트 잡 (수동 실행)
- [ ] 번들 분석 잡 (선택적)
