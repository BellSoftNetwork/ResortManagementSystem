package net.bellsoft.rms.controller.v1.auth

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.service.auth.dto.UserDetailDto
import org.springframework.http.ResponseEntity
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "인증", description = "인증 API")
@Validated
@RestController
@RequestMapping("/api/v1/auth")
interface AuthController {
    @Operation(summary = "회원가입", description = "회원가입 처리")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @PostMapping("/register")
    fun registerUser(
        @RequestBody @Valid
        userRegistrationRequest: UserRegistrationRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>>
}
