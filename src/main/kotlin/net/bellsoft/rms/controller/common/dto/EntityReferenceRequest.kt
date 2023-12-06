package net.bellsoft.rms.controller.common.dto

import io.swagger.v3.oas.annotations.media.Schema

@Schema(description = "엔티티 참조 ID 정보")
data class EntityReferenceRequest(
    @Schema(description = "엔티티 ID", example = "1")
    val id: Long,
)
