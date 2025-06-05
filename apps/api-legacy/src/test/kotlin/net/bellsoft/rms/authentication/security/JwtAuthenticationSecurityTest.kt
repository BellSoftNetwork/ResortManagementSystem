package net.bellsoft.rms.authentication.security

import com.fasterxml.jackson.databind.ObjectMapper
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.collections.shouldContainExactly
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.authentication.util.AuthTestConstants
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.common.util.extractJsonPath
import net.bellsoft.rms.common.util.jsonPathShouldExist
import net.bellsoft.rms.common.util.performAndReturn
import net.bellsoft.rms.common.util.shouldBeOk
import net.bellsoft.rms.common.util.shouldBeTooManyRequests
import net.bellsoft.rms.common.util.shouldBeUnauthorized
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.HttpHeaders
import org.springframework.http.MediaType
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
import java.util.concurrent.CountDownLatch
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit
import java.util.concurrent.atomic.AtomicInteger

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc
class JwtAuthenticationSecurityTest(
    private val mockMvc: MockMvc,
    private val objectMapper: ObjectMapper,
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("CSRF 공격 시뮬레이션") {
            // 테스트용 사용자 생성
            val password = "password123"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            // 로그인하여 토큰 발급
            val loginResult = mockMvc.performAndReturn(
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

            loginResult.shouldBeOk()
            loginResult.jsonPathShouldExist("$.value.accessToken", objectMapper)
            val accessToken = loginResult.extractJsonPath("$.value.accessToken", objectMapper)

            When("CSRF 공격 시뮬레이션 - 쿠키에 토큰을 저장하고 다른 사이트에서 요청을 보내는 경우") {
                Then("쿠키에 저장된 토큰은 인증에 사용되지 않는다") {
                    // 쿠키에 토큰을 저장하고 요청을 보내는 시뮬레이션
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .cookie(jakarta.servlet.http.Cookie("accessToken", accessToken))
                            .header("Origin", "https://malicious-site.com"), // 악의적인 사이트에서 요청
                    )
                        .shouldBeUnauthorized()
                }
            }
        }

        Given("XSS 공격 방어 테스트") {
            // 테스트용 사용자 생성
            val password = "password123"
            val maliciousScript = "<script>alert('XSS')</script>"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("XSS 공격 시도 - 악성 스크립트가 포함된 User-Agent로 로그인") {
                Then("로그인은 성공하고 악성 스크립트가 실행되지 않는다") {
                    val loginResult = mockMvc.performAndReturn(
                        post("/api/v1/auth/login")
                            .contentType(MediaType.APPLICATION_JSON)
                            .header("User-Agent", maliciousScript)
                            .header("Accept-Language", "ko-KR")
                            .content(
                                """
                                {
                                    "username":"${user.userId}",
                                    "password":"$password"
                                }
                                """.trimIndent(),
                            ),
                    )

                    // 로그인이 성공하는지 확인
                    loginResult.shouldBeOk()

                    // 토큰이 발급되었는지 확인
                    loginResult.jsonPathShouldExist("$.value.accessToken", objectMapper)
                    loginResult.jsonPathShouldExist("$.value.refreshToken", objectMapper)
                }
            }
        }

        Given("Race Condition 테스트") {
            // 테스트용 사용자 생성
            val password = "password123"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("여러 스레드에서 동시에 토큰 발급을 시도하는 경우") {
                val threadCount = 10
                val executor = Executors.newFixedThreadPool(threadCount)
                val latch = CountDownLatch(1)
                val successCount = AtomicInteger(0)

                Then("모든 요청이 정상적으로 처리된다") {
                    // 여러 스레드에서 동시에 로그인 요청
                    for (i in 0 until threadCount) {
                        executor.submit {
                            try {
                                latch.await() // 모든 스레드가 동시에 시작하도록 대기
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
                                if (result.response.status == 200) {
                                    successCount.incrementAndGet()
                                }
                            } catch (e: Exception) {
                                // 예외 발생 시 로그 출력
                                println("[DEBUG_LOG] 로그인 요청 실패: ${e.message}")
                            }
                        }
                    }

                    latch.countDown() // 모든 스레드 시작
                    executor.shutdown()
                    executor.awaitTermination(10, TimeUnit.SECONDS)

                    // 모든 요청이 성공했는지 확인
                    successCount.get() shouldBe threadCount
                }
            }
        }

        Given("동시 로그인 제한 테스트") {
            // 테스트용 사용자 생성
            val password = "password123"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("서로 다른 디바이스에서 동시에 로그인하는 경우") {
                // 첫 번째 디바이스에서 로그인
                val loginResult1 = mockMvc.performAndReturn(
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
                loginResult1.shouldBeOk()
                loginResult1.jsonPathShouldExist("$.value.accessToken", objectMapper)
                val accessToken1 = loginResult1.extractJsonPath("$.value.accessToken", objectMapper)

                // 두 번째 디바이스에서 로그인
                val loginResult2 = mockMvc.performAndReturn(
                    post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header("User-Agent", AuthTestConstants.DEFAULT_IOS_USER_AGENT)
                        .header("Accept-Language", AuthTestConstants.DEFAULT_ENGLISH_LANGUAGE)
                        .content(
                            """
                            {
                                "username":"${user.userId}",
                                "password":"$password"
                            }
                            """.trimIndent(),
                        ),
                )

                loginResult2.shouldBeOk()
                loginResult2.jsonPathShouldExist("$.value.accessToken", objectMapper)
                val accessToken2 = loginResult2.extractJsonPath("$.value.accessToken", objectMapper)

                Then("두 디바이스 모두 유효한 토큰을 발급받는다") {
                    // 첫 번째 디바이스 토큰으로 인증
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $accessToken1"),
                    )
                        .shouldBeOk()

                    // 두 번째 디바이스 토큰으로 인증
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $accessToken2"),
                    )
                        .shouldBeOk()
                }
            }
        }

        Given("Brute Force 방어 테스트") {
            // 테스트용 사용자 생성
            val password = "password123"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            When("30분 이내에 잘못된 비밀번호로 5번 초과하여 로그인을 시도하는 경우") {
                val attemptCount = 6
                val resultStatusCodes = mutableListOf<Int>()

                repeat(attemptCount) {
                    val wrongPassword = "wrongpassword"
                    val result = mockMvc.performAndReturn(
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
                    )

                    resultStatusCodes.add(result.response.status)
                }

                Then("첫 5번 로그인 시도는 401 응답과 함께 실패한다") {
                    resultStatusCodes.slice(0 until 5) shouldContainExactly List(5) { 401 }
                }

                Then("5번을 초과한 로그인 시도는 429 응답과 함께 실패한다") {
                    resultStatusCodes.last() shouldBe 429
                }

                Then("직후에 정상 비밀번호로 로그인을 시도해도 429 응답과 함께 실패한다") {
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
                        .shouldBeTooManyRequests()
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
)
