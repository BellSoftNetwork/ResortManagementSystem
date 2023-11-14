package net.bellsoft.rms.domain.room.event

import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface RoomEventRepository : JpaRepository<RoomEvent, Long> {
    fun findAllByRoomId(roomId: Long, pageable: Pageable): Page<RoomEvent>
}
