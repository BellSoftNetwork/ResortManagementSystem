package net.bellsoft.rms.common.util

import jakarta.persistence.EntityManager
import jakarta.persistence.PersistenceContext
import org.springframework.beans.factory.InitializingBean
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional

@Component
class TestDatabaseSupport : InitializingBean {
    @PersistenceContext
    private lateinit var entityManager: EntityManager

    private var tableNames: List<String> = listOf()

    override fun afterPropertiesSet() {
        tableNames = entityManager
            .createNativeQuery("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = SCHEMA()")
            .resultList
            .map { it.toString() }
            .filterNot { it.startsWith("database_changelog") }
    }

    @Transactional
    fun clear() {
        entityManager.clear()
        entityManager.createNativeQuery("SET REFERENTIAL_INTEGRITY FALSE").executeUpdate()
        tableNames.forEach { entityManager.createNativeQuery("TRUNCATE TABLE $it").executeUpdate() }
        entityManager.createNativeQuery("SET REFERENTIAL_INTEGRITY TRUE").executeUpdate()
    }
}
