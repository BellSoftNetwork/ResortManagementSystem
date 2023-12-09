package net.bellsoft.rms.room.dto.service

import net.bellsoft.rms.room.dto.request.RoomCreateRequest
import net.bellsoft.rms.room.type.RoomStatus

data class RoomCreateDto(
    val number: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val note: String,
    val status: RoomStatus,
) {
    companion object {
        fun of(dto: RoomCreateRequest) = RoomCreateDto(
            number = dto.number,
            peekPrice = dto.peekPrice,
            offPeekPrice = dto.offPeekPrice,
            description = dto.description,
            note = dto.note,
            status = dto.status,
        )
    }
}
