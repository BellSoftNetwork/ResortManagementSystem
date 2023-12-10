package net.bellsoft.rms.room.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.response.RoomDetailDto
import net.bellsoft.rms.room.dto.service.RoomCreateDto
import net.bellsoft.rms.room.dto.service.RoomPatchDto
import org.springframework.data.domain.Pageable

interface RoomService {
    fun findAll(pageable: Pageable, filter: RoomFilterDto): EntityListDto<RoomDetailDto>

    fun find(id: Long): RoomDetailDto

    fun create(request: RoomCreateDto): RoomDetailDto

    fun update(id: Long, request: RoomPatchDto): RoomDetailDto

    fun delete(id: Long)

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<RoomDetailDto>>
}
