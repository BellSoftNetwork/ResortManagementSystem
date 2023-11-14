package net.bellsoft.rms.component.history.dto

import io.swagger.v3.oas.annotations.media.Schema
import net.bellsoft.rms.component.history.type.HistoryType
import net.bellsoft.rms.domain.base.dto.RevisionDetails
import java.time.LocalDateTime

@Schema(description = "엔티티 데이터 이력 정보")
data class EntityHistoryDto<T>(
    @Schema(description = "해당 이력의 엔티티 상태")
    val entity: T,

    @Schema(description = "이력 타입", example = "CREATED")
    val historyType: HistoryType,

    @Schema(description = "이력 생성일", example = "2020-01-01 00:00:00")
    val historyCreatedAt: LocalDateTime,

    @Schema(description = "변경된 필드", example = "name")
    val updatedFields: HashSet<String>,
) {
    companion object {
        fun <ENTITY : Any, DTO> of(revisionDetails: RevisionDetails<ENTITY>, converter: (ENTITY) -> DTO) =
            EntityHistoryDto(
                entity = converter(revisionDetails.entity),
                historyType = HistoryType.of(revisionDetails.revisionType),
                historyCreatedAt = revisionDetails.revisionInfo.createdAt,
                updatedFields = revisionDetails.modifiedFields,
            )
    }
}
