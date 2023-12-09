package net.bellsoft.rms.room.dto.filter

import net.bellsoft.rms.room.type.RoomStatus
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
