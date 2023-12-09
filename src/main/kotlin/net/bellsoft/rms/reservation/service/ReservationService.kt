package net.bellsoft.rms.reservation.service

import net.bellsoft.rms.authentication.component.SecuritySupport
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.mapper.ReservationMapper
import net.bellsoft.rms.reservation.repository.ReservationRepository
import net.bellsoft.rms.reservation.repository.ReservationRoomRepository
import net.bellsoft.rms.revision.component.EntityRevisionComponent
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.room.repository.RoomRepository
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ReservationService(
    private val securitySupport: SecuritySupport,
    private val entityRevisionComponent: EntityRevisionComponent,
    private val reservationRepository: ReservationRepository,
    private val reservationRoomRepository: ReservationRoomRepository,
    private val roomRepository: RoomRepository,
    private val reservationMapper: ReservationMapper,
) {
    fun findAll(pageable: Pageable, filter: ReservationFilterDto): EntityListDto<ReservationDetailDto> {
        // TODO: N + 1 쿼리 해결 필요
        return EntityListDto.of(
            reservationRepository.getFilteredReservations(pageable, filter),
            reservationMapper::toDto,
        )
    }

    fun findById(id: Long): ReservationDetailDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        return reservationMapper.toDto(reservation)
    }

    @Transactional
    fun create(reservationCreateDto: ReservationCreateDto): ReservationDetailDto {
        val reservation = reservationRepository.save(
            reservationMapper.toEntity(reservationCreateDto).apply { rooms = mutableListOf() },
        )
        val requestRoomIds = reservationCreateDto.rooms.map { it.id }.toSet()

        updateReservationRooms(reservation, requestRoomIds)

        return reservationMapper.toDto(reservationRepository.saveAndFlush(reservation))
    }

    @Transactional
    fun update(id: Long, reservationPatchDto: ReservationPatchDto): ReservationDetailDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        reservationMapper.updateEntity(reservationPatchDto, reservation)
        if (reservationPatchDto.rooms.isPresent) {
            val requestRoomIds = reservationPatchDto.rooms.get().map { it.id }.toSet()
            updateReservationRooms(reservation, requestRoomIds)
        }

        return reservationMapper.toDto(reservationRepository.saveAndFlush(reservation))
    }

    @Transactional
    fun updateReservationRooms(reservation: Reservation, requestRoomIds: Set<Long>) {
        reservationRoomRepository.deleteAllByReservation(reservation)
        reservationRoomRepository.flush()

        val requiredRooms = roomRepository.findByIdInOrderByNumberAsc(requestRoomIds)

        reservation.updateRooms(requiredRooms)

        reservationRoomRepository.saveAll(reservation.rooms)
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

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<ReservationDetailDto>> =
        EntityListDto
            .of(entityRevisionComponent.findAllHistory(Reservation::class, reservationMapper::toDto, id, pageable))
}
