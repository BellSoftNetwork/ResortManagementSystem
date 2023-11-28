package net.bellsoft.rms.controller.v1.reservation.dto

import java.time.LocalDate

data class ReservationRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val searchText: String?,
)
