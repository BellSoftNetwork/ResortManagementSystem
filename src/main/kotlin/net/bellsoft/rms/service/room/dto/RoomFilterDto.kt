package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.controller.v1.room.dto.RoomRequestFilter
import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDate

data class RoomFilterDto(
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val status: RoomStatus? = null,
) {
    companion object {
        fun of(dto: RoomRequestFilter) = RoomFilterDto(
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            status = dto.status,
        )
    }
}
