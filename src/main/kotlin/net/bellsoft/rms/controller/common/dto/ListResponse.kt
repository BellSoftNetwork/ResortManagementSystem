package net.bellsoft.rms.controller.common.dto

import io.swagger.v3.oas.annotations.media.Schema
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.common.dto.PageDto
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity

@Schema(description = "리스트 응답 정보")
data class ListResponse<T>(
    @Schema(description = "페이지 정보")
    val page: PageDto,

    @Schema(description = "필터 조건")
    val filter: Any?,

    @Schema(description = "아이템 값 리스트")
    val values: Collection<T>,
) {
    fun toResponseEntity(status: HttpStatus = HttpStatus.OK) = ResponseEntity.status(status).body(this)

    companion object {
        fun <T : EntityListDto<V>, V> of(entityListDto: T, filter: Any? = null) = ListResponse(
            page = entityListDto.page,
            filter = filter,
            values = entityListDto.values,
        )
    }
}
