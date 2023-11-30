package net.bellsoft.rms.controller.v1.admin

import mu.KLogging
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.admin.dto.AdminUserCreateRequest
import net.bellsoft.rms.controller.v1.admin.dto.AdminUserPatchRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.PermissionRequiredDataException
import net.bellsoft.rms.service.auth.AuthService
import net.bellsoft.rms.service.auth.dto.UserCreateDto
import net.bellsoft.rms.service.auth.dto.UserDetailDto
import net.bellsoft.rms.service.auth.dto.UserPatchDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AccountControllerImpl(
    private val authService: AuthService,
) : AccountController {
    override fun getAccounts(pageable: Pageable) = ListResponse
        .of((authService.findAll(pageable)))
        .toResponseEntity()

    override fun createAccount(
        user: User,
        request: AdminUserCreateRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>> {
        if (request.role >= UserRole.ADMIN && user.role != UserRole.SUPER_ADMIN)
            throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

        return SingleResponse
            .of(authService.register(UserCreateDto.of(request)))
            .toResponseEntity(HttpStatus.CREATED)
    }

    override fun updateAccount(
        id: Long,
        user: User,
        request: AdminUserPatchRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>> {
        request.role.orElse(null)?.let {
            if (it >= UserRole.ADMIN && user.role != UserRole.SUPER_ADMIN)
                throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")
        }
        if (!authService.isUpdatableAccount(user, id))
            throw PermissionRequiredDataException("동일 또는 상위 권한 계정 정보 수정 불가")

        return SingleResponse
            .of(authService.updateAccount(id, UserPatchDto.of(request)))
            .toResponseEntity()
    }

    companion object : KLogging()
}
