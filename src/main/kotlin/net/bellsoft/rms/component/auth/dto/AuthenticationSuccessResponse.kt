package net.bellsoft.rms.component.auth.dto

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.auth.dto.UserDto

data class AuthenticationSuccessResponse(
    val message: String,
    val value: UserDto,
) {
    companion object {
        fun of(user: User, message: String) = AuthenticationSuccessResponse(
            message = message,
            value = UserDto.of(user),
        )
    }
}
