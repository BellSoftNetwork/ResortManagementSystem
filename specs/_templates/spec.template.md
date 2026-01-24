---
id: {spec-id}
title: "{제목}"
status: draft  # draft | approved | completed | deprecated
type: {type}   # migration | product | infra | refactor
version: 1.0.0
created: {date}
updated: {date}

depends_on: []
related: []
replaces: null
replaced_by: null

impact: both  # backend | frontend | both | infra
risk: medium  # low | medium | high
effort: medium  # small | medium | large

changelog:
  - date: {date}
    description: "초기 스펙 작성"
---

# {제목}

> 한 줄 요약: 이 스펙이 해결하려는 문제와 목표

---

## 1. 개요

### 1.1 배경 및 문제

현재 상황과 해결해야 할 문제를 설명합니다.

### 1.2 목표

이 스펙을 통해 달성하려는 목표:
- 목표 1
- 목표 2

### 1.3 비목표 (Non-Goals)

이 스펙에서 다루지 않는 범위:
- 비목표 1

---

## 2. 상세 설계

### 2.1 주요 변경 사항

변경되는 내용을 상세히 설명합니다.

### 2.2 데이터 구조

```typescript
// 또는 Go, Kotlin 등 해당 언어로 작성
interface Example {
  field: string
}
```

### 2.3 API 변경 (해당되는 경우)

| Method | Path | Description | 변경 사항 |
|--------|------|-------------|----------|
| GET | `/api/v1/example` | 예시 | 신규 |

---

## 3. 호환성 고려사항

### 3.1 하위 호환성

기존 시스템과의 호환성 유지 방안.

### 3.2 마이그레이션

기존 데이터/코드 마이그레이션 필요 여부와 방법.

---

## 4. 테스트 계획

### 4.1 단위 테스트

- [ ] 테스트 케이스 1
- [ ] 테스트 케이스 2

### 4.2 통합 테스트

- [ ] 테스트 케이스 1

---

## 5. 완료 조건

- [ ] 조건 1
- [ ] 조건 2
- [ ] 테스트 통과
- [ ] 코드 리뷰 완료

---

## 6. 참고 자료

- [관련 문서 링크]()
