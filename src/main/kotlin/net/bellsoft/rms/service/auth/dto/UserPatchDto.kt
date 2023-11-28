package net.bellsoft.rms.service.auth.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.security.crypto.password.PasswordEncoder

data class UserPatchDto(
    val password: JsonNullable<String>,
    val name: JsonNullable<String>,
    val isLock: JsonNullable<Boolean>,
    val role: JsonNullable<UserRole>,
) {
    fun updateEntity(entity: User, passwordEncoder: PasswordEncoder) {
        password.orElse(null)?.let { entity.changePassword(passwordEncoder, it) }
        name.orElse(null)?.let { entity.name = it }
        isLock.orElse(null)?.let { entity.status = if (it) UserStatus.INACTIVE else UserStatus.ACTIVE }
        role.orElse(null)?.let { entity.role = it }
    }
}
