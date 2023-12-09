package net.bellsoft.rms.revision.component

import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.revision.repository.RevisionDetailsRepository
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component
import kotlin.reflect.KClass

@Component
class EntityRevisionComponent(
    private val revisionDetailsRepository: RevisionDetailsRepository,
) {
    fun <ENTITY : Any, DTO> findAllHistory(
        clazz: KClass<*>,
        converter: (ENTITY) -> DTO,
        id: Long,
    ) = revisionDetailsRepository
        .findAllByIdToRevisionInfo<ENTITY>(clazz, id)
        .map { EntityRevisionDto.of(it, converter) }

    fun <ENTITY : Any, DTO> findAllHistory(
        clazz: KClass<*>,
        converter: (ENTITY) -> DTO,
        id: Long,
        pageable: Pageable,
    ) = revisionDetailsRepository
        .findAllByIdToRevisionInfo<ENTITY>(clazz, id, pageable)
        .map { EntityRevisionDto.of(it, converter) }
}
