package net.bellsoft.rms.room.dto.response

import net.bellsoft.rms.room.type.RoomStatus
import net.bellsoft.rms.user.dto.response.UserSummaryDto
import java.time.LocalDateTime

data class RoomDetailDto(
    val id: Long,
    val number: String,
    // TODO: 추후 Not Null 로 변경 필요
    val roomGroup: RoomGroupSummaryDto?,
    val note: String,
    val status: RoomStatus,
    val createdAt: LocalDateTime,
    val createdBy: UserSummaryDto,
    val updatedAt: LocalDateTime,
    val updatedBy: UserSummaryDto,
)
