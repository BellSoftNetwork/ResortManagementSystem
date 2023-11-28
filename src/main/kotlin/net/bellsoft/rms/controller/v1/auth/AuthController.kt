package net.bellsoft.rms.controller.v1.auth

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import mu.KLogging
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.mapper.model.UserMapper
import net.bellsoft.rms.service.auth.AuthService
import org.springframework.http.HttpStatus
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "인증", description = "인증 API")
@Validated
@RestController
@RequestMapping("/api/v1/auth")
class AuthController(
    private val authService: AuthService,
    private val userMapper: UserMapper,
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
    ) = SingleResponse
        .of((authService.register(userMapper.toDto(userRegistrationRequest))))
        .toResponseEntity(HttpStatus.CREATED)

    companion object : KLogging()
}
