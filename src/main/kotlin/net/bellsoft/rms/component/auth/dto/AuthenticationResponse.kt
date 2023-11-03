package net.bellsoft.rms.component.auth.dto

import org.springframework.security.core.GrantedAuthority

data class AuthenticationResponse(
    val email: String,
    val message: String,
    val authorities: Collection<GrantedAuthority> = listOf(),
)
