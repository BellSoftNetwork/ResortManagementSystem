package net.bellsoft.rms.room.repository

import net.bellsoft.rms.room.entity.RoomGroup
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.repository.history.RevisionRepository
import org.springframework.stereotype.Repository

@Repository
interface RoomGroupRepository :
    JpaRepository<RoomGroup, Long>,
    RevisionRepository<RoomGroup, Long, Long>,
    RoomGroupCustomRepository
