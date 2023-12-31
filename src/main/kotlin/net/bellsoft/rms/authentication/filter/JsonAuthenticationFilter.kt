package net.bellsoft.rms.authentication.filter

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import net.bellsoft.rms.authentication.dto.request.LoginRequest
import net.bellsoft.rms.authentication.exception.InvalidAuthHeaderException
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.Authentication
import org.springframework.security.web.authentication.AbstractAuthenticationProcessingFilter
import org.springframework.security.web.util.matcher.AntPathRequestMatcher
import org.springframework.util.StringUtils

class JsonAuthenticationFilter(requestMatcher: AntPathRequestMatcher, private val objectMapper: ObjectMapper) :
    AbstractAuthenticationProcessingFilter(requestMatcher) {
    override fun attemptAuthentication(
        request: HttpServletRequest,
        response: HttpServletResponse,
    ): Authentication {
        if (!isValidRequest(request))
            throw InvalidAuthHeaderException("'Content-Type: application/json' 설정 필요")

        val signInRequest = objectMapper.readValue(request.reader, LoginRequest::class.java)
        request.setAttribute("username", signInRequest.username)

        val token = signInRequest.let {
            if (StringUtils.hasText(it.username).not() || StringUtils.hasText(it.password).not())
                throw AuthenticationCredentialsNotFoundException("'username' 및 'password' 정보 필요")

            UsernamePasswordAuthenticationToken(it.username, it.password)
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
