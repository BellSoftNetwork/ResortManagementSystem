package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodCreateRequest
import net.bellsoft.rms.domain.reservation.method.ReservationMethodStatus

data class ReservationMethodCreateDto(
    val name: String,
    val commissionRate: Double = 0.0,
    val requireUnpaidAmountCheck: Boolean = false,
    val status: ReservationMethodStatus = ReservationMethodStatus.ACTIVE,
) {
    companion object {
        fun of(dto: ReservationMethodCreateRequest) = ReservationMethodCreateDto(
            name = dto.name,
            commissionRate = dto.commissionRate,
            requireUnpaidAmountCheck = dto.requireUnpaidAmountCheck,
        )
    }
}
