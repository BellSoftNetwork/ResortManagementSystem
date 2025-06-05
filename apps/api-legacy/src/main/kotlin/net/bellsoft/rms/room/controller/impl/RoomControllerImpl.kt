package net.bellsoft.rms.room.controller.impl

import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.room.controller.RoomController
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.filter.RoomRequestFilter
import net.bellsoft.rms.room.dto.request.RoomCreateRequest
import net.bellsoft.rms.room.dto.request.RoomPatchRequest
import net.bellsoft.rms.room.dto.service.RoomCreateDto
import net.bellsoft.rms.room.dto.service.RoomPatchDto
import net.bellsoft.rms.room.service.RoomService
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable
import org.springframework.web.bind.annotation.RestController

@RestController
class RoomControllerImpl(
    private val roomService: RoomService,
) : RoomController {
    override fun getRooms(pageable: Pageable, filter: RoomRequestFilter) = ListResponse
        .of(roomService.findAll(pageable, RoomFilterDto.of(filter)), filter)
        .toResponseEntity()

    override fun getRoom(id: Long) = SingleResponse
        .of(roomService.find(id))
        .toResponseEntity()

    override fun createRoom(
        user: User,
        request: RoomCreateRequest,
    ) = SingleResponse
        .of(roomService.create(RoomCreateDto.of(request)))

    override fun updateRoom(
        user: User,
        id: Long,
        request: RoomPatchRequest,
    ) = SingleResponse
        .of(roomService.update(id, RoomPatchDto.of(request)))
        .toResponseEntity()

    override fun deleteRoom(
        user: User,
        id: Long,
    ) {
        roomService.delete(id)
    }

    override fun getRoomHistory(
        id: Long,
        pageable: Pageable,
    ) = ListResponse
        .of(roomService.findHistory(id, pageable))
        .toResponseEntity()
}
