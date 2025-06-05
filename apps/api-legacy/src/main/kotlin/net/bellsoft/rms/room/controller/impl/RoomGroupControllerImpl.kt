package net.bellsoft.rms.room.controller.impl

import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.room.controller.RoomGroupController
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.filter.RoomRequestFilter
import net.bellsoft.rms.room.dto.request.RoomGroupCreateRequest
import net.bellsoft.rms.room.dto.request.RoomGroupPatchRequest
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.bellsoft.rms.room.dto.service.RoomGroupPatchDto
import net.bellsoft.rms.room.service.RoomGroupService
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable
import org.springframework.web.bind.annotation.RestController

@RestController
class RoomGroupControllerImpl(
    private val roomGroupService: RoomGroupService,
) : RoomGroupController {
    override fun getRoomGroups(pageable: Pageable) = ListResponse
        .of(roomGroupService.findAll(pageable))
        .toResponseEntity()

    override fun getRoomGroup(id: Long, filter: RoomRequestFilter) = SingleResponse
        .of(roomGroupService.find(id, RoomFilterDto.of(filter)))
        .toResponseEntity()

    override fun createRoomGroup(
        user: User,
        request: RoomGroupCreateRequest,
    ) = SingleResponse
        .of(roomGroupService.create(RoomGroupCreateDto.of(request)))

    override fun updateRoomGroup(
        user: User,
        id: Long,
        request: RoomGroupPatchRequest,
    ) = SingleResponse
        .of(roomGroupService.update(id, RoomGroupPatchDto.of(request)))
        .toResponseEntity()

    override fun deleteRoomGroup(
        user: User,
        id: Long,
    ) {
        roomGroupService.delete(id)
    }
}
