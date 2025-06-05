package net.bellsoft.rms.user.repository

import net.bellsoft.rms.user.entity.User
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface UserRepository : JpaRepository<User, Long>, UserCustomRepository {
    fun findByUserId(userId: String): User?
    fun findByEmail(email: String): User?
    fun existsByEmail(email: String): Boolean
}
