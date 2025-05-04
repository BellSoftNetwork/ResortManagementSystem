package net.bellsoft.rms.common.controller

import io.jsonwebtoken.ExpiredJwtException
import io.jsonwebtoken.MalformedJwtException
import io.jsonwebtoken.UnsupportedJwtException
import io.jsonwebtoken.security.SignatureException
import io.swagger.v3.oas.annotations.media.Content
import io.swagger.v3.oas.annotations.media.Schema
import io.swagger.v3.oas.annotations.responses.ApiResponse
import jakarta.validation.ConstraintViolationException
import net.bellsoft.rms.authentication.exception.InvalidRefreshTokenException
import net.bellsoft.rms.authentication.exception.InvalidTokenException
import net.bellsoft.rms.authentication.exception.TooManyRequestsException
import net.bellsoft.rms.common.dto.response.ErrorResponse
import net.bellsoft.rms.common.exception.BadRequestException
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.PermissionRequiredDataException
import net.bellsoft.rms.common.exception.UnprocessableEntityException
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.http.ResponseEntity
import org.springframework.http.converter.HttpMessageNotReadableException
import org.springframework.security.authentication.BadCredentialsException
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.validation.BindException
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice

@RestControllerAdvice
interface ExceptionController {
    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(MethodArgumentNotValidException::class)
    fun handleMethodArgumentNotValidException(
        ex: MethodArgumentNotValidException,
    ): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(BindException::class)
    fun handleBindException(ex: BindException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(ConstraintViolationException::class)
    fun handleConstraintViolationException(
        ex: ConstraintViolationException,
    ): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(HttpMessageNotReadableException::class)
    fun handleHttpMessageNotReadableException(ex: HttpMessageNotReadableException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(BadRequestException::class)
    fun handleBadRequestException(ex: BadRequestException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "422", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(UnprocessableEntityException::class)
    fun handleUnprocessableEntityException(ex: UnprocessableEntityException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "422", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(DataIntegrityViolationException::class)
    fun handleDataIntegrityViolationException(ex: DataIntegrityViolationException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(UsernameNotFoundException::class)
    fun handleUsernameNotFoundException(ex: UsernameNotFoundException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(InvalidTokenException::class)
    fun handleInvalidTokenException(ex: InvalidTokenException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(IllegalArgumentException::class)
    fun handleIllegalArgumentException(ex: IllegalArgumentException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "400", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(IllegalStateException::class)
    fun handleIllegalStateException(ex: IllegalStateException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "404", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(DataNotFoundException::class)
    fun handleDataNotFoundException(ex: DataNotFoundException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "403", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(PermissionRequiredDataException::class)
    fun handlePermissionRequiredDataException(ex: PermissionRequiredDataException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "429", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(TooManyRequestsException::class)
    fun handleTooManyRequestsException(ex: TooManyRequestsException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(InvalidRefreshTokenException::class)
    fun handleInvalidRefreshTokenException(ex: InvalidRefreshTokenException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(ExpiredJwtException::class)
    fun handleExpiredJwtException(ex: ExpiredJwtException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(MalformedJwtException::class)
    fun handleMalformedJwtException(ex: MalformedJwtException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(UnsupportedJwtException::class)
    fun handleUnsupportedJwtException(ex: UnsupportedJwtException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(SignatureException::class)
    fun handleSignatureException(ex: SignatureException): ResponseEntity<ErrorResponse>

    @ApiResponse(responseCode = "401", content = [Content(schema = Schema(implementation = ErrorResponse::class))])
    @ExceptionHandler(BadCredentialsException::class)
    fun handleBadCredentialsException(ex: BadCredentialsException): ResponseEntity<ErrorResponse>
}
