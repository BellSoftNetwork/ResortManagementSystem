package net.bellsoft.rms.mapper.common

import jakarta.persistence.EntityManager
import org.mapstruct.TargetType
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.stereotype.Component

@Component
class ReferenceMapper(
    private val entityManager: EntityManager,
) {
    @IdToReference
    fun <T> idToReference(id: Long?, @TargetType type: Class<T>): T? =
        id?.let { entityManager.getReference(type, it) }

    @IdToReference
    fun <T> idToReference(id: JsonNullable<Long?>, @TargetType type: Class<T>): T? =
        id.orElse(null)?.let { entityManager.getReference(type, it) }
}
