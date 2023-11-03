package net.bellsoft.rms.component.auth

import mu.KLogging
import org.springframework.security.authentication.AuthenticationProvider
import org.springframework.security.authentication.BadCredentialsException
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.Authentication
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Component

@Component
class JsonAuthenticationProvider(
    private val userDetailsService: UserDetailsService,
    private val passwordEncoder: PasswordEncoder,
) : AuthenticationProvider {
    override fun authenticate(authentication: Authentication): Authentication {
        val user = userDetailsService.loadUserByUsername(authentication.name)

        if (passwordEncoder.matches(authentication.credentials as String, user.password))
            return UsernamePasswordAuthenticationToken(user, user.password, user.authorities)

        throw BadCredentialsException("비밀번호 불일치")
    }

    override fun supports(authentication: Class<*>): Boolean {
        return UsernamePasswordAuthenticationToken::class.java.isAssignableFrom(authentication)
    }

    companion object : KLogging()
}
