---
id: api-response-compat
title: "API 응답 완전 호환성 검증"
status: completed
type: migration
version: 1.1.0
created: 2026-01-07
updated: 2026-01-11

depends_on: [history-api-compat]
related: [production-cutover]
replaces: null
replaced_by: null

impact: both
risk: high
effort: large

changelog:
  - date: 2026-01-11
    description: "호환성 테스트 스크립트 완성, 40/40 엔드포인트 검증 통과"
  - date: 2026-01-07
    description: "TODO.md Phase 2에서 마이그레이션"
---

# API 응답 완전 호환성 검증

> api-core의 모든 API 응답이 api-legacy와 완전히 동일한지 검증하고 수정

---

## 1. 개요

### 1.1 배경 및 문제

프론트엔드가 api-legacy 응답을 기준으로 개발되어 있으므로, 모든 응답이 완전히 동일해야 함:
- JSON 키 존재 여부
- 값 타입 (int vs int64)
- null 처리 차이
- 날짜/시간 형식

### 1.2 목표

- 양쪽 API 응답을 자동으로 비교하는 스크립트 개선
- 모든 차이점 식별 및 수정
- Golden file 기반 호환성 테스트

### 1.3 비목표 (Non-Goals)

- API 기능 추가/변경 (호환성만 검증)

---

## 2. 비교 대상 API

### 2.1 전체 API 목록

| API | Path | 상태 |
|-----|------|------|
| 객실 목록 | `GET /api/v1/rooms` | ✅ 검증 완료 |
| 객실 상세 | `GET /api/v1/rooms/{id}` | ✅ 검증 완료 |
| 객실 이력 | `GET /api/v1/rooms/{id}/histories` | ✅ 검증 완료 |
| 객실 그룹 목록 | `GET /api/v1/room-groups` | ✅ 검증 완료 |
| 객실 그룹 상세 | `GET /api/v1/room-groups/{id}` | ✅ 검증 완료 |
| 예약 목록 | `GET /api/v1/reservations` | ✅ 검증 완료 |
| 예약 상세 | `GET /api/v1/reservations/{id}` | ✅ 검증 완료 |
| 예약 이력 | `GET /api/v1/reservations/{id}/histories` | ✅ 검증 완료 |
| 예약 통계 | `GET /api/v1/reservation-statistics` | ✅ 검증 완료 |
| 결제 수단 | `GET /api/v1/payment-methods` | ✅ 검증 완료 |
| 현재 사용자 | `GET /api/v1/my` | ✅ 검증 완료 |
| 사용자 목록 | `GET /api/v1/admin/accounts` | ✅ 검증 완료 |
| 앱 설정 | `GET /api/v1/config` | ✅ 검증 완료 |
| 환경 정보 | `GET /api/v1/env` | ✅ 검증 완료 |

**검증 도구**: `python3 scripts/api-compatibility-test.py`
**결과**: 40/40 엔드포인트 정상 동작 (2026-01-11)

---

## 3. 알려진 호환성 이슈

### 3.1 JSON 키 누락 이슈

Go의 `omitempty` 태그로 인해 null/빈 값 필드가 누락될 수 있음:

```go
// ❌ 문제: 값이 없으면 키 자체가 사라짐
Note *string `json:"note,omitempty"`

// ✅ 해결: omitempty 제거
Note *string `json:"note"`
```

### 3.2 날짜/시간 형식

```
api-legacy: 2024-01-15T10:30:00 (타임존 없음)
api-core:   동일 형식 확인 필요
```

### 3.3 숫자 타입

int vs int64 등의 차이로 인한 JSON 직렬화 차이 확인 필요

### 3.4 페이지네이션 메타데이터

Spring Boot 페이지네이션 응답 형식과 완전 일치 필요

---

## 4. 검증 방법

### 4.1 비교 스크립트

```bash
# scripts/compare-api-responses.py
python3 scripts/compare-api-responses.py
```

### 4.2 Golden File 테스트

api-legacy 실제 응답을 golden file로 저장 후 api-core 응답과 비교

---

## 5. 완료 조건

- [x] `scripts/api-compatibility-test.py` 스크립트 완성
- [x] 모든 API 응답 비교 완료 (40/40 엔드포인트)
- [x] JSON 키 누락 이슈 수정 (omitempty 제거)
- [x] 날짜/시간 형식 통일 (ISO 8601, 타임존 없음)
- [x] 숫자 타입 통일
- [x] 페이지네이션 형식 통일 (Spring Boot 호환)
- [x] Golden file 기반 테스트 추가 (25개 golden file 저장)

## 6. 관련 파일

- `scripts/api-compatibility-test.py`: 호환성 테스트 스크립트
- `tests/contracts/api-endpoints-manifest.yaml`: API 엔드포인트 매니페스트
- `tests/contracts/golden/`: Golden file 저장 디렉토리 (25개 파일)
