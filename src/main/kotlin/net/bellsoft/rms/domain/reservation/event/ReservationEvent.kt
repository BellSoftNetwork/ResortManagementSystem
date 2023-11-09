package net.bellsoft.rms.domain.reservation.event

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
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.user.User
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where

@Entity
@Table(name = "reservation_event")
@SQLDelete(sql = "UPDATE reservation_event SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
class ReservationEvent(
    user: User,
    reservation: Reservation,
    detail: String,

    @Column(name = "type", nullable = false, columnDefinition = "TINYINT")
    var type: ReservationEventType,
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
    @JoinColumn(name = "reservation_id", nullable = false)
    var reservation: Reservation = reservation
        private set

    @Column(name = "detail", nullable = false, length = 1000)
    var detail: String = detail
        private set
}
