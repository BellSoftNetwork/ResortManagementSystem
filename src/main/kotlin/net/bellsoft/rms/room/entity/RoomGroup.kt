package net.bellsoft.rms.room.entity

import jakarta.persistence.CascadeType
import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.OneToMany
import jakarta.persistence.OrderBy
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseMustAuditEntity
import net.bellsoft.rms.common.entity.BaseTimeEntity
import org.hibernate.annotations.Comment
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import java.io.Serial
import java.io.Serializable

@Entity
@AuditTable("room_group_history")
@Table(
    name = "room_group",
    uniqueConstraints = [
        UniqueConstraint(name = "uc_room_group_name", columnNames = ["name", "deleted_at"]),
    ],
)
@SQLDelete(sql = "UPDATE room_group SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
@Comment("객실 그룹")
class RoomGroup(
    @Column(name = "name", nullable = false, length = 20)
    @Comment("객실 그룹명")
    var name: String,

    @Column(name = "peek_price", nullable = false)
    @Comment("성수기 가격")
    var peekPrice: Int = 0,

    @Column(name = "off_peek_price", nullable = false)
    @Comment("비성수기 가격")
    var offPeekPrice: Int = 0,

    @Column(name = "description", nullable = false, length = 200)
    @Comment("객실 그룹 설명")
    var description: String = "",

    @OneToMany(mappedBy = "roomGroup", cascade = [CascadeType.PERSIST])
    @OrderBy("id ASC")
    @Comment("객실")
    var rooms: MutableList<Room> = mutableListOf(),
) : Serializable, BaseMustAuditEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = -3767055553695226176L
    }
}
