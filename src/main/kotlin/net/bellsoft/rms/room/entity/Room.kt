package net.bellsoft.rms.room.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseMustAuditEntity
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.room.type.RoomStatus
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited
import java.io.Serial
import java.io.Serializable

@Entity
@Audited(withModifiedFlag = true)
@AuditTable("room_history")
@Table(
    name = "room",
    uniqueConstraints = [
        UniqueConstraint(name = "uc_room_number", columnNames = ["number", "deleted_at"]),
    ],
)
@SQLDelete(sql = "UPDATE room SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
class Room(
    @Column(name = "number", nullable = false, length = 10)
    var number: String,

    @Column(name = "peek_price", nullable = false)
    var peekPrice: Int = 0,

    @Column(name = "off_peek_price", nullable = false)
    var offPeekPrice: Int = 0,

    @Column(name = "description", nullable = false, length = 200)
    var description: String = "",

    @Column(name = "note", nullable = false, length = 200)
    var note: String = "",

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    var status: RoomStatus = RoomStatus.INACTIVE,
) : Serializable, BaseMustAuditEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = -4336707694159308855L
    }
}
