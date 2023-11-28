package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.method.ReservationMethodStatus

data class ReservationMethodCreateDto(
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean = false,
    val status: ReservationMethodStatus = ReservationMethodStatus.ACTIVE,
)
