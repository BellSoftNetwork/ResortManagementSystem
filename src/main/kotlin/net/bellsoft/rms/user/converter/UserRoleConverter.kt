package net.bellsoft.rms.user.converter

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter
import net.bellsoft.rms.user.type.UserRole

@Converter(autoApply = true)
class UserRoleConverter : AttributeConverter<UserRole, Int> {
    override fun convertToDatabaseColumn(attribute: UserRole): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): UserRole = UserRole.values().first { it.value == dbData }
}
