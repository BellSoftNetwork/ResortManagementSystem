package net.bellsoft.rms.authentication.exception

import org.springframework.http.HttpStatus
import org.springframework.security.core.AuthenticationException
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.web.bind.annotation.ResponseStatus

class UserNotFoundException(message: String = "존재하지 않는 사용자") : UsernameNotFoundException(message)
class PasswordMismatchException(message: String = "일치하지 않는 비밀번호") : UsernameNotFoundException(message)
class InvalidAuthHeaderException(message: String = "유효하지 않는 인증 헤더") : AuthenticationException(message)

@ResponseStatus(HttpStatus.UNAUTHORIZED)
open class InvalidTokenException(message: String = "유효하지 않는 토큰") : AuthenticationException(message)
class TokenNotFoundException(message: String = "존재하지 않는 토큰") : InvalidTokenException(message)

@ResponseStatus(HttpStatus.UNAUTHORIZED)
class InvalidRefreshTokenException(message: String = "유효하지 않은 리프레시 토큰입니다.") : InvalidTokenException(message)

@ResponseStatus(HttpStatus.UNAUTHORIZED)
class InvalidAccessTokenException(message: String = "유효하지 않은 액세스 토큰입니다.") : InvalidTokenException(message)
