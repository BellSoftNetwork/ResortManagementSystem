package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.controller.v1.room.dto.RoomCreateRequest
import net.bellsoft.rms.controller.v1.room.dto.RoomRequestFilter
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomDetailDto
import net.bellsoft.rms.service.room.dto.RoomFilterDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.MappingTarget
import org.mapstruct.Mappings
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    uses = [JsonNullableMapper::class, ReferenceMapper::class],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
interface RoomMapper {
    @Mappings(
        Mapping(target = "createdBy", source = "createdBy.email"),
        Mapping(target = "updatedBy", source = "updatedBy.email"),
    )
    fun toDto(entity: Room): RoomDetailDto

    fun toDto(dto: RoomCreateRequest): RoomCreateDto
    fun toDto(dto: RoomRequestFilter): RoomFilterDto

    fun toEntity(dto: RoomCreateDto): Room

    fun updateEntity(dto: RoomPatchDto, @MappingTarget entity: Room)
}
