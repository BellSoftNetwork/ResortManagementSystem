package net.bellsoft.rms.domain.room

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface RoomRepository : JpaRepository<Room, Long> {
    fun findByNumber(name: String): Room?
}
