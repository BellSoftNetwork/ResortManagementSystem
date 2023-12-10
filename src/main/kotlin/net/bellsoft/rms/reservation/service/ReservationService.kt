package net.bellsoft.rms.reservation.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import org.springframework.data.domain.Pageable

interface ReservationService {
    fun findAll(pageable: Pageable, filter: ReservationFilterDto): EntityListDto<ReservationDetailDto>

    fun findById(id: Long): ReservationDetailDto

    fun create(reservationCreateDto: ReservationCreateDto): ReservationDetailDto

    fun update(id: Long, reservationPatchDto: ReservationPatchDto): ReservationDetailDto

    fun delete(id: Long)

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<ReservationDetailDto>>
}
