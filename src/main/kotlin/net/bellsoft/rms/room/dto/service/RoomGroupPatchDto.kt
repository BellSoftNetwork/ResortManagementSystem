package net.bellsoft.rms.room.dto.service

import net.bellsoft.rms.room.dto.request.RoomGroupPatchRequest
import org.openapitools.jackson.nullable.JsonNullable

data class RoomGroupPatchDto(
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val peekPrice: JsonNullable<Int> = JsonNullable.undefined(),
    val offPeekPrice: JsonNullable<Int> = JsonNullable.undefined(),
    val description: JsonNullable<String> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: RoomGroupPatchRequest) = RoomGroupPatchDto(
            name = dto.name,
            peekPrice = dto.peekPrice,
            offPeekPrice = dto.offPeekPrice,
            description = dto.description,
        )
    }
}
