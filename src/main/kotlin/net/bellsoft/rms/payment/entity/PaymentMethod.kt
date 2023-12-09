package net.bellsoft.rms.payment.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.payment.type.PaymentMethodStatus
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import java.io.Serial
import java.io.Serializable

@Entity
@Table(
    name = "payment_method",
    uniqueConstraints = [
        UniqueConstraint(name = "uc_payment_method_name", columnNames = ["name", "deleted_at"]),
    ],
)
@SQLDelete(sql = "UPDATE payment_method SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
class PaymentMethod(
    @Column(name = "name", nullable = false, length = 20)
    var name: String,

    @Column(name = "commission_rate", nullable = false)
    var commissionRate: Double,

    @Column(name = "required_unpaid_amount_check", nullable = false)
    var requireUnpaidAmountCheck: Boolean = false,

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    var status: PaymentMethodStatus = PaymentMethodStatus.INACTIVE,
) : Serializable, BaseTimeEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = 6340491803785384426L
    }
}
