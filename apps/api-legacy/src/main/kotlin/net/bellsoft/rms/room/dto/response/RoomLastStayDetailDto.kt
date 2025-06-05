package net.bellsoft.rms.room.dto.response

import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto

data class RoomLastStayDetailDto(
    val room: RoomSummaryDto,
    val lastReservation: ReservationDetailDto?,
)
