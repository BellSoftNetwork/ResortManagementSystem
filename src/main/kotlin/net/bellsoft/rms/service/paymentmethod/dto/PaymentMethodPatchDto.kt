package net.bellsoft.rms.service.paymentmethod.dto

import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodPatchRequest
import org.openapitools.jackson.nullable.JsonNullable

data class PaymentMethodPatchDto(
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val commissionRate: JsonNullable<Double> = JsonNullable.undefined(),
    val requireUnpaidAmountCheck: JsonNullable<Boolean> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: PaymentMethodPatchRequest) = PaymentMethodPatchDto(
            name = dto.name,
            commissionRate = dto.commissionRate,
            requireUnpaidAmountCheck = dto.requireUnpaidAmountCheck,
        )
    }
}
