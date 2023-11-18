package net.bellsoft.rms.controller.v1.my.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import net.bellsoft.rms.service.auth.dto.AccountPatchDto

@Schema(description = "내 계정 수정 요청 정보")
data class MyPatchRequest(
    @Schema(description = "비밀번호", example = "password!@#")
    @field:Size(min = 8, max = 20)
    val password: String? = null,
) {
    fun toDto() = AccountPatchDto(
        password = password,
    )
}
