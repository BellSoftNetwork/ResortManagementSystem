package net.bellsoft.rms.authentication.handler

import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.authentication.exception.TokenExpiredException
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.security.core.AuthenticationException
import org.springframework.security.web.AuthenticationEntryPoint
import org.springframework.stereotype.Component

/**
 * JWT 인증 실패 처리 핸들러
 *
 * 인증 실패 시 401 응답을 반환한다.
 */
@Component
class JwtAuthenticationEntryPoint : AuthenticationEntryPoint {
    override fun commence(
        request: HttpServletRequest,
        response: HttpServletResponse,
        authException: AuthenticationException,
    ) {
        logger.error { "인증 실패: ${authException.message}" }

        response.status = HttpStatus.UNAUTHORIZED.value()
        response.contentType = MediaType.APPLICATION_JSON_VALUE

        val message = when (authException) {
            is TokenExpiredException -> "토큰이 만료되었습니다."
            else -> "인증에 실패했습니다."
        }

        response.writer.write("{\"message\":\"$message\",\"status\":401}")
    }

    companion object : KLogging()
}
