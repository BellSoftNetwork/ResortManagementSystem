package net.bellsoft.rms.room.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.response.RoomGroupDetailDto
import net.bellsoft.rms.room.dto.response.RoomGroupSummaryDto
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.bellsoft.rms.room.dto.service.RoomGroupPatchDto
import org.springframework.data.domain.Pageable

interface RoomGroupService {
    fun findAll(pageable: Pageable): EntityListDto<RoomGroupSummaryDto>

    fun find(id: Long, filter: RoomFilterDto): RoomGroupDetailDto

    fun create(request: RoomGroupCreateDto): RoomGroupSummaryDto

    fun update(id: Long, request: RoomGroupPatchDto): RoomGroupSummaryDto

    fun delete(id: Long)
}
