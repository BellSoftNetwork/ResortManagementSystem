package net.bellsoft.rms.service.auth

import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.exception.UnprocessableEntityException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.service.auth.dto.UserDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service

@Service
class AuthService(
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
) : UserDetailsService {
    @Throws(UsernameNotFoundException::class)
    override fun loadUserByUsername(username: String): UserDetails {
        return userRepository.findByEmail(username)
            ?: throw UserNotFoundException("$username 은 존재하지 않는 사용자입니다")
    }

    fun register(userRegistrationRequest: UserRegistrationRequest): UserDto {
        try {
            return UserDto.of(userRepository.save(userRegistrationRequest.toEntity(passwordEncoder)))
        } catch (ex: DataIntegrityViolationException) {
            throw UnprocessableEntityException("${userRegistrationRequest.email} 로 가입할 수 없습니다")
        }
    }
}
