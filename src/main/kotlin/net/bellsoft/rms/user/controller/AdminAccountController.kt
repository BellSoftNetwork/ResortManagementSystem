package net.bellsoft.rms.user.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.dto.request.AdminUserCreateRequest
import net.bellsoft.rms.user.dto.request.AdminUserPatchRequest
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable
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
interface AdminAccountController {
    @Operation(summary = "계정 리스트", description = "계정 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getAccounts(pageable: Pageable): ResponseEntity<ListResponse<UserDetailDto>>

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
        request: AdminUserCreateRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>>

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
        request: AdminUserPatchRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>>
}
