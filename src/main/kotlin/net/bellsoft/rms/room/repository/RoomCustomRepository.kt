package net.bellsoft.rms.room.repository

import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.entity.Room
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface RoomCustomRepository {
    fun getFilteredRooms(pageable: Pageable, filter: RoomFilterDto): Page<Room>
}
