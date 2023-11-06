package net.bellsoft.rms.controller.v1.admin.dto

import net.bellsoft.rms.domain.user.UserRole

data class AccountPatchRequest(
    val password: String? = null,
    val name: String? = null,
    val isLock: Boolean? = null,
    val role: UserRole? = null,
)
