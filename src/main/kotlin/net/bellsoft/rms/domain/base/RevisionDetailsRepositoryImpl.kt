package net.bellsoft.rms.domain.base

import jakarta.persistence.EntityManager
import net.bellsoft.rms.domain.base.dto.RevisionDetails
import org.hibernate.envers.AuditReaderFactory
import org.hibernate.envers.query.AuditEntity
import org.hibernate.envers.query.AuditQuery
import org.hibernate.envers.query.order.AuditOrder
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import org.springframework.stereotype.Repository
import org.springframework.transaction.annotation.Transactional
import kotlin.reflect.KClass

@Transactional(readOnly = true)
@Repository
class RevisionDetailsRepositoryImpl(
    private val entityManager: EntityManager,
) : RevisionDetailsRepository {
    override fun <T : Any> findAllByIdToRevisionInfo(clazz: KClass<*>, id: Long): List<RevisionDetails<T>> {
        return createBaseChangeQuery(clazz, id)
            .resultList
            .map { RevisionDetails.of(it) }
    }

    override fun <T : Any> findAllByIdToRevisionInfo(
        clazz: KClass<*>,
        id: Long,
        pageable: Pageable,
    ): Page<RevisionDetails<T>> {
        val count = createBaseQuery(clazz, id)
            .addProjection(AuditEntity.revisionNumber().count())
            .singleResult as Long

        val baseChangeQuery = createBaseChangeQuery(clazz, id)

        mapPropertySort(pageable.sort).forEach(baseChangeQuery::addOrder)

        val resultList = baseChangeQuery
            .setFirstResult(pageable.offset.toInt())
            .setMaxResults(pageable.pageSize)
            .resultList
            .map { RevisionDetails.of<T>(it) }

        return PageImpl(resultList, pageable, count)
    }

    private fun mapPropertySort(sort: Sort): List<AuditOrder> {
        if (sort.isEmpty)
            return listOf(AuditEntity.revisionNumber().asc())

        val result: MutableList<AuditOrder> = mutableListOf()

        for (order in sort) {
            val property = AuditEntity.property(order.property)

            result.add(
                if (order.direction.isAscending)
                    property.asc()
                else
                    property.desc(),
            )
        }

        return result
    }

    private fun createBaseQuery(clazz: KClass<*>, id: Long): AuditQuery {
        return auditReader().createQuery()
            .forRevisionsOfEntity(clazz.java, false, true)
            .add(AuditEntity.id().eq(id))
    }

    private fun createBaseChangeQuery(clazz: KClass<*>, id: Long): AuditQuery {
        return auditReader().createQuery()
            .forRevisionsOfEntityWithChanges(clazz.java, true)
            .add(AuditEntity.id().eq(id))
    }

    private fun auditReader() = AuditReaderFactory.get(entityManager)
}
