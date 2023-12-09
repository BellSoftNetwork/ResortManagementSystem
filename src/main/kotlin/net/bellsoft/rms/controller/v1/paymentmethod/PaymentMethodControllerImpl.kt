package net.bellsoft.rms.controller.v1.paymentmethod

import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodCreateRequest
import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodPatchRequest
import net.bellsoft.rms.service.paymentmethod.PaymentMethodService
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodCreateDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodPatchDto
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
