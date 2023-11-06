package net.bellsoft.rms.component.auth.dto

import net.bellsoft.rms.domain.user.User
import java.time.LocalDateTime

data class AuthenticationSuccessResponse(
    val id: Long,
    val email: String,
    val name: String,
    val role: String,
    val createdAt: LocalDateTime,
    val message: String,
) {
    companion object {
        fun of(user: User, message: String) = AuthenticationSuccessResponse(
            id = user.id,
            email = user.email,
            name = user.name,
            role = user.role.name,
            createdAt = user.createdAt,
            message = message,
        )
    }
}
