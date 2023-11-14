package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.event.RoomEvent
import net.bellsoft.rms.domain.room.event.RoomEventType
import net.bellsoft.rms.domain.user.User
import java.time.LocalDateTime

data class RoomEventDto(
    val id: Long,
    val user: User,
    val detail: String,
    val type: RoomEventType,
    val createdAt: LocalDateTime,
) {
    companion object {
        fun of(roomEvent: RoomEvent) = RoomEventDto(
            id = roomEvent.id,
            user = roomEvent.user,
            detail = roomEvent.detail,
            type = roomEvent.type,
            createdAt = roomEvent.createdAt,
        )
    }
}
