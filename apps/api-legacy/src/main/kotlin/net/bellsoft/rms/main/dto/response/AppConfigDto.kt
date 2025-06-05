package net.bellsoft.rms.main.dto.response

import io.swagger.v3.oas.annotations.media.Schema

@Schema(description = "앱 설정 정보")
data class AppConfigDto(
    @Schema(description = "일반 회원 신규 가입 가능 여부", example = "false")
    val isAvailableRegistration: Boolean,
)
