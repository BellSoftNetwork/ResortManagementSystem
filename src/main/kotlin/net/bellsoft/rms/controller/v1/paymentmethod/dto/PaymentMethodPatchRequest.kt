package net.bellsoft.rms.controller.v1.paymentmethod.dto

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import org.hibernate.validator.constraints.Range
import org.openapitools.jackson.nullable.JsonNullable

@Schema(description = "결제 수단 수정 요청 정보")
data class PaymentMethodPatchRequest(
    @Schema(description = "결제 수단명", example = "네이버")
    @field:Size(min = 2, max = 20)
    val name: JsonNullable<String> = JsonNullable.undefined(),

    @Schema(description = "수수료율", example = "0.2")
    @field:Range(min = 0, max = 1)
    val commissionRate: JsonNullable<Double> = JsonNullable.undefined(),

    @Schema(description = "미수금 금액 알림", example = "false")
    val requireUnpaidAmountCheck: JsonNullable<Boolean> = JsonNullable.undefined(),
)
