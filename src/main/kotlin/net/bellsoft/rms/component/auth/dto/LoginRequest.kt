package net.bellsoft.rms.component.auth.dto

data class LoginRequest(
    val username: String,
    val password: String,
)
