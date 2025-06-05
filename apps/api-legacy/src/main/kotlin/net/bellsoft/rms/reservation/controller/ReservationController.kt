package net.bellsoft.rms.reservation.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.reservation.dto.filter.ReservationRequestFilter
import net.bellsoft.rms.reservation.dto.request.ReservationCreateRequest
import net.bellsoft.rms.reservation.dto.request.ReservationPatchRequest
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.annotation.Secured
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
interface ReservationController {
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
    ): ResponseEntity<ListResponse<ReservationDetailDto>>

    @Operation(summary = "예약 조회", description = "예약 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getReservation(@PathVariable("id") id: Long): ResponseEntity<SingleResponse<ReservationDetailDto>>

    @Operation(summary = "예약 생성", description = "예약 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PostMapping
    fun createReservation(
        @AuthenticationPrincipal user: User,

        @RequestBody @Valid
        request: ReservationCreateRequest,
    ): SingleResponse<ReservationDetailDto>

    @Operation(summary = "예약 수정", description = "기존 예약 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PatchMapping("/{id}")
    fun updateReservation(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,

        @RequestBody @Valid
        request: ReservationPatchRequest,
    ): ResponseEntity<SingleResponse<ReservationDetailDto>>

    @Operation(summary = "예약 삭제", description = "기존 예약 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteReservation(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,
    )

    @Operation(summary = "예약 이력", description = "예약 정보 변경 이력 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @GetMapping("/{id}/histories")
    fun getReservationHistory(
        @PathVariable("id") id: Long,
        pageable: Pageable,
    ): ResponseEntity<ListResponse<EntityRevisionDto<ReservationDetailDto>>>
}
