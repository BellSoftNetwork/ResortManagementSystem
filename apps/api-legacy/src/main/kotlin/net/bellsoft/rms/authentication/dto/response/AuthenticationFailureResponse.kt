package net.bellsoft.rms.authentication.dto.response

data class AuthenticationFailureResponse(
    val username: String,
    val message: String,
)
