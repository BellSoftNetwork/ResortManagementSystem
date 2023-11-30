package net.bellsoft.rms.controller.v1.reservation

import mu.KLogging
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationCreateRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationPatchRequest
import net.bellsoft.rms.controller.v1.reservation.dto.ReservationRequestFilter
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.reservation.ReservationService
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import net.bellsoft.rms.service.reservation.dto.ReservationPatchDto
import org.springframework.data.domain.Pageable
import org.springframework.web.bind.annotation.RestController

@RestController
class ReservationControllerImpl(
    private val reservationService: ReservationService,
) : ReservationController {
    override fun getReservations(
        pageable: Pageable,
        filter: ReservationRequestFilter,
    ) = ListResponse
        .of(reservationService.findAll(pageable, ReservationFilterDto.of(filter)), filter)
        .toResponseEntity()

    override fun getReservation(id: Long) = SingleResponse
        .of(reservationService.findById(id))
        .toResponseEntity()

    override fun createReservation(
        user: User,
        request: ReservationCreateRequest,
    ) = SingleResponse
        .of(reservationService.create(ReservationCreateDto.of(request)))

    override fun updateReservation(
        user: User,
        id: Long,
        request: ReservationPatchRequest,
    ) = SingleResponse
        .of(reservationService.update(id, ReservationPatchDto.of(request)))
        .toResponseEntity()

    override fun deleteReservation(
        user: User,
        id: Long,
    ) {
        reservationService.delete(id)
    }

    override fun getReservationHistory(
        id: Long,
        pageable: Pageable,
    ) = ListResponse
        .of(reservationService.findHistory(id, pageable))
        .toResponseEntity()

    companion object : KLogging()
}
