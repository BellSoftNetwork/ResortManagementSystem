package net.bellsoft.rms.authentication.dto.request

import jakarta.validation.constraints.NotBlank

/**
 * 리프레시 토큰 요청 DTO
 *
 * @property refreshToken 리프레시 토큰
 */
data class RefreshTokenRequest(
    @field:NotBlank(message = "리프레시 토큰은 필수입니다.")
    val refreshToken: String,
)
