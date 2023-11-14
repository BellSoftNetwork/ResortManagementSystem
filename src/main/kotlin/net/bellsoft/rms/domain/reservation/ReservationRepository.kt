package net.bellsoft.rms.domain.reservation

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.repository.history.RevisionRepository
import org.springframework.stereotype.Repository

@Repository
interface ReservationRepository : JpaRepository<Reservation, Long>, RevisionRepository<Reservation, Long, Long>
