package net.bellsoft.rms.controller.v1.my

import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.my.dto.MyPatchRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.mapper.model.UserMapper
import net.bellsoft.rms.service.auth.AuthService
import net.bellsoft.rms.service.auth.dto.UserPatchDto
import org.springframework.web.bind.annotation.RestController

@RestController
class MyControllerImpl(
    private val authService: AuthService,
    private val userMapper: UserMapper,
) : MyController {
    override fun displayMySelf(user: User) = SingleResponse
        .of(userMapper.toDto(user))
        .toResponseEntity()

    override fun updateMySelf(user: User, request: MyPatchRequest) = SingleResponse
        .of(authService.updateAccount(user.id, UserPatchDto.of(request)))
        .toResponseEntity()
}
