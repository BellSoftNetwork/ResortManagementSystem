# Hibernate Envers 테이블 구조

> api-legacy에서 Hibernate Envers가 생성한 히스토리 테이블 구조

---

## revision_info 테이블

```sql
CREATE TABLE revision_info (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  timestamp BIGINT NOT NULL,  -- Unix timestamp (밀리초)
  user_id BIGINT              -- 변경한 사용자 ID (nullable)
);
```

---

## *_history 테이블 공통 구조

예: `room_history`

```sql
CREATE TABLE room_history (
  rev BIGINT NOT NULL,           -- revision_info.id 참조
  revtype TINYINT,               -- 0=ADD, 1=MOD, 2=DEL
  id BIGINT NOT NULL,            -- 원본 엔티티 ID
  
  -- 원본 테이블의 모든 컬럼들
  number VARCHAR(10),
  note VARCHAR(200),
  status TINYINT,
  
  -- 변경 플래그 컬럼들 (*_mod)
  number_mod BIT(1),
  note_mod BIT(1),
  status_mod BIT(1),
  
  PRIMARY KEY (rev, id),
  FOREIGN KEY (rev) REFERENCES revision_info(id)
);
```

---

## RevisionType 매핑

| revtype | Envers Enum | 프론트엔드 HistoryType |
|---------|-------------|----------------------|
| 0 | ADD | CREATED |
| 1 | MOD | UPDATED |
| 2 | DEL | DELETED |

---

## 기존 히스토리 테이블

api-legacy에서 사용하는 히스토리 테이블:

- `room_history`
- `reservation_history`
- `reservation_room_history`

---

## api-core audit_logs 테이블

api-core에서는 Envers 대신 `audit_logs` 테이블 사용:

```sql
CREATE TABLE audit_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  entity_type VARCHAR(100) NOT NULL,    -- "room", "reservation" 등
  entity_id BIGINT NOT NULL,
  action VARCHAR(20) NOT NULL,          -- "CREATE", "UPDATE", "DELETE"
  old_values JSON,                       -- 변경 전 값
  new_values JSON,                       -- 변경 후 값
  changed_fields JSON,                   -- 변경된 필드 목록
  user_id BIGINT,
  username VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_audit_logs_entity (entity_type, entity_id)
);
```

---

## 마이그레이션 결정

**결정**: 기존 Envers 히스토리 데이터 포기, 새로운 `audit_logs` 테이블만 사용

- 기존 `*_history` 테이블 데이터는 마이그레이션하지 않음
- 마이그레이션 시점 이후부터 새로 쌓이는 데이터만 사용
- 향후 필요시 별도 마이그레이션 스크립트로 처리 가능

자세한 내용은 [history-api-compat 스펙](../../specs/migration/history-api-compat/spec.md) 참조.
