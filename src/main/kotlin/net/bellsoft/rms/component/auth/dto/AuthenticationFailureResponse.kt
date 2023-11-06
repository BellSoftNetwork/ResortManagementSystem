package net.bellsoft.rms.component.auth.dto

data class AuthenticationFailureResponse(
    val email: String,
    val message: String,
)
