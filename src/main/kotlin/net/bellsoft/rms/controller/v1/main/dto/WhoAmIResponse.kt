package net.bellsoft.rms.controller.v1.main.dto

import net.bellsoft.rms.domain.user.User
import java.time.LocalDateTime

data class WhoAmIResponse(
    val id: Long,
    val email: String,
    val name: String,
    val role: String,
    val createdAt: LocalDateTime,
) {
    companion object {
        fun of(user: User) = WhoAmIResponse(
            id = user.id,
            email = user.email,
            name = user.name,
            role = user.role.name,
            createdAt = user.createdAt,
        )
    }
}
