package net.bellsoft.rms.authentication.handler

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.authentication.dto.response.AuthenticationSuccessResponse
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.mapper.UserMapper
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.security.core.Authentication
import org.springframework.security.web.authentication.AuthenticationSuccessHandler
import org.springframework.stereotype.Component

@Component
class JsonAuthenticationSuccessHandler(
    private val objectMapper: ObjectMapper,
    private val userMapper: UserMapper,
) : AuthenticationSuccessHandler {
    override fun onAuthenticationSuccess(
        request: HttpServletRequest,
        response: HttpServletResponse,
        authentication: Authentication,
    ) {
        val user = authentication.principal as User
        logger.info { "인증 성공 (username: ${user.username}, authorities: ${user.authorities})" }

        response.apply {
            this.status = HttpStatus.OK.value()
            this.contentType = MediaType.APPLICATION_JSON_VALUE
        }

        objectMapper.writeValue(
            response.writer,
            AuthenticationSuccessResponse(
                message = "인증 성공",
                value = userMapper.toDto(user),
            ),
        )
    }

    companion object : KLogging()
}
