package net.bellsoft.rms.service.dto

import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

data class PageDto(
    val index: Int,
    val size: Int,
    val totalElements: Long? = null,
    val totalPages: Int? = null,
) {
    companion object {
        fun of(pageable: Pageable): PageDto {
            return PageDto(
                index = pageable.pageNumber,
                size = pageable.pageSize,
            )
        }

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
