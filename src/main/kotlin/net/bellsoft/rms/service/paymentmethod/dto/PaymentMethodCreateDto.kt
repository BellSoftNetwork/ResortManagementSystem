package net.bellsoft.rms.service.paymentmethod.dto

import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodCreateRequest
import net.bellsoft.rms.domain.paymentmethod.PaymentMethodStatus

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
