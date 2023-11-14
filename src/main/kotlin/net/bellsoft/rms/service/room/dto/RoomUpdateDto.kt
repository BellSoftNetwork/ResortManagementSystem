package net.bellsoft.rms.service.room.dto

import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomStatus

data class RoomUpdateDto(
    val number: String? = null,
    val peekPrice: Int? = null,
    val offPeekPrice: Int? = null,
    val description: String? = null,
    val note: String? = null,
    val status: RoomStatus? = null,
) {
    fun updateEntity(entity: Room) {
        number?.let { entity.number = it }
        peekPrice?.let { entity.peekPrice = it }
        offPeekPrice?.let { entity.offPeekPrice = it }
        description?.let { entity.description = it }
        note?.let { entity.note = it }
        status?.let { entity.status = it }
    }
}
