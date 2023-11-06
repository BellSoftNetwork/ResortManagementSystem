package net.bellsoft.rms.service.admin

import net.bellsoft.rms.controller.v1.admin.dto.AccountCreateRequest
import net.bellsoft.rms.controller.v1.admin.dto.AccountPatchRequest
import net.bellsoft.rms.controller.v1.admin.dto.AccountResponse
import net.bellsoft.rms.controller.v1.admin.dto.AccountResponses
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.domain.user.UserStatus
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.PermissionRequiredDataException
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
    fun findAll(pageable: Pageable): AccountResponses {
        val users = userRepository.findAll(pageable)

        return AccountResponses.of(users)
    }

    @Transactional
    fun createAccount(
        requestUser: User,
        accountCreateRequest: AccountCreateRequest,
    ): AccountResponse {
        if (requestUser.role < UserRole.ADMIN)
            throw PermissionRequiredDataException("관리자 권한 필요")
        if (accountCreateRequest.role >= UserRole.ADMIN && requestUser.role != UserRole.SUPER_ADMIN)
            throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

        return AccountResponse.of(userRepository.save(accountCreateRequest.toEntity()))
    }

    @Transactional
    fun updateAccount(
        requestUser: User,
        updateUserId: Long,
        accountPatchRequest: AccountPatchRequest,
    ): AccountResponse {
        if (requestUser.role < UserRole.ADMIN)
            throw PermissionRequiredDataException("관리자 권한 필요")

        val updateUser = userRepository.findByIdOrNull(updateUserId)
            ?: throw DataNotFoundException("존재하지 않는 사용자")

        if (requestUser.role == UserRole.ADMIN && updateUser.role >= UserRole.ADMIN)
            throw PermissionRequiredDataException("동일 또는 상위 권한 계정 정보 수정 불가")

        return AccountResponse.of(
            updateUser.apply {
                accountPatchRequest.password?.let { changePassword(passwordEncoder, it) }
                accountPatchRequest.name?.let { name = it }
                accountPatchRequest.isLock?.let { status = if (it) UserStatus.INACTIVE else UserStatus.ACTIVE }
                accountPatchRequest.role?.let {
                    if (it >= UserRole.ADMIN && requestUser.role != UserRole.SUPER_ADMIN)
                        throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

                    role = it
                }
            },
        )
    }
}
