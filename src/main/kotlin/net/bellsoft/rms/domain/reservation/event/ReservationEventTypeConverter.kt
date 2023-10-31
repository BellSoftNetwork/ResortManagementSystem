package net.bellsoft.rms.domain.reservation.event

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class ReservationEventTypeConverter : AttributeConverter<ReservationEventType, Int> {
    override fun convertToDatabaseColumn(attribute: ReservationEventType): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): ReservationEventType =
        ReservationEventType.values().first { it.value == dbData }
}
