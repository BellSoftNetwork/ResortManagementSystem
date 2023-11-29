package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomDetailDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
import org.mapstruct.Mapper
import org.mapstruct.MappingTarget
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    uses = [JsonNullableMapper::class, ReferenceMapper::class, UserMapper::class],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
interface RoomMapper {
    fun toDto(entity: Room): RoomDetailDto

    fun toEntity(dto: RoomCreateDto): Room

    fun updateEntity(dto: RoomPatchDto, @MappingTarget entity: Room)
}
