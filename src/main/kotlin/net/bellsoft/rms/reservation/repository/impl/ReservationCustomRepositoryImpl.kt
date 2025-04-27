package net.bellsoft.rms.reservation.repository.impl

import com.querydsl.core.types.Order
import com.querydsl.core.types.OrderSpecifier
import com.querydsl.core.types.Projections
import com.querydsl.core.types.dsl.Expressions
import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationStatisticsDto
import net.bellsoft.rms.reservation.dto.response.StatisticsDataDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.entity.QReservation
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.repository.ReservationCustomRepository
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.support.PageableExecutionUtils
import org.springframework.stereotype.Repository
import java.time.LocalDate
import java.time.format.DateTimeFormatter

@Repository
class ReservationCustomRepositoryImpl(
    private val jpaQueryFactory: JPAQueryFactory,
) : ReservationCustomRepository {
    override fun getFilteredReservations(pageable: Pageable, filter: ReservationFilterDto): Page<Reservation> {
        val result = getFilteredReservationsBaseQuery(filter)
            .select(QReservation.reservation)
            .orderBy(
                *pageable.sort.map {
                    OrderSpecifier(
                        if (it.isAscending) Order.ASC else Order.DESC,
                        Expressions.path(String::class.java, QReservation.reservation, it.property),
                    )
                }.toList().toTypedArray<OrderSpecifier<String>?>(),
            )
            .offset(pageable.offset)
            .limit(pageable.pageSize.toLong())
            .fetch()

        return PageableExecutionUtils.getPage(result, pageable) {
            getFilteredReservationsBaseQuery(filter)
                .select(QReservation.reservation.count())
                .fetchOne()!!
        }
    }

    override fun getReservationStatistics(
        startDate: LocalDate,
        endDate: LocalDate,
        periodType: StatisticsPeriodType,
    ): ReservationStatisticsDto {
        // 통계 데이터 조회
        val stats = when (periodType) {
            StatisticsPeriodType.DAILY -> getDailyStatistics(startDate, endDate)
            StatisticsPeriodType.WEEKLY -> getWeeklyStatistics(startDate, endDate)
            StatisticsPeriodType.MONTHLY -> getMonthlyStatistics(startDate, endDate)
            StatisticsPeriodType.YEARLY -> getYearlyStatistics(startDate, endDate)
        }

        return ReservationStatisticsDto(
            periodType = periodType,
            stats = stats,
        )
    }

    private fun getMonthlyStatistics(startDate: LocalDate, endDate: LocalDate): List<StatisticsDataDto> {
        // 기본 데이터 조회 - 연, 월, 가격 합계, 건수, 인원 합계
        data class MonthData(
            val year: Int,
            val month: Int,
            val totalPrice: Long,
            val count: Int,
            val totalPeopleCount: Int,
        )

        val results = jpaQueryFactory
            .from(QReservation.reservation)
            .where(
                QReservation.reservation.stayStartAt.between(startDate, endDate)
                    .or(QReservation.reservation.stayEndAt.between(startDate, endDate)),
            )
            .groupBy(
                QReservation.reservation.stayStartAt.year(),
                QReservation.reservation.stayStartAt.month(),
            )
            .select(
                Projections.constructor(
                    MonthData::class.java,
                    QReservation.reservation.stayStartAt.year(),
                    QReservation.reservation.stayStartAt.month(),
                    QReservation.reservation.price.sum().longValue(),
                    QReservation.reservation.count().intValue(),
                    QReservation.reservation.peopleCount.sum().intValue(),
                ),
            )
            .orderBy(
                QReservation.reservation.stayStartAt.year().asc(),
                QReservation.reservation.stayStartAt.month().asc(),
            )
            .fetch()

        // 결과를 원하는 형태로 변환
        return results.map { monthData ->
            val periodString = String.format("%04d-%02d", monthData.year, monthData.month)
            StatisticsDataDto(
                period = periodString,
                totalSales = monthData.totalPrice,
                totalReservations = monthData.count,
                totalGuests = monthData.totalPeopleCount,
            )
        }
    }

    private fun getDailyStatistics(startDate: LocalDate, endDate: LocalDate): List<StatisticsDataDto> {
        // 기본 데이터 조회 - 날짜, 가격 합계, 건수, 인원 합계
        data class DailyData(
            val date: LocalDate,
            val totalPrice: Long,
            val count: Int,
            val totalPeopleCount: Int,
        )

        val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd")

        return jpaQueryFactory
            .from(QReservation.reservation)
            .where(
                QReservation.reservation.stayStartAt.between(startDate, endDate)
                    .or(QReservation.reservation.stayEndAt.between(startDate, endDate)),
            )
            .groupBy(QReservation.reservation.stayStartAt)
            .select(
                Projections.constructor(
                    DailyData::class.java,
                    QReservation.reservation.stayStartAt,
                    QReservation.reservation.price.sum().longValue(),
                    QReservation.reservation.count().intValue(),
                    QReservation.reservation.peopleCount.sum().intValue(),
                ),
            )
            .orderBy(QReservation.reservation.stayStartAt.asc())
            .fetch()
            .map { dailyData ->
                StatisticsDataDto(
                    period = dailyData.date.format(formatter),
                    totalSales = dailyData.totalPrice,
                    totalReservations = dailyData.count,
                    totalGuests = dailyData.totalPeopleCount,
                )
            }
    }

    private fun getWeeklyStatistics(startDate: LocalDate, endDate: LocalDate): List<StatisticsDataDto> {
        // 기본 데이터 조회 - 날짜, 가격 합계, 건수, 인원 합계
        data class DailyData(
            val date: LocalDate,
            val totalPrice: Long,
            val count: Int,
            val totalPeopleCount: Int,
        )

        val results = jpaQueryFactory
            .from(QReservation.reservation)
            .where(
                QReservation.reservation.stayStartAt.between(startDate, endDate)
                    .or(QReservation.reservation.stayEndAt.between(startDate, endDate)),
            )
            .groupBy(QReservation.reservation.stayStartAt)
            .select(
                Projections.constructor(
                    DailyData::class.java,
                    QReservation.reservation.stayStartAt,
                    QReservation.reservation.price.sum().longValue(),
                    QReservation.reservation.count().intValue(),
                    QReservation.reservation.peopleCount.sum().intValue(),
                ),
            )
            .orderBy(QReservation.reservation.stayStartAt.asc())
            .fetch()

        // 주 단위로 그룹화하여 데이터 집계
        return results
            .groupBy { dailyData ->
                val year = dailyData.date.year
                val week = dailyData.date.get(java.time.temporal.WeekFields.ISO.weekOfWeekBasedYear())
                "$year-${week.toString().padStart(2, '0')}"
            }
            .map { (weekKey, dailyDataList) ->
                StatisticsDataDto(
                    period = weekKey,
                    totalSales = dailyDataList.sumOf { it.totalPrice },
                    totalReservations = dailyDataList.sumOf { it.count },
                    totalGuests = dailyDataList.sumOf { it.totalPeopleCount },
                )
            }
            .sortedBy { it.period }
    }

    private fun getYearlyStatistics(startDate: LocalDate, endDate: LocalDate): List<StatisticsDataDto> {
        // 기본 데이터 조회 - 연도, 가격 합계, 건수, 인원 합계
        data class YearData(
            val year: Int,
            val totalPrice: Long,
            val count: Int,
            val totalPeopleCount: Int,
        )

        return jpaQueryFactory
            .from(QReservation.reservation)
            .where(
                QReservation.reservation.stayStartAt.between(startDate, endDate)
                    .or(QReservation.reservation.stayEndAt.between(startDate, endDate)),
            )
            .groupBy(QReservation.reservation.stayStartAt.year())
            .select(
                Projections.constructor(
                    YearData::class.java,
                    QReservation.reservation.stayStartAt.year(),
                    QReservation.reservation.price.sum().longValue(),
                    QReservation.reservation.count().intValue(),
                    QReservation.reservation.peopleCount.sum().intValue(),
                ),
            )
            .orderBy(QReservation.reservation.stayStartAt.year().asc())
            .fetch()
            .map { yearData ->
                StatisticsDataDto(
                    period = yearData.year.toString(),
                    totalSales = yearData.totalPrice,
                    totalReservations = yearData.count,
                    totalGuests = yearData.totalPeopleCount,
                )
            }
    }

    private fun getFilteredReservationsBaseQuery(filter: ReservationFilterDto) = jpaQueryFactory
        .from(QReservation.reservation)
        .where(
            goeStayStartAt(filter.stayStartAt)?.or(goeStayEndAt(filter.stayStartAt)),
            loeStayStartAt(filter.stayEndAt)?.or(loeStayEndAt(filter.stayEndAt)),
            likeSearchText(filter.searchText),
            eqStatus(filter.status),
            eqType(filter.type),
        )

    private fun goeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.goe(it) }

    private fun loeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.loe(it) }

    private fun goeStayEndAt(stayEndAt: LocalDate?) =
        stayEndAt?.let { QReservation.reservation.stayEndAt.goe(it) }

    private fun loeStayEndAt(stayEndAt: LocalDate?) =
        stayEndAt?.let { QReservation.reservation.stayEndAt.loe(it) }

    private fun likeSearchText(searchText: String?) =
        searchText?.let {
            QReservation.reservation.name.like("%$it%")
                .or(QReservation.reservation.phone.like("%$it%"))
        }

    private fun eqStatus(status: ReservationStatus?) =
        status?.let { QReservation.reservation.status.eq(it) }

    private fun eqType(status: ReservationType?) =
        status?.let { QReservation.reservation.type.eq(it) }
}
