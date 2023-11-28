package net.bellsoft.rms.controller.v1.room.dto

import net.bellsoft.rms.domain.room.RoomStatus
import java.time.LocalDate

data class RoomRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val status: RoomStatus?,
)
