package net.bellsoft.rms.domain.room

import net.bellsoft.rms.service.room.dto.RoomFilterDto
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface RoomCustomRepository {
    fun getFilteredRooms(pageable: Pageable, filter: RoomFilterDto): Page<Room>
}
