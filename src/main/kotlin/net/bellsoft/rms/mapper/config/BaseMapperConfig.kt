package net.bellsoft.rms.mapper.config

import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import org.mapstruct.MapperConfig
import org.mapstruct.NullValuePropertyMappingStrategy

@MapperConfig(
    uses = [
        JsonNullableMapper::class,
        ReferenceMapper::class,
    ],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
interface BaseMapperConfig
