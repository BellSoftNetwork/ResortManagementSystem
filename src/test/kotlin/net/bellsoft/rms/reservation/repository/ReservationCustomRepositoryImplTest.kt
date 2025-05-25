package net.bellsoft.rms.reservation.repository

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.common.util.SecurityTestSupport
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.user.entity.User
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
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

            When("stayStartAt 기준 오름차순으로 정렬하면") {
                val pageable = PageRequest.of(0, 10, Sort.by(Sort.Direction.ASC, "stayStartAt"))
                val filter = ReservationFilterDto()
                val result = reservationRepository.getFilteredReservations(pageable, filter)

                Then("입실일이 빠른 순서대로 정렬된다") {
                    result.content.size shouldBe 2
                    result.content[0].stayStartAt shouldBe LocalDate.of(2024, 5, 1)
                    result.content[1].stayStartAt shouldBe LocalDate.of(2024, 6, 1)
                }
            }

            When("stayStartAt 기준 내림차순으로 정렬하면") {
                val pageable = PageRequest.of(0, 10, Sort.by(Sort.Direction.DESC, "stayStartAt"))
                val filter = ReservationFilterDto()
                val result = reservationRepository.getFilteredReservations(pageable, filter)

                Then("입실일이 늦은 순서대로 정렬된다") {
                    result.content.size shouldBe 2
                    result.content[0].stayStartAt shouldBe LocalDate.of(2024, 6, 1)
                    result.content[1].stayStartAt shouldBe LocalDate.of(2024, 5, 1)
                }
            }

            When("오름차순과 내림차순 정렬 결과를 비교하면") {
                val ascPageable = PageRequest.of(0, 10, Sort.by(Sort.Direction.ASC, "stayStartAt"))
                val descPageable = PageRequest.of(0, 10, Sort.by(Sort.Direction.DESC, "stayStartAt"))
                val filter = ReservationFilterDto()

                val ascResult = reservationRepository.getFilteredReservations(ascPageable, filter)
                val descResult = reservationRepository.getFilteredReservations(descPageable, filter)

                Then("정렬 결과가 서로 반대 순서여야 한다") {
                    ascResult.content.size shouldBe 2
                    descResult.content.size shouldBe 2

                    ascResult.content[0].id shouldBe descResult.content[1].id
                    ascResult.content[1].id shouldBe descResult.content[0].id
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
