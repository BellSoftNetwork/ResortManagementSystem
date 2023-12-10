package net.bellsoft.rms.reservation.service.impl

import net.bellsoft.rms.authentication.component.SecuritySupport
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.ReservationDetailDto
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.exception.UnavailableRoomException
import net.bellsoft.rms.reservation.mapper.ReservationMapper
import net.bellsoft.rms.reservation.repository.ReservationRepository
import net.bellsoft.rms.reservation.repository.ReservationRoomRepository
import net.bellsoft.rms.reservation.service.ReservationService
import net.bellsoft.rms.revision.component.EntityRevisionComponent
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.repository.RoomRepository
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ReservationServiceImpl(
    private val securitySupport: SecuritySupport,
    private val entityRevisionComponent: EntityRevisionComponent,
    private val reservationRepository: ReservationRepository,
    private val reservationRoomRepository: ReservationRoomRepository,
    private val roomRepository: RoomRepository,
    private val reservationMapper: ReservationMapper,
) : ReservationService {
    override fun findAll(pageable: Pageable, filter: ReservationFilterDto): EntityListDto<ReservationDetailDto> {
        // TODO: N + 1 쿼리 해결 필요
        return EntityListDto.of(
            reservationRepository.getFilteredReservations(pageable, filter),
            reservationMapper::toDto,
        )
    }

    override fun findById(id: Long): ReservationDetailDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        return reservationMapper.toDto(reservation)
    }

    @Transactional
    override fun create(reservationCreateDto: ReservationCreateDto): ReservationDetailDto {
        val reservation = reservationRepository.save(
            reservationMapper.toEntity(reservationCreateDto).apply { rooms = mutableListOf() },
        )
        val requestRoomIds = reservationCreateDto.rooms.map { it.id }.toSet()

        updateReservationRooms(reservation, requestRoomIds)

        return reservationMapper.toDto(reservationRepository.saveAndFlush(reservation))
    }

    @Transactional
    override fun update(id: Long, reservationPatchDto: ReservationPatchDto): ReservationDetailDto {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")

        reservationMapper.updateEntity(reservationPatchDto, reservation)
        if (reservationPatchDto.rooms.isPresent) {
            val requestRoomIds = reservationPatchDto.rooms.get().map { it.id }.toSet()
            updateReservationRooms(reservation, requestRoomIds)
        }

        return reservationMapper.toDto(reservationRepository.saveAndFlush(reservation))
    }

    private fun updateReservationRooms(reservation: Reservation, requestRoomIds: Set<Long>) {
        val reservedRooms = roomRepository.getReservedRooms(
            RoomFilterDto(
                stayStartAt = reservation.stayStartAt,
                stayEndAt = reservation.stayEndAt,
                excludeReservationId = reservation.id,
            ),
            requestRoomIds,
        )

        if (reservedRooms.isNotEmpty()) {
            val reservedRoomNumbers = reservedRooms.joinToString { it.number }
            throw UnavailableRoomException("변경하려는 객실이 이미 다른 예약 건에 배정되어 있음 ($reservedRoomNumbers)")
        }

        reservationRoomRepository.deleteAllByReservation(reservation)
        reservationRoomRepository.flush()

        val requiredRooms = roomRepository.findByIdInOrderByNumberAsc(requestRoomIds)

        reservation.updateRooms(requiredRooms)

        reservationRoomRepository.saveAll(reservation.rooms)
    }

    @Transactional
    override fun delete(id: Long) {
        val reservation = reservationRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약")
        val user = securitySupport.getCurrentUser()
            ?: throw UserNotFoundException("로그인 필요")

        reservationRepository.save(reservation.apply { updatedBy = user })
        reservationRepository.flush()
        reservationRepository.delete(reservation)
    }

    override fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<ReservationDetailDto>> =
        EntityListDto
            .of(entityRevisionComponent.findAllHistory(Reservation::class, reservationMapper::toDto, id, pageable))
}
