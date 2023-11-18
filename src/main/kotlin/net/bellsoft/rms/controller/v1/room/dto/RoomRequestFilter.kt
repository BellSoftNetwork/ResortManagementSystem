package net.bellsoft.rms.controller.v1.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import net.bellsoft.rms.service.room.dto.RoomFilterDto
import java.time.LocalDate

data class RoomRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val status: RoomStatus?,
) {
    fun toDto() = RoomFilterDto(
        stayStartAt = stayStartAt,
        stayEndAt = stayEndAt,
        status = status,
    )
}
