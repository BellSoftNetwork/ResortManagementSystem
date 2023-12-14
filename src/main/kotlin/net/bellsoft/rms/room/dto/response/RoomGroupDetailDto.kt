package net.bellsoft.rms.room.dto.response

import net.bellsoft.rms.user.dto.response.UserSummaryDto
import java.time.LocalDateTime

data class RoomGroupDetailDto(
    val id: Long,
    val name: String,
    val peekPrice: Int,
    val offPeekPrice: Int,
    val description: String,
    val rooms: List<RoomLastStayDetailDto>,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
