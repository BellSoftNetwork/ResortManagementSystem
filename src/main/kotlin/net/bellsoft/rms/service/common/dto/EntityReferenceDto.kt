package net.bellsoft.rms.service.common.dto

import net.bellsoft.rms.controller.common.dto.EntityReferenceRequest
import net.bellsoft.rms.domain.base.BaseEntity

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
