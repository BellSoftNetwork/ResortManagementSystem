package net.bellsoft.rms.controller.v1.reservation.dto

import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto

data class ReservationMethodCreateRequest(
    val name: String,
    val commissionRate: Double,
) {
    fun toDto() = ReservationMethodCreateDto(
        name = name,
        commissionRate = commissionRate,
    )
}
