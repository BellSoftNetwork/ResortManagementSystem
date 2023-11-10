package net.bellsoft.rms.controller.v1.admin.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.service.admin.dto.AccountPatchDto

@Schema(description = "계정 수정 요청 정보")
data class AccountPatchRequest(
    @Schema(description = "비밀번호", example = "password!@#")
    @field:Size(min = 8, max = 20)
    val password: String? = null,

    @Schema(description = "이름", example = "방울")
    @field:Size(min = 2, max = 20)
    val name: String? = null,

    @Schema(description = "계정 잠금 상태", example = "false")
    val isLock: Boolean? = null,

    @Schema(description = "계정 권한", example = "NORMAL")
    val role: UserRole? = null,
) {
    fun toDto() = AccountPatchDto(
        password = password,
        name = name,
        isLock = isLock,
        role = role,
    )
}
