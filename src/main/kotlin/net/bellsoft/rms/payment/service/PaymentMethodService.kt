package net.bellsoft.rms.payment.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.payment.dto.response.PaymentMethodDetailDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodCreateDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodPatchDto
import org.springframework.data.domain.Pageable

interface PaymentMethodService {
    fun findAll(pageable: Pageable): EntityListDto<PaymentMethodDetailDto>

    fun find(id: Long): PaymentMethodDetailDto

    fun create(paymentMethodCreateDto: PaymentMethodCreateDto): PaymentMethodDetailDto

    fun update(id: Long, paymentMethodPatchDto: PaymentMethodPatchDto): PaymentMethodDetailDto

    fun delete(id: Long): Boolean
}
