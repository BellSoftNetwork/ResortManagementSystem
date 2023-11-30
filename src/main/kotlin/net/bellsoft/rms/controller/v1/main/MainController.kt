package net.bellsoft.rms.controller.v1.main

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.main.dto.EnvResponse
import net.bellsoft.rms.service.config.dto.AppConfigDto
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "메인", description = "메인 API")
@RestController
@RequestMapping("/api/v1")
interface MainController {
    @Operation(summary = "백엔드 환경 정보", description = "백엔드 환경 정보 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/env")
    fun displayEnv(): SingleResponse<EnvResponse>

    @Operation(summary = "애플리케이션 설정 정보", description = "앱 설정 정보 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/config")
    fun displayConfig(): ResponseEntity<SingleResponse<AppConfigDto>>
}
