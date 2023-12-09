package net.bellsoft.rms.payment.dto.service

import net.bellsoft.rms.payment.dto.request.PaymentMethodCreateRequest
import net.bellsoft.rms.payment.type.PaymentMethodStatus

data class PaymentMethodCreateDto(
    val name: String,
    val commissionRate: Double = 0.0,
    val requireUnpaidAmountCheck: Boolean = false,
    val status: PaymentMethodStatus = PaymentMethodStatus.ACTIVE,
) {
    companion object {
        fun of(dto: PaymentMethodCreateRequest) = PaymentMethodCreateDto(
            name = dto.name,
            commissionRate = dto.commissionRate,
            requireUnpaidAmountCheck = dto.requireUnpaidAmountCheck,
        )
    }
}
