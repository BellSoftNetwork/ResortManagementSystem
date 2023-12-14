package net.bellsoft.rms.room.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.room.dto.response.RoomGroupSummaryDto
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.bellsoft.rms.room.dto.service.RoomGroupPatchDto
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.entity.RoomGroup
import net.bellsoft.rms.user.mapper.UserMapper
import org.mapstruct.BeanMapping
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.MappingTarget
import org.mapstruct.Mappings
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    config = BaseMapperConfig::class,
    uses = [UserMapper::class],
)
interface RoomGroupMapper {
    fun toDto(entity: RoomGroup): RoomGroupSummaryDto

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.SET_TO_DEFAULT)
    @Mappings(
        Mapping(target = "updatedBy", ignore = true),
        Mapping(target = "rooms", source = "rooms"),
    )
    fun toEntity(dto: RoomGroupCreateDto, rooms: MutableList<Room> = mutableListOf()): RoomGroup

    @Mappings(
        Mapping(target = "updatedBy", ignore = true),
        Mapping(target = "rooms", ignore = true),
    )
    fun updateEntity(dto: RoomGroupPatchDto, @MappingTarget entity: RoomGroup)
}
