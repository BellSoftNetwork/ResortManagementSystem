package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.reservation.method.ReservationMethodStatus

data class ReservationMethodCreateDto(
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean = false,
) {
    fun toEntity() = ReservationMethod(
        name = name,
        commissionRate = commissionRate,
        requireUnpaidAmountCheck = requireUnpaidAmountCheck,
        status = ReservationMethodStatus.ACTIVE,
    )
}
