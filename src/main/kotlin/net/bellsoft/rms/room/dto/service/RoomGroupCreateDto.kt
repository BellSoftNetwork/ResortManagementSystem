package net.bellsoft.rms.room.dto.service

import net.bellsoft.rms.room.dto.request.RoomGroupCreateRequest

data class RoomGroupCreateDto(
    val name: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
) {
    companion object {
        fun of(dto: RoomGroupCreateRequest) = RoomGroupCreateDto(
            name = dto.name,
            peekPrice = dto.peekPrice,
            offPeekPrice = dto.offPeekPrice,
            description = dto.description,
        )
    }
}
