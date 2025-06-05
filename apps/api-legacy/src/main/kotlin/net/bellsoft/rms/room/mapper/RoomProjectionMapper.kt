package net.bellsoft.rms.room.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.reservation.mapper.ReservationMapper
import net.bellsoft.rms.room.dto.projection.RoomLastReservationProjection
import net.bellsoft.rms.room.dto.response.RoomGroupDetailDto
import net.bellsoft.rms.room.dto.response.RoomLastStayDetailDto
import net.bellsoft.rms.room.entity.RoomGroup
import net.bellsoft.rms.user.mapper.UserMapper
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Mappings

@Mapper(
    config = BaseMapperConfig::class,
    uses = [UserMapper::class, RoomMapper::class, ReservationMapper::class],
)
interface RoomProjectionMapper {
    @Mappings(
        Mapping(target = "rooms", source = "roomProjections"),
    )
    fun toDto(entity: RoomGroup, roomProjections: List<RoomLastReservationProjection>): RoomGroupDetailDto

    @Mappings(
        Mapping(target = "room", source = "room"),
        Mapping(target = "lastReservation", source = "lastReservation"),
    )
    fun toDto(entity: RoomLastReservationProjection): RoomLastStayDetailDto
}
