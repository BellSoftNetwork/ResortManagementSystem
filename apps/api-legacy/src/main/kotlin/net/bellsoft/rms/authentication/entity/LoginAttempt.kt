package net.bellsoft.rms.authentication.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Index
import jakarta.persistence.Table
import net.bellsoft.rms.common.entity.BaseEntity
import org.hibernate.annotations.Comment
import java.time.LocalDateTime

@Entity
@Table(
    name = "login_attempts",
    indexes = [
        Index(name = "idx_login_attempts_username_attempt_at", columnList = "username, attempt_at"),
        Index(name = "idx_login_attempts_ip_address_attempt_at", columnList = "ip_address, attempt_at"),
        Index(
            name = "idx_login_attempts_username_ip_address_attempt_at",
            columnList = "username, ip_address, attempt_at",
        ),
    ],
)
class LoginAttempt(
    @Column(name = "username", nullable = false, length = 50)
    @Comment("계정 ID")
    var username: String,

    @Column(name = "ip_address", nullable = false, length = 50)
    @Comment("IP 주소")
    var ipAddress: String,

    @Column(name = "successful", nullable = false)
    @Comment("로그인 성공 여부")
    var successful: Boolean,

    @Column(name = "attempt_at", nullable = false)
    @Comment("로그인 시도 시각")
    var attemptAt: LocalDateTime = LocalDateTime.now(),

    @Column(name = "os_info", length = 50)
    @Comment("운영체제 정보")
    var osInfo: String? = null,

    @Column(name = "language_info", length = 50)
    @Comment("언어 설정 정보")
    var languageInfo: String? = null,

    @Column(name = "user_agent", length = 500)
    @Comment("사용자 에이전트 정보")
    var userAgent: String? = null,

    @Column(name = "device_fingerprint", length = 50)
    @Comment("디바이스 정보 해시")
    var deviceFingerprint: String? = null,
) : BaseEntity() {
    companion object {
        private const val serialVersionUID: Long = -903809232978889145L
    }
}
