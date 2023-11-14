package net.bellsoft.rms.component.auth

import net.bellsoft.rms.domain.user.User
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Component

@Component
class SecuritySupport {
    fun getCurrentUser(): User? {
        val authentication = SecurityContextHolder.getContext().authentication

        return if (authentication?.isAuthenticated == true)
            authentication.principal as User
        else
            null
    }
}
