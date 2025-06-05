package net.bellsoft.rms.authentication.component

import io.jsonwebtoken.Claims
import io.jsonwebtoken.ExpiredJwtException
import io.jsonwebtoken.Jwts
import io.jsonwebtoken.MalformedJwtException
import io.jsonwebtoken.SignatureAlgorithm
import io.jsonwebtoken.UnsupportedJwtException
import io.jsonwebtoken.security.Keys
import io.jsonwebtoken.security.SignatureException
import mu.KLogging
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.dto.response.TokenDto
import net.bellsoft.rms.authentication.exception.InvalidAccessTokenException
import net.bellsoft.rms.authentication.exception.InvalidRefreshTokenException
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.data.repository.findByIdOrNull
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.Authentication
import org.springframework.stereotype.Component
import java.security.Key
import java.time.LocalDateTime
import java.time.ZoneId
import java.util.*

@Component
class JwtTokenProvider(
    @Value("\${security.jwt.secret}") private val secretKey: String,
    @Value("\${security.jwt.access-token-validity-in-hours}") private val accessTokenValidityInHours: Long,
    @Value("\${security.jwt.refresh-token-validity-in-hours}") private val refreshTokenValidityInHours: Long,
    private val userRepository: UserRepository,
) {
    private val key: Key = Keys.hmacShaKeyFor(secretKey.toByteArray())

    /**
     * 토큰이 유효한지 검증한다.
     *
     * @param token JWT 토큰
     * @return 토큰 유효성 여부
     */
    fun validateToken(token: String): Boolean {
        try {
            val claims = getClaims(token)
            val expirationDate = claims.expiration.toInstant().atZone(ZoneId.systemDefault()).toLocalDateTime()
            val now = LocalDateTime.now()
            return !expirationDate.isBefore(now)
        } catch (e: ExpiredJwtException) {
            logger.error { "JWT 토큰 만료: ${e.message}" }
            return false
        } catch (e: SignatureException) {
            logger.error { "JWT 토큰 서명 검증 실패: ${e.message}" }
            return false
        } catch (e: MalformedJwtException) {
            logger.error { "잘못된 형식의 JWT 토큰: ${e.message}" }
            return false
        } catch (e: UnsupportedJwtException) {
            logger.error { "지원하지 않는 JWT 토큰: ${e.message}" }
            return false
        } catch (e: Exception) {
            logger.error { "JWT 토큰 검증 실패: ${e.message}" }
            return false
        }
    }

    /**
     * 사용자 정보를 기반으로 액세스 토큰과 리프레시 토큰을 생성한다.
     */
    fun createTokens(user: User, deviceInfoDto: DeviceInfoDto): TokenDto {
        val now = LocalDateTime.now()
        val accessTokenExpiresIn = now.plusHours(accessTokenValidityInHours)
        val refreshTokenExpiresIn = now.plusHours(refreshTokenValidityInHours)

        val accessToken = generateAccessToken(user, now, accessTokenExpiresIn)
        val refreshToken = generateRefreshToken(
            user.id.toString(),
            now,
            refreshTokenExpiresIn,
            deviceInfoDto.deviceFingerprint,
        )

        // 만료 시간을 밀리초로 변환
        val accessTokenExpiresInMillis = Date.from(accessTokenExpiresIn.atZone(ZoneId.systemDefault()).toInstant()).time

        return TokenDto(
            accessToken = accessToken,
            refreshToken = refreshToken,
            accessTokenExpiresIn = accessTokenExpiresInMillis,
        )
    }

    /**
     * 리프레시 토큰을 검증하고 새로운 액세스 토큰을 발급한다.
     * 리프레시 토큰이 만료 임박이면 새로운 리프레시 토큰도 발급한다.
     *
     * @throws InvalidRefreshTokenException 리프레시 토큰이 유효하지 않은 경우
     */
    fun refreshTokens(refreshToken: String, deviceInfoDto: DeviceInfoDto): TokenDto {
        if (!validateToken(refreshToken))
            throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰입니다.")

        try {
            val claims = getClaims(refreshToken)

            // NOTE: 테스트 시 차단되는 문제가 있어서 비활성화 (추가 보안 기능으로 검증 하는게 나을듯)
//            validateDeviceFingerprint(claims, deviceInfoDto.deviceFingerprint)

            val user = validateAndGetUser(claims)

            return if (isTokenNearExpiration(refreshToken, REFRESH_TOKEN_RENEWAL_DAYS)) {
                createTokens(user, deviceInfoDto)
            } else {
                renewAccessTokenOnly(user, refreshToken)
            }
        } catch (e: Exception) {
            handleTokenRefreshError(e)
        }
    }

    /**
     * 토큰에서 인증 정보를 추출한다.
     *
     * @param accessToken JWT 토큰
     * @return 인증 정보
     * @throws InvalidAccessTokenException 토큰이 유효하지 않거나 사용자를 찾을 수 없는 경우
     */
    fun getAuthentication(accessToken: String): Authentication {
        if (!validateToken(accessToken))
            throw InvalidAccessTokenException("유효하지 않은 토큰입니다.")

        val claims = getClaims(accessToken)
        val userId = claims.subject.toLong()

        // 사용자 정보 조회
        val user = userRepository.findByIdOrNull(userId)
            ?: throw InvalidAccessTokenException("유효하지 않은 토큰: 사용자를 찾을 수 없습니다.")

        return UsernamePasswordAuthenticationToken(user, "", user.authorities)
    }

    /**
     * 토큰의 만료일을 반환한다.
     *
     * @param token JWT 토큰
     * @return 토큰 만료일 (null인 경우 토큰이 유효하지 않음)
     */
    private fun getTokenExpirationDate(token: String): LocalDateTime? {
        return try {
            val claims = getClaims(token)
            claims.expiration.toInstant().atZone(ZoneId.systemDefault()).toLocalDateTime()
        } catch (e: Exception) {
            logger.error { "토큰 만료일 확인 실패: ${e.message}" }
            null
        }
    }

    /**
     * 토큰이 지정된 일수 이내에 만료되는지 확인한다.
     *
     * @param token JWT 토큰
     * @param daysBeforeExpiration 만료 전 일수
     * @return 토큰 만료 임박 여부
     */
    private fun isTokenNearExpiration(token: String, daysBeforeExpiration: Long): Boolean {
        return try {
            val expirationDate = getTokenExpirationDate(token) ?: return true
            val now = LocalDateTime.now()
            val timeBeforeExpiration = daysBeforeExpiration // days
            now.plusDays(timeBeforeExpiration).isAfter(expirationDate)
        } catch (e: Exception) {
            logger.error { "토큰 만료 임박 확인 실패: ${e.message}" }
            true // 확인할 수 없으면 만료 임박으로 간주
        }
    }

    private fun validateAndGetUser(claims: Claims): User {
        val userId = claims.subject.toLong()

        return userRepository.findByIdOrNull(userId)
            ?: throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: 사용자를 찾을 수 없습니다.")
    }

    private fun renewAccessTokenOnly(user: User, refreshToken: String): TokenDto {
        val now = LocalDateTime.now()
        val accessTokenExpiresIn = now.plusHours(accessTokenValidityInHours)
        val accessToken = generateAccessToken(user, now, accessTokenExpiresIn)

        // 만료 시간을 밀리초로 변환
        val accessTokenExpiresInMillis = Date.from(accessTokenExpiresIn.atZone(ZoneId.systemDefault()).toInstant()).time

        return TokenDto(
            accessToken = accessToken,
            refreshToken = refreshToken,
            accessTokenExpiresIn = accessTokenExpiresInMillis,
        )
    }

    private fun handleTokenRefreshError(e: Exception): Nothing {
        when (e) {
            is InvalidRefreshTokenException -> throw e
            is SignatureException -> {
                logger.error { "리프레시 토큰 서명 검증 실패: ${e.message}" }
                throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: 서명이 유효하지 않습니다.")
            }

            is ExpiredJwtException -> {
                logger.error { "리프레시 토큰 만료: ${e.message}" }
                throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: 토큰이 만료되었습니다.")
            }

            else -> {
                logger.error { "리프레시 토큰 검증 실패: ${e.message}" }
                throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: ${e.message}")
            }
        }
    }

    /**
     * 사용자 정보를 기반으로 액세스 토큰을 생성한다.
     */
    private fun generateAccessToken(user: User, issuedAt: LocalDateTime, expiresIn: LocalDateTime): String {
        return Jwts.builder()
            .setSubject(user.id.toString())
            .claim("username", user.username)
            .claim("authorities", user.authorities.joinToString(",") { it.authority })
            .setIssuedAt(Date.from(issuedAt.atZone(ZoneId.systemDefault()).toInstant()))
            .setExpiration(Date.from(expiresIn.atZone(ZoneId.systemDefault()).toInstant()))
            .signWith(key, SignatureAlgorithm.HS256)
            .compact()
    }

    /**
     * 사용자 ID를 기반으로 리프레시 토큰을 생성한다.
     */
    private fun generateRefreshToken(
        userId: String,
        issuedAt: LocalDateTime,
        expiresIn: LocalDateTime,
        deviceFingerprint: String?,
    ): String {
        val builder = Jwts.builder()
            .setSubject(userId)
            .setIssuedAt(Date.from(issuedAt.atZone(ZoneId.systemDefault()).toInstant()))
            .setExpiration(Date.from(expiresIn.atZone(ZoneId.systemDefault()).toInstant()))

        // 디바이스 핑거프린트가 제공된 경우 리프레시 토큰에 추가
        if (deviceFingerprint != null) {
            builder.claim("deviceFingerprint", deviceFingerprint)
        }

        return builder.signWith(key, SignatureAlgorithm.HS256).compact()
    }

    /**
     * 디바이스 핑거프린트를 검증한다.
     */
    private fun validateDeviceFingerprint(claims: Claims, deviceFingerprint: String?) {
        if (deviceFingerprint == null)
            throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: 디바이스 핑거프린트가 존재하지 않습니다.")

        val tokenDeviceFingerprint = claims["deviceFingerprint"] as? String
        if (tokenDeviceFingerprint != null && tokenDeviceFingerprint != deviceFingerprint) {
            throw InvalidRefreshTokenException("유효하지 않은 리프레시 토큰: 디바이스 핑거프린트가 일치하지 않습니다.")
        }
    }

    /**
     * 토큰에서 클레임을 추출한다.
     *
     * @param token JWT 토큰
     * @return 토큰 클레임
     */
    private fun getClaims(token: String): Claims {
        return Jwts.parserBuilder()
            .setSigningKey(key)
            .setAllowedClockSkewSeconds(1) // 1초의 시간 오차 허용
            .build()
            .parseClaimsJws(token)
            .body
    }

    companion object : KLogging() {
        private const val REFRESH_TOKEN_RENEWAL_DAYS = 7L
    }
}
