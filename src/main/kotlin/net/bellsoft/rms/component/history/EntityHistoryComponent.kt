package net.bellsoft.rms.component.history

import net.bellsoft.rms.component.history.dto.EntityHistoryDto
import net.bellsoft.rms.domain.revision.RevisionDetailsRepository
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component
import kotlin.reflect.KClass

@Component
class EntityHistoryComponent(
    private val revisionDetailsRepository: RevisionDetailsRepository,
) {
    fun <ENTITY : Any, DTO> findAllHistory(
        clazz: KClass<*>,
        converter: (ENTITY) -> DTO,
        id: Long,
    ) = revisionDetailsRepository
        .findAllByIdToRevisionInfo<ENTITY>(clazz, id)
        .map { EntityHistoryDto.of(it, converter) }

    fun <ENTITY : Any, DTO> findAllHistory(
        clazz: KClass<*>,
        converter: (ENTITY) -> DTO,
        id: Long,
        pageable: Pageable,
    ) = revisionDetailsRepository
        .findAllByIdToRevisionInfo<ENTITY>(clazz, id, pageable)
        .map { EntityHistoryDto.of(it, converter) }
}
