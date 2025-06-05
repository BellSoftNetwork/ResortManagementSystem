package net.bellsoft.rms.payment.dto.service

import net.bellsoft.rms.payment.dto.request.PaymentMethodPatchRequest
import org.openapitools.jackson.nullable.JsonNullable

data class PaymentMethodPatchDto(
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val commissionRate: JsonNullable<Double> = JsonNullable.undefined(),
    val requireUnpaidAmountCheck: JsonNullable<Boolean> = JsonNullable.undefined(),
    @get:JvmName("getIsDefaultSelect")
    val isDefaultSelect: JsonNullable<Boolean> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: PaymentMethodPatchRequest) = PaymentMethodPatchDto(
            name = dto.name,
            commissionRate = dto.commissionRate,
            requireUnpaidAmountCheck = dto.requireUnpaidAmountCheck,
            isDefaultSelect = dto.isDefaultSelect,
        )
    }
}
