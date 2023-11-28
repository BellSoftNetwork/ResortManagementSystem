package net.bellsoft.rms.component.auth.dto

import net.bellsoft.rms.service.auth.dto.UserDetailDto

data class AuthenticationSuccessResponse(
    val message: String,
    val value: UserDetailDto,
)
