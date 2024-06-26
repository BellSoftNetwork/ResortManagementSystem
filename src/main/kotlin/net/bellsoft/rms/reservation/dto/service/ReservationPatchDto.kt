package net.bellsoft.rms.reservation.dto.service

import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.common.util.convert
import net.bellsoft.rms.reservation.dto.request.ReservationPatchRequest
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import org.openapitools.jackson.nullable.JsonNullable
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationPatchDto(
    val paymentMethod: JsonNullable<EntityReferenceDto> = JsonNullable.undefined(),
    val rooms: JsonNullable<Set<EntityReferenceDto>> = JsonNullable.undefined(),
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val phone: JsonNullable<String> = JsonNullable.undefined(),
    val peopleCount: JsonNullable<Int> = JsonNullable.undefined(),
    val stayStartAt: JsonNullable<LocalDate> = JsonNullable.undefined(),
    val stayEndAt: JsonNullable<LocalDate> = JsonNullable.undefined(),
    val checkInAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),
    val checkOutAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),
    val price: JsonNullable<Int> = JsonNullable.undefined(),
    val deposit: JsonNullable<Int> = JsonNullable.undefined(),
    val paymentAmount: JsonNullable<Int> = JsonNullable.undefined(),
    val refundAmount: JsonNullable<Int> = JsonNullable.undefined(),
    val brokerFee: JsonNullable<Int> = JsonNullable.undefined(),
    val note: JsonNullable<String> = JsonNullable.undefined(),
    val canceledAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),
    val status: JsonNullable<ReservationStatus> = JsonNullable.undefined(),
    val type: JsonNullable<ReservationType> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: ReservationPatchRequest) = ReservationPatchDto(
            paymentMethod = dto.paymentMethod,
            rooms = dto.rooms.convert(EntityReferenceDto::of),
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
