package net.bellsoft.rms.service.reservation.dto

import org.openapitools.jackson.nullable.JsonNullable

data class ReservationMethodPatchDto(
    val name: JsonNullable<String>,
    val commissionRate: JsonNullable<Double>,
    val requireUnpaidAmountCheck: JsonNullable<Boolean>,
)
