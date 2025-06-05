package net.bellsoft.rms.common.controller.impl

import io.jsonwebtoken.ExpiredJwtException
import io.jsonwebtoken.MalformedJwtException
import io.jsonwebtoken.UnsupportedJwtException
import io.jsonwebtoken.security.SignatureException
import jakarta.validation.ConstraintViolationException
import mu.KLogging
import net.bellsoft.rms.authentication.exception.InvalidRefreshTokenException
import net.bellsoft.rms.authentication.exception.InvalidTokenException
import net.bellsoft.rms.authentication.exception.TooManyRequestsException
import net.bellsoft.rms.common.controller.ExceptionController
import net.bellsoft.rms.common.dto.response.ErrorResponse
import net.bellsoft.rms.common.exception.BadRequestException
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.PermissionRequiredDataException
import net.bellsoft.rms.common.exception.UnprocessableEntityException
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.http.converter.HttpMessageNotReadableException
import org.springframework.security.authentication.BadCredentialsException
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.validation.BindException
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.bind.annotation.RestControllerAdvice

@RestControllerAdvice
class ExceptionControllerImpl : ExceptionController {
    override fun handleMethodArgumentNotValidException(
        ex: MethodArgumentNotValidException,
    ): ResponseEntity<ErrorResponse> {
        return handleBindException(ex)
    }

    override fun handleBindException(ex: BindException): ResponseEntity<ErrorResponse> {
        val fieldErrors = ex.fieldErrors.map {
            "'${it.field}'은(는) ${it.defaultMessage} (요청 값: ${it.rejectedValue})"
        }

        val globalErrors = ex.globalErrors.map { it.defaultMessage.toString() }

        return ResponseEntity
            .badRequest()
            .body(
                ErrorResponse(
                    message = "잘못된 요청",
                    errors = globalErrors,
                    fieldErrors = fieldErrors,
                ),
            )
    }

    override fun handleConstraintViolationException(
        ex: ConstraintViolationException,
    ): ResponseEntity<ErrorResponse> {
        val fieldErrors = ex.constraintViolations.map {
            "'${it.propertyPath}'은(는) ${it.message} (요청 값: ${it.invalidValue})"
        }

        return ResponseEntity
            .badRequest()
            .body(
                ErrorResponse(
                    message = "잘못된 요청",
                    errors = listOf(ex.message ?: ""),
                    fieldErrors = fieldErrors,
                ),
            )
    }

    override fun handleHttpMessageNotReadableException(
        ex: HttpMessageNotReadableException,
    ): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .badRequest()
            .body(ErrorResponse("JSON 파싱 실패"))
    }

    override fun handleBadRequestException(ex: BadRequestException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .badRequest()
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleUnprocessableEntityException(ex: UnprocessableEntityException): ResponseEntity<ErrorResponse> {
        logger.error(ex.message.toString())

        return ResponseEntity
            .unprocessableEntity()
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleDataIntegrityViolationException(
        ex: DataIntegrityViolationException,
    ): ResponseEntity<ErrorResponse> {
        logger.error(ex.message.toString())

        return ResponseEntity
            .unprocessableEntity()
            .body(ErrorResponse("처리할 수 없는 데이터"))
    }

    override fun handleUsernameNotFoundException(ex: UsernameNotFoundException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("유효하지 않은 계정"))
    }

    override fun handleInvalidTokenException(ex: InvalidTokenException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleIllegalArgumentException(ex: IllegalArgumentException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.BAD_REQUEST)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleIllegalStateException(ex: IllegalStateException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.BAD_REQUEST)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleDataNotFoundException(ex: DataNotFoundException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.NOT_FOUND)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handlePermissionRequiredDataException(
        ex: PermissionRequiredDataException,
    ): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.FORBIDDEN)
            .body(ErrorResponse(ex.message.toString()))
    }

    companion object : KLogging()

    override fun handleTooManyRequestsException(ex: TooManyRequestsException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.TOO_MANY_REQUESTS)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleInvalidRefreshTokenException(ex: InvalidRefreshTokenException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse(ex.message.toString()))
    }

    override fun handleExpiredJwtException(ex: ExpiredJwtException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("토큰이 만료되었습니다."))
    }

    override fun handleMalformedJwtException(ex: MalformedJwtException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("잘못된 형식의 토큰입니다."))
    }

    override fun handleUnsupportedJwtException(ex: UnsupportedJwtException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("지원하지 않는 토큰입니다."))
    }

    override fun handleSignatureException(ex: SignatureException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("토큰 서명이 유효하지 않습니다."))
    }

    override fun handleBadCredentialsException(ex: BadCredentialsException): ResponseEntity<ErrorResponse> {
        logger.info(ex.message.toString())

        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse("아이디 또는 비밀번호가 일치하지 않습니다."))
    }
}
