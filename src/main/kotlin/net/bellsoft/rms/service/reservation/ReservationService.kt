package net.bellsoft.rms.service.reservation

import net.bellsoft.rms.component.auth.SecuritySupport
import net.bellsoft.rms.component.history.EntityHistoryComponent
import net.bellsoft.rms.component.history.dto.EntityHistoryDto
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationRepository
import net.bellsoft.rms.domain.reservation.ReservationRoomRepository
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.mapper.model.ReservationMapper
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationDetailDto
import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import net.bellsoft.rms.service.reservation.dto.ReservationPatchDto
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ReservationService(
    private val securitySupport: SecuritySupport,
    private val entityHistoryComponent: EntityHistoryComponent,
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

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityHistoryDto<ReservationDetailDto>> = EntityListDto
        .of(entityHistoryComponent.findAllHistory(Reservation::class, reservationMapper::toDto, id, pageable))
}
