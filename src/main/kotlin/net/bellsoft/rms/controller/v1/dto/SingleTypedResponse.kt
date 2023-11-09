package net.bellsoft.rms.controller.v1.dto

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity

// TODO: SingleResponse 를 이 클래스로 대체 예정 (Issue: #23)
data class SingleTypedResponse<T>(
    val value: T,
) {
    fun toResponseEntity(status: HttpStatus) = ResponseEntity.status(status).body(value)

    companion object {
        fun <T> of(value: T) = SingleTypedResponse(value)
    }
}
