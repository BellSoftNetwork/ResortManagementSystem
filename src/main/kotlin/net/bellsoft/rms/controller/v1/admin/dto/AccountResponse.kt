package net.bellsoft.rms.controller.v1.admin.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import java.time.LocalDateTime

data class AccountResponse(
    val id: Long,
    val email: String,
    val name: String,
    val role: UserRole,
    val createdAt: LocalDateTime,
) {
    companion object {
        fun of(user: User) = AccountResponse(
            id = user.id,
            email = user.email,
            name = user.name,
            role = user.role,
            createdAt = user.createdAt,
        )
    }
}
