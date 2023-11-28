package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodCreateRequest
import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodDetailDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodPatchDto
import org.mapstruct.Mapper
import org.mapstruct.MappingTarget
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    uses = [JsonNullableMapper::class, ReferenceMapper::class],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
interface ReservationMethodMapper {
    fun toDto(entity: ReservationMethod): ReservationMethodDetailDto

    fun toDto(dto: ReservationMethodCreateRequest): ReservationMethodCreateDto

    fun toEntity(dto: ReservationMethodCreateDto): ReservationMethod

    fun updateEntity(dto: ReservationMethodPatchDto, @MappingTarget entity: ReservationMethod)
}