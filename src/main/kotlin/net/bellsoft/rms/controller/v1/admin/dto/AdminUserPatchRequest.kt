package net.bellsoft.rms.controller.v1.admin.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.user.UserRole
import org.openapitools.jackson.nullable.JsonNullable

@Schema(description = "계정 수정 요청 정보")
data class AdminUserPatchRequest(
    @Schema(description = "계정 ID", example = "bell")
    @field:Size(min = 4, max = 30)
    val userId: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "이메일", example = "bell@bellsoft.net")
    @field:Email
    val email: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "비밀번호", example = "password!@#")
    @field:Size(min = 8, max = 20)
    val password: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "이름", example = "방울")
    @field:Size(min = 2, max = 20)
    val name: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "계정 잠금 상태", example = "false")
    val isLock: JsonNullable<Boolean> = JsonNullable.undefined(),

    @Schema(description = "계정 권한", example = "NORMAL")
    val role: JsonNullable<UserRole> = JsonNullable.undefined(),
)
