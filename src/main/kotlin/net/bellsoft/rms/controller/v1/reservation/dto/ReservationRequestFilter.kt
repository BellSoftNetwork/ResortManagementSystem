package net.bellsoft.rms.controller.v1.reservation.dto

import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import java.time.LocalDate

data class ReservationRequestFilter(
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val searchText: String?,
) {
    fun toDto() = ReservationFilterDto(
        stayStartAt = stayStartAt,
        stayEndAt = stayEndAt,
        searchText = searchText,
    )
}
