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
import net.bellsoft.rms.domain.base.BaseTime
import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.user.User
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import java.time.LocalDateTime
import java.util.*

@Entity
@Table(name = "reservation")
@SQLDelete(sql = "UPDATE reservation SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
class Reservation(
    user: User,
    reservationMethod: ReservationMethod,
    room: Room? = null,
    name: String,
    phone: String? = null,
    peopleCount: Int? = null,
    stayStartAt: Date,
    stayEndAt: Date,
    checkInAt: LocalDateTime? = null,
    checkOutAt: LocalDateTime? = null,
    price: Int,
    paymentAmount: Int = 0,
    refundAmount: Int = 0,
    reservationFee: Int = 0,
    brokerFee: Int = 0,
    note: String = "",
    canceledAt: LocalDateTime? = null,

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    var status: ReservationStatus = ReservationStatus.PENDING,
) : BaseTime() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "user_id", nullable = false)
    var user: User = user
        private set

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "reservation_method_id", nullable = false)
    var reservationMethod: ReservationMethod = reservationMethod
        private set

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "room_id")
    var room: Room? = room
        private set

    @Column(name = "name", nullable = false, length = 30)
    var name: String = name

    @Column(name = "phone", length = 20)
    var phone: String? = phone
        private set

    @Column(name = "people_count")
    var peopleCount: Int? = peopleCount
        private set

    @Column(name = "stay_start_at", nullable = false)
    var stayStartAt: Date = stayStartAt
        private set

    @Column(name = "stay_end_at", nullable = false)
    var stayEndAt: Date = stayEndAt
        private set

    @Column(name = "check_in_at")
    var checkInAt: LocalDateTime? = checkInAt
        private set

    @Column(name = "check_out_at")
    var checkOutAt: LocalDateTime? = checkOutAt
        private set

    @Column(name = "price", nullable = false)
    var price: Int = price
        private set

    @Column(name = "payment_amount", nullable = false)
    var paymentAmount: Int = paymentAmount
        private set

    @Column(name = "refund_amount", nullable = false)
    var refundAmount: Int = refundAmount
        private set

    @Column(name = "reservation_fee", nullable = false)
    var reservationFee: Int = reservationFee
        private set

    @Column(name = "broker_fee", nullable = false)
    var brokerFee: Int = brokerFee
        private set

    @Column(name = "note", nullable = false, length = 200)
    var note: String = note
        private set

    @Column(name = "canceled_at")
    var canceledAt: LocalDateTime? = canceledAt
        private set
}
