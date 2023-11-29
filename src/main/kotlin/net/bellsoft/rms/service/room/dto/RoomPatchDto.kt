package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.controller.v1.room.dto.RoomPatchRequest
import net.bellsoft.rms.domain.room.RoomStatus
import org.openapitools.jackson.nullable.JsonNullable

data class RoomPatchDto(
    val number: JsonNullable<String> = JsonNullable.undefined(),
    val peekPrice: JsonNullable<Int> = JsonNullable.undefined(),
    val offPeekPrice: JsonNullable<Int> = JsonNullable.undefined(),
    val description: JsonNullable<String> = JsonNullable.undefined(),
    val note: JsonNullable<String> = JsonNullable.undefined(),
    val status: JsonNullable<RoomStatus> = JsonNullable.undefined(),
) {
    companion object {
        fun of(dto: RoomPatchRequest) = RoomPatchDto(
            number = dto.number,
            peekPrice = dto.peekPrice,
            offPeekPrice = dto.offPeekPrice,
            description = dto.description,
            note = dto.note,
            status = dto.status,
        )
    }
}
