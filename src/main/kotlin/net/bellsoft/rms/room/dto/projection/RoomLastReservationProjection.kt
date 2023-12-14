package net.bellsoft.rms.room.dto.projection

import com.querydsl.core.annotations.QueryProjection
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.room.entity.Room

data class RoomLastReservationProjection @QueryProjection constructor(
    val room: Room,
    val lastReservation: Reservation?,
)
