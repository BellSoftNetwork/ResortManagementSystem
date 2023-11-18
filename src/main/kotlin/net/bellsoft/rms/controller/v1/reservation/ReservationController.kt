package net.bellsoft.rms.controller.v1.reservation

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import mu.KLogging
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationCreateRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationRequestFilter
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationUpdateRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.reservation.ReservationService
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController

@Tag(name = "예약", description = "예약 관리 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/reservations")
class ReservationController(
    private val reservationService: ReservationService,
) {
    @Operation(summary = "예약 리스트", description = "예약 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getReservations(
        pageable: Pageable,
        filter: ReservationRequestFilter,
    ) = ListResponse
        .of(reservationService.findAll(pageable, filter.toDto()), filter)
        .toResponseEntity()

    @Operation(summary = "예약 조회", description = "예약 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getReservation(@PathVariable("id") id: Long) = SingleResponse
        .of(reservationService.findById(id))
        .toResponseEntity()

    @Operation(summary = "예약 생성", description = "예약 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @PostMapping
    fun createReservation(
        @AuthenticationPrincipal user: User,

        @RequestBody @Valid
        request: ReservationCreateRequest,
    ) = SingleResponse
        .of(reservationService.create(request.toDto()))

    @Operation(summary = "예약 수정", description = "기존 예약 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @PatchMapping("/{id}")
    fun updateReservation(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,

        @RequestBody @Valid
        request: ReservationUpdateRequest,
    ) = SingleResponse
        .of(reservationService.update(id, request.toDto()))
        .toResponseEntity()

    @Operation(summary = "예약 삭제", description = "기존 예약 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteReservation(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,
    ) {
        reservationService.delete(id)
    }

    @Operation(summary = "예약 이력", description = "예약 정보 변경 이력 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}/histories")
    fun getReservationHistory(
        @PathVariable("id") id: Long,
        pageable: Pageable,
    ) = ListResponse
        .of(reservationService.findHistory(id, pageable))
        .toResponseEntity()

    companion object : KLogging()
}
