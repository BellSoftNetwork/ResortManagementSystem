package net.bellsoft.rms.domain.user

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import io.mockk.every
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import org.springframework.jdbc.core.JdbcTemplate
import java.time.LocalDateTime

@JpaEntityTest
internal class UserRepositoryTest(
    private val userRepository: UserRepository,
    private val jdbcTemplate: JdbcTemplate,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("한 유저가 생성된 상황에서") {
            val user = userRepository.save(fixture())

            When("한 유저를 같은 트랜잭션 안에서 여러번 조회하면") {
                Then("캐싱된 값이 나온다") {
                    val selectedUser1 = userRepository.findByIdOrNull(user.id)
                    val selectedUser2 = userRepository.findByIdOrNull(user.id)

                    selectedUser1 shouldNotBe null
                    selectedUser1 shouldBe selectedUser2
                }
            }

            When("등록된 이메일로 유저를 조회하면") {
                val selectedUser = userRepository.findByEmail(user.email)

                Then("정상적으로 조회된다") {
                    selectedUser?.id shouldBe user.id
                }
            }

            When("등록되지 않은 이메일로 유저를 조회하면") {
                val unknownUser = userRepository.findByEmail("NON-EXIST-EMAIL@bellsoft.net")

                Then("빈 값이 조회된다") {
                    unknownUser shouldBe null
                }
            }

            When("유저를 삭제하면") {
                val userId = user.id
                userRepository.delete(user)

                Then("유저가 조회되지 않는다") {
                    userRepository.findByIdOrNull(userId) shouldBe null
                }

                Then("실제 DB에 순수 쿼리로 조회했을 때 deleted_at 값과 함께 존재한다 (JdbcTemplate)") {
                    val localDateTime = jdbcTemplate.queryForObject(
                        "SELECT deleted_at FROM user WHERE id = $userId",
                        LocalDateTime::class.java,
                    )

                    localDateTime shouldNotBe null
                }
            }
        }

        Given("2022년 가입 유저 2명, 2023년 가입 유저가 1명 생성된 상황에서") {
            val createdAtLocalDateTimes = listOf(
                LocalDateTime.of(2022, 1, 1, 0, 0),
                LocalDateTime.of(2022, 12, 30, 23, 23),
                LocalDateTime.of(2023, 1, 1, 0, 0),
            )

            createdAtLocalDateTimes.forEach {
                every { LocalDateTime.now() } returns it
                userRepository.save(fixture())
            }

            When("2022년 가입 유저를 조회하면") {
                val users = userRepository.getYearCreatedUsers(2022)

                Then("2명이 조회된다") {
                    users.size shouldBe 2
                }
            }
        }
    },
)
