package net.bellsoft.rms.component.auth.dto

data class LoginRequest(
    val email: String,
    val password: String,
)
