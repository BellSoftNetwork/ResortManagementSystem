package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDateTime

data class RoomDto(
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
) {
    companion object {
        fun of(room: Room) = RoomDto(
            id = room.id,
            number = room.number,
            peekPrice = room.peekPrice,
            offPeekPrice = room.offPeekPrice,
            description = room.description,
            note = room.note,
            status = room.status,
            createdAt = room.createdAt,
            createdBy = room.createdBy.email,
            updatedAt = room.updatedAt,
            updatedBy = room.updatedBy.email,
        )
    }
}
