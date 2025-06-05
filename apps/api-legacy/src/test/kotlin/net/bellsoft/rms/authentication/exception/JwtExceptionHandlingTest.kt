package net.bellsoft.rms.authentication.exception

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.authentication.component.JwtTokenProvider
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.common.util.performAndReturn
import net.bellsoft.rms.common.util.shouldBeUnauthorized
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.HttpHeaders
import org.springframework.http.MediaType
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
import java.lang.reflect.Field
import java.security.Key
import java.util.*

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc
class JwtExceptionHandlingTest(
    private val mockMvc: MockMvc,
    private val userRepository: UserRepository,
    private val jwtTokenProvider: JwtTokenProvider,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("토큰 형식 오류 처리") {
            When("잘못된 형식의 JWT 토큰으로 인증을 시도하면") {
                val malformedToken = "invalid.jwt.token"

                Then("토큰 검증이 실패하고 false를 반환한다") {
                    jwtTokenProvider.validateToken(malformedToken) shouldBe false
                }

                Then("API 요청 시 401 Unauthorized 응답을 반환한다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $malformedToken"),
                    )
                        .shouldBeUnauthorized()
                }
            }

            When("지원하지 않는 JWT 토큰으로 인증을 시도하면") {
                // 헤더와 페이로드만 있고 서명이 없는 토큰
                val unsupportedToken = "eyJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIn0"

                Then("토큰 검증이 실패하고 false를 반환한다") {
                    jwtTokenProvider.validateToken(unsupportedToken) shouldBe false
                }

                Then("API 요청 시 401 Unauthorized 응답을 반환한다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $unsupportedToken"),
                    )
                        .shouldBeUnauthorized()
                }
            }
        }

        Given("토큰 서명 오류 처리") {
            When("잘못된 서명의 JWT 토큰으로 인증을 시도하면") {
                // 유효한 형식이지만 서명이 다른 토큰
                val invalidSignatureToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
                    "eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
                    "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

                Then("토큰 검증이 실패하고 false를 반환한다") {
                    jwtTokenProvider.validateToken(invalidSignatureToken) shouldBe false
                }

                Then("API 요청 시 401 Unauthorized 응답을 반환한다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $invalidSignatureToken"),
                    )
                        .shouldBeUnauthorized()
                }
            }
        }

        Given("토큰 만료 오류 처리") {
            // 테스트용 사용자 생성
            val user = userRepository.save(fixture())

            // 만료된 토큰 생성
            val expiredToken = createExpiredToken(user, jwtTokenProvider)

            When("만료된 JWT 토큰으로 인증을 시도하면") {
                Then("토큰 검증이 실패하고 false를 반환한다") {
                    jwtTokenProvider.validateToken(expiredToken) shouldBe false
                }

                Then("API 요청 시 401 Unauthorized 응답을 반환한다") {
                    val result = mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $expiredToken"),
                    )

                    result.shouldBeUnauthorized()
                }
            }
        }

        Given("리프레시 토큰 예외 처리") {
            val deviceInfoDto: DeviceInfoDto = fixture()

            When("유효하지 않은 리프레시 토큰으로 갱신을 시도하면") {
                val invalidRefreshToken = "invalid.refresh.token"
                val deviceFingerprint = "test-device-fingerprint"

                Then("InvalidRefreshTokenException이 발생한다") {
                    shouldThrow<InvalidRefreshTokenException> {
                        jwtTokenProvider.refreshTokens(invalidRefreshToken, deviceInfoDto)
                    }
                }

                Then("API 요청 시 401 Unauthorized 응답을 반환한다") {
                    mockMvc.performAndReturn(
                        post("/api/v1/auth/refresh")
                            .contentType(MediaType.APPLICATION_JSON)
                            .content(
                                """
                                {
                                    "refreshToken":"$invalidRefreshToken",
                                    "deviceFingerprint":"$deviceFingerprint"
                                }
                                """.trimIndent(),
                            ),
                    )
                        .shouldBeUnauthorized()
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
) {
    companion object {
        // 만료된 토큰 생성 함수
        private fun createExpiredToken(user: User, jwtTokenProvider: JwtTokenProvider): String {
            // JwtTokenProvider의 private 필드에 접근하기 위한 리플렉션
            val keyField: Field = JwtTokenProvider::class.java.getDeclaredField("key")
            keyField.isAccessible = true
            val key = keyField.get(jwtTokenProvider) as Key

            val now = Date()
            val pastDate = Date(now.time - 1000) // 1초 전

            // 만료된 액세스 토큰 생성
            return io.jsonwebtoken.Jwts.builder()
                .setSubject(user.id.toString())
                .claim("username", user.username)
                .claim("authorities", user.authorities.joinToString(",") { it.authority })
                .setIssuedAt(pastDate)
                .setExpiration(pastDate)
                .signWith(key, io.jsonwebtoken.SignatureAlgorithm.HS256)
                .compact()
        }
    }
}
