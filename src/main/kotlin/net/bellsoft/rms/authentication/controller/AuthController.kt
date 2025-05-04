package net.bellsoft.rms.authentication.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.authentication.dto.request.LoginRequest
import net.bellsoft.rms.authentication.dto.request.RefreshTokenRequest
import net.bellsoft.rms.authentication.dto.request.UserRegistrationRequest
import net.bellsoft.rms.authentication.dto.response.TokenDto
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.dto.response.UserDetailDto
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

    @Operation(summary = "로그인", description = "로그인 처리 및 JWT 토큰 발급")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
            ApiResponse(responseCode = "401", description = "인증 실패"),
        ],
    )
    @PostMapping("/login")
    fun login(
        @RequestBody @Valid
        loginRequest: LoginRequest,
    ): ResponseEntity<SingleResponse<TokenDto>>

    @Operation(summary = "토큰 갱신", description = "리프레시 토큰을 사용하여 액세스 토큰 갱신")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
            ApiResponse(responseCode = "401", description = "유효하지 않은 리프레시 토큰"),
        ],
    )
    @PostMapping("/refresh")
    fun refreshToken(
        @RequestBody @Valid
        refreshTokenRequest: RefreshTokenRequest,
    ): ResponseEntity<SingleResponse<TokenDto>>
}
