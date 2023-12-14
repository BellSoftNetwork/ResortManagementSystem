package net.bellsoft.rms.room.dto.service

import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.room.dto.request.RoomPatchRequest
import net.bellsoft.rms.room.type.RoomStatus
import org.openapitools.jackson.nullable.JsonNullable

data class RoomPatchDto(
    val roomGroup: JsonNullable<EntityReferenceDto> = JsonNullable.undefined(),
    val number: JsonNullable<String> = JsonNullable.undefined(),
    val note: JsonNullable<String> = JsonNullable.undefined(),
    val status: JsonNullable<RoomStatus> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: RoomPatchRequest) = RoomPatchDto(
            roomGroup = dto.roomGroup,
            number = dto.number,
            note = dto.note,
            status = dto.status,
        )
    }
}
