package net.bellsoft.rms.room.repository

import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.projection.RoomLastReservationProjection
import net.bellsoft.rms.room.entity.RoomGroup

interface RoomGroupCustomRepository {
    fun getFilteredRoomsOrderByLastStayAt(
        roomGroup: RoomGroup,
        filter: RoomFilterDto,
    ): List<RoomLastReservationProjection>
}
