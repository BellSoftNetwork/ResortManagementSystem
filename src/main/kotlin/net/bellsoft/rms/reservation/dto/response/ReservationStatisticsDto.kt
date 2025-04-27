package net.bellsoft.rms.reservation.dto.response

/**
 * 통계 기간 타입
 */
enum class StatisticsPeriodType {
    DAILY, WEEKLY, MONTHLY, YEARLY
}

/**
 * 통계 데이터 DTO
 */
data class StatisticsDataDto(
    /**
     * 기간 식별자 (일별: YYYY-MM-DD, 월별: YYYY-MM, 주별: YYYY-WW, 연별: YYYY)
     */
    val period: String,

    /**
     * 총 매출액
     */
    val totalSales: Long,

    /**
     * 총 예약 건수
     */
    val totalReservations: Int,

    /**
     * 총 방문 인원
     */
    val totalGuests: Int,
)

/**
 * 예약 통계 응답 DTO
 */
data class ReservationStatisticsDto(
    /**
     * 통계 기간 타입
     */
    val periodType: StatisticsPeriodType = StatisticsPeriodType.MONTHLY,

    /**
     * 통계 데이터
     */
    val stats: List<StatisticsDataDto> = emptyList(),

    /**
     * 월별 통계 데이터 (이전 버전 호환성)
     */
    val monthlyStats: List<MonthlyStatDto> = emptyList(),
) {
    /**
     * 월별 통계 데이터 (이전 버전 호환성)
     */
    data class MonthlyStatDto(
        /**
         * 년월 (YYYY-MM 형식)
         */
        val yearMonth: String,

        /**
         * 총 매출액
         */
        val totalSales: Long,

        /**
         * 총 예약 건수
         */
        val totalReservations: Int,

        /**
         * 총 방문 인원
         */
        val totalGuests: Int,
    )

    /**
     * StatisticsDataDto를 MonthlyStatDto로 변환
     */
    constructor(stats: List<StatisticsDataDto>) : this(
        periodType = StatisticsPeriodType.MONTHLY,
        stats = stats,
        monthlyStats = stats.map {
            MonthlyStatDto(
                yearMonth = it.period,
                totalSales = it.totalSales,
                totalReservations = it.totalReservations,
                totalGuests = it.totalGuests,
            )
        },
    )
}
