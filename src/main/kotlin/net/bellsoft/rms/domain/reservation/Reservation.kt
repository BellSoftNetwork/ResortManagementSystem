package net.bellsoft.rms.domain.reservation

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.OneToMany
import jakarta.persistence.OrderBy
import jakarta.persistence.Table
import net.bellsoft.rms.domain.base.BaseMustAuditEntity
import net.bellsoft.rms.domain.base.BaseTimeEntity
import net.bellsoft.rms.domain.paymentmethod.PaymentMethod
import net.bellsoft.rms.domain.room.Room
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited
import org.hibernate.envers.RelationTargetAuditMode
import java.io.Serial
import java.io.Serializable
import java.time.LocalDate
import java.time.LocalDateTime

@Entity
@Audited(withModifiedFlag = true)
@AuditTable("reservation_history")
@Table(name = "reservation")
@SQLDelete(sql = "UPDATE reservation SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
class Reservation(
    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "payment_method_id_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "payment_method_id", nullable = false)
    var paymentMethod: PaymentMethod,

    @OneToMany(mappedBy = "reservation", fetch = FetchType.LAZY)
    @OrderBy("id ASC")
    var rooms: MutableList<ReservationRoom> = mutableListOf(),

    @Column(name = "name", nullable = false, length = 30)
    var name: String,

    @Column(name = "phone", nullable = false, length = 20)
    var phone: String,

    @Column(name = "people_count", nullable = false)
    var peopleCount: Int = 0,

    @Column(name = "stay_start_at", nullable = false)
    var stayStartAt: LocalDate,

    @Column(name = "stay_end_at", nullable = false)
    var stayEndAt: LocalDate,

    @Column(name = "check_in_at")
    var checkInAt: LocalDateTime? = null,

    @Column(name = "check_out_at")
    var checkOutAt: LocalDateTime? = null,

    @Column(name = "price", nullable = false)
    var price: Int,

    @Column(name = "payment_amount", nullable = false)
    var paymentAmount: Int = 0,

    @Column(name = "refund_amount", nullable = false)
    var refundAmount: Int = 0,

    @Column(name = "broker_fee", nullable = false)
    var brokerFee: Int = 0,

    @Column(name = "note", nullable = false, length = 200)
    var note: String = "",

    @Column(name = "canceled_at")
    var canceledAt: LocalDateTime? = null,

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    var status: ReservationStatus = ReservationStatus.PENDING,
) : Serializable, BaseMustAuditEntity() {
    fun addRoom(room: Room) {
        rooms.add(ReservationRoom(reservation = this, room = room))
    }

    fun removeRoom(room: Room) {
        rooms.removeIf { it.room == room }
    }

    fun updateRooms(rooms: List<Room>) {
        this.rooms.clear()
        this.rooms.addAll(rooms.map { ReservationRoom(reservation = this, room = it) })
    }

    companion object {
        @Serial
        private const val serialVersionUID: Long = 737436737363715200L
    }
}
