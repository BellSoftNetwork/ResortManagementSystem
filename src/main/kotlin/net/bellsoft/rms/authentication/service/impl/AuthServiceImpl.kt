package net.bellsoft.rms.authentication.service.impl

import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.authentication.service.AuthService
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class AuthServiceImpl(
    private val userRepository: UserRepository,
) : AuthService {
    override fun loadUserByUsername(username: String): UserDetails {
        val user = if (username.contains("@"))
            userRepository.findByEmail(username)
        else
            userRepository.findByUserId(username)

        return user ?: throw UserNotFoundException("$username 은 존재하지 않는 사용자입니다")
    }
}
