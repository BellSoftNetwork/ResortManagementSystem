package net.bellsoft.rms.controller.v1.main

import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.main.dto.EnvResponse
import net.bellsoft.rms.service.config.ConfigService
import org.springframework.beans.factory.annotation.Value
import org.springframework.web.bind.annotation.RestController

@RestController
class MainControllerImpl(
    @Value("\${application.deploy.commit_sha}") private val commitSha: String,
    @Value("\${application.deploy.commit_short_sha}") private val commitShortSha: String,
    @Value("\${application.deploy.commit_title}") private val commitTitle: String,
    @Value("\${application.deploy.commit_timestamp}") private val commitTimestamp: String,
    private val configService: ConfigService,
) : MainController {
    override fun displayEnv() = SingleResponse
        .of(
            EnvResponse.of(
                applicationFullName = "Resort Management System",
                applicationShortName = "RMS",
                commitSha = commitSha,
                commitShortSha = commitShortSha,
                commitTitle = commitTitle,
                commitTimestamp = commitTimestamp,
            ),
        )

    override fun displayConfig() = SingleResponse
        .of(configService.getAppConfig())
        .toResponseEntity()
}
