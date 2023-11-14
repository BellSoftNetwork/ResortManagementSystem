package net.bellsoft.rms.domain

import net.bellsoft.rms.config.P6SpyFormatter
import net.bellsoft.rms.config.QueryDslConfig
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest
import org.springframework.context.annotation.Import

@DataJpaTest
@Import(QueryDslConfig::class, P6SpyFormatter::class)
annotation class JpaEntityTest
