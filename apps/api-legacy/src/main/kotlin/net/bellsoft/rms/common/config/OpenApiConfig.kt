package net.bellsoft.rms.common.config

import io.swagger.v3.oas.annotations.enums.SecuritySchemeIn
import io.swagger.v3.oas.annotations.enums.SecuritySchemeType
import io.swagger.v3.oas.annotations.security.SecurityScheme
import io.swagger.v3.oas.models.ExternalDocumentation
import io.swagger.v3.oas.models.OpenAPI
import io.swagger.v3.oas.models.info.Info
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
@SecurityScheme(name = "basicAuth", scheme = "basic", type = SecuritySchemeType.APIKEY, `in` = SecuritySchemeIn.COOKIE)
class OpenApiConfig {
    @Bean
    fun customOpenAPI(): OpenAPI = OpenAPI()
        .info(
            Info().apply {
                title = "Resort Management System API"
                description = "Resort Management System 공식 API 문서입니다"
                version = "v0.0.1"
            },
        )
        .externalDocs(
            ExternalDocumentation()
                .description("Resort Management System API Git Repository")
                .url("https://gitlab.bellsoft.net/resort-assistant/resort-management-system.git"),
        )
}
