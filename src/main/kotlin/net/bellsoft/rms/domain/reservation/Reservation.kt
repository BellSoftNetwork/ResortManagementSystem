package net.bellsoft.rms.domain.reservation

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.Table
import net.bellsoft.rms.domain.base.BaseMustAudit
import net.bellsoft.rms.domain.base.BaseTime
import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.room.Room
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.hibernate.envers.AuditTable
import org.hibernate.envers.Audited
import org.hibernate.envers.RelationTargetAuditMode
import java.io.Serializable
import java.time.LocalDate
import java.time.LocalDateTime

@Entity
@Audited(withModifiedFlag = true)
@AuditTable("reservation_history")
@Table(name = "reservation")
@SQLDelete(sql = "UPDATE reservation SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
class Reservation(
    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "reservation_method_id_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "reservation_method_id", nullable = false)
    var reservationMethod: ReservationMethod,

    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "room_id_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "room_id")
    var room: Room? = null,

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

    @Column(name = "reservation_fee", nullable = false)
    var reservationFee: Int = 0,

    @Column(name = "broker_fee", nullable = false)
    var brokerFee: Int = 0,

    @Column(name = "note", nullable = false, length = 200)
    var note: String = "",

    @Column(name = "canceled_at")
    var canceledAt: LocalDateTime? = null,

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    var status: ReservationStatus = ReservationStatus.PENDING,
) : Serializable, BaseMustAudit() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set
}
