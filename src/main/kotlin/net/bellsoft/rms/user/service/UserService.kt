package net.bellsoft.rms.user.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.UnprocessableEntityException
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.dto.service.UserPatchDto
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.mapper.UserMapper
import net.bellsoft.rms.user.repository.UserRepository
import net.bellsoft.rms.user.type.UserRole
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class UserService(
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val userMapper: UserMapper,
) {
    @Transactional
    fun register(userCreateDto: UserCreateDto): UserDetailDto {
        try {
            return userMapper.toDto(userRepository.save(userCreateDto.toEntity(passwordEncoder)))
        } catch (ex: DataIntegrityViolationException) {
            throw UnprocessableEntityException("${userCreateDto.userId}(${userCreateDto.email}) 로 가입할 수 없습니다")
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
