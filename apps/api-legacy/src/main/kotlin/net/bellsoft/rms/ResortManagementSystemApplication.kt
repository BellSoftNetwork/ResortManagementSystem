package net.bellsoft.rms

import net.bellsoft.rms.common.annotation.ExcludeFromJacocoGeneratedReport
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.data.envers.repository.support.EnversRevisionRepositoryFactoryBean
import org.springframework.data.jpa.repository.config.EnableJpaRepositories

@SpringBootApplication
@EnableJpaRepositories(repositoryFactoryBeanClass = EnversRevisionRepositoryFactoryBean::class)
class ResortManagementSystemApplication

@ExcludeFromJacocoGeneratedReport
fun main(args: Array<String>) {
    runApplication<ResortManagementSystemApplication>(*args)
}
