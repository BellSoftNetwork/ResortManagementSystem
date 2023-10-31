package net.bellsoft.rms.domain.user

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class UserStatusConverter : AttributeConverter<UserStatus, Int> {
    override fun convertToDatabaseColumn(attribute: UserStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): UserStatus = UserStatus.values().first { it.value == dbData }
}
