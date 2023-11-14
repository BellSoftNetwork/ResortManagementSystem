package net.bellsoft.rms.service.admin

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.service.admin.dto.AccountCreateDto
import net.bellsoft.rms.service.admin.dto.AccountPatchDto
import net.bellsoft.rms.service.auth.dto.UserDto
import net.bellsoft.rms.service.common.dto.EntityListDto
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class AdminAccountService(
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
) {
    fun findAll(pageable: Pageable): EntityListDto<UserDto> {
        return EntityListDto.of(userRepository.findAll(pageable), UserDto::of)
    }

    @Transactional
    fun createAccount(accountCreateDto: AccountCreateDto): UserDto {
        return UserDto.of(userRepository.save(accountCreateDto.toEntity(passwordEncoder)))
    }

    @Transactional
    fun updateAccount(
        updateUserId: Long,
        accountPatchDto: AccountPatchDto,
    ): UserDto {
        val updateUser = userRepository.findByIdOrNull(updateUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        accountPatchDto.updateEntity(updateUser, passwordEncoder)

        return UserDto.of(userRepository.save(updateUser))
    }

    fun isUpdatableAccount(requestUser: User, targetUserId: Long): Boolean {
        if (requestUser.role == UserRole.SUPER_ADMIN) return true

        val targetUser = userRepository.findByIdOrNull(targetUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        return targetUser.role < UserRole.ADMIN
    }
}
