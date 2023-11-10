package net.bellsoft.rms.service.admin.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import org.springframework.security.crypto.password.PasswordEncoder

data class AccountPatchDto(
    val password: String? = null,
    val name: String? = null,
    val isLock: Boolean? = null,
    val role: UserRole? = null,
) {
    fun updateEntity(entity: User, passwordEncoder: PasswordEncoder) {
        password?.let { entity.changePassword(passwordEncoder, it) }
        name?.let { entity.name = it }
        isLock?.let { entity.status = if (it) UserStatus.INACTIVE else UserStatus.ACTIVE }
        role?.let { entity.role = it }
    }
}
