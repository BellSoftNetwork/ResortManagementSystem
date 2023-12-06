package net.bellsoft.rms.mapper.common

import jakarta.persistence.EntityManager
import net.bellsoft.rms.service.common.dto.EntityReferenceDto
import org.mapstruct.TargetType
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.stereotype.Component

@Component
class ReferenceMapper(
    private val entityManager: EntityManager,
) {
    @IdToReference
    fun <T> longIdToReference(id: Long, @TargetType type: Class<T>): T =
        id.let { entityManager.getReference(type, it) }

    @IdToReference
    fun <T> nullableIdToReference(nullable: JsonNullable<Long?>, @TargetType type: Class<T>) =
        nullable.orElse(null)?.let { entityManager.getReference(type, it) }

    fun <T> refIdToReference(dto: EntityReferenceDto, @TargetType type: Class<T>): T =
        entityManager.getReference(type, dto.id)

    fun <T> nullableRefIdToReference(id: JsonNullable<EntityReferenceDto?>, @TargetType type: Class<T>) =
        id.orElse(null)?.let { entityManager.getReference(type, it.id) }

    fun <T> refIdsToReference(
        ids: Collection<EntityReferenceDto>,
        @TargetType type: Class<T>,
    ): Collection<T> =
        ids.map { entityManager.getReference(type, it.id) }

    fun <T> nullableRefIdsToReference(
        ids: JsonNullable<Collection<EntityReferenceDto>>,
        @TargetType type: Class<T>,
    ): Collection<T> = ids.get().map { entityManager.getReference(type, it.id) }
}
