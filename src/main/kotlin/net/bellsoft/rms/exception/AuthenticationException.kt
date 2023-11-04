package net.bellsoft.rms.exception

import org.springframework.security.core.AuthenticationException
import org.springframework.security.core.userdetails.UsernameNotFoundException

class UserNotFoundException(message: String = "존재하지 않는 사용자") : UsernameNotFoundException(message)
class PasswordMismatchException(message: String = "일치하지 않는 비밀번호") : UsernameNotFoundException(message)
class InvalidAuthHeaderException(message: String = "유효하지 않는 인증 헤더") : AuthenticationException(message)
open class InvalidTokenException(message: String = "유효하지 않는 토큰") : AuthenticationException(message)
class TokenNotFoundException(message: String = "존재하지 않는 토큰") : InvalidTokenException(message)
