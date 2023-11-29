package net.bellsoft.rms.controller.v1.my

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.my.dto.MyPatchRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.mapper.model.UserMapper
import net.bellsoft.rms.service.auth.AuthService
import net.bellsoft.rms.service.auth.dto.UserPatchDto
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
class MyController(
    private val authService: AuthService,
    private val userMapper: UserMapper,
) {
    @Operation(summary = "로그인 계정 정보", description = "로그인 계정 정보 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @RequestMapping(method = [RequestMethod.GET, RequestMethod.POST])
    fun displayMySelf(@AuthenticationPrincipal user: User) = SingleResponse
        .of(userMapper.toDto(user))
        .toResponseEntity()

    @Operation(summary = "내 계정 정보 수정", description = "현재 로그인 된 계정 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @PatchMapping
    fun updateMySelf(@AuthenticationPrincipal user: User, @RequestBody request: MyPatchRequest) = SingleResponse
        .of(authService.updateAccount(user.id, UserPatchDto.of(request)))
        .toResponseEntity()
}
