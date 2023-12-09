package net.bellsoft.rms.user.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.dto.request.MyPatchRequest
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.entity.User
import org.springframework.http.ResponseEntity
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestMethod
import org.springframework.web.bind.annotation.RestController

@Tag(name = "내 정보", description = "내 정보 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/my")
interface MyController {
    @Operation(summary = "로그인 계정 정보", description = "로그인 계정 정보 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @RequestMapping(method = [RequestMethod.GET, RequestMethod.POST])
    fun displayMySelf(@AuthenticationPrincipal user: User): ResponseEntity<SingleResponse<UserDetailDto>>

    @Operation(summary = "내 계정 정보 수정", description = "현재 로그인 된 계정 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @PatchMapping
    fun updateMySelf(
        @AuthenticationPrincipal user: User,
        @RequestBody request: MyPatchRequest,
    ): ResponseEntity<SingleResponse<UserDetailDto>>
}
