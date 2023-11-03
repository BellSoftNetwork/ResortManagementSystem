package net.bellsoft.rms.controller.v1.auth.dto

import io.swagger.v3.oas.annotations.media.Schema
import net.bellsoft.rms.controller.v1.dto.SingleResponse
import net.bellsoft.rms.domain.user.User

@Schema(description = "회원가입 정보")
data class RegisteredUserResponse(
    @Schema(description = "사용자 이메일", example = "bell@softbell.net")
    val email: String,
) : SingleResponse() {
    companion object {
        fun of(user: User): RegisteredUserResponse {
            return RegisteredUserResponse(user.email)
        }
    }
}
