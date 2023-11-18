package net.bellsoft.rms.controller.v1.reservation

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodCreateRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodUpdateRequest
import net.bellsoft.rms.service.reservation.ReservationMethodService
import net.bellsoft.rms.service.reservation.dto.ReservationMethodDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
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

@Tag(name = "예약 수단", description = "예약 수단 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/reservation-methods")
class ReservationMethodController(
    private val reservationMethodService: ReservationMethodService,
) {
    @Operation(summary = "예약 수단 리스트", description = "예약 수단 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getReservations(pageable: Pageable) = ListResponse
        .of(reservationMethodService.findAll(pageable))
        .toResponseEntity()

    @Operation(summary = "예약 수단 조회", description = "예약 수단 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getReservation(@PathVariable("id") id: Long): ResponseEntity<SingleResponse<ReservationMethodDto>> {
        return SingleResponse
            .of(reservationMethodService.find(id))
            .toResponseEntity(HttpStatus.OK)
    }

    @Operation(summary = "예약 수단 생성", description = "예약 수단 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @PostMapping
    fun createReservation(
        @RequestBody @Valid
        request: ReservationMethodCreateRequest,
    ) = SingleResponse
        .of(reservationMethodService.create(request.toDto()))
        .toResponseEntity(HttpStatus.CREATED)

    @Operation(summary = "예약 수단 수정", description = "기존 예약 수단 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @PatchMapping("/{id}")
    fun updateReservation(
        @PathVariable("id") id: Long,
        @RequestBody @Valid
        request: ReservationMethodUpdateRequest,
    ) = SingleResponse
        .of(reservationMethodService.update(id, request.toDto()))
        .toResponseEntity(HttpStatus.OK)

    @Operation(summary = "예약 수단 삭제", description = "기존 예약 수단 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteReservation(@PathVariable("id") id: Long) {
        reservationMethodService.delete(id)
    }
}