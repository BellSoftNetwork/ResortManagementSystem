package net.bellsoft.rms.domain.reservation.event

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface ReservationEventRepository : JpaRepository<ReservationEvent, Long>
