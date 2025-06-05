package net.bellsoft.rms.user.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Table
import jakarta.persistence.UniqueConstraint
import net.bellsoft.rms.common.annotation.ExcludeFromJacocoGeneratedReport
import net.bellsoft.rms.common.entity.BaseTimeEntity
import net.bellsoft.rms.user.type.UserRole
import net.bellsoft.rms.user.type.UserStatus
import org.hibernate.annotations.Comment
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.crypto.password.PasswordEncoder
import java.io.Serial
import java.io.Serializable

@Entity
@Table(
    name = "user",
    uniqueConstraints = [
        UniqueConstraint(name = "uc_user_user_id", columnNames = ["user_id", "deleted_at"]),
        UniqueConstraint(name = "uc_user_email", columnNames = ["email", "deleted_at"]),
    ],
)
@SQLDelete(sql = "UPDATE user SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTimeEntity.SOFT_DELETE_CONDITION)
@Comment("사용자")
class User(
    @Column(name = "user_id", nullable = false, length = 30)
    @Comment("계정 ID")
    var userId: String,

    @Column(name = "email", nullable = true, length = 100)
    @Comment("이메일")
    var email: String?,

    @Column(name = "name", nullable = false, length = 20)
    @Comment("이름")
    var name: String,

    @Column(name = "password", nullable = false, length = 100)
    @Comment("비밀번호")
    private var password: String,

    @Column(name = "status", nullable = false, columnDefinition = "TINYINT")
    @Comment("계정 상태")
    var status: UserStatus = UserStatus.INACTIVE,

    @Column(name = "role", nullable = false, columnDefinition = "TINYINT")
    @Comment("계정 권한")
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
        return userId
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
