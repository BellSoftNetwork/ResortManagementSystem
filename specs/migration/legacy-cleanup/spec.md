---
id: legacy-cleanup
title: "레거시 정리"
status: draft
type: migration
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

depends_on: [production-cutover]
related: []
replaces: null
replaced_by: null

impact: both
risk: low
effort: small

changelog:
  - date: 2026-01-07
    description: "TODO.md Phase 5에서 마이그레이션"
---

# 레거시 정리

> 전환 완료 후 api-legacy 관련 코드 및 리소스 정리

---

## 1. 개요

### 1.1 배경 및 문제

api-core로 완전 전환 후 불필요한 레거시 자원 정리 필요:
- api-legacy 코드
- CI/CD 파이프라인
- 문서
- DB 테이블

### 1.2 목표

- api-legacy 완전 제거
- 문서 정리
- DB 정리 (선택)

### 1.3 비목표 (Non-Goals)

- 새로운 기능 추가

---

## 2. 정리 대상

### 2.1 api-legacy 제거

- [ ] api-legacy 컨테이너/Pod 제거
- [ ] api-legacy CI/CD 파이프라인 비활성화/제거
- [ ] `apps/api-legacy/` 디렉토리 제거 (또는 archive 브랜치로 이동)

### 2.2 문서 업데이트

- [ ] `README.md` - api-legacy 관련 내용 제거
- [ ] `CLAUDE.md` → `AGENTS.md` - Go API만 설명
- [ ] `docs/API_COMPARISON.md` 제거 또는 아카이브

### 2.3 DB 정리 (선택)

Liquibase 관련 테이블:
- [ ] `databasechangelog` 제거 결정
- [ ] `databasechangeloglock` 제거 결정

Hibernate Envers 히스토리 테이블:
- [ ] `room_history` 유지/제거 결정
- [ ] `reservation_history` 유지/제거 결정
- [ ] `reservation_room_history` 유지/제거 결정

**참고**: 새 `audit_logs` 테이블로 완전 대체 시 제거 가능

---

## 3. 완료 조건

- [ ] api-legacy 코드 제거/아카이브
- [ ] CI/CD 정리
- [ ] 문서 업데이트
- [ ] DB 정리 결정 및 실행 (선택)
