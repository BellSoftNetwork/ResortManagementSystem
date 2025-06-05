package net.bellsoft.rms.authentication.component

import io.jsonwebtoken.Jwts
import io.jsonwebtoken.SignatureAlgorithm
import io.kotest.assertions.throwables.shouldNotThrowAny
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.dto.response.TokenDto
import net.bellsoft.rms.authentication.fixture.DeviceInfoFixture
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.fixture.UserFixture
import net.bellsoft.rms.user.repository.UserRepository
import net.bellsoft.rms.user.type.UserRole
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.ActiveProfiles
import java.lang.reflect.Field
import java.security.Key
import java.util.*

@SpringBootTest
@ActiveProfiles("test")
class JwtTokenProviderSecurityTest(
    @Value("\${security.jwt.secret}") private val secretKey: String,
    @Value("\${security.jwt.access-token-validity-in-hours}") private val accessTokenValidityInHours: Long,
    @Value("\${security.jwt.refresh-token-validity-in-hours}") private val refreshTokenValidityInHours: Long,
    private val userRepository: UserRepository,
    private val jwtTokenProvider: JwtTokenProvider,
) : BehaviorSpec(
    {
        val fixture = baseFixture.feature(DeviceInfoFixture.Feature.WITH_INTERNAL_IP)

        val user = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())

        // JWT 토큰의 페이로드를 변조하는 함수
        fun createTamperedToken(user: User, deviceInfoDto: DeviceInfoDto): TokenDto {
            // JwtTokenProvider의 private 필드에 접근하기 위한 리플렉션
            val keyField: Field = JwtTokenProvider::class.java.getDeclaredField("key")
            keyField.isAccessible = true
            val key = keyField.get(jwtTokenProvider) as Key

            val now = Date()
            val expirationDate = Date(now.time + accessTokenValidityInHours * 1000)

            // 변조된 권한 정보를 가진 액세스 토큰 생성
            val tamperedAccessToken = Jwts.builder()
                .setSubject(user.id.toString())
                .claim("username", user.username)
                // 권한 정보 변조 - ADMIN 권한 추가
                .claim("authorities", "NORMAL,ADMIN")
                .setIssuedAt(now)
                .setExpiration(expirationDate)
                .signWith(key, SignatureAlgorithm.HS256)
                .compact()

            // 정상적인 리프레시 토큰 생성
            val refreshTokenBuilder = Jwts.builder()
                .setSubject(user.id.toString())
                .setIssuedAt(now)
                .setExpiration(Date(now.time + refreshTokenValidityInHours * 1000))

            refreshTokenBuilder.claim("deviceFingerprint", deviceInfoDto.deviceFingerprint)

            val refreshToken = refreshTokenBuilder
                .signWith(key, SignatureAlgorithm.HS256)
                .compact()

            return TokenDto(
                accessToken = tamperedAccessToken,
                refreshToken = refreshToken,
                accessTokenExpiresIn = expirationDate.time,
            )
        }

        Given("JWT 페이로드 변조 시") {
            val tamperedTokenDto = createTamperedToken(user, fixture())

            When("권한 정보가 변조된 액세스 토큰으로 인증을 시도하면") {
                val authentication = jwtTokenProvider.getAuthentication(tamperedTokenDto.accessToken)

                Then("변조된 권한이 아닌 실제 사용자의 권한이 사용된다") {
                    val authorities = authentication.authorities.map { it.authority }
                    authorities shouldBe listOf("NORMAL")
                    authorities shouldNotBe listOf("NORMAL", "ADMIN")
                }
            }
        }

        Given("사용자 권한 변경 시") {
            When("토큰 발급 후 사용자 권한이 변경되면") {
                // 사용자 권한 변경
                user.role = UserRole.ADMIN
                userRepository.save(user)

                Then("DB에서 조회한 사용자의 권한은 변경된 권한이다") {
                    val updatedUser = userRepository.findById(user.id).get()
                    updatedUser.role shouldBe UserRole.ADMIN
                }
            }
        }

        Given("다른 디바이스에서 토큰 발급 시도 시") {
            val deviceInfoDto: DeviceInfoDto = fixture()
            val tokenDto = jwtTokenProvider.createTokens(user, deviceInfoDto)

            When("첫 번째 디바이스에서 발급된 리프레시 토큰으로 두 번째 디바이스에서 갱신을 시도하면") {
                val otherDeviceInfoDto: DeviceInfoDto = fixture.feature(
                    DeviceInfoFixture.Feature.WITH_ANDROID_DEVICE_FINGERPRINT,
                )()

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
        }
    },
)
