package net.bellsoft.rms.reservation.dto.response

import net.bellsoft.rms.payment.dto.response.PaymentMethodDetailDto
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import net.bellsoft.rms.room.dto.response.RoomDetailDto
import net.bellsoft.rms.user.dto.response.UserSummaryDto
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationDetailDto(
    val id: Long,
    val paymentMethod: PaymentMethodDetailDto,
    val rooms: List<RoomDetailDto>,
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
    val type: ReservationType,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
