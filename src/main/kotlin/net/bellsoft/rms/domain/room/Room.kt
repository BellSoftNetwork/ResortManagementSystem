package net.bellsoft.rms.domain.room

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.domain.base.BaseMustAudit
import net.bellsoft.rms.domain.base.BaseTime
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited

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
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
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
) : BaseMustAudit() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set
}
