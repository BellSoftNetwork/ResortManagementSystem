package net.bellsoft.rms.domain.room

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.repository.history.RevisionRepository
import org.springframework.stereotype.Repository

@Repository
interface RoomRepository : JpaRepository<Room, Long>, RevisionRepository<Room, Long, Long>, RoomCustomRepository {
    fun findByIdInOrderByNumberAsc(ids: Collection<Long>): List<Room>

    fun findByNumber(name: String): Room?
}
