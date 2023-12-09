package net.bellsoft.rms.user.controller.impl

import mu.KLogging
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.common.exception.PermissionRequiredDataException
import net.bellsoft.rms.user.controller.AdminAccountController
import net.bellsoft.rms.user.dto.request.AdminUserCreateRequest
import net.bellsoft.rms.user.dto.request.AdminUserPatchRequest
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.dto.service.UserPatchDto
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.service.UserService
import net.bellsoft.rms.user.type.UserRole
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminAccountControllerImpl(
    private val userService: UserService,
) : AdminAccountController {
    override fun getAccounts(pageable: Pageable) = ListResponse
        .of((userService.findAll(pageable)))
        .toResponseEntity()

    override fun createAccount(
        user: User,
        request: AdminUserCreateRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>> {
        if (request.role >= UserRole.ADMIN && user.role != UserRole.SUPER_ADMIN)
            throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

        return SingleResponse
            .of(userService.register(UserCreateDto.of(request)))
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
        if (!userService.isUpdatableAccount(user, id))
            throw PermissionRequiredDataException("동일 또는 상위 권한 계정 정보 수정 불가")

        return SingleResponse
            .of(userService.updateAccount(id, UserPatchDto.of(request)))
            .toResponseEntity()
    }

    companion object : KLogging()
}
