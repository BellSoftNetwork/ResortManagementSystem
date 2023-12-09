package net.bellsoft.rms.reservation.dto.filter

import net.bellsoft.rms.reservation.type.ReservationStatus
import java.time.LocalDate

data class ReservationFilterDto(
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val searchText: String? = null,
    val status: ReservationStatus? = null,
) {
    companion object {
        fun of(dto: ReservationRequestFilter) = ReservationFilterDto(
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            searchText = dto.searchText,
            status = dto.status,
        )
    }
}
