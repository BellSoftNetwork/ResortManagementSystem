package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDateTime

data class RoomDetailDto(
    val id: Long,
    val number: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val note: String,
    val status: RoomStatus,
    val createdAt: LocalDateTime,
    val createdBy: String,
    val updatedAt: LocalDateTime,
    val updatedBy: String,
)
