package net.bellsoft.rms.util

import org.openapitools.jackson.nullable.JsonNullable

fun <T, D> JsonNullable<T>.convert(converter: (T) -> D): JsonNullable<D> {
    if (!this.isPresent)
        return JsonNullable.undefined()

    return JsonNullable.of(converter(this.get()))
}
