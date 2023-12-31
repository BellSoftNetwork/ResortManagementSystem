package net.bellsoft.rms.payment.controller.impl

import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.payment.controller.PaymentMethodController
import net.bellsoft.rms.payment.dto.request.PaymentMethodCreateRequest
import net.bellsoft.rms.payment.dto.request.PaymentMethodPatchRequest
import net.bellsoft.rms.payment.dto.service.PaymentMethodCreateDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodPatchDto
import net.bellsoft.rms.payment.service.PaymentMethodService
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.RestController

@RestController
class PaymentMethodControllerImpl(
    private val paymentMethodService: PaymentMethodService,
) : PaymentMethodController {
    override fun getPaymentMethods(pageable: Pageable) = ListResponse
        .of(paymentMethodService.findAll(pageable))
        .toResponseEntity()

    override fun getPaymentMethod(id: Long) = SingleResponse
        .of(paymentMethodService.find(id))
        .toResponseEntity(HttpStatus.OK)

    override fun createPaymentMethod(
        request: PaymentMethodCreateRequest,
    ) = SingleResponse
        .of(paymentMethodService.create(PaymentMethodCreateDto.of(request)))
        .toResponseEntity(HttpStatus.CREATED)

    override fun updatePaymentMethod(
        id: Long,
        request: PaymentMethodPatchRequest,
    ) = SingleResponse
        .of(paymentMethodService.update(id, PaymentMethodPatchDto.of(request)))
        .toResponseEntity(HttpStatus.OK)

    override fun deletePaymentMethod(id: Long) {
        paymentMethodService.delete(id)
    }
}
