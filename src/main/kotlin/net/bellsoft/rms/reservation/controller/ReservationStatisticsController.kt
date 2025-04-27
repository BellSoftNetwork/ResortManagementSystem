package net.bellsoft.rms.reservation.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.reservation.dto.response.ReservationStatisticsDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import org.springframework.http.ResponseEntity
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestParam

@Tag(name = "예약 통계", description = "예약 통계 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RequestMapping("/api/v1/reservation-statistics")
interface ReservationStatisticsController {
    @Operation(summary = "예약 통계", description = "지정된 기간 내의 예약 통계 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getReservationStatistics(
        @RequestParam("startDate") startDate: String,
        @RequestParam("endDate") endDate: String,
        @RequestParam(
            "periodType",
            required = false,
            defaultValue = "MONTHLY",
        ) periodType: StatisticsPeriodType = StatisticsPeriodType.MONTHLY,
    ): ResponseEntity<SingleResponse<ReservationStatisticsDto>>
}
