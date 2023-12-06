package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.mapper.common.IdToReference
import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationDetailDto
import net.bellsoft.rms.service.reservation.dto.ReservationPatchDto
import org.mapstruct.DecoratedWith
import org.mapstruct.IterableMapping
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.MappingTarget
import org.mapstruct.Mappings
import org.mapstruct.NullValueMappingStrategy
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    uses = [
        JsonNullableMapper::class,
        ReferenceMapper::class,
        UserMapper::class,
        RoomMapper::class,
        ReservationMethodMapper::class,
    ],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
@DecoratedWith(ReservationMapperDecorator::class)
interface ReservationMapper {
    fun toDto(entity: Reservation): ReservationDetailDto

    @IterableMapping(nullValueMappingStrategy = NullValueMappingStrategy.RETURN_DEFAULT)
    @Mappings(
        Mapping(target = "reservationMethod", source = "reservationMethodId", qualifiedBy = [IdToReference::class]),
        Mapping(target = "rooms", ignore = true),
    )
    fun toEntity(dto: ReservationCreateDto): Reservation

    @Mappings(
        Mapping(target = "reservationMethod", source = "reservationMethodId", qualifiedBy = [IdToReference::class]),
        Mapping(target = "rooms", ignore = true),
    )
    fun updateEntity(dto: ReservationPatchDto, @MappingTarget entity: Reservation)
}
