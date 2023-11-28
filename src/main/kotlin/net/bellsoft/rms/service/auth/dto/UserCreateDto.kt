package net.bellsoft.rms.service.auth.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import org.springframework.security.crypto.password.PasswordEncoder

data class UserCreateDto(
    val name: String,
    val email: String,
    val password: String,
    val status: UserStatus = UserStatus.ACTIVE,
    val role: UserRole = UserRole.NORMAL,
) {
    fun toEntity(passwordEncoder: PasswordEncoder): User {
        return User(
            email = email,
            name = name,
            password = passwordEncoder.encode(password),
            status = status,
            role = role,
        )
    }
}
