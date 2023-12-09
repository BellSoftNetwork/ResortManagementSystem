package net.bellsoft.rms.revision.repository

import net.bellsoft.rms.revision.dto.RevisionDetails
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import kotlin.reflect.KClass

@Repository
interface RevisionDetailsRepository {
    fun <T : Any> findAllByIdToRevisionInfo(clazz: KClass<*>, id: Long): List<RevisionDetails<T>>
    fun <T : Any> findAllByIdToRevisionInfo(clazz: KClass<*>, id: Long, pageable: Pageable): Page<RevisionDetails<T>>
}
