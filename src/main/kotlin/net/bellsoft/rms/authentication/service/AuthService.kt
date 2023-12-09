package net.bellsoft.rms.authentication.service

import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class AuthService(
    private val userRepository: UserRepository,
) : UserDetailsService {
    @Throws(UsernameNotFoundException::class)
    override fun loadUserByUsername(username: String): UserDetails {
        val user = if (username.contains("@"))
            userRepository.findByEmail(username)
        else
            userRepository.findByUserId(username)

        return user ?: throw UserNotFoundException("$username 은 존재하지 않는 사용자입니다")
    }
}
