package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDate

data class RoomFilterDto(
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val status: RoomStatus? = null,
)
