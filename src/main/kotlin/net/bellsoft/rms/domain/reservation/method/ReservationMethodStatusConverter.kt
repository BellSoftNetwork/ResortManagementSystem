package net.bellsoft.rms.domain.reservation.method

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class ReservationMethodStatusConverter : AttributeConverter<ReservationMethodStatus, Int> {
    override fun convertToDatabaseColumn(attribute: ReservationMethodStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): ReservationMethodStatus =
        ReservationMethodStatus.values().first { it.value == dbData }
}
