package net.bellsoft.rms.service.dto

import org.springframework.data.domain.Page

data class EntityListDto<T>(
    val page: PageDto,
    val values: Collection<T>,
) {
    companion object {
        fun <ENTITY> of(page: Page<ENTITY>) = EntityListDto(
            page = PageDto.of(page),
            values = page.content,
        )

        fun <ENTITY, DTO> of(page: Page<ENTITY>, values: Collection<DTO>) = EntityListDto(
            page = PageDto.of(page),
            values = values,
        )

        fun <ENTITY, DTO> of(page: Page<ENTITY>, converter: (ENTITY) -> DTO) = EntityListDto(
            page = PageDto.of(page),
            values = page.content.map(converter),
        )

        fun <T> of(values: Collection<T>) = EntityListDto(
            page = PageDto(
                index = 0,
                size = values.size,
                totalElements = values.size.toLong(),
                totalPages = 1,
            ),
            values = values,
        )
    }
}
