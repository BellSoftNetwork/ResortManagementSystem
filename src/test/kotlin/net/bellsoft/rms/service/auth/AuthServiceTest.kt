package net.bellsoft.rms.service.auth

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.util.TestDatabaseSupport
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class AuthServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val authService: AuthService,
) : BehaviorSpec(
    {
        val fixture = baseFixture.new {
            property(UserRegistrationRequest::password) { "password" }
        }

        Given("가입한 사용자가 없는 상황에서") {
            When("신규 회원 가입 시도 시") {
                val userRegistrationRequest: UserRegistrationRequest = fixture {
                    property(UserRegistrationRequest::email) { "userId@mail.com" }
                }
                val user = authService.register(userRegistrationRequest)

                Then("정상적으로 가입된다") {
                    user.email shouldBe userRegistrationRequest.email
                }
            }

            When("존재하지 않는 계정 아이디로 유저 로드 시도 시") {
                val email = "NOT_EXISTS_USER_ID@mail.com"

                Then("유저 로드에 실패한다") {
                    shouldThrow<UsernameNotFoundException> {
                        authService.loadUserByUsername(email)
                    }.message shouldBe "$email 은 존재하지 않는 사용자입니다"
                }
            }
        }

        Given("기존에 가입한 사용자가 있는 상황에서") {
            val userRegistrationRequest: UserRegistrationRequest = fixture()
            authService.register(userRegistrationRequest)

            When("기존 사용자 ID 와 동일한 ID 로 가입 요청 시") {
                Then("가입이 거부된다") {
                    shouldThrow<DataIntegrityViolationException> {
                        authService.register(userRegistrationRequest)
                    }
                }
            }

            When("기존 사용자 ID 와 다른 ID 로 가입 요청 시") {
                val newUserRegistrationRequest: UserRegistrationRequest = fixture {
                    property(UserRegistrationRequest::email) { "new@mail.com" }
                }
                val newUser = authService.register(newUserRegistrationRequest)

                Then("정상적으로 가입된다") {
                    newUser.email shouldBe newUserRegistrationRequest.email
                }
            }

            When("유효한 계정 아이디로 유저 로드 시도 시") {
                val email = userRegistrationRequest.email

                Then("유저 엔티티가 정상적으로 로드된다") {
                    authService.loadUserByUsername(email).username shouldBe email
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
