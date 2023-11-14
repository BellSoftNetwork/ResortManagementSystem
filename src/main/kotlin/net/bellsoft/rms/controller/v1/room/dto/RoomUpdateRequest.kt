package net.bellsoft.rms.controller.v1.room.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.room.RoomStatus
import net.bellsoft.rms.service.room.dto.RoomUpdateDto

@Schema(description = "객실 수정 요청 정보")
data class RoomUpdateRequest(
    @Schema(description = "객실 번호", example = "101")
    @field:Size(min = 2, max = 20)
    val number: String? = null,

    @Schema(description = "성수기 가격", example = "100000")
    val peekPrice: Int? = null,

    @Schema(description = "비성수기 가격", example = "80000")
    val offPeekPrice: Int? = null,

    @Schema(description = "객실 설명", example = "와이파이 사용 가능")
    val description: String? = null,

    @Schema(description = "객실 노트", example = "수도 공사 중")
    val note: String? = null,

    @Schema(description = "객실 상태", example = "NORMAL")
    val status: RoomStatus? = null,
) {
    fun toDto() = RoomUpdateDto(
        number = number,
        peekPrice = peekPrice,
        offPeekPrice = offPeekPrice,
        description = description,
        note = note,
        status = status,
    )
}
