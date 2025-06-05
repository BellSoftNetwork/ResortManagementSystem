package net.bellsoft.rms.reservation.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.reservation.dto.response.ReservationStatisticsDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import org.springframework.data.domain.Pageable
import java.time.LocalDate

interface ReservationService {
    fun findAll(pageable: Pageable, filter: ReservationFilterDto): EntityListDto<ReservationDetailDto>

    fun findById(id: Long): ReservationDetailDto

    fun create(reservationCreateDto: ReservationCreateDto): ReservationDetailDto

    fun update(id: Long, reservationPatchDto: ReservationPatchDto): ReservationDetailDto

    fun delete(id: Long)

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<ReservationDetailDto>>

    /**
     * 지정된 기간 내의 예약 통계를 조회합니다.
     *
     * @param startDate 시작 날짜
     * @param endDate 종료 날짜
     * @param periodType 통계 기간 타입 (기본값: MONTHLY)
     * @return 예약 통계 데이터
     */
    fun getStatistics(
        startDate: LocalDate,
        endDate: LocalDate,
        periodType: StatisticsPeriodType = StatisticsPeriodType.MONTHLY,
    ): ReservationStatisticsDto
}
