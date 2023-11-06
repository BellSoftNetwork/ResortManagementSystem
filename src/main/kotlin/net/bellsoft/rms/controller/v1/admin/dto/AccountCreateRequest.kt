package net.bellsoft.rms.controller.v1.admin.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole

data class AccountCreateRequest(
    val name: String,
    val email: String,
    val password: String,
    val role: UserRole,
) {
    fun toEntity(): User {
        return User(
            email = email,
            name = name,
            password = password,
            role = role,
        )
    }
}
