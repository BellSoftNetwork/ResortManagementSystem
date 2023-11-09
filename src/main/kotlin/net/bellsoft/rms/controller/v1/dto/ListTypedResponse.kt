package net.bellsoft.rms.controller.v1.dto

import net.bellsoft.rms.service.dto.EntityListDto
import net.bellsoft.rms.service.dto.PageDto
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity

// TODO: ListResponse 를 이 클래스로 대체 예정 (Issue: #23)
data class ListTypedResponse<T>(
    val page: PageDto,
    val values: Collection<T>,
) {
    fun toResponseEntity(status: HttpStatus = HttpStatus.OK) = ResponseEntity.status(status).body(this)

    companion object {
        fun <T : EntityListDto<V>, V> of(entityListDto: T) = ListTypedResponse(
            page = entityListDto.page,
            values = entityListDto.values,
        )
    }
}
