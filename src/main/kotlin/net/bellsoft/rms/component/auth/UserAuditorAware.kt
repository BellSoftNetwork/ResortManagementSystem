package net.bellsoft.rms.component.auth

import net.bellsoft.rms.domain.user.User
import org.springframework.data.domain.AuditorAware
import org.springframework.stereotype.Component
import java.util.*

@Component
class UserAuditorAware(
    private val securitySupport: SecuritySupport,
) : AuditorAware<User> {
    override fun getCurrentAuditor(): Optional<User> {
        return Optional.of(
            securitySupport.getCurrentUser()
                ?: return Optional.empty(),
        )
    }
}
