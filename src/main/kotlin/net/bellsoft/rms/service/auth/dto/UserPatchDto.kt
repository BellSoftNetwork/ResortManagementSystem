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
    val isLock: JsonNullable<Boolean> = JsonNullable.undefined(),
    val role: JsonNullable<UserRole> = JsonNullable.undefined(),
) {
    fun updateEntity(entity: User, passwordEncoder: PasswordEncoder) {
        password.orElse(null)?.let { entity.changePassword(passwordEncoder, it) }
        name.orElse(null)?.let { entity.name = it }
        isLock.orElse(null)?.let { entity.status = if (it) UserStatus.INACTIVE else UserStatus.ACTIVE }
        role.orElse(null)?.let { entity.role = it }
    }

    companion object {
        fun of(dto: AdminUserPatchRequest) = UserPatchDto(
            password = dto.password,
            name = dto.name,
            isLock = dto.isLock,
            role = dto.role,
        )

        fun of(dto: MyPatchRequest) = UserPatchDto(
            password = dto.password,
        )
    }
}
