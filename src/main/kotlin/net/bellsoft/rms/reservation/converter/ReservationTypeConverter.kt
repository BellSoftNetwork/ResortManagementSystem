package net.bellsoft.rms.reservation.converter

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter
import net.bellsoft.rms.reservation.type.ReservationType

@Converter(autoApply = true)
class ReservationTypeConverter : AttributeConverter<ReservationType, Int> {
    override fun convertToDatabaseColumn(attribute: ReservationType): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): ReservationType =
        ReservationType.values().first { it.value == dbData }
}
