package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomStatus

data class RoomCreateDto(
    val number: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val note: String,
    val status: RoomStatus,
) {
    fun toEntity() = Room(
        number = number,
        peekPrice = peekPrice,
        offPeekPrice = offPeekPrice,
        description = description,
        note = note,
        status = status,
    )
}
