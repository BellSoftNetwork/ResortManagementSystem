package net.bellsoft.rms.service.admin

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.domain.user.UserFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.service.admin.dto.AccountCreateDto
import net.bellsoft.rms.service.admin.dto.AccountPatchDto
import net.bellsoft.rms.util.TestDatabaseSupport
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class AdminAccountServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val adminAccountService: AdminAccountService,
    private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("계정이 생성되지 않은 상태에서") {
            When("계정 리스트 조회 시") {
                val result = adminAccountService.findAll(PageRequest.of(0, 10))

                Then("빈 리스트가 반환된다") {
                    result.page.totalElements shouldBe 0
                }
            }

            When("신규 계정 등록 시") {
                val accountCreateDto = fixture<AccountCreateDto>()
                val result = adminAccountService.createAccount(accountCreateDto)

                Then("정상적으로 등록된다") {
                    result.email shouldBe accountCreateDto.email
                    result.name shouldBe accountCreateDto.name
                    result.role shouldBe accountCreateDto.role
                }
            }

            When("존재하지 않는 계정 수정 시도 시") {
                val exception = shouldThrow<DataNotFoundException> {
                    adminAccountService.updateAccount(-1, fixture())
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 사용자"
                }
            }
        }

        Given("50개의 계정이 생성된 상태에서") {
            val accounts = userRepository.saveAll(
                (1..50).map { fixture<User> { property(User::email) { "find_all_test$it@mail.com" } } },
            )

            When("계정 리스트 조회 시") {
                val accountsDto = adminAccountService.findAll(PageRequest.of(0, 20, Sort.by("id").descending()))

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
                val accountCreateDto = fixture<AccountCreateDto>()
                val result = adminAccountService.createAccount(accountCreateDto)

                Then("정상적으로 등록된다") {
                    result.email shouldBe accountCreateDto.email
                    result.name shouldBe accountCreateDto.name
                    result.role shouldBe accountCreateDto.role
                }
            }

            When("기존 계정 수정 시도 시") {
                val account = accounts.first()
                val accountPatchDto = AccountPatchDto(name = "변경된 이름")
                val result = adminAccountService.updateAccount(account.id, accountPatchDto)

                Then("정상적으로 수정된다") {
                    result.name shouldBe accountPatchDto.name
                }
            }
        }

        Given("일반, 관리자, 최고 관리자 계정이 등록된 상태에서") {
            val normalUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())
            val adminUser = userRepository.save(fixture.feature(UserFixture.Feature.ADMIN)())
            val superAdminUser = userRepository.save(fixture.feature(UserFixture.Feature.SUPER_ADMIN)())

            When("최고 관리자가 최고 관리자를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(superAdminUser, superAdminUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("최고 관리자가 일반 관리자를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(superAdminUser, adminUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("최고 관리자가 일반 유저를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(superAdminUser, normalUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("일반 관리자가 최고 관리자를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(adminUser, superAdminUser.id)

                Then("수정 불가능한 상태로 반환된다") {
                    result shouldBe false
                }
            }

            When("일반 관리자가 일반 관리자를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(adminUser, adminUser.id)

                Then("수정 불가능한 상태로 반환된다") {
                    result shouldBe false
                }
            }

            When("일반 관리자가 일반 유저를 수정할 수 있는지 확인하면") {
                val result = adminAccountService.isUpdatableAccount(adminUser, normalUser.id)

                Then("수정 가능한 상태로 반환된다") {
                    result shouldBe true
                }
            }

            When("일반 관리자가 존재하지 않는 유저를 수정할 수 있는지 확인하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    adminAccountService.isUpdatableAccount(adminUser, -1)
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
