package net.bellsoft.rms.controller.v1.reservation

import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodCreateRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationMethodPatchRequest
import net.bellsoft.rms.service.reservation.ReservationMethodService
import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodPatchDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.RestController

@RestController
class ReservationMethodControllerImpl(
    private val reservationMethodService: ReservationMethodService,
) : ReservationMethodController {
    override fun getReservations(pageable: Pageable) = ListResponse
        .of(reservationMethodService.findAll(pageable))
        .toResponseEntity()

    override fun getReservation(id: Long) = SingleResponse
        .of(reservationMethodService.find(id))
        .toResponseEntity(HttpStatus.OK)

    override fun createReservation(
        request: ReservationMethodCreateRequest,
    ) = SingleResponse
        .of(reservationMethodService.create(ReservationMethodCreateDto.of(request)))
        .toResponseEntity(HttpStatus.CREATED)

    override fun updateReservation(
        id: Long,
        request: ReservationMethodPatchRequest,
    ) = SingleResponse
        .of(reservationMethodService.update(id, ReservationMethodPatchDto.of(request)))
        .toResponseEntity(HttpStatus.OK)

    override fun deleteReservation(id: Long) {
        reservationMethodService.delete(id)
    }
}
