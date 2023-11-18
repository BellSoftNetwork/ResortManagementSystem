package net.bellsoft.rms.domain.reservation

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.util.SecurityTestSupport
import net.bellsoft.rms.util.TestDatabaseSupport
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.repository.findByIdOrNull
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDateTime

@SpringBootTest
@ActiveProfiles("test")
internal class ReservationRepositoryTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val securityTestSupport: SecurityTestSupport,
    private val reservationMethodRepository: ReservationMethodRepository,
    private val reservationRepository: ReservationRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture
        val loginUser: User = fixture()

        mockkStatic(LocalDateTime::class)

        beforeContainer {
            if (it.descriptor.isRootTest())
                securityTestSupport.login(loginUser)
        }

        Given("예약 정보가 생성된 상황에서") {
            val reservationMethod = reservationMethodRepository.save(fixture())
            val reservation = reservationRepository.save(
                fixture {
                    property(Reservation::reservationMethod) { reservationMethod }
                    property(Reservation::room) { null }
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

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
