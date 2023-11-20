package net.bellsoft.rms.domain.user

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Table
import net.bellsoft.rms.annotation.ExcludeFromJacocoGeneratedReport
import net.bellsoft.rms.domain.base.BaseTimeEntity
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.crypto.password.PasswordEncoder
import java.io.Serial
import java.io.Serializable

@Entity
@Table(name = "user")
@SQLDelete(sql = "UPDATE user SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
class User(
    @Column(name = "email", nullable = false, unique = true, length = 100)
    var email: String,

    @Column(name = "name", nullable = false, length = 20)
    var name: String,

    @Column(name = "password", nullable = false, length = 100)
    private var password: String,

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    var status: UserStatus = UserStatus.INACTIVE,

    @Column(name = "role", nullable = false, columnDefinition = "TINYINT")
    var role: UserRole = UserRole.NORMAL,
) : Serializable, BaseTimeEntity(), UserDetails {
    @ExcludeFromJacocoGeneratedReport
    override fun getAuthorities(): MutableCollection<out GrantedAuthority> {
        return mutableListOf<GrantedAuthority>(SimpleGrantedAuthority(this.role.name))
    }

    override fun getPassword(): String {
        return password
    }

    fun changePassword(passwordEncoder: PasswordEncoder, rawPassword: String) {
        this.password = passwordEncoder.encode(rawPassword)
    }

    override fun getUsername(): String {
        return email
    }

    @ExcludeFromJacocoGeneratedReport
    override fun isAccountNonExpired(): Boolean {
        return true
    }

    @ExcludeFromJacocoGeneratedReport
    override fun isAccountNonLocked(): Boolean {
        return true
    }

    @ExcludeFromJacocoGeneratedReport
    override fun isCredentialsNonExpired(): Boolean {
        return true
    }

    @ExcludeFromJacocoGeneratedReport
    override fun isEnabled(): Boolean {
        return true
    }

    override fun toString(): String {
        return "User(id=$id, email='$email', password='$password', name='$name', role='$role', status='$status', " +
            "createdAt='$createdAt', updatedAt='$updatedAt', deletedAt=$deletedAt)"
    }

    companion object {
        @Serial
        private const val serialVersionUID: Long = -8921744413467798119L
    }
}
