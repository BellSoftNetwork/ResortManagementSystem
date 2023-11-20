package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import java.time.LocalDateTime

data class ReservationMethodDto(
    val id: Long,
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime,
) {
    companion object {
        fun of(reservationMethod: ReservationMethod) = ReservationMethodDto(
            id = reservationMethod.id,
            name = reservationMethod.name,
            commissionRate = reservationMethod.commissionRate,
            requireUnpaidAmountCheck = reservationMethod.requireUnpaidAmountCheck,
            createdAt = reservationMethod.createdAt,
            updatedAt = reservationMethod.updatedAt,
        )
    }
}
