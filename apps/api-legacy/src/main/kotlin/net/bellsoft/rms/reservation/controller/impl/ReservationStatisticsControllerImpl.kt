package net.bellsoft.rms.reservation.controller.impl

import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.reservation.controller.ReservationStatisticsController
import net.bellsoft.rms.reservation.dto.response.ReservationStatisticsDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.service.ReservationService
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.time.LocalDate

@RestController
class ReservationStatisticsControllerImpl(
    private val reservationService: ReservationService,
) : ReservationStatisticsController {
    override fun getReservationStatistics(
        startDate: String,
        endDate: String,
        periodType: StatisticsPeriodType,
    ): ResponseEntity<SingleResponse<ReservationStatisticsDto>> {
        return SingleResponse
            .of(
                reservationService.getStatistics(
                    LocalDate.parse(startDate),
                    LocalDate.parse(endDate),
                    periodType,
                ),
            )
            .toResponseEntity()
    }
}
