package net.bellsoft.rms.controller.v1.main

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import net.bellsoft.rms.controller.v1.main.dto.EnvResponse
import net.bellsoft.rms.controller.v1.main.dto.WhoAmIResponse
import net.bellsoft.rms.domain.user.User
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestMethod
import org.springframework.web.bind.annotation.RestController

@Tag(name = "메인", description = "메인 API")
@RestController
@RequestMapping("/api/v1")
class MainController {
    @Operation(summary = "백엔드 환경 정보", description = "백엔드 환경 정보 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/env")
    fun displayEnv() = EnvResponse.of("Resort Management System", "RMS", "v1.0")

    @Operation(summary = "로그인 계정 정보", description = "로그인 계정 정보 조회")
    @SecurityRequirement(name = "basicAuth")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @RequestMapping("/whoami", method = [RequestMethod.GET, RequestMethod.POST])
    fun displayMySelf(@AuthenticationPrincipal user: User) = WhoAmIResponse.of(user)
}
