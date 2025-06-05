package net.bellsoft.rms.authentication.controller

import com.fasterxml.jackson.databind.ObjectMapper
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.authentication.entity.LoginAttempt
import net.bellsoft.rms.authentication.fixture.LoginAttemptFixture
import net.bellsoft.rms.authentication.repository.LoginAttemptRepository
import net.bellsoft.rms.authentication.util.AuthTestConstants
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.common.util.jsonPathShouldExist
import net.bellsoft.rms.common.util.performAndReturn
import net.bellsoft.rms.common.util.shouldBeOk
import net.bellsoft.rms.common.util.shouldBeTooManyRequests
import net.bellsoft.rms.common.util.shouldBeUnauthorized
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.MediaType
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc
class AuthControllerBruteForceTest(
    private val mockMvc: MockMvc,
    private val objectMapper: ObjectMapper,
    private val userRepository: UserRepository,
    private val loginAttemptRepository: LoginAttemptRepository,
    private val passwordEncoder: PasswordEncoder,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("브루트 포스 공격 방어 테스트") {
            // 테스트용 사용자 생성
            val password = "password123"
            val wrongPassword = "wrongpassword"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })
            val loginAttemptFixture = fixture.new {
                property(LoginAttempt::username) { user.userId }
            }

            When("여러 번 잘못된 비밀번호로 로그인을 시도하면") {
                // 5번의 실패한 로그인 시도
                repeat(5) {
                    mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("X-Forwarded-For", "127.0.0.1")
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(
                                """
                                {
                                    "username":"${user.userId}",
                                    "password":"$wrongPassword"
                                }
                                """.trimIndent(),
                            ),
                    ).shouldBeUnauthorized()
                }

                Then("로그인 시도가 데이터베이스에 기록된다") {
                    val attempts = loginAttemptRepository.findAll()
                    attempts.size shouldBe 5
                    attempts.all { !it.successful } shouldBe true
                }

                Then("추가 로그인 시도는 429 Too Many Requests 응답을 반환한다") {
                    mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(
                                """
                                {
                                    "username":"${user.userId}",
                                    "password":"$wrongPassword"
                                }
                                """.trimIndent(),
                            ),
                    ).shouldBeTooManyRequests()
                }
            }

            When("시간이 지나 이전 실패 시도가 윈도우를 벗어나면") {
                // 과거의 실패한 로그인 시도 기록 (윈도우 밖)
                repeat(5) {
                    loginAttemptRepository.save(
                        loginAttemptFixture
                            .feature(LoginAttemptFixture.Feature.PAST_FAILED)
                            .feature(LoginAttemptFixture.Feature.WITH_SPECIFIC_IP)<LoginAttempt>(),
                    )
                }

                Then("로그인 시도가 허용된다") {
                    val result = mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                            .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                            .content(
                                """
                                {
                                    "username":"${user.userId}",
                                    "password":"$password"
                                }
                                """.trimIndent(),
                            ),
                    )

                    result.shouldBeOk()
                    result.jsonPathShouldExist("$.value.accessToken", objectMapper)
                }
            }

            When("OS가 변경되면") {
                // Windows에서 성공한 로그인 시도
                mockMvc.performAndReturn(
                    post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header("User-Agent", AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT)
                        .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                        .content(
                            """
                            {
                                "username":"${user.userId}",
                                "password":"$password"
                            }
                            """.trimIndent(),
                        ),
                )
                    .shouldBeOk()

                // Android에서 로그인 시도 (다른 OS)
                mockMvc.performAndReturn(
                    post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header("User-Agent", AuthTestConstants.DEFAULT_ANDROID_USER_AGENT)
                        .header("Accept-Language", AuthTestConstants.DEFAULT_KOREAN_LANGUAGE)
                        .content(
                            """
                            {
                                "username":"${user.userId}",
                                "password":"$password"
                            }
                            """.trimIndent(),
                        ),
                ).shouldBeOk()

                Then("로그인은 성공하지만 디바이스 변경이 감지된다") {
                    // 로그에서 디바이스 변경 감지 메시지를 확인할 수 있음
                    val attempts = loginAttemptRepository.findAll()
                    attempts.size shouldBe 2

                    // OS 정보가 올바르게 저장되었는지 확인
                    val firstAttempt = attempts.find { it.osInfo == "Windows" }
                    val secondAttempt = attempts.find { it.osInfo == "Android" }

                    firstAttempt shouldNotBe null
                    secondAttempt shouldNotBe null
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
)
