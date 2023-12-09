package net.bellsoft.rms.reservation.repository

import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.entity.Reservation
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface ReservationCustomRepository {
    fun getFilteredReservations(pageable: Pageable, filter: ReservationFilterDto): Page<Reservation>
}
