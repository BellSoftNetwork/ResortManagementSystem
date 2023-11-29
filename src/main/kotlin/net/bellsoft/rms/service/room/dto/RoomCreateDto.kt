package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.controller.v1.room.dto.RoomCreateRequest
import net.bellsoft.rms.domain.room.RoomStatus

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
