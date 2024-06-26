package net.bellsoft.rms.reservation.dto.service

import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.reservation.dto.request.ReservationCreateRequest
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationCreateDto(
    val paymentMethod: EntityReferenceDto,
    val rooms: Set<EntityReferenceDto> = emptySet(),
    val name: String,
    val phone: String = "",
    val peopleCount: Int = 0,
    val stayStartAt: LocalDate,
    val stayEndAt: LocalDate,
    val checkInAt: LocalDateTime? = null,
    val checkOutAt: LocalDateTime? = null,
    val price: Int = 0,
    val deposit: Int = 0,
    val paymentAmount: Int = 0,
    val refundAmount: Int = 0,
    val brokerFee: Int = 0,
    val note: String = "",
    val canceledAt: LocalDateTime? = null,
    val status: ReservationStatus = ReservationStatus.PENDING,
    val type: ReservationType = ReservationType.STAY,
) {
    companion object {
        fun of(dto: ReservationCreateRequest) = ReservationCreateDto(
            paymentMethod = dto.paymentMethod,
            rooms = EntityReferenceDto.of(dto.rooms),
            name = dto.name,
            phone = dto.phone,
            peopleCount = dto.peopleCount,
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            checkInAt = dto.checkInAt,
            checkOutAt = dto.checkOutAt,
            price = dto.price,
            deposit = dto.deposit,
            paymentAmount = dto.paymentAmount,
            refundAmount = dto.refundAmount,
            brokerFee = dto.brokerFee,
            note = dto.note,
            canceledAt = dto.canceledAt,
            status = dto.status,
            type = dto.type,
        )
    }
}
