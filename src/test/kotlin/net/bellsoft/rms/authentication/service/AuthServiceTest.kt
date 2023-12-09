package net.bellsoft.rms.authentication.service

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class AuthServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val authService: AuthService,
    private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("가입한 사용자가 없는 상황에서") {
            When("존재하지 않는 계정 아이디로 유저 로드 시도 시") {
                val userId = "NOT_EXISTS_USER_ID"

                Then("유저 로드에 실패한다") {
                    shouldThrow<UsernameNotFoundException> {
                        authService.loadUserByUsername(userId)
                    }.message shouldBe "$userId 은 존재하지 않는 사용자입니다"
                }
            }
        }

        Given("기존에 가입한 사용자가 있는 상황에서") {
            userRepository.save(
                fixture {
                    property(User::userId) { "userId" }
                    property(User::email) { "userId@mail.com" }
                },
            )

            When("유효한 계정 아이디로 유저 로드 시도 시") {
                Then("유저 엔티티가 정상적으로 로드된다") {
                    authService.loadUserByUsername("userId@mail.com").username shouldBe "userId"
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
