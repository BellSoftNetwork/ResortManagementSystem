package net.bellsoft.rms.user.controller.impl

import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.controller.MyController
import net.bellsoft.rms.user.dto.request.MyPatchRequest
import net.bellsoft.rms.user.dto.service.UserPatchDto
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.mapper.UserMapper
import net.bellsoft.rms.user.service.UserService
import org.springframework.web.bind.annotation.RestController

@RestController
class MyControllerImpl(
    private val userService: UserService,
    private val userMapper: UserMapper,
) : MyController {
    override fun displayMySelf(user: User) = SingleResponse
        .of(userMapper.toDto(user))
        .toResponseEntity()

    override fun updateMySelf(user: User, request: MyPatchRequest) = SingleResponse
        .of(userService.updateAccount(user.id, UserPatchDto.of(request)))
        .toResponseEntity()
}
