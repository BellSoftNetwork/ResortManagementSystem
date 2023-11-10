package net.bellsoft.rms.service.admin

import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.PermissionRequiredDataException
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
    fun createAccount(
        requestUser: User,
        accountCreateDto: AccountCreateDto,
    ): UserDto {
        if (requestUser.role < UserRole.ADMIN)
            throw PermissionRequiredDataException("관리자 권한 필요")
        if (accountCreateDto.role >= UserRole.ADMIN && requestUser.role != UserRole.SUPER_ADMIN)
            throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

        return UserDto.of(userRepository.save(accountCreateDto.toEntity(passwordEncoder)))
    }

    @Transactional
    fun updateAccount(
        requestUser: User,
        updateUserId: Long,
        accountPatchDto: AccountPatchDto,
    ): UserDto {
        if (requestUser.role < UserRole.ADMIN)
            throw PermissionRequiredDataException("관리자 권한 필요")

        val updateUser = userRepository.findByIdOrNull(updateUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        if (requestUser.role == UserRole.ADMIN && updateUser.role >= UserRole.ADMIN)
            throw PermissionRequiredDataException("동일 또는 상위 권한 계정 정보 수정 불가")

        accountPatchDto.role?.let {
            if (it >= UserRole.ADMIN && requestUser.role != UserRole.SUPER_ADMIN)
                throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")
        }

        accountPatchDto.updateEntity(updateUser, passwordEncoder)

        return UserDto.of(userRepository.save(updateUser))
    }
}
