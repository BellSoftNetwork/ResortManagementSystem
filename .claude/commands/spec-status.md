# /spec-status - 스펙 상태 확인

모든 스펙의 현재 상태를 확인합니다.

## 사용법

```
/spec-status
/spec-status <category>
```

## 파라미터

- `category` (선택): migration | frontend | infra | api-core

## 출력 형식

```
📊 스펙 상태 요약

## migration/
| 스펙 | 제목 | spec.md | plan.md |
|------|------|:-------:|:-------:|
| history-api-compat | History API 완전 호환 | ✅ completed | - |
| api-response-compat | API 응답 완전 호환성 | approved | not-started |
| db-schema-unification | DB 스키마 통합 | draft | not-started |

## frontend/
| 스펙 | 제목 | spec.md | plan.md |
|------|------|:-------:|:-------:|
| unit-testing | 유닛 테스트 작성 | draft | not-started |

## infra/
| 스펙 | 제목 | spec.md | plan.md |
|------|------|:-------:|:-------:|
| gitlab-ci-restructure | GitLab CI 재구성 | approved | 🚧 in-progress |

---
총 7개 스펙: 1 completed, 2 approved, 4 draft
```

## 실행 절차

1. `specs/` 하위 모든 `spec.md` 파일 스캔
2. YAML frontmatter에서 status 추출
3. 해당 디렉토리에 `plan.md` 있으면 status 추출
4. 카테고리별로 그룹화하여 테이블 출력
5. 요약 통계 출력
