package net.bellsoft.rms.service.reservation

import net.bellsoft.rms.component.auth.SecuritySupport
import net.bellsoft.rms.component.history.EntityHistoryComponent
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationRepository
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationDto
import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import net.bellsoft.rms.service.reservation.dto.ReservationUpdateDto
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class ReservationService(
    private val securitySupport: SecuritySupport,
    private val entityHistoryComponent: EntityHistoryComponent,
    private val reservationRepository: ReservationRepository,
    private val reservationMethodRepository: ReservationMethodRepository,
    private val roomRepository: RoomRepository,
) {
    fun findAll(pageable: Pageable, filter: ReservationFilterDto): EntityListDto<ReservationDto> {
        return EntityListDto.of(
            reservationRepository.getFilteredReservations(pageable, filter),
            ReservationDto::of,
        )
    }

    fun findById(id: Long): ReservationDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        return ReservationDto.of(reservation)
    }

    @Transactional
    fun create(reservationCreateDto: ReservationCreateDto): ReservationDto {
        reservationCreateDto.loadProxyEntities(
            reservationMethodRepository = reservationMethodRepository,
            roomRepository = roomRepository,
        )

        return ReservationDto.of(reservationRepository.save(reservationCreateDto.toEntity()))
    }

    @Transactional
    fun update(id: Long, reservationUpdateDto: ReservationUpdateDto): ReservationDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        reservationUpdateDto.loadProxyEntities(
            reservationMethodRepository = reservationMethodRepository,
            roomRepository = roomRepository,
        )
        reservationUpdateDto.updateEntity(reservation)

        return ReservationDto.of(reservationRepository.save(reservation))
    }

    @Transactional
    fun delete(id: Long) {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")
        val user = securitySupport.getCurrentUser()
            ?: throw UserNotFoundException("로그인 필요")

        reservationRepository.save(reservation.apply { updatedBy = user })
        reservationRepository.flush()
        reservationRepository.delete(reservation)
    }

    fun findHistory(id: Long, pageable: Pageable) = EntityListDto
        .of(entityHistoryComponent.findAllHistory(Reservation::class, ReservationDto::of, id, pageable))
}
