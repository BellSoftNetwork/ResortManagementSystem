package net.bellsoft.rms.common.dto.service

import net.bellsoft.rms.common.dto.request.EntityReferenceRequest
import net.bellsoft.rms.common.entity.BaseEntity

data class EntityReferenceDto(
    val id: Long,
) {
    companion object {
        fun of(entity: BaseEntity) = EntityReferenceDto(
            id = entity.id,
        )

        fun of(dto: EntityReferenceRequest) = EntityReferenceDto(
            id = dto.id,
        )

        fun of(dto: Collection<EntityReferenceRequest>) = dto.map { of(it) }.toSet()
    }
}
