package net.bellsoft.rms.controller.v1.admin

import mu.KLogging
import net.bellsoft.rms.controller.v1.admin.dto.AccountCreateRequest
import net.bellsoft.rms.controller.v1.admin.dto.AccountPatchRequest
import net.bellsoft.rms.controller.v1.admin.dto.AccountResponse
import net.bellsoft.rms.controller.v1.admin.dto.AccountResponses
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.admin.AdminAccountService
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/api/v1/admin/accounts")
class AccountController(
    private val adminAccountService: AdminAccountService,
) {
    @GetMapping("")
    fun getAccounts(pageable: Pageable): ResponseEntity<AccountResponses> {
        val response = adminAccountService.findAll(pageable)

        return ResponseEntity.ok(response)
    }

    @PostMapping("")
    fun createAccount(
        @AuthenticationPrincipal user: User,
        @RequestBody request: AccountCreateRequest,
    ): ResponseEntity<AccountResponse> {
        val response = adminAccountService.createAccount(user, request)

        return ResponseEntity.status(HttpStatus.CREATED).body(response)
    }

    @PatchMapping("/{id}")
    fun updateAccount(
        @PathVariable("id") id: Long,
        @AuthenticationPrincipal user: User,
        @RequestBody request: AccountPatchRequest,
    ): ResponseEntity<AccountResponse> {
        return ResponseEntity.ok(adminAccountService.updateAccount(user, id, request))
    }

    companion object : KLogging()
}
