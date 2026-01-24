---
spec_id: gitlab-ci-restructure
status: in-progress
started_at: 2026-01-07
completed_at: null
assignee: null

blockers: []
dependencies: []
---

# GitLab CI 빌드 파이프라인 재구성 - 구현 계획

> spec.md의 요구사항을 구현하기 위한 실행 계획

---

## 1. 구현 단계

### Phase 1: 기본 잡 구성 ✅ 완료

- [x] api-core test job 구성
- [x] api-core build job 구성
- [x] api-legacy test job 구성
- [x] api-legacy build job 구성
- [x] frontend-web test job 구성
- [x] frontend-web build job 구성

### Phase 2: 조건부 실행 ✅ 완료

- [x] 앱별 변경 감지 rules 설정
- [x] 병렬 실행 구성

### Phase 3: 테스트 환경 개선 🚧 진행 중

- [ ] MySQL 서비스 포함 테스트
- [ ] Redis 서비스 포함 테스트
- [ ] 통합 테스트 잡 추가

### Phase 4: 선택적 기능

- [ ] 번들 크기 분석 잡
- [ ] 성능 테스트 잡

---

## 2. 기술 결정 사항

### 결정 1: 이미지 태그 형식

**결정**: `{app}/{branch}:{commit-hash}` 형식 사용

**이유**: 앱별 독립 배포 및 버전 추적 용이

### 결정 2: Vue 앱 빌드 분리

**결정**: 기존 Spring Boot 앱 내 static 포함 방식 제거, nginx 기반 별도 컨테이너

**이유**: k8s ingress URL prefix로 분리 운영 예정

---

## 3. 진행 로그

### 2026-01-07

- Phase 1, 2 완료
- Phase 3 진행 중: DB 서비스 포함 테스트 환경 구성 필요

---

## 4. 검증 체크리스트

- [x] api-core 테스트 잡 실행 성공
- [x] api-legacy 테스트 잡 실행 성공
- [x] frontend-web 테스트 잡 실행 성공
- [x] 변경된 앱만 빌드되는지 확인
- [ ] DB 포함 통합 테스트 실행
