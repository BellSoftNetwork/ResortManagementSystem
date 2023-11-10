package net.bellsoft.rms.service.admin

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.ints.shouldBeGreaterThanOrEqual
import io.kotest.matchers.longs.shouldBeGreaterThanOrEqual
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.PermissionRequiredDataException
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.controller.v1.admin.AccountCreateDtoFixture
import net.bellsoft.rms.fixture.domain.user.UserFixture
import net.bellsoft.rms.fixture.util.feature
import net.bellsoft.rms.service.admin.dto.AccountCreateDto
import net.bellsoft.rms.service.admin.dto.AccountPatchDto
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class AdminAccountServiceTest(
    private val adminAccountService: AdminAccountService,
    private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

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
                            page.totalPages shouldBeGreaterThanOrEqual 3
                            page.totalElements shouldBeGreaterThanOrEqual 50
                            values.size shouldBe 20
                            values.first().id shouldBe accountIds.first()
                            values.last().id shouldBe accountIds[19]
                        }
                    }
                }
            }
        }

        Given("최고 관리자가") {
            val user = userRepository.save(fixture.feature(UserFixture.Feature.SUPER_ADMIN)())

            When("최고 관리자 권한 계정 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture
                    .feature(AccountCreateDtoFixture.Feature.SUPER_ADMIN)()
                val createdUser = adminAccountService.createAccount(user, accountCreateDto)

                Then("정상적으로 추가된다") {
                    createdUser.email shouldBe accountCreateDto.email
                }
            }

            When("관리자 권한 계정 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture
                    .feature(AccountCreateDtoFixture.Feature.ADMIN)()
                val createdUser = adminAccountService.createAccount(user, accountCreateDto)

                Then("정상적으로 추가된다") {
                    createdUser.email shouldBe accountCreateDto.email
                }
            }

            When("동일한 계정을 반복해서 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture()
                adminAccountService.createAccount(user, accountCreateDto)

                Then("중복 등록으로 유저 추가에 실패한다") {
                    shouldThrow<DataIntegrityViolationException> {
                        adminAccountService.createAccount(user, accountCreateDto)
                    }
                }
            }

            When("일반 권한 계정의 정보 변경 시도 시") {
                val createdUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())
                val updatedUser = adminAccountService.updateAccount(
                    requestUser = user,
                    updateUserId = createdUser.id,
                    accountPatchDto = AccountPatchDto(
                        password = "CHANGE_PASSWORD",
                        name = "CHANGE_NAME",
                        isLock = true,
                        role = UserRole.ADMIN,
                    ),
                )

                Then("정상적으로 변경된다") {
                    updatedUser.run {
                        name shouldBe "CHANGE_NAME"
                        role shouldBe UserRole.ADMIN
                    }
                }
            }

            When("존재하지 않는 계정의 정보 변경 시도 시") {
                Then("조회 불가 예외가 발생한다") {
                    shouldThrow<DataNotFoundException> {
                        adminAccountService.updateAccount(
                            requestUser = user,
                            updateUserId = -1,
                            accountPatchDto = AccountPatchDto(
                                name = "CHANGE_NAME",
                            ),
                        )
                    }
                }
            }
        }

        Given("일반 관리자가") {
            val user = userRepository.save(fixture.feature(UserFixture.Feature.ADMIN)())

            When("관리자 권한 계정 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture
                    .feature(AccountCreateDtoFixture.Feature.ADMIN)()

                Then("권한 부족으로 계정 추가에 실패한다") {
                    shouldThrow<PermissionRequiredDataException> {
                        adminAccountService.createAccount(user, accountCreateDto)
                    }
                }
            }

            When("일반 권한 계정 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture
                    .feature(AccountCreateDtoFixture.Feature.NORMAL)()
                val createdUser = adminAccountService.createAccount(user, accountCreateDto)

                Then("정상적으로 추가된다") {
                    createdUser.email shouldBe accountCreateDto.email
                }
            }

            When("일반 권한 계정의 정보 변경 시도 시") {
                val createdUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())
                val updatedUser = adminAccountService.updateAccount(
                    requestUser = user,
                    updateUserId = createdUser.id,
                    accountPatchDto = AccountPatchDto(
                        name = "CHANGE_NAME",
                    ),
                )

                Then("정상적으로 변경된다") {
                    updatedUser.run {
                        name shouldBe "CHANGE_NAME"
                    }
                }
            }

            When("관리자 권한 계정의 정보 변경 시도 시") {
                val createdUser = userRepository.save(fixture.feature(UserFixture.Feature.ADMIN)())

                Then("변경할 수 없다") {
                    shouldThrow<PermissionRequiredDataException> {
                        adminAccountService.updateAccount(
                            requestUser = user,
                            updateUserId = createdUser.id,
                            accountPatchDto = AccountPatchDto(
                                name = "CHANGE_NAME",
                            ),
                        )
                    }
                }
            }

            When("일반 권한 계정의 권한을 관리자로 변경 시도 시") {
                val createdUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())

                Then("변경할 수 없다") {
                    shouldThrow<PermissionRequiredDataException> {
                        adminAccountService.updateAccount(
                            requestUser = user,
                            updateUserId = createdUser.id,
                            accountPatchDto = AccountPatchDto(
                                role = UserRole.ADMIN,
                            ),
                        )
                    }
                }
            }
        }

        Given("일반 유저가") {
            val user = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())

            When("일반 권한 계정 추가 시도 시") {
                val accountCreateDto: AccountCreateDto = fixture
                    .feature(AccountCreateDtoFixture.Feature.NORMAL)()

                Then("권한 부족으로 계정 추가에 실패한다") {
                    shouldThrow<PermissionRequiredDataException> {
                        adminAccountService.createAccount(user, accountCreateDto)
                    }
                }
            }

            When("일반 권한 계정의 정보 변경 시도 시") {
                val createdUser = userRepository.save(fixture.feature(UserFixture.Feature.NORMAL)())

                Then("변경할 수 없다") {
                    shouldThrow<PermissionRequiredDataException> {
                        adminAccountService.updateAccount(
                            requestUser = user,
                            updateUserId = createdUser.id,
                            accountPatchDto = AccountPatchDto(
                                role = UserRole.ADMIN,
                            ),
                        )
                    }
                }
            }
        }
    },
)
