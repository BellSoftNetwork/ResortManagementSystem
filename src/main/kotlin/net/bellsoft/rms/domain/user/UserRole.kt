package net.bellsoft.rms.domain.user

enum class UserRole(val value: Int) {
    NORMAL(0),
    ADMIN(100),
    SUPER_ADMIN(127),
}
