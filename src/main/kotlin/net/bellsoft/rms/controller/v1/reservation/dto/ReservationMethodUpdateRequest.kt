package net.bellsoft.rms.controller.v1.reservation.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.service.reservation.dto.ReservationMethodUpdateDto
import org.hibernate.validator.constraints.Range

@Schema(description = "예약 수단 수정 요청 정보")
data class ReservationMethodUpdateRequest(
    @Schema(description = "예약 수단명", example = "네이버")
    @field:Size(min = 2, max = 20)
    val name: String? = null,

    @Schema(description = "수수료율", example = "0.2")
    @field:Range(min = 0, max = 1)
    val commissionRate: Double? = null,
) {
    fun toDto() = ReservationMethodUpdateDto(
        name = name,
        commissionRate = commissionRate,
    )
}
