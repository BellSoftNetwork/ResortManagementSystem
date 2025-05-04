package net.bellsoft.rms.authentication.dto.response

/**
 * JWT 토큰 정보를 담는 DTO
 *
 * @property accessToken 액세스 토큰
 * @property refreshToken 리프레시 토큰
 * @property accessTokenExpiresIn 액세스 토큰 만료 시간 (밀리초)
 */
data class TokenDto(
    val accessToken: String,
    val refreshToken: String,
    val accessTokenExpiresIn: Long,
)
