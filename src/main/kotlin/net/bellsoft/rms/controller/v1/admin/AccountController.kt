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
import net.bellsoft.rms.service.admin.AdminAccountService
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@Tag(name = "계정", description = "계정 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/admin/accounts")
class AccountController(
    private val adminAccountService: AdminAccountService,
) {
    @Operation(summary = "계정 리스트", description = "계정 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getAccounts(pageable: Pageable) = ListResponse
        .of((adminAccountService.findAll(pageable)))
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
    ) = SingleResponse
        .of(adminAccountService.createAccount(user, request.toDto()))
        .toResponseEntity(HttpStatus.CREATED)

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
    ) = SingleResponse
        .of(adminAccountService.updateAccount(user, id, request.toDto()))
        .toResponseEntity()

    companion object : KLogging()
}
