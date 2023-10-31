package net.bellsoft.rms.domain.room

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class RoomStatusConverter : AttributeConverter<RoomStatus, Int> {
    override fun convertToDatabaseColumn(attribute: RoomStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): RoomStatus = RoomStatus.values().first { it.value == dbData }
}
