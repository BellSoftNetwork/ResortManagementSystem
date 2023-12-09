package net.bellsoft.rms.service.paymentmethod

import net.bellsoft.rms.domain.paymentmethod.PaymentMethodRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.mapper.model.PaymentMethodMapper
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodCreateDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodDetailDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodPatchDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class PaymentMethodService(
    private val paymentMethodRepository: PaymentMethodRepository,
    private val paymentMethodMapper: PaymentMethodMapper,
) {
    fun findAll(pageable: Pageable): EntityListDto<PaymentMethodDetailDto> {
        return EntityListDto.of(
            paymentMethodRepository.findAll(pageable),
            paymentMethodMapper::toDto,
        )
    }

    fun find(id: Long): PaymentMethodDetailDto {
        val paymentMethod = paymentMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 결제 수단")

        return paymentMethodMapper.toDto(paymentMethod)
    }

    @Transactional
    fun create(paymentMethodCreateDto: PaymentMethodCreateDto): PaymentMethodDetailDto {
        try {
            val paymentMethod = paymentMethodMapper.toEntity(paymentMethodCreateDto)

            return paymentMethodMapper.toDto(paymentMethodRepository.save(paymentMethod))
        } catch (e: DataIntegrityViolationException) {
            throw DuplicateDataException("이미 존재하는 결제 수단")
        }
    }

    @Transactional
    fun update(id: Long, paymentMethodPatchDto: PaymentMethodPatchDto): PaymentMethodDetailDto {
        val paymentMethod = paymentMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 결제 수단")

        paymentMethodMapper.updateEntity(paymentMethodPatchDto, paymentMethod)

        return paymentMethodMapper.toDto(paymentMethodRepository.save(paymentMethod))
    }

    @Transactional
    fun delete(id: Long): Boolean {
        if (!paymentMethodRepository.existsById(id))
            throw DataNotFoundException("존재하지 않는 결제 수단")

        paymentMethodRepository.deleteById(id)

        return true
    }
}
