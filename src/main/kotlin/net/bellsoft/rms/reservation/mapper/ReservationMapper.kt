package net.bellsoft.rms.reservation.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.payment.mapper.PaymentMethodMapper
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.room.mapper.RoomMapper
import net.bellsoft.rms.user.mapper.UserMapper
import org.mapstruct.DecoratedWith
import org.mapstruct.IterableMapping
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.MappingTarget
import org.mapstruct.Mappings
import org.mapstruct.NullValueMappingStrategy

@Mapper(
    config = BaseMapperConfig::class,
    uses = [
        UserMapper::class,
        RoomMapper::class,
        PaymentMethodMapper::class,
    ],
)
@DecoratedWith(ReservationMapperDecorator::class)
interface ReservationMapper {
    fun toDto(entity: Reservation): ReservationDetailDto

    @IterableMapping(nullValueMappingStrategy = NullValueMappingStrategy.RETURN_DEFAULT)
    @Mappings(
        Mapping(target = "rooms", ignore = true),
        Mapping(target = "updatedBy", ignore = true),
    )
    fun toEntity(dto: ReservationCreateDto): Reservation

    @Mappings(
        Mapping(target = "rooms", ignore = true),
        Mapping(target = "updatedBy", ignore = true),
    )
    fun updateEntity(dto: ReservationPatchDto, @MappingTarget entity: Reservation)
}
