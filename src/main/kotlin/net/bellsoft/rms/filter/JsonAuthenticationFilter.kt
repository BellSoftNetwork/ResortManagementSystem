package net.bellsoft.rms.filter

import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import net.bellsoft.rms.component.auth.dto.LoginRequest
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.Authentication
import org.springframework.security.web.authentication.AbstractAuthenticationProcessingFilter
import org.springframework.security.web.util.matcher.AntPathRequestMatcher
import org.springframework.util.StringUtils

class JsonAuthenticationFilter(requestMatcher: AntPathRequestMatcher) :
    AbstractAuthenticationProcessingFilter(requestMatcher) {
    private val objectMapper = jacksonObjectMapper()

    override fun attemptAuthentication(
        request: HttpServletRequest,
        response: HttpServletResponse,
    ): Authentication {
        check(isValidRequest(request)) { "'Content-Type: application/json' 설정 필요" }

        val signInRequest = objectMapper.readValue(request.reader, LoginRequest::class.java)
        request.setAttribute("email", signInRequest.email)

        val token = signInRequest.let {
            require(StringUtils.hasText(it.email) && StringUtils.hasText(it.password)) { "'email' 및 'password' 정보 필요'" }

            UsernamePasswordAuthenticationToken(it.email, it.password)
        }

        return authenticationManager.authenticate(token)
    }

    private fun isValidRequest(request: HttpServletRequest): Boolean {
        if (request.method != HttpMethod.POST.name())
            return false

        if (request.contentType != MediaType.APPLICATION_JSON_VALUE)
            return false

        return true
    }
}
