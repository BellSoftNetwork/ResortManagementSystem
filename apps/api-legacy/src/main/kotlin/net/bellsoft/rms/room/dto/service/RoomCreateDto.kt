package net.bellsoft.rms.room.dto.service

import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.room.dto.request.RoomCreateRequest
import net.bellsoft.rms.room.type.RoomStatus

data class RoomCreateDto(
    val roomGroup: EntityReferenceDto,
    val number: String,
    val note: String,
    val status: RoomStatus,
) {
    companion object {
        fun of(dto: RoomCreateRequest) = RoomCreateDto(
            roomGroup = dto.roomGroup,
            number = dto.number,
            note = dto.note,
            status = dto.status,
        )
    }
}
