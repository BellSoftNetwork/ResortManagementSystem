package net.bellsoft.rms.controller.v1.reservation

import net.bellsoft.rms.controller.v1.dto.ListTypedResponse
import net.bellsoft.rms.controller.v1.dto.SingleTypedResponse
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodCreateRequest
import net.bellsoft.rms.service.reservation.ReservationMethodService
import net.bellsoft.rms.service.reservation.dto.ReservationMethodDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodUpdateDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/api/v1/reservation-methods")
class ReservationMethodController(
    private val reservationMethodService: ReservationMethodService,
) {
    @GetMapping
    fun getReservations(pageable: Pageable): ResponseEntity<ListTypedResponse<ReservationMethodDto>> {
        return ListTypedResponse.of(reservationMethodService.findAll(pageable)).toResponseEntity()
    }

    @GetMapping("/{id}")
    fun getReservation(@PathVariable("id") id: Long): ResponseEntity<ReservationMethodDto> {
        return SingleTypedResponse.of(reservationMethodService.find(id)).toResponseEntity(HttpStatus.OK)
    }

    @PostMapping
    fun createReservation(@RequestBody request: ReservationMethodCreateRequest): ResponseEntity<ReservationMethodDto> {
        return SingleTypedResponse.of(reservationMethodService.create(request.toDto()))
            .toResponseEntity(HttpStatus.CREATED)
    }

    @PatchMapping("/{id}")
    fun updateReservation(
        @PathVariable("id") id: Long,
        @RequestBody reservationMethodUpdateDto: ReservationMethodUpdateDto,
    ): ResponseEntity<ReservationMethodDto> {
        return SingleTypedResponse.of(reservationMethodService.update(id, reservationMethodUpdateDto))
            .toResponseEntity(HttpStatus.OK)
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteReservation(@PathVariable("id") id: Long) {
        reservationMethodService.delete(id)
    }
}
