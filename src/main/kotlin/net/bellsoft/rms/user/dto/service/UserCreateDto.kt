package net.bellsoft.rms.user.dto.service

import net.bellsoft.rms.authentication.dto.request.UserRegistrationRequest
import net.bellsoft.rms.user.dto.request.AdminUserCreateRequest
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.type.UserRole
import net.bellsoft.rms.user.type.UserStatus
import org.springframework.security.crypto.password.PasswordEncoder

data class UserCreateDto(
    val name: String,
    val userId: String,
    val email: String?,
    val password: String,
    val status: UserStatus = UserStatus.ACTIVE,
    val role: UserRole = UserRole.NORMAL,
) {
    fun toEntity(passwordEncoder: PasswordEncoder): User {
        return User(
            userId = userId,
            email = email,
            name = name,
            password = passwordEncoder.encode(password),
            status = status,
            role = role,
        )
    }

    companion object {
        fun of(dto: AdminUserCreateRequest) = UserCreateDto(
            name = dto.name,
            userId = dto.userId,
            email = dto.email,
            password = dto.password,
            role = dto.role,
        )

        fun of(dto: UserRegistrationRequest) = UserCreateDto(
            name = dto.name,
            userId = dto.userId,
            email = dto.email,
            password = dto.password,
        )
    }
}
