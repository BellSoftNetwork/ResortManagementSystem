package net.bellsoft.rms.controller.v1.auth

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import mu.KLogging
import net.bellsoft.rms.controller.v1.auth.dto.RegisteredUserResponse
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.service.AuthService
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "인증", description = "인증 API")
@RestController
@Validated
@RequestMapping("/api/v1/auth")
class AuthController(
    val authService: AuthService,
) {
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
    ): ResponseEntity<RegisteredUserResponse> {
        return ResponseEntity
            .status(HttpStatus.CREATED)
            .body(authService.register(userRegistrationRequest))
    }

    companion object : KLogging()
}
