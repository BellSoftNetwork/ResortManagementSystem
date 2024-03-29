package net.bellsoft.rms.reservation.dto.filter

import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import java.time.LocalDate

data class ReservationRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val searchText: String?,
    val status: ReservationStatus?,
    val type: ReservationType?,
)
