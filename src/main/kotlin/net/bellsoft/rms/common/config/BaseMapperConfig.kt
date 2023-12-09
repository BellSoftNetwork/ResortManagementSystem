package net.bellsoft.rms.common.config

import net.bellsoft.rms.common.mapper.JsonNullableMapper
import net.bellsoft.rms.common.mapper.ReferenceMapper
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
