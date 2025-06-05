package net.bellsoft.rms.user.service

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.dto.service.UserPatchDto
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.fixture.UserFixture
import net.bellsoft.rms.user.repository.UserRepository
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class UserServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val userService: UserService,
    private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture.new {
            property(UserCreateDto::password) { "password" }
        }

        Given("가입한 사용자가 없는 상황에서") {
            When("신규 회원 가입 시도 시") {
                val userCreateDto: UserCreateDto = fixture {
                    property(UserCreateDto::userId) { "userId" }
                }
                val user = userService.register(userCreateDto)

                Then("정상적으로 가입된다") {
                    user.userId shouldBe userCreateDto.userId
                }
            }

            When("계정 리스트 조회 시") {
                val result = userService.findAll(PageRequest.of(0, 10))

                Then("빈 리스트가 반환된다") {
                    result.page.totalElements shouldBe 0
                }
            }

            When("신규 계정 등록 시") {
                val userCreateDto = fixture<UserCreateDto>()
                val result = userService.register(userCreateDto)

                Then("정상적으로 등록된다") {
                    result.userId shouldBe userCreateDto.userId
                    result.email shouldBe userCreateDto.email
                    result.name shouldBe userCreateDto.name
                    result.role shouldBe userCreateDto.role
                }
            }

            When("존재하지 않는 계정 수정 시도 시") {
                val exception = shouldThrow<DataNotFoundException> {
                    userService.updateAccount(-1, fixture())
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 사용자"
                }
            }
        }

        Given("기존에 가입한 사용자가 있는 상황에서") {
            val userCreateDto: UserCreateDto = fixture {
                property(UserCreateDto::userId) { "userId" }
                property(UserCreateDto::email) { "userId@mail.com" }
            }
            userService.register(userCreateDto)

            When("기존 사용자 ID 와 동일한 ID 로 가입 요청 시") {
                Then("가입이 거부된다") {
                    shouldThrow<DataIntegrityViolationException> {
                        userService.register(userCreateDto)
                    }
                }
            }

            When("기존 사용자 ID 와 다른 ID 로 가입 요청 시") {
                val newUserCreateDto: UserCreateDto = fixture {
                    property(UserCreateDto::userId) { "newId" }
                    property(UserCreateDto::email) { "newId@mail.com" }
                }
                val newUser = userService.register(newUserCreateDto)

                Then("정상적으로 가입된다") {
                    newUser.userId shouldBe newUserCreateDto.userId
                    newUser.email shouldBe newUserCreateDto.email
                }
            }
        }

        Given("50개의 계정이 생성된 상태에서") {
            val accounts = userRepository.saveAll(
                (1..50).map { fixture<User> { property(User::userId) { "find_all_test$it" } } },
            )

            When("계정 리스트 조회 시") {
                val accountsDto = userService.findAll(PageRequest.of(0, 20, Sort.by("id").descending()))

                Then("정상적으로 모두 조회된다") {
                    val accountIds = accounts.map { it.id }.sortedDescending()

                    assertSoftly {
                        accountsDto.run {
                            page.index shouldBe 0
                            page.size shouldBe 20
                            page.totalPages shouldBe 3
                            page.totalElements shouldBe 50
                            values.size shouldBe 20
                            values.first().id shouldBe accountIds.first()
                            values.last().id shouldBe accountIds[19]
                        }
                    }
                }
            }

            When("신규 계정 등록 시") {
                val userCreateDto = fixture<UserCreateDto>()
                val result = userService.register(userCreateDto)

                Then("정상적으로 등록된다") {
                    result.userId shouldBe userCreateDto.userId
                    result.email shouldBe userCreateDto.email
                    result.name shouldBe userCreateDto.name
                    result.role shouldBe userCreateDto.role
                }
            }

            When("기존 계정 수정 시도 시") {
                val account = accounts.first()
                val userPatchDto = UserPatchDto(
                    name = JsonNullable.of("변경된 이름"),
                )
                val result = userService.updateAccount(account.id, userPatchDto)

                Then("정상적으로 수정된다") {
                    result.name shouldBe "변경된 이름"
                }
            }
        }

        Given("일반, 관리자, 최고 관리자 계정이 등록된 상태에서") {
            val normalUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())
            val adminUser = userRepository.save(fixture.feature(UserFixture.Feature.ADMIN)())
            val superAdminUser = userRepository.save(fixture.feature(UserFixture.Feature.SUPER_ADMIN)())

            When("최고 관리자가 최고 관리자를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(superAdminUser, superAdminUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("최고 관리자가 일반 관리자를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(superAdminUser, adminUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("최고 관리자가 일반 유저를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(superAdminUser, normalUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("일반 관리자가 최고 관리자를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(adminUser, superAdminUser.id)

                Then("수정 불가능한 상태로 반환된다") {
                    result shouldBe false
                }
            }

            When("일반 관리자가 일반 관리자를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(adminUser, adminUser.id)

                Then("수정 불가능한 상태로 반환된다") {
                    result shouldBe false
                }
            }

            When("일반 관리자가 일반 유저를 수정할 수 있는지 확인하면") {
                val result = userService.isUpdatableAccount(adminUser, normalUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("일반 관리자가 존재하지 않는 유저를 수정할 수 있는지 확인하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    userService.isUpdatableAccount(adminUser, -1)
                }

                Then("사용자를 찾을 수 없다는 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 사용자"
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
