---
id: api-legacy-room-group-management
title: "api-legacy 객실 그룹 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-legacy 객실 그룹 관리

> 객실 그룹(카테고리) CRUD

---

## 1. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/room-groups` | 목록 | USER |
| GET | `/api/v1/room-groups/{id}` | 상세 | USER |
| POST | `/api/v1/room-groups` | 생성 | ADMIN |
| PATCH | `/api/v1/room-groups/{id}` | 수정 | ADMIN |
| DELETE | `/api/v1/room-groups/{id}` | 삭제 | ADMIN |

---

## 2. RoomGroup 엔티티

```kotlin
@Entity
@Audited
@AuditTable("room_group_history")
@Table(
    uniqueConstraints = [
        UniqueConstraint(columnNames = ["name", "deleted_at"])
    ]
)
@SQLDelete(sql = "UPDATE room_group SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at = '1970-01-01 00:00:00'")
@Comment("객실 그룹")
class RoomGroup(
    @Column(nullable = false, length = 20)
    var name: String,
    
    @Column(name = "peek_price", nullable = false)
    var peekPrice: Int,
    
    @Column(name = "off_peek_price", nullable = false)
    var offPeekPrice: Int,
    
    @Column(nullable = false, length = 200)
    var description: String = "",
    
    @OneToMany(mappedBy = "roomGroup", cascade = [PERSIST])
    @OrderBy("id ASC")
    val rooms: MutableList<Room> = mutableListOf()
) : BaseMustAuditEntity()
```

---

## 3. 서비스

```kotlin
interface RoomGroupService {
    fun getAll(pageable: Pageable): EntityListDto<RoomGroup, RoomGroupSummaryDto>
    fun getById(id: Long, filter: RoomGroupDetailFilterRequest?): RoomGroupDetailDto
    fun create(request: RoomGroupCreateRequest): RoomGroupSummaryDto
    fun update(id: Long, request: RoomGroupPatchRequest): RoomGroupSummaryDto
    fun delete(id: Long)
}
```

---

## 4. DTO

### 4.1 Summary (목록용)

```kotlin
data class RoomGroupSummaryDto(
    val id: Long,
    val name: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto?,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto?
)
```

### 4.2 Detail (상세용)

```kotlin
data class RoomGroupDetailDto(
    val id: Long,
    val name: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val rooms: List<RoomLastStayDetailDto>,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto?,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto?
)

data class RoomLastStayDetailDto(
    val room: RoomDetailDto,
    val lastReservation: ReservationDetailDto?
)
```

---

## 5. 삭제 제약

```kotlin
@Transactional
override fun delete(id: Long) {
    val roomGroup = roomGroupRepository.findByIdOrThrow(id)
    
    if (roomGroup.rooms.isNotEmpty()) {
        throw ConflictException("객실이 포함된 그룹은 삭제할 수 없습니다")
    }
    
    roomGroupRepository.delete(roomGroup)
}
```
