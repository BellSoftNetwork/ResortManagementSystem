package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus

data class RoomCreateDto(
    val number: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val note: String,
    val status: RoomStatus,
)
