package net.bellsoft.rms.authentication.controller.impl

import jakarta.servlet.http.HttpServletRequest
import mu.KLogging
import net.bellsoft.rms.authentication.component.JwtTokenProvider
import net.bellsoft.rms.authentication.controller.AuthController
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.dto.request.LoginRequest
import net.bellsoft.rms.authentication.dto.request.RefreshTokenRequest
import net.bellsoft.rms.authentication.dto.request.UserRegistrationRequest
import net.bellsoft.rms.authentication.exception.InvalidRefreshTokenException
import net.bellsoft.rms.authentication.service.LoginAttemptService
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.service.UserService
import org.springframework.http.HttpStatus
import org.springframework.security.authentication.AuthenticationManager
import org.springframework.security.authentication.BadCredentialsException
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.web.bind.annotation.RestController

@RestController
class AuthControllerImpl(
    private val userService: UserService,
    private val jwtTokenProvider: JwtTokenProvider,
    private val authenticationManager: AuthenticationManager,
    private val loginAttemptService: LoginAttemptService,
    private val request: HttpServletRequest,
) : AuthController {
    override fun registerUser(
        userRegistrationRequest: UserRegistrationRequest,
    ) = SingleResponse
        .of((userService.register(UserCreateDto.of(userRegistrationRequest))))
        .toResponseEntity(HttpStatus.CREATED)

    override fun login(loginRequest: LoginRequest) = try {
        val username = loginRequest.username
        val deviceInfoDto = DeviceInfoDto.fromRequest(request)

        // 로그인 시도 횟수 검사
        loginAttemptService.checkLoginAttempts(username, deviceInfoDto)

        // 디바이스 변경 감지
        val isDeviceChanged = loginAttemptService.isDeviceChanged(username, deviceInfoDto)
        if (isDeviceChanged) {
            logger.warn { "디바이스 변경 감지: $username, 이전 디바이스와 다른 디바이스에서 로그인 시도" }
            // 여기서 추가 검증 로직을 구현할 수 있음 (예: 이메일 알림, 2FA 요구 등)
        }

        try {
            // 인증 처리
            val authentication = authenticationManager.authenticate(
                UsernamePasswordAuthenticationToken(username, loginRequest.password),
            )
            SecurityContextHolder.getContext().authentication = authentication

            // 성공한 로그인 시도 기록
            loginAttemptService.recordLoginAttempt(
                username = username,
                deviceInfoDto = deviceInfoDto,
                successful = true,
            )

            // JWT 토큰 생성
            val user = authentication.principal as net.bellsoft.rms.user.entity.User
            val tokenDto = jwtTokenProvider.createTokens(user, deviceInfoDto)

            SingleResponse.of(tokenDto).toResponseEntity(HttpStatus.OK)
        } catch (e: BadCredentialsException) {
            // 실패한 로그인 시도 기록
            loginAttemptService.recordLoginAttempt(
                username = username,
                deviceInfoDto = deviceInfoDto,
                successful = false,
            )
            logger.warn { "로그인 실패 (잘못된 자격 증명): $username, IP: ${deviceInfoDto.ipAddress}" }
            throw e
        }
    } catch (e: Exception) {
        logger.error(e) { "로그인 실패: ${e.message}" }
        throw e
    }

    override fun refreshToken(refreshTokenRequest: RefreshTokenRequest) = try {
        val deviceInfoDto = DeviceInfoDto.fromRequest(request)

        try {
            // 리프레시 토큰 검증
            if (!jwtTokenProvider.validateToken(refreshTokenRequest.refreshToken)) {
                throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰입니다.")
            }

            // 새 액세스 토큰 발급
            val tokenDto = jwtTokenProvider
                .refreshTokens(refreshTokenRequest.refreshToken, deviceInfoDto)

            SingleResponse.of(tokenDto).toResponseEntity(HttpStatus.OK)
        } catch (e: InvalidRefreshTokenException) {
            logger.warn { "토큰 갱신 실패 (유효하지 않은 리프레시 토큰): IP: ${deviceInfoDto.ipAddress}" }
            throw e
        }
    } catch (e: Exception) {
        logger.error(e) { "토큰 갱신 실패: ${e.message}" }
        throw e
    }

    companion object : KLogging()
}
