---
id: api-legacy-room-management
title: "api-legacy 객실 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-legacy 객실 관리

> 객실 CRUD 및 Envers 히스토리

---

## 1. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/rooms` | 객실 목록 | USER |
| GET | `/api/v1/rooms/{id}` | 객실 상세 | USER |
| POST | `/api/v1/rooms` | 객실 생성 | ADMIN |
| PATCH | `/api/v1/rooms/{id}` | 객실 수정 | ADMIN |
| DELETE | `/api/v1/rooms/{id}` | 객실 삭제 | ADMIN |
| GET | `/api/v1/rooms/{id}/histories` | 변경 이력 | ADMIN |

---

## 2. Room 엔티티

```kotlin
@Entity
@Audited(withModifiedFlag = true)
@AuditTable("room_history")
@Table(
    uniqueConstraints = [
        UniqueConstraint(columnNames = ["number", "deleted_at"])
    ]
)
@SQLDelete(sql = "UPDATE room SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at = '1970-01-01 00:00:00'")
@Comment("객실")
class Room(
    @Column(nullable = false, length = 10)
    var number: String,
    
    @Audited(withModifiedFlag = true, targetAuditMode = NOT_AUDITED)
    @ManyToOne(fetch = EAGER, optional = false)
    @JoinColumn(nullable = false)
    var roomGroup: RoomGroup,
    
    @Column(nullable = false, length = 200)
    var note: String = "",
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var status: RoomStatus = RoomStatus.NORMAL
) : BaseMustAuditEntity()
```

---

## 3. 객실 상태

```kotlin
enum class RoomStatus(val value: Int) {
    DAMAGED(-10),
    CONSTRUCTION(-1),
    INACTIVE(0),
    NORMAL(1)
}
```

---

## 4. 서비스

```kotlin
interface RoomService {
    fun getAll(pageable: Pageable, filter: RoomFilterRequest): EntityListDto<Room, RoomDetailDto>
    fun getById(id: Long): RoomDetailDto
    fun create(request: RoomCreateRequest): RoomDetailDto
    fun update(id: Long, request: RoomPatchRequest): RoomDetailDto
    fun delete(id: Long)
    fun getHistories(id: Long, pageable: Pageable): EntityListDto<Triple<Room, RevisionDetails, Set<String>>, EntityRevisionDto<RoomDetailDto>>
}
```

---

## 5. 필터

```kotlin
data class RoomFilterRequest(
    val roomGroupId: Long?,
    val status: RoomStatus?,
    val search: String?
)
```

---

## 6. QueryDSL 쿼리

```kotlin
class RoomCustomRepositoryImpl(
    private val queryFactory: JPAQueryFactory
) : RoomCustomRepository {
    
    override fun findAllWithFilter(
        pageable: Pageable,
        filter: RoomFilterRequest
    ): Page<Room> {
        val query = queryFactory
            .selectFrom(room)
            .leftJoin(room.roomGroup, roomGroup).fetchJoin()
            .where(
                roomGroupIdEq(filter.roomGroupId),
                statusEq(filter.status),
                numberContains(filter.search)
            )
        
        val content = query
            .offset(pageable.offset)
            .limit(pageable.pageSize.toLong())
            .orderBy(getOrderSpecifier(pageable.sort))
            .fetch()
        
        val total = queryFactory
            .select(room.count())
            .from(room)
            .where(/* same conditions */)
            .fetchOne() ?: 0L
        
        return PageImpl(content, pageable, total)
    }
}
```

---

## 7. 히스토리 조회

```kotlin
@GetMapping("/{id}/histories")
@PreAuthorize("hasAnyRole('ADMIN', 'SUPER_ADMIN')")
fun getHistories(
    @PathVariable id: Long,
    pageable: Pageable
): ListResponse<EntityRevisionDto<RoomDetailDto>> {
    return ListResponse(
        roomService.getHistories(id, pageable)
    )
}
```
