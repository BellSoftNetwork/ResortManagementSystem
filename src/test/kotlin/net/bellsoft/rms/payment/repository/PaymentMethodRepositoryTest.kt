package net.bellsoft.rms.payment.repository

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.common.annotation.JpaEntityTest
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class PaymentMethodRepositoryTest(
    private val paymentMethodRepository: PaymentMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("결제 수단이 생성된 상황에서") {
            val paymentMethod = paymentMethodRepository.save(fixture())

            When("등록된 id로 조회하면") {
                val selectedPaymentMethod = paymentMethodRepository.findByIdOrNull(paymentMethod.id)

                Then("정상적으로 조회된다") {
                    selectedPaymentMethod?.id shouldBe paymentMethod.id
                }
            }

            When("등록되지 않은 id로 조회하면") {
                val findPaymentMethod = paymentMethodRepository.findByIdOrNull(-1)

                Then("빈 값이 조회된다") {
                    findPaymentMethod shouldBe null
                }
            }

            When("결제 수단을 삭제하면") {
                val paymentMethodId = paymentMethod.id
                paymentMethodRepository.delete(paymentMethod)

                Then("결제 수단 정보가 조회되지 않는다") {
                    paymentMethodRepository.findByIdOrNull(paymentMethodId) shouldBe null
                }
            }
        }
    },
)
