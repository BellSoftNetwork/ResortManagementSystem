package net.bellsoft.rms.service.auth.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import org.springframework.security.crypto.password.PasswordEncoder

data class AccountCreateDto(
    val name: String,
    val email: String,
    val password: String,
    val role: UserRole,
) {
    fun toEntity(passwordEncoder: PasswordEncoder): User {
        return User(
            email = email,
            name = name,
            password = passwordEncoder.encode(password),
            role = role,
        )
    }
}
