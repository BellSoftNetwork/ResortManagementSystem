package net.bellsoft.rms.room.dto.request

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import org.hibernate.validator.constraints.Range

@Schema(description = "객실 그룹 생성 요청 정보")
data class RoomGroupCreateRequest(
    @Schema(description = "객실 그룹명", example = "20평형")
    @field:Size(min = 2, max = 20)
    val name: String,

    @Schema(description = "성수기 가격", example = "100000")
    @field:Range(min = 0, max = 100000000)
    val peekPrice: Int,

    @Schema(description = "비성수기 가격", example = "80000")
    @field:Range(min = 0, max = 100000000)
    val offPeekPrice: Int,

    @Schema(description = "객실 설명", example = "와이파이 사용 가능")
    @field:Size(min = 0, max = 200)
    val description: String,
)
