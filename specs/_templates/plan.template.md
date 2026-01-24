---
spec_id: {spec-id}
status: not-started  # not-started | in-progress | blocked | done
started_at: null
completed_at: null
assignee: null

blockers: []
dependencies: []
---

# {제목} - 구현 계획

> spec.md의 요구사항을 구현하기 위한 실행 계획

---

## 1. 구현 단계

### Phase 1: 준비

- [ ] 관련 코드 분석
- [ ] 테스트 환경 준비

### Phase 2: 핵심 구현

- [ ] 구현 항목 1
- [ ] 구현 항목 2

### Phase 3: 테스트 및 검증

- [ ] 단위 테스트 작성
- [ ] 통합 테스트 실행
- [ ] 기존 기능 회귀 테스트

### Phase 4: 문서화 및 정리

- [ ] 코드 주석 정리
- [ ] README/문서 업데이트
- [ ] 불필요한 코드 제거

---

## 2. 기술 결정 사항

### 결정 1: {제목}

**문제**: 해결해야 할 기술적 문제

**옵션**:
1. 옵션 A - 장단점
2. 옵션 B - 장단점

**결정**: 선택한 옵션과 이유

---

## 3. 리스크 및 대응

| 리스크 | 영향도 | 대응 방안 |
|--------|--------|----------|
| 리스크 1 | 높음 | 대응 방안 |

---

## 4. 진행 로그

### {날짜}

- 진행 내용 기록
- 발견된 이슈 및 해결 방안

---

## 5. 검증 체크리스트

- [ ] `docker compose exec api-core make test` 통과
- [ ] `docker compose exec api-core make lint` 통과
- [ ] `docker compose exec frontend yarn lint` 통과
- [ ] `docker compose exec frontend yarn build` 통과
- [ ] 기존 기능 정상 동작 확인
