package net.bellsoft.rms.service.reservation.dto

import java.time.LocalDate

data class ReservationFilterDto(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val searchText: String?,
)
