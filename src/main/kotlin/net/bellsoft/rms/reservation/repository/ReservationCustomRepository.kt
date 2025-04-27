package net.bellsoft.rms.reservation.repository

import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationStatisticsDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.entity.Reservation
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.time.LocalDate

interface ReservationCustomRepository {
    fun getFilteredReservations(pageable: Pageable, filter: ReservationFilterDto): Page<Reservation>

    /**
     * 지정된 기간 내의 예약 통계를 조회합니다.
     *
     * @param startDate 시작 날짜
     * @param endDate 종료 날짜
     * @param periodType 통계 기간 타입 (기본값: MONTHLY)
     * @return 예약 통계 데이터
     */
    fun getReservationStatistics(
        startDate: LocalDate,
        endDate: LocalDate,
        periodType: StatisticsPeriodType = StatisticsPeriodType.MONTHLY,
    ): ReservationStatisticsDto
}
