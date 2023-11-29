package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import net.bellsoft.rms.service.auth.dto.UserSummaryDto
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
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
