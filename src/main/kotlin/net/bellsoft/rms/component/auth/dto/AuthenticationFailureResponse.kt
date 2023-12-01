package net.bellsoft.rms.component.auth.dto

data class AuthenticationFailureResponse(
    val username: String,
    val message: String,
)
