package net.bellsoft.rms.room.dto.response

import net.bellsoft.rms.user.dto.response.UserSummaryDto
import java.time.LocalDateTime

data class RoomGroupSummaryDto(
    val id: Long,
    val name: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
