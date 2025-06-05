package net.bellsoft.rms.authentication.service

import mu.KLogging
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.entity.LoginAttempt
import net.bellsoft.rms.authentication.exception.TooManyRequestsException
import net.bellsoft.rms.authentication.repository.LoginAttemptRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Service
import java.time.LocalDateTime.now

/**
 * 로그인 시도 서비스
 *
 * 로그인 시도를 기록하고 브루트 포스 공격을 방지하는 기능을 제공한다.
 */
@Service
class LoginAttemptService(
    private val loginAttemptRepository: LoginAttemptRepository,
    @Value("\${security.login-attempt.max-attempts}") private val maxAttempts: Int,
    @Value("\${security.login-attempt.window-minutes}") private val windowMinutes: Int,
) {
    /**
     * 로그인 시도를 기록한다.
     *
     * @param username 사용자 ID
     * @param deviceInfoDto 디바이스 정보 DTO
     * @param successful 성공 여부
     */
    fun recordLoginAttempt(
        username: String,
        deviceInfoDto: DeviceInfoDto,
        successful: Boolean,
    ) {
        loginAttemptRepository.save(
            LoginAttempt(
                username = username,
                ipAddress = deviceInfoDto.ipAddress,
                successful = successful,
                userAgent = deviceInfoDto.userAgent,
                osInfo = deviceInfoDto.osInfo,
                languageInfo = deviceInfoDto.languageInfo,
                deviceFingerprint = deviceInfoDto.deviceFingerprint,
            ),
        )
    }

    /**
     * 로그인 시도가 허용되는지 확인한다. (DeviceInfo DTO 사용)
     * 최대 시도 횟수를 초과하면 예외를 발생시킨다.
     * 마지막 성공한 로그인 이후의 실패 시도만 카운트한다.
     *
     * @param username 사용자 ID
     * @param deviceInfoDto 디바이스 정보 DTO
     * @throws TooManyRequestsException 최대 시도 횟수를 초과한 경우
     */
    fun checkLoginAttempts(username: String, deviceInfoDto: DeviceInfoDto) {
        val windowStart = now().minusMinutes(windowMinutes.toLong())

        // 사용자 ID와 IP 주소 조합 기반 검사 (마지막 성공 이후 실패만 카운트)
        val combinedFailedAttempts = loginAttemptRepository.countFailedAttemptsAfterLastSuccessfulByUsernameAndIp(
            username,
            deviceInfoDto.ipAddress,
            windowStart,
        )
        if (combinedFailedAttempts >= maxAttempts) {
            throw TooManyRequestsException("너무 많은 로그인 시도가 있었습니다. $windowMinutes 분 후에 다시 시도해주세요.")
        }

        // IP 주소 기반 검사 (마지막 성공 이후 실패만 카운트)
        val ipFailedAttempts =
            loginAttemptRepository.countFailedAttemptsAfterLastSuccessfulByIp(deviceInfoDto.ipAddress, windowStart)
        if (ipFailedAttempts >= maxAttempts * 3) { // IP 기반은 더 많은 시도 허용
            throw TooManyRequestsException("이 IP 주소에서 너무 많은 로그인 시도가 있었습니다. $windowMinutes 분 후에 다시 시도해주세요.")
        }
    }

    /**
     * 디바이스 변경 여부를 확인한다. (DeviceInfo DTO 사용)
     * 마지막 성공한 로그인과 디바이스 핑거프린트가 다르면 true를 반환한다.
     *
     * @param username 사용자 ID
     * @param deviceInfoDto 디바이스 정보 DTO
     * @return 디바이스 변경 여부
     */
    fun isDeviceChanged(username: String, deviceInfoDto: DeviceInfoDto): Boolean {
        val lastSuccessfulLogin = loginAttemptRepository
            .findTopByUsernameAndSuccessfulTrueOrderByAttemptAtDesc(username)
            ?: return false

        return lastSuccessfulLogin.deviceFingerprint != null &&
            lastSuccessfulLogin.deviceFingerprint != deviceInfoDto.deviceFingerprint
    }

    companion object : KLogging()
}
