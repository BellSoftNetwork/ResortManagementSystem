package net.bellsoft.rms.authentication.controller.impl

import mu.KLogging
import net.bellsoft.rms.authentication.controller.AuthController
import net.bellsoft.rms.authentication.dto.request.UserRegistrationRequest
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.service.UserService
import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.RestController

@RestController
class AuthControllerImpl(
    private val userService: UserService,
) : AuthController {
    override fun registerUser(
        userRegistrationRequest: UserRegistrationRequest,
    ) = SingleResponse
        .of((userService.register(UserCreateDto.of(userRegistrationRequest))))
        .toResponseEntity(HttpStatus.CREATED)

    companion object : KLogging()
}
