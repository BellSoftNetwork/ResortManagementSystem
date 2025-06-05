package net.bellsoft.rms.user.repository

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import io.mockk.every
import io.mockk.mockkStatic
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseNullFixture
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.type.UserStatus
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDateTime

@SpringBootTest
@ActiveProfiles("test")
internal class UserTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        val fixture = baseNullFixture

        mockkStatic(LocalDateTime::class)
        every { LocalDateTime.now() } returns TEST_LOCAL_DATE_TIME

        Given("유저가 없는 상황에서") {
            When("최소 파라미터로 유저 엔티티를 생성하면") {
                val originalUser: User = fixture()
                val createdUser = userRepository.save(originalUser)

                Then("설정한 값이 정상적으로 등록된다") {
                    assertSoftly {
                        createdUser.password shouldBe originalUser.password
                        createdUser.email shouldBe originalUser.email
                        createdUser.name shouldBe originalUser.name
                    }
                }

                Then("값을 제공하지 않은 칼럼에 기본 값이 등록된다") {
                    assertSoftly {
                        createdUser.id shouldNotBe null
                        createdUser.status shouldBe UserStatus.INACTIVE
                    }
                }

                Then("생성 및 수정 시각이 현재로 등록된다") {
                    assertSoftly {
                        createdUser.createdAt shouldBe TEST_LOCAL_DATE_TIME
                        createdUser.updatedAt shouldBe TEST_LOCAL_DATE_TIME
                    }
                }
            }
        }

        Given("유저가 생성된 상황에서") {
            val user = userRepository.save(fixture())

            When("기존에 존재하는 이메일로 유저 엔티티를 생성하면") {
                val duplicatedUser: User = fixture {
                    property(User::email) { user.email }
                }

                Then("DataIntegrityViolationException 예외가 발생된다") {
                    shouldThrow<DataIntegrityViolationException> {
                        userRepository.save(duplicatedUser)
                    }
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
) {
    companion object {
        val TEST_LOCAL_DATE_TIME: LocalDateTime = LocalDateTime.of(2022, 10, 9, 17, 30)
    }
}
