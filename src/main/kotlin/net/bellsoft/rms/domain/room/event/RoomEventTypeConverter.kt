package net.bellsoft.rms.domain.room.event

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class RoomEventTypeConverter : AttributeConverter<RoomEventType, Int> {
    override fun convertToDatabaseColumn(attribute: RoomEventType): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): RoomEventType =
        RoomEventType.values().first { it.value == dbData }
}
