package net.bellsoft.rms.service.auth.dto

import net.bellsoft.rms.controller.v1.admin.dto.AdminUserPatchRequest
import net.bellsoft.rms.controller.v1.my.dto.MyPatchRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.security.crypto.password.PasswordEncoder

data class UserPatchDto(
    val password: JsonNullable<String> = JsonNullable.undefined(),
    val name: JsonNullable<String> = JsonNullable.undefined(),
    val userId: JsonNullable<String> = JsonNullable.undefined(),
    val email: JsonNullable<String?> = JsonNullable.undefined(),
    val isLock: JsonNullable<Boolean> = JsonNullable.undefined(),
    val role: JsonNullable<UserRole> = JsonNullable.undefined(),
) {
    fun updateEntity(entity: User, passwordEncoder: PasswordEncoder) {
        if (password.isPresent) entity.changePassword(passwordEncoder, password.get())
        if (name.isPresent) entity.name = name.get()
        if (userId.isPresent) entity.userId = userId.get()
        if (email.isPresent) entity.email = email.get()
        if (isLock.isPresent) entity.status = if (isLock.get()) UserStatus.INACTIVE else UserStatus.ACTIVE
        if (role.isPresent) entity.role = role.get()
    }

    companion object {
        fun of(dto: AdminUserPatchRequest) = UserPatchDto(
            userId = dto.userId,
            email = dto.email,
            password = dto.password,
            name = dto.name,
            isLock = dto.isLock,
            role = dto.role,
        )

        fun of(dto: MyPatchRequest) = UserPatchDto(
            email = dto.email,
            password = dto.password,
        )
    }
}
