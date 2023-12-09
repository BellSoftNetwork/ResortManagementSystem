package net.bellsoft.rms.reservation.converter

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter
import net.bellsoft.rms.reservation.type.ReservationStatus

@Converter(autoApply = true)
class ReservationStatusConverter : AttributeConverter<ReservationStatus, Int> {
    override fun convertToDatabaseColumn(attribute: ReservationStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): ReservationStatus =
        ReservationStatus.values().first { it.value == dbData }
}
