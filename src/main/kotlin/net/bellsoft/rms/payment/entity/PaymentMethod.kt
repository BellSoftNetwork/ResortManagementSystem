package net.bellsoft.rms.payment.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.payment.type.PaymentMethodStatus
import org.hibernate.annotations.Comment
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
@Comment("결제 수단")
class PaymentMethod(
    @Column(name = "name", nullable = false, length = 20)
    @Comment("결제 수단명")
    var name: String,

    @Column(name = "commission_rate", nullable = false)
    @Comment("수수료율")
    var commissionRate: Double,

    @Column(name = "required_unpaid_amount_check", nullable = false)
    @Comment("미수금 금액 알림")
    var requireUnpaidAmountCheck: Boolean = false,

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    @Comment("결제수단 상태")
    var status: PaymentMethodStatus = PaymentMethodStatus.INACTIVE,
) : Serializable, BaseTimeEntity() {
    companion object {
        @Serial
        private const val serialVersionUID: Long = 6340491803785384426L
    }
}
