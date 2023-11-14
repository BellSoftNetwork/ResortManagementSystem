package net.bellsoft.rms.util

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Component

@Component
class SecurityTestSupport(
    private val userRepository: UserRepository,
) {
    private val fixture = baseFixture

    fun login(user: User? = null) = userRepository.save(user ?: fixture()).also {
        SecurityContextHolder.getContext()
            .authentication = UsernamePasswordAuthenticationToken(it, it.password, it.authorities)
    }

    fun logout() {
        SecurityContextHolder.getContext().authentication = null
    }
}
