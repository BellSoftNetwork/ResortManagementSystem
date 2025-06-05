package net.bellsoft.rms.room.dto.request

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import org.openapitools.jackson.nullable.JsonNullable

@Schema(description = "객실 그룹 수정 요청 정보")
data class RoomGroupPatchRequest(
    @Schema(description = "객실 그룹명", example = "20평형")
    @field:Size(min = 2, max = 20)
    val name: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "성수기 가격", example = "100000")
    val peekPrice: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "비성수기 가격", example = "80000")
    val offPeekPrice: JsonNullable<Int> = JsonNullable.undefined(),

    @Schema(description = "객실 설명", example = "와이파이 사용 가능")
    val description: JsonNullable<String> = JsonNullable.undefined(),
)
