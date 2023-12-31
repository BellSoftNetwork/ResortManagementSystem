package net.bellsoft.rms.reservation.entity

import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseMustAuditEntity
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.room.entity.Room
import org.hibernate.annotations.Comment
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited
import java.io.Serial
import java.io.Serializable

@Entity
@Audited
@AuditTable("reservation_room_history")
@Table(
    name = "reservation_room",
    uniqueConstraints = [
        UniqueConstraint(
            name = "uc_reservation_room_reservation_id_and_room_id",
            columnNames = ["reservation_id", "room_id", "deleted_at"],
        ),
    ],
)
@SQLDelete(sql = "UPDATE reservation_room SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
@Comment("예약 객실")
class ReservationRoom(
    @ManyToOne(fetch = FetchType.EAGER, optional = false)
    @JoinColumn(name = "reservation_id", nullable = false)
    @Comment("예약 정보")
    var reservation: Reservation,

    @ManyToOne(fetch = FetchType.EAGER, optional = false)
    @JoinColumn(name = "room_id", nullable = false)
    @Comment("배정 객실")
    var room: Room,
) : Serializable, BaseMustAuditEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = 737436737363715200L
    }
}
