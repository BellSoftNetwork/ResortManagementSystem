package net.bellsoft.rms.room.dto.response

import net.bellsoft.rms.room.type.RoomStatus
import net.bellsoft.rms.user.dto.response.UserSummaryDto
import java.time.LocalDateTime

data class RoomSummaryDto(
    val id: Long,
    val number: String,
    val note: String,
    val status: RoomStatus,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
