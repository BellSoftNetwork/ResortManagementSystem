package net.bellsoft.rms.domain.reservation.method

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.domain.base.BaseTime
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import java.io.Serializable

@Entity
@Table(
    name = "reservation_method",
    uniqueConstraints = [
        UniqueConstraint(name = "uc_reservation_method_name", columnNames = ["name", "deleted_at"]),
    ],
)
@SQLDelete(sql = "UPDATE reservation_method SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
class ReservationMethod(
    @Column(name = "name", nullable = false, length = 20)
    var name: String,

    @Column(name = "commission_rate", nullable = false)
    var commissionRate: Double,

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    var status: ReservationMethodStatus = ReservationMethodStatus.INACTIVE,
) : Serializable, BaseTime() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set
}
