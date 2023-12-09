package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.domain.reservation.ReservationRoom
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.mapper.config.BaseMapperConfig
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomDetailDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
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
