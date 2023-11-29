package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodPatchRequest
import org.openapitools.jackson.nullable.JsonNullable

data class ReservationMethodPatchDto(
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val commissionRate: JsonNullable<Double> = JsonNullable.undefined(),
    val requireUnpaidAmountCheck: JsonNullable<Boolean> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: ReservationMethodPatchRequest) = ReservationMethodPatchDto(
            name = dto.name,
            commissionRate = dto.commissionRate,
            requireUnpaidAmountCheck = dto.requireUnpaidAmountCheck,
        )
    }
}
