package net.bellsoft.rms.domain.reservation

import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface ReservationCustomRepository {
    fun getFilteredReservations(pageable: Pageable, filter: ReservationFilterDto): Page<Reservation>
}
