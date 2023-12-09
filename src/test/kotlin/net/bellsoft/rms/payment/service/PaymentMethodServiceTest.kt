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
                    ),
                )

                Then("정상적으로 수정된다") {
                    updatePaymentMethod.run {
                        name shouldBe "BSN"
                        commissionRate shouldBe 0.2
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

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
