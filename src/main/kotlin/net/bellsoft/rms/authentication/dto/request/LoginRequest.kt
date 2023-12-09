package net.bellsoft.rms.authentication.dto.request

data class LoginRequest(
    val username: String,
    val password: String,
)
