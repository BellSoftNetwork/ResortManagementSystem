package net.bellsoft.rms.service.reservation

import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodUpdateDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class ReservationMethodService(
    private val reservationMethodRepository: ReservationMethodRepository,
) {
    fun findAll(pageable: Pageable): EntityListDto<ReservationMethodDto> {
        return EntityListDto.of(
            reservationMethodRepository.findAll(pageable),
            ReservationMethodDto::of,
        )
    }

    fun find(id: Long): ReservationMethodDto {
        val reservationMethod = reservationMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약 수단")

        return ReservationMethodDto.of(reservationMethod)
    }

    @Transactional
    fun create(reservationMethodCreateDto: ReservationMethodCreateDto): ReservationMethodDto {
        try {
            return ReservationMethodDto.of(reservationMethodRepository.save(reservationMethodCreateDto.toEntity()))
        } catch (e: DataIntegrityViolationException) {
            throw DuplicateDataException("이미 존재하는 예약 수단")
        }
    }

    @Transactional
    fun update(id: Long, reservationMethodUpdateDto: ReservationMethodUpdateDto): ReservationMethodDto {
        val reservationMethod = reservationMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 예약 수단")

        reservationMethodUpdateDto.updateEntity(reservationMethod)

        return ReservationMethodDto.of(reservationMethodRepository.save(reservationMethod))
    }

    @Transactional
    fun delete(id: Long): Boolean {
        if (!reservationMethodRepository.existsById(id))
            throw DataNotFoundException("존재하지 않는 예약 수단")

        reservationMethodRepository.deleteById(id)

        return true
    }
}
