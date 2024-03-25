package net.bellsoft.rms.reservation.entity

import jakarta.persistence.CascadeType
import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.FetchType
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.OneToMany
import jakarta.persistence.OrderBy
import jakarta.persistence.Table
import net.bellsoft.rms.common.entity.BaseMustAuditEntity
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.payment.entity.PaymentMethod
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import net.bellsoft.rms.room.entity.Room
import org.hibernate.annotations.Comment
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
@Comment("예약 정보")
class Reservation(
    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "payment_method_id_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "payment_method_id", nullable = false)
    @Comment("결제 수단")
    var paymentMethod: PaymentMethod,

    @OneToMany(mappedBy = "reservation", cascade = [CascadeType.ALL], orphanRemoval = true)
    @OrderBy("id ASC")
    @Comment("배정 객실")
    var rooms: MutableList<ReservationRoom> = mutableListOf(),

    @Column(name = "name", nullable = false, length = 30)
    @Comment("예약자명")
    var name: String,

    @Column(name = "phone", nullable = false, length = 20)
    @Comment("예약자 전화번호")
    var phone: String,

    @Column(name = "people_count", nullable = false)
    @Comment("예약 인원")
    var peopleCount: Int = 0,

    @Column(name = "stay_start_at", nullable = false)
    @Comment("입실일")
    var stayStartAt: LocalDate,

    @Column(name = "stay_end_at", nullable = false)
    @Comment("퇴실일")
    var stayEndAt: LocalDate,

    @Column(name = "check_in_at")
    @Comment("체크인 시각")
    var checkInAt: LocalDateTime? = null,

    @Column(name = "check_out_at")
    @Comment("체크아웃 시각")
    var checkOutAt: LocalDateTime? = null,

    @Column(name = "price", nullable = false)
    @Comment("예약 가격")
    var price: Int,

    @Column(name = "payment_amount", nullable = false)
    @Comment("현재 총 지불 금액")
    var paymentAmount: Int = 0,

    @Column(name = "refund_amount", nullable = false)
    @Comment("환불 금액")
    var refundAmount: Int = 0,

    @Column(name = "broker_fee", nullable = false)
    @Comment("플랫폼 수수료")
    var brokerFee: Int = 0,

    @Column(name = "note", nullable = false, length = 200)
    @Comment("메모")
    var note: String = "",

    @Column(name = "canceled_at")
    @Comment("예약 취소 시각")
    var canceledAt: LocalDateTime? = null,

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    @Comment("예약 상태")
    var status: ReservationStatus = ReservationStatus.PENDING,

    @Column(name = "type", nullable = false, columnDefinition = "TINYINT")
    @Comment("예약 구분")
    var type: ReservationType = ReservationType.STAY,
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
