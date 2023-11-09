package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.reservation.method.ReservationMethodStatus

data class ReservationMethodCreateDto(
    val name: String,
    val commissionRate: Double,
) {
    fun toEntity() = ReservationMethod(
        name = name,
        commissionRate = commissionRate,
        status = ReservationMethodStatus.ACTIVE,
    )
}
