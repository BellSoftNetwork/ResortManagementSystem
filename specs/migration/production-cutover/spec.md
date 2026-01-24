---
id: production-cutover
title: "운영 환경 전환 준비"
status: approved
type: migration
version: 1.2.0
created: 2026-01-07
updated: 2026-01-12

depends_on: [api-response-compat, db-schema-unification]
related: [legacy-cleanup]
replaces: null
replaced_by: null

impact: infra
risk: high
effort: large

changelog:
  - date: 2026-01-12
    description: "프론트엔드 연동 검증 완료, 블루-그린 실행 계획 추가 (plan.md)"
  - date: 2026-01-11
    description: "api-response-compat 완료로 상태 approved 전환"
  - date: 2026-01-07
    description: "TODO.md Phase 4에서 마이그레이션"
---

# 운영 환경 전환 준비

> api-legacy를 내리고 api-core로 완전 전환하기 위한 준비

---

## 1. 개요

### 1.1 배경 및 문제

두 API가 동시에 운영 중이며, api-core로 완전 전환 필요:
- 트래픽 전환 계획
- 롤백 전략
- 모니터링 설정

### 1.2 목표

- K8s 배포 구성 업데이트
- 무중단 전환 계획 수립
- 모니터링 및 알림 설정

### 1.3 비목표 (Non-Goals)

- api-legacy 코드 제거 (별도 스펙)

---

## 2. K8s 배포 구성

### 2.1 디렉토리 구조

```
deploy/example/base/
├── api-core/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
└── ingress.yaml
```

### 2.2 Ingress 라우팅

```yaml
# /api/v1/* → api-core 서비스
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /api/v1
            pathType: Prefix
            backend:
              service:
                name: api-core
                port:
                  number: 8080
```

---

## 3. 무중단 전환 계획

### 3.1 전환 순서

1. api-core를 api-legacy와 동일한 클러스터에 배포
2. 일부 트래픽을 api-core로 라우팅 (Canary)
3. 모니터링 및 문제 확인
4. 전체 트래픽 api-core로 전환
5. api-legacy 제거

### 3.2 롤백 계획

- Ingress 설정 복구로 즉시 롤백 가능
- api-legacy Pod 유지 (전환 안정화까지)

---

## 4. 모니터링

### 4.1 헬스체크

```yaml
livenessProbe:
  httpGet:
    path: /actuator/health/liveness
    port: 8080
readinessProbe:
  httpGet:
    path: /actuator/health/readiness
    port: 8080
```

### 4.2 메트릭

- 에러율
- 응답시간 (p50, p95, p99)
- 요청 수

---

## 5. 완료 조건

- [ ] K8s Deployment/Service/ConfigMap 확인
- [ ] Ingress 라우팅 규칙 수정
- [ ] 환경 변수 및 Secret 설정
- [ ] Canary 배포 전략 구현
- [ ] 롤백 계획 수립
- [x] 헬스체크 엔드포인트 연결 확인 (api-core /actuator/health/* 동작 확인)
- [ ] 로그 수집 설정
- [ ] 모니터링 대시보드 구성

---

## 6. 선행 조건 상태

| 조건 | 상태 | 비고 |
|------|:----:|------|
| API 응답 호환성 (api-response-compat) | ✅ | 40/40 엔드포인트 검증 완료 |
| DB 스키마 통합 (db-schema-unification) | ✅ | 동일 DB 직접 사용 전략 확정 |
| 호환성 테스트 스크립트 | ✅ | `scripts/api-compatibility-test.py` |
| Golden file 테스트 | ✅ | 25개 golden file 저장 |
| **프론트엔드 연동 검증** | ✅ | 브라우저 테스트 완료, 에러 없음 |
| **블루-그린 실행 계획** | ✅ | `plan.md` 작성 완료 |

---

## 7. 검증 결과 요약 (2026-01-12)

### 프론트엔드 연동 테스트 결과

| 페이지 | 상태 | 확인 내용 |
|--------|:----:|----------|
| 로그인 | ✅ | testadmin 계정 로그인 성공 |
| 대시보드 | ✅ | 예약 캘린더, 테이블 정상 표시 |
| 객실 관리 | ✅ | 58개 객실 목록 정상 |
| 예약 관리 | ✅ | 18개 예약 목록 정상 |
| 통계 | ✅ | 매출/예약/인원 차트 정상 |
| 결제 수단 | ✅ | 5개 결제 수단 정상 |
| 계정 관리 | ✅ | 5개 계정 정상 |
| 콘솔 에러 | ✅ | 없음 |

### 알려진 차이점 (무해)

| 항목 | api-legacy | api-core | 영향 |
|------|------------|----------|------|
| /actuator/health 응답 | `{"status":"UP","groups":[...]}` | `{"status":"UP","components":{...}}` | K8s 헬스체크 영향 없음 |
| History API | Hibernate Envers 데이터 | 신규 audit_logs 데이터 | 기존 히스토리 미표시 (의도적)
