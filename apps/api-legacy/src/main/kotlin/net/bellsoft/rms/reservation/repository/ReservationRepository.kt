package net.bellsoft.rms.reservation.repository

import net.bellsoft.rms.reservation.entity.Reservation
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.repository.history.RevisionRepository
import org.springframework.stereotype.Repository

@Repository
interface ReservationRepository :
    JpaRepository<Reservation, Long>,
    RevisionRepository<Reservation, Long, Long>,
    ReservationCustomRepository
