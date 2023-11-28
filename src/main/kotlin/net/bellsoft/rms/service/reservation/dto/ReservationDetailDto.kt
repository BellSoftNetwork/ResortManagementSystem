package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.service.room.dto.RoomDetailDto
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationDetailDto(
    val id: Long,
    val reservationMethod: ReservationMethodDetailDto,
    val room: RoomDetailDto?,
    val name: String,
    val phone: String,
    val peopleCount: Int,
    val stayStartAt: LocalDate,
    val stayEndAt: LocalDate,
    val checkInAt: LocalDateTime?,
    val checkOutAt: LocalDateTime?,
    val price: Int,
    val paymentAmount: Int,
    val refundAmount: Int,
    val brokerFee: Int,
    val note: String,
    val canceledAt: LocalDateTime?,
    val status: ReservationStatus,
    val createdAt: LocalDateTime,
    val createdBy: String,
    val updatedAt: LocalDateTime,
    val updatedBy: String,
)
