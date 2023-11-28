package net.bellsoft.rms.controller.v1.reservation.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.reservation.ReservationStatus
import org.hibernate.validator.constraints.Range
import java.time.LocalDate
import java.time.LocalDateTime

@Schema(description = "예약 생성 요청 정보")
data class ReservationCreateRequest(
    @Schema(description = "예약 수단 ID", example = "1")
    val reservationMethodId: Long,

    @Schema(description = "객실 ID", example = "1")
    val roomId: Long? = null,

    @Schema(description = "예약자명", example = "홍길동")
    @field:Size(min = 2, max = 30)
    val name: String,

    @Schema(description = "예약자 전화번호", example = "010-0000-0000")
    @field:Size(min = 2, max = 20)
    val phone: String = "",

    @Schema(description = "예약 인원", example = "4")
    @field:Range(min = 0, max = 1000)
    val peopleCount: Int = 0,

    @Schema(description = "입실일", example = "2023-11-15")
    val stayStartAt: LocalDate,

    @Schema(description = "퇴실일", example = "2023-11-16")
    val stayEndAt: LocalDate,

    @Schema(description = "체크인 시각", example = "2023-11-15 17:00:00")
    val checkInAt: LocalDateTime? = null,

    @Schema(description = "체크아웃 시각", example = "2023-11-16 10:00:00")
    val checkOutAt: LocalDateTime? = null,

    @Schema(description = "예약 가격", example = "100000")
    @field:Range(min = 0, max = 100000000)
    val price: Int = 0,

    @Schema(description = "현재 총 지불 금액", example = "80000")
    @field:Range(min = 0, max = 100000000)
    val paymentAmount: Int = 0,

    @Schema(description = "환불 금액", example = "0")
    @field:Range(min = 0, max = 100000000)
    val refundAmount: Int = 0,

    @Schema(description = "플랫폼 수수료", example = "5000")
    @field:Range(min = 0, max = 100000000)
    val brokerFee: Int = 0,

    @Schema(description = "메모", example = "밤 늦게 입실 예정")
    @field:Size(min = 0, max = 200)
    val note: String = "",

    @Schema(description = "예약 취소 시각", example = "2023-11-15 20:00:00")
    val canceledAt: LocalDateTime? = null,

    @Schema(description = "예약 상태", example = "NORMAL")
    val status: ReservationStatus = ReservationStatus.PENDING,
)
