package net.bellsoft.rms.domain.reservation.event

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationRepository
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class ReservationEventRepositoryTest(
    private val userRepository: UserRepository,
    private val reservationMethodRepository: ReservationMethodRepository,
    private val reservationRepository: ReservationRepository,
    private val reservationEventRepository: ReservationEventRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("예약 이벤트가 생성된 상황에서") {
            val user = userRepository.save(fixture())
            val reservationMethod = reservationMethodRepository.save(fixture())
            val reservation = reservationRepository.save(
                fixture {
                    property(Reservation::user) { user }
                    property(Reservation::reservationMethod) { reservationMethod }
                    property(Reservation::room) { null }
                },
            )
            val reservationEvent = reservationEventRepository.save(
                fixture {
                    property(ReservationEvent::user) { user }
                    property(ReservationEvent::reservation) { reservation }
                },
            )

            When("등록된 id로 조회하면") {
                val selectedReservationEvent = reservationEventRepository.findByIdOrNull(reservationEvent.id)

                Then("정상적으로 조회된다") {
                    selectedReservationEvent?.id shouldBe reservationEvent.id
                }
            }

            When("등록되지 않은 id로 조회하면") {
                val unknownReservationEvent = reservationEventRepository.findByIdOrNull(-1)

                Then("빈 값이 조회된다") {
                    unknownReservationEvent shouldBe null
                }
            }

            When("예약 이벤트를 삭제하면") {
                val reservationEventId = reservationEvent.id
                reservationEventRepository.delete(reservationEvent)

                Then("예약 이벤트 정보가 조회되지 않는다") {
                    reservationEventRepository.findByIdOrNull(reservationEventId) shouldBe null
                }
            }
        }
    },
)
