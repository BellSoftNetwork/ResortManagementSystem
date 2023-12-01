package net.bellsoft.rms.controller.v1.my.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.Size
import org.openapitools.jackson.nullable.JsonNullable

@Schema(description = "내 계정 수정 요청 정보")
data class MyPatchRequest(
    @Schema(description = "이메일", example = "bell@bellsoft.net")
    @field:Email
    val email: JsonNullable<String?> = JsonNullable.undefined(),

    @Schema(description = "비밀번호", example = "password!@#")
    @field:Size(min = 8, max = 20)
    val password: JsonNullable<String> = JsonNullable.undefined(),
)
