package net.bellsoft.rms.controller.v1.auth.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.Size

@Schema(description = "회원가입 요청 정보")
data class UserRegistrationRequest(
    @Schema(description = "계정 ID", example = "bell")
    @field:Size(min = 4, max = 30)
    val userId: String,

    @Schema(description = "이메일", example = "bell@bellsoft.net")
    @field:Email
    val email: String,

    @Schema(description = "비밀번호", example = "password!@#")
    @field:Size(min = 8, max = 20)
    val password: String,

    @Schema(description = "이름", example = "방울")
    @field:Size(min = 2, max = 20)
    val name: String,
)
