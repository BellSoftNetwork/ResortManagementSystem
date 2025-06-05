package net.bellsoft.rms.authentication.entity

import net.bellsoft.rms.user.entity.User
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.userdetails.UserDetails
import java.io.Serial
import java.io.Serializable

class UserWrapper(
    val id: Long,
    private val username: String,
    private val password: String,
    private val authorities: List<GrantedAuthority>,
    private val accountNonExpired: Boolean,
    private val accountNonLocked: Boolean,
    private val credentialsNonExpired: Boolean,
    private val enabled: Boolean,
) : Serializable, UserDetails {
    override fun getUsername(): String {
        return username
    }

    override fun getPassword(): String {
        return password
    }

    override fun getAuthorities(): List<GrantedAuthority> {
        return authorities
    }

    override fun isAccountNonExpired(): Boolean {
        return accountNonExpired
    }

    override fun isAccountNonLocked(): Boolean {
        return accountNonLocked
    }

    override fun isCredentialsNonExpired(): Boolean {
        return credentialsNonExpired
    }

    override fun isEnabled(): Boolean {
        return enabled
    }

    companion object {

        @Serial
        private const val serialVersionUID: Long = 3968529540276588739L

        fun of(user: User) = UserWrapper(
            id = user.id,
            username = user.userId,
            password = user.password,
            authorities = user.authorities.toList(),
            accountNonExpired = user.isAccountNonExpired,
            accountNonLocked = user.isAccountNonLocked,
            credentialsNonExpired = user.isCredentialsNonExpired,
            enabled = user.isEnabled,
        )
    }
}
