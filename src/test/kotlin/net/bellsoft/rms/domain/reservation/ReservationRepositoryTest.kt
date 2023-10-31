package net.bellsoft.rms.domain.reservation

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class ReservationRepositoryTest(
    private val userRepository: UserRepository,
    private val reservationMethodRepository: ReservationMethodRepository,
    private val reservationRepository: ReservationRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("예약 정보가 생성된 상황에서") {
            val user = userRepository.save(fixture())
            val reservationMethod = reservationMethodRepository.save(fixture())
            val reservation = reservationRepository.save(
                fixture {
                    property(Reservation::user) { user }
                    property(Reservation::reservationMethod) { reservationMethod }
                },
            )

            When("등록된 id로 조회하면") {
                val selectedReservation = reservationRepository.findByIdOrNull(reservation.id)

                Then("정상적으로 조회된다") {
                    selectedReservation?.id shouldBe reservation.id
                }
            }

            When("등록되지 않은 id로 조회하면") {
                val unknownReservation = reservationRepository.findByIdOrNull(-1)

                Then("빈 값이 조회된다") {
                    unknownReservation shouldBe null
                }
            }

            When("예약 정보를 삭제하면") {
                val reservationId = reservation.id
                reservationRepository.delete(reservation)

                Then("예약 정보가 조회되지 않는다") {
                    reservationRepository.findByIdOrNull(reservationId) shouldBe null
                }
            }
        }
    },
)
