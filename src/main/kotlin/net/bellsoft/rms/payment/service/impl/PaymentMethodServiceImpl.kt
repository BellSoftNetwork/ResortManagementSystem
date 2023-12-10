package net.bellsoft.rms.payment.service.impl

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.payment.dto.response.PaymentMethodDetailDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodCreateDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodPatchDto
import net.bellsoft.rms.payment.mapper.PaymentMethodMapper
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.payment.service.PaymentMethodService
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class PaymentMethodServiceImpl(
    private val paymentMethodRepository: PaymentMethodRepository,
    private val paymentMethodMapper: PaymentMethodMapper,
) : PaymentMethodService {
    override fun findAll(pageable: Pageable): EntityListDto<PaymentMethodDetailDto> {
        return EntityListDto.of(
            paymentMethodRepository.findAll(pageable),
            paymentMethodMapper::toDto,
        )
    }

    override fun find(id: Long): PaymentMethodDetailDto {
        val paymentMethod = paymentMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 결제 수단")

        return paymentMethodMapper.toDto(paymentMethod)
    }

    @Transactional
    override fun create(paymentMethodCreateDto: PaymentMethodCreateDto): PaymentMethodDetailDto {
        try {
            val paymentMethod = paymentMethodMapper.toEntity(paymentMethodCreateDto)

            return paymentMethodMapper.toDto(paymentMethodRepository.save(paymentMethod))
        } catch (e: DataIntegrityViolationException) {
            throw DuplicateDataException("이미 존재하는 결제 수단")
        }
    }

    @Transactional
    override fun update(id: Long, paymentMethodPatchDto: PaymentMethodPatchDto): PaymentMethodDetailDto {
        val paymentMethod = paymentMethodRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 결제 수단")

        paymentMethodMapper.updateEntity(paymentMethodPatchDto, paymentMethod)

        return paymentMethodMapper.toDto(paymentMethodRepository.save(paymentMethod))
    }

    @Transactional
    override fun delete(id: Long): Boolean {
        if (!paymentMethodRepository.existsById(id))
            throw DataNotFoundException("존재하지 않는 결제 수단")

        paymentMethodRepository.deleteById(id)

        return true
    }
}
