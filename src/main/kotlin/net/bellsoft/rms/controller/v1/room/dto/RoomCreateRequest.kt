package net.bellsoft.rms.controller.v1.room.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.room.RoomStatus
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import org.hibernate.validator.constraints.Range

@Schema(description = "객실 생성 요청 정보")
data class RoomCreateRequest(
    @Schema(description = "객실 번호", example = "101")
    @field:Size(min = 2, max = 20)
    val number: String,

    @Schema(description = "성수기 가격", example = "100000")
    @field:Range(min = 0, max = 100000000)
    val peekPrice: Int,

    @Schema(description = "비성수기 가격", example = "80000")
    @field:Range(min = 0, max = 100000000)
    val offPeekPrice: Int,

    @Schema(description = "객실 설명", example = "와이파이 사용 가능")
    @field:Size(min = 0, max = 200)
    val description: String,

    @Schema(description = "객실 노트", example = "수도 공사 중")
    @field:Size(min = 0, max = 200)
    val note: String,

    @Schema(description = "객실 상태", example = "NORMAL")
    val status: RoomStatus,
) {
    fun toDto() = RoomCreateDto(
        number = number,
        peekPrice = peekPrice,
        offPeekPrice = offPeekPrice,
        description = description,
        note = note,
        status = status,
    )
}
