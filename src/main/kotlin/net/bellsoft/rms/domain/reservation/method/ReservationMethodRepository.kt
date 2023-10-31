package net.bellsoft.rms.domain.reservation.method

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface ReservationMethodRepository : JpaRepository<ReservationMethod, Long>
