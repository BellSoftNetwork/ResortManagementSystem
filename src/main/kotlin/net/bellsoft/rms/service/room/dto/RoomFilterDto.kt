package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDate

data class RoomFilterDto(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val status: RoomStatus?,
)
