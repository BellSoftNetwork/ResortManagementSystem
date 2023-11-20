package net.bellsoft.rms.controller.v1.reservation.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.service.reservation.dto.ReservationUpdateDto
import org.hibernate.validator.constraints.Range
import java.time.LocalDate
import java.time.LocalDateTime

@Schema(description = "예약 수정 요청 정보")
data class ReservationUpdateRequest(
    @Schema(description = "예약 수단 ID", example = "1")
    val reservationMethodId: Long? = null,

    @Schema(description = "객실 ID", example = "1")
    val roomId: Long? = null,

    @Schema(description = "예약자명", example = "홍길동")
    @field:Size(min = 2, max = 30)
    val name: String? = null,

    @Schema(description = "예약자 전화번호", example = "010-0000-0000")
    @field:Size(min = 2, max = 20)
    val phone: String? = null,

    @Schema(description = "예약 인원", example = "4")
    @field:Range(min = 0, max = 1000)
    val peopleCount: Int? = null,

    @Schema(description = "입실일", example = "2023-11-15")
    val stayStartAt: LocalDate? = null,

    @Schema(description = "퇴실일", example = "2023-11-16")
    val stayEndAt: LocalDate? = null,

    @Schema(description = "체크인 시각", example = "2023-11-15 17:00:00")
    val checkInAt: LocalDateTime? = null,

    @Schema(description = "체크아웃 시각", example = "2023-11-16 10:00:00")
    val checkOutAt: LocalDateTime? = null,

    @Schema(description = "예약 가격", example = "100000")
    @field:Range(min = 0, max = 100000000)
    val price: Int? = null,

    @Schema(description = "현재 총 지불 금액", example = "80000")
    @field:Range(min = 0, max = 100000000)
    val paymentAmount: Int? = null,

    @Schema(description = "환불 금액", example = "0")
    @field:Range(min = 0, max = 100000000)
    val refundAmount: Int? = null,

    @Schema(description = "플랫폼 수수료", example = "5000")
    @field:Range(min = 0, max = 100000000)
    val brokerFee: Int? = null,

    @Schema(description = "메모", example = "밤 늦게 입실 예정")
    @field:Size(min = 0, max = 200)
    val note: String? = null,

    @Schema(description = "예약 취소 시각", example = "2023-11-15 20:00:00")
    val canceledAt: LocalDateTime? = null,

    @Schema(description = "예약 상태", example = "NORMAL")
    val status: ReservationStatus? = null,
) {
    fun toDto() = ReservationUpdateDto(
        reservationMethodId = reservationMethodId,
        roomId = roomId,
        name = name,
        phone = phone,
        peopleCount = peopleCount,
        stayStartAt = stayStartAt,
        stayEndAt = stayEndAt,
        checkInAt = checkInAt,
        checkOutAt = checkOutAt,
        price = price,
        paymentAmount = paymentAmount,
        refundAmount = refundAmount,
        brokerFee = brokerFee,
        note = note,
        canceledAt = canceledAt,
        status = status,
    )
}
