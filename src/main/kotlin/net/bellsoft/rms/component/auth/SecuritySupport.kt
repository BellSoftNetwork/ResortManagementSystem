package net.bellsoft.rms.component.auth

import mu.KLogging
import net.bellsoft.rms.domain.user.User
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Component

@Component
class SecuritySupport {
    fun getCurrentUser(): User? {
        val authentication = SecurityContextHolder.getContext().authentication

        return if (authentication?.isAuthenticated == true)
            if (authentication.principal is User) {
                authentication.principal as User
            } else {
                logger.error("principal is not User: $authentication")
                null
            }
        else
            null
    }

    companion object : KLogging()
}
