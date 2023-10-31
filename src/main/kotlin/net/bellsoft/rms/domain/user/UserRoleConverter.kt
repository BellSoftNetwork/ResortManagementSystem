package net.bellsoft.rms.domain.user

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class UserRoleConverter : AttributeConverter<UserRole, Int> {
    override fun convertToDatabaseColumn(attribute: UserRole): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): UserRole = UserRole.values().first { it.value == dbData }
}
