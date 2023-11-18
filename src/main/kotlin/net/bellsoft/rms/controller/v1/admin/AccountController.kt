package net.bellsoft.rms.controller.v1.admin

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import mu.KLogging
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.admin.dto.AccountCreateRequest
import net.bellsoft.rms.controller.v1.admin.dto.AccountPatchRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.exception.PermissionRequiredDataException
import net.bellsoft.rms.service.auth.AuthService
import net.bellsoft.rms.service.auth.dto.UserDto
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.annotation.Secured
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "계정 (어드민)", description = "어드민용 계정 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@Secured("ADMIN", "SUPER_ADMIN")
@RequestMapping("/api/v1/admin/accounts")
class AccountController(
    private val authService: AuthService,
) {
    @Operation(summary = "계정 리스트", description = "계정 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getAccounts(pageable: Pageable) = ListResponse
        .of((authService.findAll(pageable)))
        .toResponseEntity()

    @Operation(summary = "계정 생성", description = "신규 계정 추가")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @PostMapping
    fun createAccount(
        @AuthenticationPrincipal user: User,
        @RequestBody @Valid
        request: AccountCreateRequest,
    ): ResponseEntity<SingleResponse<UserDto>> {
        if (request.role >= UserRole.ADMIN && user.role != UserRole.SUPER_ADMIN)
            throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")

        return SingleResponse
            .of(authService.createAccount(request.toDto()))
            .toResponseEntity(HttpStatus.CREATED)
    }

    @Operation(summary = "계정 수정", description = "기존 계정 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @PatchMapping("/{id}")
    fun updateAccount(
        @PathVariable("id") id: Long,
        @AuthenticationPrincipal user: User,
        @RequestBody @Valid
        request: AccountPatchRequest,
    ): ResponseEntity<SingleResponse<UserDto>> {
        request.role?.let {
            if (it >= UserRole.ADMIN && user.role != UserRole.SUPER_ADMIN)
                throw PermissionRequiredDataException("관리자 이상 권한 설정 시 최고 관리자 권한 필요")
        }
        if (!authService.isUpdatableAccount(user, id))
            throw PermissionRequiredDataException("동일 또는 상위 권한 계정 정보 수정 불가")

        return SingleResponse
            .of(authService.updateAccount(id, request.toDto()))
            .toResponseEntity()
    }

    companion object : KLogging()
}
