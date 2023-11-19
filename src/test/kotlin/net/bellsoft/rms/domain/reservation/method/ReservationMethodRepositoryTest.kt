package net.bellsoft.rms.domain.reservation.method

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class ReservationMethodRepositoryTest(
    private val reservationMethodRepository: ReservationMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("예약 수단이 생성된 상황에서") {
            val reservationMethod = reservationMethodRepository.save(fixture())

            When("등록된 id로 조회하면") {
                val selectedReservationMethod = reservationMethodRepository.findByIdOrNull(reservationMethod.id)

                Then("정상적으로 조회된다") {
                    selectedReservationMethod?.id shouldBe reservationMethod.id
                }
            }

            When("등록되지 않은 id로 조회하면") {
                val findReservationMethod = reservationMethodRepository.findByIdOrNull(-1)

                Then("빈 값이 조회된다") {
                    findReservationMethod shouldBe null
                }
            }

            When("예약 수단을 삭제하면") {
                val reservationMethodId = reservationMethod.id
                reservationMethodRepository.delete(reservationMethod)

                Then("예약 수단 정보가 조회되지 않는다") {
                    reservationMethodRepository.findByIdOrNull(reservationMethodId) shouldBe null
                }
            }
        }
    },
)
