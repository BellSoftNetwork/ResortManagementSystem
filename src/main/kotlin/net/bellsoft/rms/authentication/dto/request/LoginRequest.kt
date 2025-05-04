package net.bellsoft.rms.authentication.dto.request

import jakarta.validation.constraints.NotBlank

/**
 * 로그인 요청 DTO
 *
 * @property username 사용자 아이디
 * @property password 비밀번호
 */
data class LoginRequest(
    @field:NotBlank(message = "사용자 아이디는 필수입니다.")
    val username: String,

    @field:NotBlank(message = "비밀번호는 필수입니다.")
    val password: String,
)
