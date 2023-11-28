package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.ReservationStatus
import org.openapitools.jackson.nullable.JsonNullable
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationPatchDto(
    val reservationMethodId: JsonNullable<Long>,
    val roomId: JsonNullable<Long?>,
    val name: JsonNullable<String>,
    val phone: JsonNullable<String>,
    val peopleCount: JsonNullable<Int>,
    val stayStartAt: JsonNullable<LocalDate>,
    val stayEndAt: JsonNullable<LocalDate>,
    val checkInAt: JsonNullable<LocalDateTime?>,
    val checkOutAt: JsonNullable<LocalDateTime?>,
    val price: JsonNullable<Int>,
    val paymentAmount: JsonNullable<Int>,
    val refundAmount: JsonNullable<Int>,
    val brokerFee: JsonNullable<Int>,
    val note: JsonNullable<String>,
    val canceledAt: JsonNullable<LocalDateTime?>,
    val status: JsonNullable<ReservationStatus>,
)
