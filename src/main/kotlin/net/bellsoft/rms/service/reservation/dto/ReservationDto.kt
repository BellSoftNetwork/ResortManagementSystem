package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.service.room.dto.RoomDto
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationDto(
    val id: Long,
    val reservationMethod: ReservationMethodDto,
    val room: RoomDto? = null,
    val name: String,
    val phone: String,
    val peopleCount: Int,
    val stayStartAt: LocalDate,
    val stayEndAt: LocalDate,
    val checkInAt: LocalDateTime? = null,
    val checkOutAt: LocalDateTime? = null,
    val price: Int,
    val paymentAmount: Int = 0,
    val refundAmount: Int = 0,
    val reservationFee: Int = 0,
    val brokerFee: Int = 0,
    val note: String = "",
    val canceledAt: LocalDateTime? = null,
    val status: ReservationStatus = ReservationStatus.PENDING,
    val createdAt: LocalDateTime,
    val createdBy: String,
    val updatedAt: LocalDateTime,
    val updatedBy: String,
) {
    companion object {
        fun of(reservation: Reservation) = ReservationDto(
            id = reservation.id,
            reservationMethod = ReservationMethodDto.of(reservation.reservationMethod),
            room = reservation.room?.let { RoomDto.of(it) },
            name = reservation.name,
            phone = reservation.phone,
            peopleCount = reservation.peopleCount,
            stayStartAt = reservation.stayStartAt,
            stayEndAt = reservation.stayEndAt,
            checkInAt = reservation.checkInAt,
            checkOutAt = reservation.checkOutAt,
            price = reservation.price,
            paymentAmount = reservation.paymentAmount,
            refundAmount = reservation.refundAmount,
            reservationFee = reservation.reservationFee,
            brokerFee = reservation.brokerFee,
            note = reservation.note,
            canceledAt = reservation.canceledAt,
            status = reservation.status,
            createdAt = reservation.createdAt,
            createdBy = reservation.createdBy.email,
            updatedAt = reservation.updatedAt,
            updatedBy = reservation.updatedBy.email,
        )
    }
}
