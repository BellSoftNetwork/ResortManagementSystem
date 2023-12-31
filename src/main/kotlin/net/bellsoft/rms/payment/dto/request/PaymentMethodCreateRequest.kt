package net.bellsoft.rms.payment.dto.request

import io.swagger.v3.oas.annotations.media.Schema
import jakarta.validation.constraints.Size
import org.hibernate.validator.constraints.Range

@Schema(description = "결제 수단 생성 요청 정보")
data class PaymentMethodCreateRequest(
    @Schema(description = "결제 수단명", example = "네이버")
    @field:Size(min = 2, max = 20)
    val name: String,

    @Schema(description = "수수료율", example = "0.2")
    @field:Range(min = 0, max = 1)
    val commissionRate: Double = 0.0,

    @Schema(description = "미수금 금액 알림", example = "false")
    val requireUnpaidAmountCheck: Boolean = false,
)
