package net.bellsoft.rms.service.reservation

import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.mapper.model.ReservationMethodMapper
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodDetailDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodPatchDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ReservationMethodService(
    private val reservationMethodRepository: ReservationMethodRepository,
    private val reservationMethodMapper: ReservationMethodMapper,
) {
    fun findAll(pageable: Pageable): EntityListDto<ReservationMethodDetailDto> {
        return EntityListDto.of(
            reservationMethodRepository.findAll(pageable),
            reservationMethodMapper::toDto,
        )
    }

    fun find(id: Long): ReservationMethodDetailDto {
        val reservationMethod = reservationMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약 수단")

        return reservationMethodMapper.toDto(reservationMethod)
    }

    @Transactional
    fun create(reservationMethodCreateDto: ReservationMethodCreateDto): ReservationMethodDetailDto {
        try {
            val reservationMethod = reservationMethodMapper.toEntity(reservationMethodCreateDto)

            return reservationMethodMapper.toDto(reservationMethodRepository.save(reservationMethod))
        } catch (e: DataIntegrityViolationException) {
            throw DuplicateDataException("이미 존재하는 예약 수단")
        }
    }

    @Transactional
    fun update(id: Long, reservationMethodPatchDto: ReservationMethodPatchDto): ReservationMethodDetailDto {
        val reservationMethod = reservationMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약 수단")

        reservationMethodMapper.updateEntity(reservationMethodPatchDto, reservationMethod)

        return reservationMethodMapper.toDto(reservationMethodRepository.save(reservationMethod))
    }

    @Transactional
    fun delete(id: Long): Boolean {
        if (!reservationMethodRepository.existsById(id))
            throw DataNotFoundException("존재하지 않는 예약 수단")

        reservationMethodRepository.deleteById(id)

        return true
    }
}
