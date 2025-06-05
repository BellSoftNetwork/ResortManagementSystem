package net.bellsoft.rms.common.mapper

import org.mapstruct.Condition
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.stereotype.Component

@Component
class JsonNullableMapper {
    fun <T> wrap(value: T): JsonNullable<T> = JsonNullable.of(value)
    fun <T> unwrap(nullable: JsonNullable<T?>) = nullable.orElse(null)

    @Condition
    fun <T> isPresent(nullable: JsonNullable<T?>) = nullable.isPresent
}
