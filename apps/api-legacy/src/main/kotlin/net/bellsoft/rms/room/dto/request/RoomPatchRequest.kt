package net.bellsoft.rms.room.dto.request

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.room.type.RoomStatus
import org.openapitools.jackson.nullable.JsonNullable

@Schema(description = "객실 수정 요청 정보")
data class RoomPatchRequest(
    @Schema(description = "객실 그룹 ID", example = "1")
    val roomGroup: JsonNullable<EntityReferenceDto> = JsonNullable.undefined(),

    @Schema(description = "객실 번호", example = "101")
    @field:Size(min = 2, max = 20)
    val number: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "객실 노트", example = "수도 공사 중")
    val note: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "객실 상태", example = "NORMAL")
    val status: JsonNullable<RoomStatus> = JsonNullable.undefined(),
)
