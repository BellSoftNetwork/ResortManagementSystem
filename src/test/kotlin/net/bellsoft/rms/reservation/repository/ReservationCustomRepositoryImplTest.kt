package net.bellsoft.rms.reservation.repository

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.common.util.SecurityTestSupport
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.user.entity.User
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDate

@SpringBootTest
@ActiveProfiles("test")
internal class ReservationCustomRepositoryImplTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val securityTestSupport: SecurityTestSupport,
    private val paymentMethodRepository: PaymentMethodRepository,
    private val reservationRepository: ReservationRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture
        val loginUser: User = fixture()

        beforeContainer {
            if (it.descriptor.isRootTest())
                securityTestSupport.login(loginUser)
        }

        Given("예약 데이터가 있는 상황에서") {
            val paymentMethod = paymentMethodRepository.save(fixture())

            // 2024년 5월 데이터
            val reservation1 = reservationRepository.save(
                fixture {
                    property(Reservation::paymentMethod) { paymentMethod }
                    property(Reservation::stayStartAt) { LocalDate.of(2024, 5, 1) }
                    property(Reservation::stayEndAt) { LocalDate.of(2024, 5, 3) }
                    property(Reservation::price) { 100000 }
                    property(Reservation::peopleCount) { 2 }
                },
            )

            // 2024년 6월 데이터
            val reservation2 = reservationRepository.save(
                fixture {
                    property(Reservation::paymentMethod) { paymentMethod }
                    property(Reservation::stayStartAt) { LocalDate.of(2024, 6, 1) }
                    property(Reservation::stayEndAt) { LocalDate.of(2024, 6, 3) }
                    property(Reservation::price) { 200000 }
                    property(Reservation::peopleCount) { 3 }
                },
            )

            When("월별 통계를 조회하면") {
                val startDate = LocalDate.of(2024, 5, 1)
                val endDate = LocalDate.of(2024, 6, 30)
                val statistics =
                    reservationRepository.getReservationStatistics(startDate, endDate, StatisticsPeriodType.MONTHLY)

                Then("통계 데이터가 정상적으로 조회된다") {
                    statistics.periodType shouldBe StatisticsPeriodType.MONTHLY
                    statistics.stats.size shouldBe 2

                    // 5월 데이터 확인
                    val mayStats = statistics.stats.find { it.period == "2024-05" }
                    mayStats shouldNotBe null
                    mayStats?.totalSales shouldBe 100000L
                    mayStats?.totalReservations shouldBe 1
                    mayStats?.totalGuests shouldBe 2

                    // 6월 데이터 확인
                    val juneStats = statistics.stats.find { it.period == "2024-06" }
                    juneStats shouldNotBe null
                    juneStats?.totalSales shouldBe 200000L
                    juneStats?.totalReservations shouldBe 1
                    juneStats?.totalGuests shouldBe 3
                }
            }

            When("일별 통계를 조회하면") {
                val startDate = LocalDate.of(2024, 5, 1)
                val endDate = LocalDate.of(2024, 6, 30)
                val statistics =
                    reservationRepository.getReservationStatistics(startDate, endDate, StatisticsPeriodType.DAILY)

                Then("통계 데이터가 정상적으로 조회된다") {
                    statistics.periodType shouldBe StatisticsPeriodType.DAILY
                    statistics.stats.size shouldBe 2

                    // 5월 1일 데이터 확인
                    val may1Stats = statistics.stats.find { it.period == "2024-05-01" }
                    may1Stats shouldNotBe null
                    may1Stats?.totalSales shouldBe 100000L
                    may1Stats?.totalReservations shouldBe 1
                    may1Stats?.totalGuests shouldBe 2

                    // 6월 1일 데이터 확인
                    val june1Stats = statistics.stats.find { it.period == "2024-06-01" }
                    june1Stats shouldNotBe null
                    june1Stats?.totalSales shouldBe 200000L
                    june1Stats?.totalReservations shouldBe 1
                    june1Stats?.totalGuests shouldBe 3
                }
            }

            When("연도별 통계를 조회하면") {
                val startDate = LocalDate.of(2024, 1, 1)
                val endDate = LocalDate.of(2024, 12, 31)
                val statistics =
                    reservationRepository.getReservationStatistics(startDate, endDate, StatisticsPeriodType.YEARLY)

                Then("통계 데이터가 정상적으로 조회된다") {
                    statistics.periodType shouldBe StatisticsPeriodType.YEARLY
                    statistics.stats.size shouldBe 1

                    // 2024년 데이터 확인
                    val year2024Stats = statistics.stats.find { it.period == "2024" }
                    year2024Stats shouldNotBe null
                    year2024Stats?.totalSales shouldBe 300000L
                    year2024Stats?.totalReservations shouldBe 2
                    year2024Stats?.totalGuests shouldBe 5
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
