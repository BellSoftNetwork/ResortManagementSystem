package net.bellsoft.rms.component.auth

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.component.auth.dto.AuthenticationSuccessResponse
import net.bellsoft.rms.domain.user.User
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.security.core.Authentication
import org.springframework.security.web.authentication.AuthenticationSuccessHandler
import org.springframework.stereotype.Component

@Component
class JsonAuthenticationSuccessHandler(private val objectMapper: ObjectMapper) : AuthenticationSuccessHandler {
    override fun onAuthenticationSuccess(
        request: HttpServletRequest,
        response: HttpServletResponse,
        authentication: Authentication,
    ) {
        val user = authentication.principal as User
        logger.info { "인증 성공 (email: ${user.username}, authorities: ${user.authorities})" }

        response.apply {
            this.status = HttpStatus.OK.value()
            this.contentType = MediaType.APPLICATION_JSON_VALUE
        }

        objectMapper.writeValue(
            response.writer,
            AuthenticationSuccessResponse.of(
                user = user,
                message = "인증 성공",
            ),
        )
    }

    companion object : KLogging()
}
