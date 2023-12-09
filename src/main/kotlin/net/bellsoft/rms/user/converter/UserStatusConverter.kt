package net.bellsoft.rms.user.converter

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter
import net.bellsoft.rms.user.type.UserStatus

@Converter(autoApply = true)
class UserStatusConverter : AttributeConverter<UserStatus, Int> {
    override fun convertToDatabaseColumn(attribute: UserStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): UserStatus = UserStatus.values().first { it.value == dbData }
}
