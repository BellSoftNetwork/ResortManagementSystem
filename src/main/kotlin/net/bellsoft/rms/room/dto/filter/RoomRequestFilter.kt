package net.bellsoft.rms.room.dto.filter

import net.bellsoft.rms.room.type.RoomStatus
import java.time.LocalDate

data class RoomRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val status: RoomStatus?,
)
