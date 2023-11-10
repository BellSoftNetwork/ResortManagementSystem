package net.bellsoft.rms.service.common.dto

import io.swagger.v3.oas.annotations.media.Schema
import org.springframework.data.domain.Page

@Schema(description = "페이지 정보")
data class PageDto(
    @Schema(description = "페이지 번호", example = "0")
    val index: Int,

    @Schema(description = "페이지 별 아이템 개수", example = "20")
    val size: Int,

    @Schema(description = "아이템 총 개수", example = "48")
    val totalElements: Long,

    @Schema(description = "총 페이지 개수", example = "3")
    val totalPages: Int,
) {
    companion object {
        fun <T> of(page: Page<T>): PageDto {
            return PageDto(
                index = page.number,
                size = page.size,
                totalElements = page.totalElements,
                totalPages = page.totalPages,
            )
        }
    }
}
