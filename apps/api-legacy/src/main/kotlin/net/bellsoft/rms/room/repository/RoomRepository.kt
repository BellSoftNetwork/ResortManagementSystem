package net.bellsoft.rms.room.repository

import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.entity.RoomGroup
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.repository.history.RevisionRepository
import org.springframework.stereotype.Repository

@Repository
interface RoomRepository : JpaRepository<Room, Long>, RevisionRepository<Room, Long, Long>, RoomCustomRepository {
    fun findByIdInOrderByNumberAsc(ids: Collection<Long>): List<Room>

    fun findByNumber(name: String): Room?
    fun existsByRoomGroup(roomGroup: RoomGroup): Boolean
}
