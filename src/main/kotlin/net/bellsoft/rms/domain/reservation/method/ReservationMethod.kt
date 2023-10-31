package net.bellsoft.rms.domain.reservation.method

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.Table
import net.bellsoft.rms.domain.base.BaseTime
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where

@Entity
@Table(name = "reservation_method")
@SQLDelete(sql = "UPDATE reservation_method SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at IS NULL")
class ReservationMethod(
    name: String,
    commissionRate: Double,
    status: ReservationMethodStatus = ReservationMethodStatus.INACTIVE,
) : BaseTime() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set

    @Column(name = "name", nullable = false, length = 20)
    var name: String = name

    @Column(name = "commission_rate", nullable = false)
    var commissionRate: Double = commissionRate
        private set

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    var status: ReservationMethodStatus = status
        private set
}
