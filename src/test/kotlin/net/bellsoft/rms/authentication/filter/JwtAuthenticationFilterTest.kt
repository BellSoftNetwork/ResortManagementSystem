package net.bellsoft.rms.authentication.filter

import com.fasterxml.jackson.databind.ObjectMapper
import io.kotest.core.spec.style.BehaviorSpec
import jakarta.servlet.http.Cookie
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.common.util.extractJsonPath
import net.bellsoft.rms.common.util.jsonPathShouldExist
import net.bellsoft.rms.common.util.performAndReturn
import net.bellsoft.rms.common.util.shouldBeOk
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

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc
class JwtAuthenticationFilterTest(
    private val mockMvc: MockMvc,
    private val objectMapper: ObjectMapper,
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("JWT 인증 필터가 요청을 처리할 때") {
            // 테스트용 사용자 생성
            val password = "password123"

            val user: User = fixture()
            userRepository.save(user.apply { changePassword(passwordEncoder, password) })

            // 로그인하여 토큰 발급
            val loginResult = mockMvc.performAndReturn(
                post("/api/v1/auth/login")
                    .contentType(MediaType.APPLICATION_JSON)
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

            When("Authorization 헤더에 유효한 Bearer 토큰이 있으면") {
                Then("인증된 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Bearer $accessToken"),
                    ).shouldBeOk()
                }
            }

            When("Authorization 헤더가 없으면") {
                Then("인증되지 않은 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my"),
                    ).shouldBeUnauthorized()
                }
            }

            When("Authorization 헤더가 Bearer 형식이 아니면") {
                Then("인증되지 않은 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .header(HttpHeaders.AUTHORIZATION, "Token $accessToken"),
                    ).shouldBeUnauthorized()
                }
            }

            When("쿠키에 토큰이 있지만 헤더에는 없으면") {
                Then("인증되지 않은 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .cookie(Cookie("accessToken", accessToken)),
                    ).shouldBeUnauthorized()
                }
            }

            When("요청 파라미터에 토큰이 있지만 헤더에는 없으면") {
                Then("인증되지 않은 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        get("/api/v1/my")
                            .param("token", accessToken),
                    ).shouldBeUnauthorized()
                }
            }

            When("요청 바디에 토큰이 있지만 헤더에는 없으면") {
                Then("인증되지 않은 요청으로 처리된다") {
                    mockMvc.performAndReturn(
                        post("/api/v1/my")
                            .contentType(MediaType.APPLICATION_JSON)
                            .content("""{"token":"$accessToken"}"""),
                    ).shouldBeUnauthorized()
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
)
