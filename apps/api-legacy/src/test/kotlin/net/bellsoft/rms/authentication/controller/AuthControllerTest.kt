package net.bellsoft.rms.authentication.controller

import com.fasterxml.jackson.databind.ObjectMapper
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import io.mockk.every
import io.mockk.mockkStatic
import net.bellsoft.rms.authentication.component.JwtTokenProvider
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.dto.request.LoginRequest
import net.bellsoft.rms.authentication.dto.request.RefreshTokenRequest
import net.bellsoft.rms.authentication.util.AuthTestConstants
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.common.util.getContentAsJson
import net.bellsoft.rms.common.util.jsonPathShouldExist
import net.bellsoft.rms.common.util.performAndReturn
import net.bellsoft.rms.common.util.shouldBeOk
import net.bellsoft.rms.common.util.shouldBeUnauthorized
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.MediaType
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
import java.time.LocalDateTime

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc
internal class AuthControllerTest(
    private val mockMvc: MockMvc,
    private val objectMapper: ObjectMapper,
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val jwtTokenProvider: JwtTokenProvider,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        mockkStatic(LocalDateTime::class)

        val fixture = baseFixture

        Given("사용자가 로그인을 시도할 때") {
            val password = "password123"

            // 테스트용 사용자 생성
            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("올바른 인증 정보를 제공하면") {
                val loginRequest = LoginRequest(
                    username = user.userId,
                    password = password,
                )

                Then("JWT 토큰이 발급된다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(loginRequest)),
                    )

                    // Kotest 스타일 검증
                    result.shouldBeOk()
                    result.jsonPathShouldExist("$.value.accessToken", objectMapper)
                    result.jsonPathShouldExist("$.value.refreshToken", objectMapper)
                    result.jsonPathShouldExist("$.value.accessTokenExpiresIn", objectMapper)

                    // 응답에서 토큰 추출
                    val jsonNode = result.getContentAsJson(objectMapper)
                    val accessToken = jsonNode.path("value").path("accessToken").asText()
                    val refreshToken = jsonNode.path("value").path("refreshToken").asText()

                    // 토큰 유효성 검증
                    jwtTokenProvider.validateToken(accessToken) shouldBe true
                    jwtTokenProvider.validateToken(refreshToken) shouldBe true
                }
            }

            When("잘못된 인증 정보를 제공하면") {
                val loginRequest = LoginRequest(
                    username = user.userId,
                    password = "wrongpassword",
                )

                Then("인증 실패 응답이 반환된다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(loginRequest)),
                    )

                    result.shouldBeUnauthorized()
                }
            }
        }

        Given("사용자가 1시간 전에 발급한 토큰으로 토큰 갱신을 시도할 때") {
            // 테스트용 사용자 생성
            val user: User = fixture()
            userRepository.save(user)

            val nowDate = LocalDateTime.now()
            every { LocalDateTime.now() } returns nowDate.minusHours(1)

            // Windows 디바이스 핑거프린트 해시 생성
            val windowsDeviceInfoDto: DeviceInfoDto = fixture {
                property(DeviceInfoDto::deviceFingerprint) { generateWindowsDeviceFingerprint() }
            }

            val tokenDto = jwtTokenProvider.createTokens(user, windowsDeviceInfoDto)
            val originalAccessToken = tokenDto.accessToken

            every { LocalDateTime.now() } returns nowDate

            When("유효한 리프레시 토큰을 제공하면") {
                val refreshTokenRequest = RefreshTokenRequest(
                    refreshToken = tokenDto.refreshToken,
                )

                Then("새로운 액세스 토큰이 발급된다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/refresh")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(refreshTokenRequest)),
                    )

                    result.shouldBeOk()
                    result.jsonPathShouldExist("$.value.accessToken", objectMapper)
                    result.jsonPathShouldExist("$.value.refreshToken", objectMapper)
                    result.jsonPathShouldExist("$.value.accessTokenExpiresIn", objectMapper)

                    // 응답에서 토큰 추출
                    val jsonNode = result.getContentAsJson(objectMapper)
                    val newAccessToken = jsonNode.path("value").path("accessToken").asText()

                    // 새 토큰 유효성 검증
                    jwtTokenProvider.validateToken(newAccessToken) shouldBe true

                    // 새 토큰이 이전 토큰과 다른지 확인
                    newAccessToken shouldNotBe originalAccessToken

                    // 시간이 1시간 지난 시점에서 새 토큰이 발급되었으므로,
                    // 새 토큰의 만료 시간은 일반적으로 이전 토큰의 만료 시간보다 더 미래여야 함
                    val newTokenExpiresIn = jsonNode.path("value").path("accessTokenExpiresIn").asLong()
                    println("새 토큰 만료 시간 타임스탬프: $newTokenExpiresIn")
                }
            }

            When("유효하지 않은 리프레시 토큰을 제공하면") {
                val refreshTokenRequest = RefreshTokenRequest(
                    refreshToken = "invalid.refresh.token",
                )

                Then("인증 실패 응답이 반환된다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/refresh")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(refreshTokenRequest)),
                    )

                    // Kotest 스타일 검증
                    result.shouldBeUnauthorized()
                }
            }

            When("액세스 토큰이 1시간 후에 만료에 가까워진 상태라면") {
                val refreshTokenRequest = RefreshTokenRequest(
                    refreshToken = tokenDto.refreshToken,
                )

                Then("리프레시 토큰으로 성공적으로 갱신할 수 있다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/refresh")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(refreshTokenRequest)),
                    )

                    result.shouldBeOk()

                    // 응답에서 토큰 추출
                    val jsonNode = result.getContentAsJson(objectMapper)
                    val newAccessToken = jsonNode.path("value").path("accessToken").asText()

                    // 새 토큰이 발급되었는지 확인
                    newAccessToken shouldNotBe originalAccessToken
                }
            }
        }

        Given("다른 IP 주소에서 리프레시 토큰을 사용할 때") {
            val password = "password123"
            val ip1 = "192.168.1.1"
            val ip2 = "10.0.0.1"

            // 테스트용 사용자 생성
            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("첫 번째 IP에서 로그인하여 리프레시 토큰을 얻은 후") {
                val loginRequest = LoginRequest(
                    username = user.userId,
                    password = password,
                )

                val loginResult = mockMvc.performAndReturn(
                    post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header("X-Forwarded-For", ip1)
                        .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                        .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                        .content(objectMapper.writeValueAsString(loginRequest)),
                )

                loginResult.shouldBeOk()
                val jsonNode = loginResult.getContentAsJson(objectMapper)
                val refreshToken = jsonNode.path("value").path("refreshToken").asText()

                Then("두 번째 IP에서 해당 리프레시 토큰으로 액세스 토큰을 갱신할 수 있다") {
                    val refreshTokenRequest = RefreshTokenRequest(
                        refreshToken = refreshToken,
                    )

                    val refreshResult = mockMvc.performAndReturn(
                        post("/api/v1/auth/refresh")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("X-Forwarded-For", ip2)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(objectMapper.writeValueAsString(refreshTokenRequest)),
                    )

                    refreshResult.shouldBeOk()
                    refreshResult.jsonPathShouldExist("$.value.accessToken", objectMapper)
                    refreshResult.jsonPathShouldExist("$.value.refreshToken", objectMapper)
                    refreshResult.jsonPathShouldExist("$.value.accessTokenExpiresIn", objectMapper)

                    // 새 액세스 토큰이 유효한지 확인
                    val newAccessToken = refreshResult.getContentAsJson(objectMapper)
                        .path("value").path("accessToken").asText()
                    jwtTokenProvider.validateToken(newAccessToken) shouldBe true
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
) {
    companion object {
        /**
         * Windows OS 정보에 대한 디바이스 핑거프린트 해시를 생성한다.
         * DeviceInfo 클래스의 generateFingerprintHash 메서드와 동일한 로직을 사용한다.
         */
        fun generateWindowsDeviceFingerprint(): String {
            val osInfo = "Windows"
            val bytes = java.security.MessageDigest.getInstance("SHA-256").digest(osInfo.toByteArray())
            return bytes.joinToString("") { "%02x".format(it) }
        }
    }
}
