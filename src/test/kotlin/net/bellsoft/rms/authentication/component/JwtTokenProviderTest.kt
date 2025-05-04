package net.bellsoft.rms.authentication.component

import io.jsonwebtoken.Jwts
import io.jsonwebtoken.SignatureAlgorithm
import io.jsonwebtoken.security.Keys
import io.kotest.assertions.throwables.shouldNotThrowAny
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.dto.response.TokenDto
import net.bellsoft.rms.authentication.exception.InvalidRefreshTokenException
import net.bellsoft.rms.authentication.fixture.DeviceInfoFixture
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.ActiveProfiles
import java.lang.reflect.Field
import java.security.Key
import java.util.*

@SpringBootTest
@ActiveProfiles("test")
class JwtTokenProviderTest(
    @Value("\${security.jwt.secret}") private val secretKey: String,
    @Value("\${security.jwt.access-token-validity-in-hours}") private val accessTokenValidityInHours: Long,
    @Value("\${security.jwt.refresh-token-validity-in-hours}") private val refreshTokenValidityInHours: Long,
    private val userRepository: UserRepository,
    private val jwtTokenProvider: JwtTokenProvider,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture.feature(DeviceInfoFixture.Feature.WITH_INTERNAL_IP)

        val deviceInfoDto: DeviceInfoDto = fixture()
        val user = userRepository.save(fixture())

        // 만료된 토큰 생성 함수
        fun createExpiredToken(user: User, deviceInfoDto: DeviceInfoDto): TokenDto {
            // JwtTokenProvider의 private 필드에 접근하기 위한 리플렉션
            val keyField: Field = JwtTokenProvider::class.java.getDeclaredField("key")
            keyField.isAccessible = true
            val key = keyField.get(jwtTokenProvider) as Key

            val now = Date()
            val pastDate = Date(now.time - 1000) // 1초 전

            // 만료된 액세스 토큰 생성
            val expiredAccessToken = Jwts.builder()
                .setSubject(user.id.toString())
                .claim("username", user.username)
                .claim("authorities", user.authorities.joinToString(",") { it.authority })
                .setIssuedAt(pastDate)
                .setExpiration(pastDate)
                .signWith(key, SignatureAlgorithm.HS256)
                .compact()

            // 만료된 리프레시 토큰 생성
            val refreshTokenBuilder = Jwts.builder()
                .setSubject(user.id.toString())
                .setIssuedAt(pastDate)
                .setExpiration(pastDate)

            refreshTokenBuilder.claim("deviceFingerprint", deviceInfoDto.deviceFingerprint)

            val expiredRefreshToken = refreshTokenBuilder
                .signWith(key, SignatureAlgorithm.HS256)
                .compact()

            return TokenDto(
                accessToken = expiredAccessToken,
                refreshToken = expiredRefreshToken,
                accessTokenExpiresIn = pastDate.time,
            )
        }

        // 유효하지 않은 서명의 토큰 생성 함수
        fun createInvalidSignatureToken(user: User, deviceInfoDto: DeviceInfoDto): TokenDto {
            // 다른 키로 서명된 토큰 생성
            val differentKey = Keys.hmacShaKeyFor("different-secret-key-for-testing-invalid-signature".toByteArray())

            val now = Date()
            val expirationDate = Date(now.time + accessTokenValidityInHours * 1000)

            // 다른 키로 서명된 액세스 토큰
            val invalidAccessToken = Jwts.builder()
                .setSubject(user.id.toString())
                .claim("username", user.username)
                .claim("authorities", user.authorities.joinToString(",") { it.authority })
                .setIssuedAt(now)
                .setExpiration(expirationDate)
                .signWith(differentKey, SignatureAlgorithm.HS256)
                .compact()

            // 다른 키로 서명된 리프레시 토큰
            val refreshTokenBuilder = Jwts.builder()
                .setSubject(user.id.toString())
                .setIssuedAt(now)
                .setExpiration(Date(now.time + refreshTokenValidityInHours * 1000))

            refreshTokenBuilder.claim("deviceFingerprint", deviceInfoDto.deviceFingerprint)

            val invalidRefreshToken = refreshTokenBuilder
                .signWith(differentKey, SignatureAlgorithm.HS256)
                .compact()

            return TokenDto(
                accessToken = invalidAccessToken,
                refreshToken = invalidRefreshToken,
                accessTokenExpiresIn = expirationDate.time,
            )
        }

        Given("JWT 토큰 생성 시") {
            When("유효한 사용자 정보와 디바이스 핑거프린트로 토큰을 생성하면") {
                val tokenDto = jwtTokenProvider.createTokens(user, deviceInfoDto)

                Then("액세스 토큰과 리프레시 토큰이 생성된다") {
                    tokenDto.accessToken shouldNotBe null
                    tokenDto.refreshToken shouldNotBe null
                    tokenDto.accessTokenExpiresIn shouldNotBe null
                }

                Then("생성된 토큰은 유효하다") {
                    jwtTokenProvider.validateToken(tokenDto.accessToken) shouldBe true
                    jwtTokenProvider.validateToken(tokenDto.refreshToken) shouldBe true
                }
            }
        }

        Given("디바이스 핑거프린트 변경 시") {
            val tokenDto = jwtTokenProvider.createTokens(user, deviceInfoDto)
            val otherDeviceInfoDto: DeviceInfoDto = fixture
                .feature(DeviceInfoFixture.Feature.WITH_ANDROID_DEVICE_FINGERPRINT)()

            When("다른 디바이스 핑거프린트로 리프레시 토큰을 사용하면") {
                Then("이상 없이 갱신에 성공한다") {
                    shouldNotThrowAny {
                        jwtTokenProvider.refreshTokens(tokenDto.refreshToken, otherDeviceInfoDto)
                    }
                }

                // NOTE: 추가 보안 기능 구현 시 별도 검증 필요
//                Then("InvalidRefreshTokenException이 발생한다") {
//                    shouldThrow<InvalidRefreshTokenException> {
//                        jwtTokenProvider.refreshTokens(tokenDto.refreshToken, otherDeviceInfoDto)
//                    }
//                }
            }

            When("디바이스 핑거프린트 없이 리프레시 토큰을 사용하면") {
                val nonFingerprintDevice: DeviceInfoDto = fixture
                    .feature(DeviceInfoFixture.Feature.WITHOUT_DEVICE_FINGERPRINT)()

                Then("이상 없이 갱신에 성공한다") {
                    shouldNotThrowAny {
                        jwtTokenProvider.refreshTokens(tokenDto.refreshToken, nonFingerprintDevice)
                    }
                }

                // NOTE: 추가 보안 기능 구현 시 별도 검증 필요
//                Then("InvalidRefreshTokenException이 발생한다") {
//                    shouldThrow<InvalidRefreshTokenException> {
//                        jwtTokenProvider.refreshTokens(tokenDto.refreshToken, nonFingerprintDevice)
//                    }
//                }
            }
        }

        Given("토큰 서명이 유효하지 않을 때") {
            val invalidTokenDto = createInvalidSignatureToken(user, deviceInfoDto)

            When("유효하지 않은 서명의 액세스 토큰으로 인증을 시도하면") {
                Then("토큰 검증에 실패한다") {
                    jwtTokenProvider.validateToken(invalidTokenDto.accessToken) shouldBe false
                }
            }

            When("유효하지 않은 서명의 리프레시 토큰으로 갱신을 시도하면") {
                Then("InvalidRefreshTokenException이 발생한다") {
                    shouldThrow<InvalidRefreshTokenException> {
                        jwtTokenProvider.refreshTokens(invalidTokenDto.refreshToken, deviceInfoDto)
                    }
                }
            }
        }

        Given("토큰이 만료되었을 때") {
            val expiredTokenDto = createExpiredToken(user, deviceInfoDto)

            When("만료된 액세스 토큰으로 인증을 시도하면") {
                Then("토큰 검증에 실패한다") {
                    jwtTokenProvider.validateToken(expiredTokenDto.accessToken) shouldBe false
                }
            }

            When("만료된 리프레시 토큰으로 갱신을 시도하면") {
                Then("InvalidRefreshTokenException이 발생한다") {
                    shouldThrow<InvalidRefreshTokenException> {
                        jwtTokenProvider.refreshTokens(expiredTokenDto.refreshToken, deviceInfoDto)
                    }
                }
            }
        }

        Given("토큰에서 인증 정보를 추출할 때") {
            val tokenDto = jwtTokenProvider.createTokens(user, deviceInfoDto)

            When("유효한 액세스 토큰으로 인증 정보를 추출하면") {
                val authentication = jwtTokenProvider.getAuthentication(tokenDto.accessToken)

                Then("사용자 정보가 포함된 인증 객체가 반환된다") {
                    authentication.isAuthenticated shouldBe true
                    val principal = authentication.principal as User
                    principal.id shouldBe user.id
                    principal.userId shouldBe user.userId
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
)
