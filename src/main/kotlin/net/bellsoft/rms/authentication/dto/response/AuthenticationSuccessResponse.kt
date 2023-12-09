package net.bellsoft.rms.authentication.dto.response

import net.bellsoft.rms.user.dto.response.UserDetailDto

data class AuthenticationSuccessResponse(
    val message: String,
    val value: UserDetailDto,
)
