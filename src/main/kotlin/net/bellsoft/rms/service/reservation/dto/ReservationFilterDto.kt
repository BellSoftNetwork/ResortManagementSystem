package net.bellsoft.rms.service.reservation.dto

import java.time.LocalDate

data class ReservationFilterDto(
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val searchText: String? = null,
)
