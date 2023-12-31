package net.bellsoft.rms.common.dto.response

import io.swagger.v3.oas.annotations.media.Schema

@Schema(description = "에러 정보")
data class ErrorResponse(
    @Schema(description = "에러 메시지")
    val message: String,

    @Schema(description = "에러")
    val errors: List<String>? = null,

    @Schema(description = "필드 에러")
    val fieldErrors: List<String>? = null,
)
