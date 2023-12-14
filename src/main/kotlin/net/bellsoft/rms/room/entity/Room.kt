package net.bellsoft.rms.room.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseMustAuditEntity
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.room.type.RoomStatus
import org.hibernate.annotations.Comment
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited
import org.hibernate.envers.RelationTargetAuditMode
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
@Comment("객실")
class Room(
    @Column(name = "number", nullable = false, length = 10)
    @Comment("객실 번호")
    var number: String,

    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "room_group_id_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.EAGER, optional = false)
    @JoinColumn(name = "room_group_id", nullable = false)
    @Comment("객실 그룹")
    var roomGroup: RoomGroup,

    @Column(name = "note", nullable = false, length = 200)
    @Comment("객실 메모")
    var note: String = "",

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    @Comment("객실 상태")
    var status: RoomStatus = RoomStatus.INACTIVE,
) : Serializable, BaseMustAuditEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = 6569216536452327330L
    }
}
