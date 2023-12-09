package net.bellsoft.rms.room.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.reservation.entity.ReservationRoom
import net.bellsoft.rms.room.dto.response.RoomDetailDto
import net.bellsoft.rms.room.dto.service.RoomCreateDto
import net.bellsoft.rms.room.dto.service.RoomPatchDto
import net.bellsoft.rms.room.entity.Room
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
interface RoomMapper {
    fun toDto(entity: Room): RoomDetailDto

    @Mapping(source = "room", target = ".")
    fun toDto(entity: ReservationRoom): RoomDetailDto

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.SET_TO_DEFAULT)
    @Mappings(
        Mapping(target = "updatedBy", ignore = true),
    )
    fun toEntity(dto: RoomCreateDto): Room

    @Mappings(
        Mapping(target = "updatedBy", ignore = true),
    )
    fun updateEntity(dto: RoomPatchDto, @MappingTarget entity: Room)
}
