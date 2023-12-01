package net.bellsoft.rms.service.auth.dto

import io.swagger.v3.oas.annotations.media.Schema
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import java.time.LocalDateTime

@Schema(description = "사용자 상세 정보")
data class UserDetailDto(
    @Schema(description = "사용자 고유 id", example = "1")
    val id: Long,

    @Schema(description = "사용자 계정", example = "bell")
    val userId: String?, // TODO: Not Null 로 변경 예정

    @Schema(description = "사용자 이메일", example = "bell@softbell.net")
    val email: String,

    @Schema(description = "사용자 이름", example = "방울")
    val name: String,

    @Schema(description = "사용자 권한", example = "NORMAL")
    val role: UserRole,

    @Schema(description = "사용자 상태", example = "ACTIVE")
    val status: UserStatus,

    @Schema(description = "프로필 이미지 주소", example = "https://gravatar.com/avatar/00000000000000000000")
    val profileImageUrl: String,

    @Schema(description = "사용자 등록 시각", example = "2020-01-01 00:00:00")
    val createdAt: LocalDateTime,

    @Schema(description = "사용자 수정 시각", example = "2020-01-01 00:00:00")
    val updatedAt: LocalDateTime,
)
