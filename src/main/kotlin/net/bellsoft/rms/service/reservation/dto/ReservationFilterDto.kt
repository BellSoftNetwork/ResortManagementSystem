package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.controller.v1.reservation.dto.ReservationRequestFilter
import java.time.LocalDate

data class ReservationFilterDto(
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val searchText: String? = null,
) {
    companion object {
        fun of(dto: ReservationRequestFilter) = ReservationFilterDto(
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            searchText = dto.searchText,
        )
    }
}
