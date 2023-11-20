package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.method.ReservationMethod

data class ReservationMethodUpdateDto(
    val name: String? = null,
    val commissionRate: Double? = null,
    val requireUnpaidAmountCheck: Boolean? = null,
) {
    fun updateEntity(entity: ReservationMethod) {
        name?.let { entity.name = it }
        commissionRate?.let { entity.commissionRate = it }
        requireUnpaidAmountCheck?.let { entity.requireUnpaidAmountCheck = it }
    }
}
