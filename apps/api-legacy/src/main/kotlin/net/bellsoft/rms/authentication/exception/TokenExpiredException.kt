package net.bellsoft.rms.authentication.exception

import org.springframework.http.HttpStatus
import org.springframework.security.core.AuthenticationException
import org.springframework.web.bind.annotation.ResponseStatus

/**
 * 토큰 만료 예외
 */
@ResponseStatus(HttpStatus.UNAUTHORIZED)
class TokenExpiredException(message: String = "토큰이 만료되었습니다.") : AuthenticationException(message)
