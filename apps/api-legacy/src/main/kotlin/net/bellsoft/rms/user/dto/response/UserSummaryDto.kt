package net.bellsoft.rms.user.dto.response

import io.swagger.v3.oas.annotations.media.Schema

@Schema(description = "사용자 요약 정보")
data class UserSummaryDto(
    @Schema(description = "사용자 고유 id", example = "1")
    val id: Long,

    @Schema(description = "사용자 ID", example = "bell")
    val userId: String,

    @Schema(description = "사용자 이메일", example = "bell@softbell.net")
    val email: String?,

    @Schema(description = "사용자 이름", example = "방울")
    val name: String,

    @Schema(description = "프로필 이미지 주소", example = "https://gravatar.com/avatar/00000000000000000000")
    val profileImageUrl: String,
)
