package net.bellsoft.rms.reservation.dto.request

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.common.dto.request.EntityReferenceRequest
import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import org.hibernate.validator.constraints.Range
import org.openapitools.jackson.nullable.JsonNullable
import java.time.LocalDate
import java.time.LocalDateTime

@Schema(description = "예약 수정 요청 정보")
data class ReservationPatchRequest(
    @Schema(description = "결제 수단 ID", example = "1")
    val paymentMethod: JsonNullable<EntityReferenceDto> = JsonNullable.undefined(),

    @Schema(description = "객실 ID", example = "1")
    val rooms: JsonNullable<Set<EntityReferenceRequest>> = JsonNullable.undefined(),

    @Schema(description = "예약자명", example = "홍길동")
    @field:Size(min = 2, max = 30)
    val name: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "예약자 전화번호", example = "010-0000-0000")
    @field:Size(max = 20)
    val phone: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "예약 인원", example = "4")
    @field:Range(min = 0, max = 1000)
    val peopleCount: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "입실일", example = "2023-11-15")
    val stayStartAt: JsonNullable<LocalDate> = JsonNullable.undefined(),

    @Schema(description = "퇴실일", example = "2023-11-16")
    val stayEndAt: JsonNullable<LocalDate> = JsonNullable.undefined(),

    @Schema(description = "체크인 시각", example = "2023-11-15 17:00:00")
    val checkInAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),

    @Schema(description = "체크아웃 시각", example = "2023-11-16 10:00:00")
    val checkOutAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),

    @Schema(description = "예약 가격", example = "100000")
    @field:Range(min = 0, max = 100000000)
    val price: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "현재 총 지불 금액", example = "80000")
    @field:Range(min = 0, max = 100000000)
    val paymentAmount: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "환불 금액", example = "0")
    @field:Range(min = 0, max = 100000000)
    val refundAmount: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "플랫폼 수수료", example = "5000")
    @field:Range(min = 0, max = 100000000)
    val brokerFee: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "메모", example = "밤 늦게 입실 예정")
    @field:Size(min = 0, max = 200)
    val note: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "예약 취소 시각", example = "2023-11-15 20:00:00")
    val canceledAt: JsonNullable<LocalDateTime?> = JsonNullable.undefined(),

    @Schema(description = "예약 상태", example = "NORMAL")
    val status: JsonNullable<ReservationStatus> = JsonNullable.undefined(),

    @Schema(description = "예약 구분", example = "STAY")
    val type: JsonNullable<ReservationType> = JsonNullable.undefined(),
)
