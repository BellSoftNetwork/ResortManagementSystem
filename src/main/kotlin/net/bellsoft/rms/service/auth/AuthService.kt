package net.bellsoft.rms.service.auth

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.UnprocessableEntityException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.mapper.model.UserMapper
import net.bellsoft.rms.service.auth.dto.UserCreateDto
import net.bellsoft.rms.service.auth.dto.UserDetailDto
import net.bellsoft.rms.service.auth.dto.UserPatchDto
import net.bellsoft.rms.service.common.dto.EntityListDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class AuthService(
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val userMapper: UserMapper,
) : UserDetailsService {
    @Throws(UsernameNotFoundException::class)
    override fun loadUserByUsername(username: String): UserDetails {
        return userRepository.findByEmail(username)
            ?: throw UserNotFoundException("$username 은 존재하지 않는 사용자입니다")
    }

    @Transactional
    fun register(userCreateDto: UserCreateDto): UserDetailDto {
        try {
            return userMapper.toDto(userRepository.save(userCreateDto.toEntity(passwordEncoder)))
        } catch (ex: DataIntegrityViolationException) {
            throw UnprocessableEntityException("${userCreateDto.email} 로 가입할 수 없습니다")
        }
    }

    fun findAll(pageable: Pageable): EntityListDto<UserDetailDto> {
        return EntityListDto.of(userRepository.findAll(pageable), userMapper::toDto)
    }

    @Transactional
    fun updateAccount(
        updateUserId: Long,
        userPatchDto: UserPatchDto,
    ): UserDetailDto {
        val updateUser = userRepository.findByIdOrNull(updateUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        userPatchDto.updateEntity(updateUser, passwordEncoder)

        return userMapper.toDto(userRepository.save(updateUser))
    }

    fun isUpdatableAccount(requestUser: User, targetUserId: Long): Boolean {
        if (requestUser.role == UserRole.SUPER_ADMIN) return true

        val targetUser = userRepository.findByIdOrNull(targetUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        return targetUser.role < UserRole.ADMIN
    }
}
