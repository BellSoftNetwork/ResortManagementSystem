package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.controller.v1.reservation.dto.ReservationCreateRequest
import net.bellsoft.rms.domain.reservation.ReservationStatus
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationCreateDto(
    val reservationMethodId: Long,
    val roomId: Long? = null,
    val name: String,
    val phone: String = "",
    val peopleCount: Int = 0,
    val stayStartAt: LocalDate,
    val stayEndAt: LocalDate,
    val checkInAt: LocalDateTime? = null,
    val checkOutAt: LocalDateTime? = null,
    val price: Int = 0,
    val paymentAmount: Int = 0,
    val refundAmount: Int = 0,
    val brokerFee: Int = 0,
    val note: String = "",
    val canceledAt: LocalDateTime? = null,
    val status: ReservationStatus = ReservationStatus.PENDING,
) {
    companion object {
        fun of(dto: ReservationCreateRequest) = ReservationCreateDto(
            reservationMethodId = dto.reservationMethodId,
            roomId = dto.roomId,
            name = dto.name,
            phone = dto.phone,
            peopleCount = dto.peopleCount,
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            checkInAt = dto.checkInAt,
            checkOutAt = dto.checkOutAt,
            price = dto.price,
            paymentAmount = dto.paymentAmount,
            refundAmount = dto.refundAmount,
            brokerFee = dto.brokerFee,
            note = dto.note,
            canceledAt = dto.canceledAt,
            status = dto.status,
        )
    }
}
