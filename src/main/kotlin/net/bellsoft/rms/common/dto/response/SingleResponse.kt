package net.bellsoft.rms.common.dto.response

import io.swagger.v3.oas.annotations.media.Schema
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity

@Schema(description = "단건 응답 정보")
data class SingleResponse<T>(
    @Schema(description = "아이템 값")
    val value: T,
) {
    fun toResponseEntity(status: HttpStatus = HttpStatus.OK) = ResponseEntity.status(status).body(this)

    companion object {
        fun <T> of(value: T) = SingleResponse(value)
    }
}
