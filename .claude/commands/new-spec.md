# /new-spec - 새 스펙 생성

새로운 기능/마이그레이션 스펙을 생성합니다.

## 사용법

```
/new-spec <category> <spec-id> "<title>"
```

## 파라미터

- `category`: migration | frontend | infra | api-core
- `spec-id`: 스펙 ID (kebab-case)
- `title`: 스펙 제목

## 예시

```
/new-spec migration db-cleanup "DB 정리"
/new-spec frontend search-component "검색 컴포넌트"
/new-spec infra monitoring-setup "모니터링 설정"
```

## 생성되는 파일

```
specs/{category}/{spec-id}/
├── spec.md    # 스펙 문서
└── plan.md    # 구현 계획 (선택)
```

## 실행 절차

1. 디렉토리 생성: `specs/{category}/{spec-id}/`
2. `specs/_templates/spec.template.md` 복사하여 `spec.md` 생성
3. 템플릿 변수 치환:
   - `{spec-id}` → 입력받은 spec-id
   - `{제목}` → 입력받은 title
   - `{date}` → 오늘 날짜 (YYYY-MM-DD)
   - `{type}` → category
4. `specs/README.md` 인덱스 업데이트
5. 생성된 파일 경로 출력
