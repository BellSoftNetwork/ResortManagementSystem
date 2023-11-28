package net.bellsoft.rms.service.reservation.dto

import java.time.LocalDateTime

data class ReservationMethodDetailDto(
    val id: Long,
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime,
)
