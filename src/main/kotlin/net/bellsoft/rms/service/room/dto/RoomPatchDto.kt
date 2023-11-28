package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import org.openapitools.jackson.nullable.JsonNullable

data class RoomPatchDto(
    val number: JsonNullable<String>,
    val peekPrice: JsonNullable<Int>,
    val offPeekPrice: JsonNullable<Int>,
    val description: JsonNullable<String>,
    val note: JsonNullable<String>,
    val status: JsonNullable<RoomStatus>,
)
