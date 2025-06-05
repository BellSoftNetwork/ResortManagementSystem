package net.bellsoft.rms.authentication.exception

import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.ResponseStatus

/**
 * 너무 많은 요청 예외
 *
 * 사용자가 짧은 시간 내에 너무 많은 요청을 보낸 경우 발생하는 예외
 */
@ResponseStatus(HttpStatus.TOO_MANY_REQUESTS)
class TooManyRequestsException(message: String) : RuntimeException(message)
