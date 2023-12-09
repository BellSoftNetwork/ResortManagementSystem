package net.bellsoft.rms.controller.v1.paymentmethod

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodCreateRequest
import net.bellsoft.rms.controller.v1.paymentmethod.dto.PaymentMethodPatchRequest
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodDetailDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.annotation.Secured
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController

@Tag(name = "결제 수단", description = "결제 수단 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/payment-methods")
interface PaymentMethodController {
    @Operation(summary = "결제 수단 리스트", description = "결제 수단 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getPaymentMethods(pageable: Pageable): ResponseEntity<ListResponse<PaymentMethodDetailDto>>

    @Operation(summary = "결제 수단 조회", description = "결제 수단 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getPaymentMethod(@PathVariable("id") id: Long): ResponseEntity<SingleResponse<PaymentMethodDetailDto>>

    @Operation(summary = "결제 수단 생성", description = "결제 수단 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PostMapping
    fun createPaymentMethod(
        @RequestBody @Valid
        request: PaymentMethodCreateRequest,
    ): ResponseEntity<SingleResponse<PaymentMethodDetailDto>>

    @Operation(summary = "결제 수단 수정", description = "기존 결제 수단 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PatchMapping("/{id}")
    fun updatePaymentMethod(
        @PathVariable("id") id: Long,
        @RequestBody @Valid
        request: PaymentMethodPatchRequest,
    ): ResponseEntity<SingleResponse<PaymentMethodDetailDto>>

    @Operation(summary = "결제 수단 삭제", description = "기존 결제 수단 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deletePaymentMethod(@PathVariable("id") id: Long)
}
