package net.bellsoft.rms.reservation.controller.impl

import mu.KLogging
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.reservation.controller.ReservationController
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.filter.ReservationRequestFilter
import net.bellsoft.rms.reservation.dto.request.ReservationCreateRequest
import net.bellsoft.rms.reservation.dto.request.ReservationPatchRequest
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.reservation.service.ReservationService
import net.bellsoft.rms.user.entity.User
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
