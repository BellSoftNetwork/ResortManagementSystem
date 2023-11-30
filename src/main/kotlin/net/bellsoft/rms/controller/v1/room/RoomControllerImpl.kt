package net.bellsoft.rms.controller.v1.room

import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.room.dto.RoomCreateRequest
import net.bellsoft.rms.controller.v1.room.dto.RoomPatchRequest
import net.bellsoft.rms.controller.v1.room.dto.RoomRequestFilter
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.room.RoomService
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomFilterDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
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
