package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.controller.v1.admin.dto.AdminUserPatchRequest
import net.bellsoft.rms.controller.v1.my.dto.MyPatchRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodPatchRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationPatchRequest
import net.bellsoft.rms.controller.v1.room.dto.RoomPatchRequest
import net.bellsoft.rms.service.auth.dto.UserPatchDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodPatchDto
import net.bellsoft.rms.service.reservation.dto.ReservationPatchDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
import org.mapstruct.Mapper
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
interface PatchDtoMapper {
    fun toDto(dto: MyPatchRequest): UserPatchDto
    fun toDto(dto: AdminUserPatchRequest): UserPatchDto
    fun toDto(dto: RoomPatchRequest): RoomPatchDto
    fun toDto(dto: ReservationPatchRequest): ReservationPatchDto
    fun toDto(dto: ReservationMethodPatchRequest): ReservationMethodPatchDto
}
