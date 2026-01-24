---
id: api-legacy-history-envers
title: "api-legacy Hibernate Envers 히스토리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: medium
---

# api-legacy Hibernate Envers 히스토리

> Hibernate Envers 기반 엔티티 변경 이력 추적

---

## 1. 개요

### 1.1 역할

- 엔티티 변경 이력 자동 추적
- 변경된 필드 감지 (withModifiedFlag)
- 리비전 기반 히스토리 조회

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `RevisionInfo.kt` | 리비전 엔티티 |
| `RevisionDetails.kt` | 리비전 상세 DTO |
| `EntityRevisionDto.kt` | 엔티티 리비전 응답 |
| `EntityRevisionComponent.kt` | 리비전 조회 컴포넌트 |
| `RevisionDetailsRepository.kt` | 리비전 리포지토리 |
| `EnversConfig.kt` | Envers 설정 |

---

## 2. 감사 대상 엔티티

| 엔티티 | 히스토리 테이블 | 변경 플래그 |
|--------|-----------------|:-----------:|
| Room | room_history | ✅ |
| RoomGroup | room_group_history | ✅ |
| Reservation | reservation_history | ✅ |
| ReservationRoom | reservation_room_history | ❌ |

---

## 3. 엔티티 설정

### 3.1 감사 대상 엔티티

```kotlin
@Entity
@Audited(withModifiedFlag = true)
@AuditTable("room_history")
class Room(
    var number: String,
    
    @Audited(withModifiedFlag = true, targetAuditMode = NOT_AUDITED)
    @ManyToOne(fetch = EAGER)
    var roomGroup: RoomGroup,
    
    var note: String,
    var status: RoomStatus
) : BaseMustAuditEntity()
```

### 3.2 RevisionInfo 엔티티

```kotlin
@Entity
@RevisionEntity
@Table(name = "revision_info")
class RevisionInfo {
    @Id
    @RevisionNumber
    @GeneratedValue(strategy = IDENTITY)
    val id: Long = 0
    
    @RevisionTimestamp
    @Column(name = "created_at")
    var createdAt: LocalDateTime = LocalDateTime.now()
}
```

---

## 4. 히스토리 조회

### 4.1 EntityRevisionComponent

```kotlin
@Component
class EntityRevisionComponent(
    @PersistenceContext
    private val entityManager: EntityManager
) {
    fun <T : Any> getRevisions(
        entityClass: KClass<T>,
        id: Long,
        pageable: Pageable
    ): EntityListDto<Triple<T, RevisionDetails, Set<String>>, EntityRevisionDto<*>> {
        val auditReader = AuditReaderFactory.get(entityManager)
        val revisions = auditReader.getRevisions(entityClass.java, id)
        // 리비전별 엔티티 스냅샷 조회
        // 변경된 필드 추출
        // DTO 변환
    }
}
```

### 4.2 변경 필드 감지

```kotlin
fun getModifiedFields(
    auditReader: AuditReader,
    entityClass: Class<*>,
    id: Long,
    revisionNumber: Number
): Set<String> {
    val query = auditReader.createQuery()
        .forRevisionsOfEntity(entityClass, false, true)
        .add(AuditEntity.id().eq(id))
        .add(AuditEntity.revisionNumber().eq(revisionNumber))
    
    val result = query.singleResult as Array<*>
    val propertyNames = result[1] as Set<String>
    return propertyNames.filter { it.endsWith("_MOD") && result[it] == true }
        .map { it.removeSuffix("_MOD") }
        .toSet()
}
```

---

## 5. DTO

### 5.1 RevisionDetails

```kotlin
data class RevisionDetails(
    val revisionId: Long,
    val revisionType: RevisionType,
    val timestamp: LocalDateTime
)
```

### 5.2 EntityRevisionDto

```kotlin
data class EntityRevisionDto<T>(
    val entity: T,
    val historyType: HistoryType,
    val historyCreatedAt: LocalDateTime,
    val updatedFields: List<String>
)
```

### 5.3 HistoryType

```kotlin
enum class HistoryType {
    CREATED,  // RevisionType.ADD
    UPDATED,  // RevisionType.MOD
    DELETED   // RevisionType.DEL
}
```

---

## 6. API 응답

### 6.1 히스토리 조회

```http
GET /api/v1/rooms/1/histories?page=0&size=20
```

### 6.2 응답

```json
{
  "data": [
    {
      "entity": {
        "id": 1,
        "number": "101",
        "roomGroupId": 1,
        "note": "수정된 메모",
        "status": "NORMAL"
      },
      "historyType": "UPDATED",
      "historyCreatedAt": "2026-01-07T14:00:00",
      "updatedFields": ["note"]
    },
    {
      "entity": { ... },
      "historyType": "CREATED",
      "historyCreatedAt": "2026-01-01T10:00:00",
      "updatedFields": []
    }
  ],
  "page": { ... }
}
```

---

## 7. 히스토리 테이블 구조

### 7.1 room_history

| 컬럼 | 설명 |
|------|------|
| id | 객실 ID |
| rev | revision_info.id 참조 |
| revtype | 0=INSERT, 1=UPDATE, 2=DELETE |
| number | 객실 번호 |
| number_mod | number 변경 여부 |
| room_group_id | 객실 그룹 ID |
| room_group_id_mod | room_group_id 변경 여부 |
| note | 메모 |
| note_mod | note 변경 여부 |
| status | 상태 |
| status_mod | status 변경 여부 |

---

## 8. api-core와의 차이점

| 항목 | api-legacy | api-core |
|------|------------|----------|
| 히스토리 저장 | Hibernate Envers | audit_logs 테이블 |
| 테이블 | *_history | audit_logs |
| 변경 필드 | _MOD 플래그 | changed_fields JSON |
| 스냅샷 | 전체 엔티티 | old/new values JSON |

---

## 9. 마이그레이션 고려사항

- 기존 Envers 히스토리 데이터는 마이그레이션하지 않음
- api-core 전환 시점 이후부터 audit_logs에 새로 쌓임
- 필요시 별도 마이그레이션 스크립트로 변환 가능
