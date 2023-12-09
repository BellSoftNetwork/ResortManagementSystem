package net.bellsoft.rms.service.paymentmethod.dto

import java.time.LocalDateTime

data class PaymentMethodDetailDto(
    val id: Long,
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime,
)
