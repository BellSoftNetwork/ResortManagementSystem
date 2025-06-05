package net.bellsoft.rms.revision.config

import jakarta.persistence.EntityManagerFactory
import org.hibernate.envers.AuditReader
import org.hibernate.envers.AuditReaderFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
class EnversConfig(
    private val entityManagerFactory: EntityManagerFactory,
) {
    @Bean
    fun auditReader(): AuditReader = AuditReaderFactory.get(entityManagerFactory.createEntityManager())
}
