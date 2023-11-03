package net.bellsoft.rms.component.auth

import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.component.auth.dto.AuthenticationResponse
import net.bellsoft.rms.domain.user.User
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.security.core.Authentication
import org.springframework.security.web.authentication.AuthenticationSuccessHandler
import org.springframework.stereotype.Component

@Component
class JsonAuthenticationSuccessHandler : AuthenticationSuccessHandler {
    private val objectMapper = jacksonObjectMapper()

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
            AuthenticationResponse(
                email = user.email,
                message = "인증 성공",
                authorities = user.authorities,
            ),
        )
    }

    companion object : KLogging()
}
