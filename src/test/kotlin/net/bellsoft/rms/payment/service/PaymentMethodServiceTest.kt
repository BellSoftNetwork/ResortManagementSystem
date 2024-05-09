package net.bellsoft.rms.payment.service

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.dto.service.PaymentMethodCreateDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodPatchDto
import net.bellsoft.rms.payment.entity.PaymentMethod
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.data.repository.findByIdOrNull
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class PaymentMethodServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val paymentMethodService: PaymentMethodService,
    private val paymentMethodRepository: PaymentMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("결제 수단이 없는 상황에서") {
            When("결제 수단을 생성하면") {
                val paymentMethod = paymentMethodService.create(
                    PaymentMethodCreateDto(
                        name = "네이버",
                        commissionRate = 0.1,
                    ),
                )

                Then("정상적으로 생성된다") {
                    paymentMethod.run {
                        name shouldBe "네이버"
                        commissionRate shouldBe 0.1
                    }
                }
            }

            When("존재하지 않는 결제 수단 정보를 조회하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    paymentMethodService.find(-1)
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 결제 수단"
                }
            }

            When("존재하지 않는 결제 수단 정보 수정을 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    paymentMethodService.update(
                        -1,
                        PaymentMethodPatchDto(),
                    )
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 결제 수단"
                }
            }

            When("존재하지 않는 결제 수단 정보 삭제를 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    paymentMethodService.delete(-1)
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 결제 수단"
                }
            }
        }

        Given("결제 수단이 등록된 상황에서") {
            val paymentMethod = paymentMethodRepository.save(fixture())

            When("동일한 결제 수단 생성을 시도하면") {
                val exception = shouldThrow<DuplicateDataException> {
                    paymentMethodService.create(
                        PaymentMethodCreateDto(
                            name = paymentMethod.name,
                            commissionRate = paymentMethod.commissionRate,
                        ),
                    )
                }

                Then("중복 생성으로 예외가 발생한다") {
                    exception.message shouldBe "이미 존재하는 결제 수단"
                }
            }

            When("결제 수단 리스트를 조회하면") {
                val paymentMethods = paymentMethodService.findAll(
                    PageRequest.of(0, 1, Sort.by("id").descending()),
                )

                Then("모든 결제 수단이 정상적으로 조회된다") {
                    paymentMethods.page.size shouldBe 1
                    paymentMethods.values.first().run {
                        name shouldBe paymentMethod.name
                        commissionRate shouldBe paymentMethod.commissionRate
                    }
                }
            }

            When("특정 결제 수단 정보를 조회하면") {
                val findPaymentMethod = paymentMethodService.find(paymentMethod.id)

                Then("정상적으로 조회된다") {
                    findPaymentMethod.run {
                        name shouldBe findPaymentMethod.name
                        commissionRate shouldBe findPaymentMethod.commissionRate
                    }
                }
            }

            When("등록한 결제 수단 정보 수정을 시도하면") {
                val updatePaymentMethod = paymentMethodService.update(
                    paymentMethod.id,
                    PaymentMethodPatchDto(
                        name = JsonNullable.of("BSN"),
                        commissionRate = JsonNullable.of(0.2),
                        isDefaultSelect = JsonNullable.of(true),
                    ),
                )

                Then("정상적으로 수정된다") {
                    updatePaymentMethod.run {
                        name shouldBe "BSN"
                        commissionRate shouldBe 0.2
                        isDefaultSelect shouldBe true
                    }
                }
            }

            When("등록한 결제 수단 정보 삭제를 시도하면") {
                val isDeleted = paymentMethodService.delete(paymentMethod.id)

                Then("정상적으로 삭제된다") {
                    isDeleted shouldBe true
                    paymentMethodRepository.findByIdOrNull(paymentMethod.id) shouldBe null
                }
            }
        }

        Given("기본 선택 옵션이 각각 활성화 및 비활성화된 결제 수단이 한 개씩 등록된 상황에서") {
            val defaultPaymentMethod = paymentMethodRepository.save(
                fixture {
                    property(PaymentMethod::isDefaultSelect) { true }
                },
            )
            val notDefaultPaymentMethod = paymentMethodRepository.save(
                fixture {
                    property(PaymentMethod::isDefaultSelect) { false }
                },
            )

            When("기본 선택 옵션이 비활성화 상태인 결제 수단을 활성화하는 변경 요청을 수행할 시") {
                paymentMethodService.update(
                    notDefaultPaymentMethod.id,
                    PaymentMethodPatchDto(
                        isDefaultSelect = JsonNullable.of(true),
                    ),
                )

                Then("변경 요청에 맞게 정상적으로 기본 선택 옵션이 활성화된다") {
                    val patchedPaymentMethod = paymentMethodRepository.findByIdOrNull(notDefaultPaymentMethod.id)

                    patchedPaymentMethod?.isDefaultSelect shouldBe true
                }

                Then("기본 선택 옵션이 활성화 되어 있던 기존 결제 수단 정보에는 기본 선택 옵션이 비활성화된다") {
                    val oldPaymentMethod = paymentMethodRepository.findByIdOrNull(defaultPaymentMethod.id)

                    oldPaymentMethod?.isDefaultSelect shouldBe false
                }
            }

            When("기본 선택 옵션이 활성화 상태인 결제 수단을 비활성화하는 변경 요청을 수행할 시") {
                paymentMethodService.update(
                    defaultPaymentMethod.id,
                    PaymentMethodPatchDto(
                        isDefaultSelect = JsonNullable.of(false),
                    ),
                )

                Then("변경 요청에 맞게 정상적으로 기본 선택 옵션이 비활성화된다") {
                    val patchedPaymentMethod = paymentMethodRepository.findByIdOrNull(notDefaultPaymentMethod.id)

                    patchedPaymentMethod?.isDefaultSelect shouldBe false
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
