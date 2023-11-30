package net.bellsoft.rms.controller.v1.auth

import mu.KLogging
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.service.auth.AuthService
import net.bellsoft.rms.service.auth.dto.UserCreateDto
import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.RestController

@RestController
class AuthControllerImpl(
    private val authService: AuthService,
) : AuthController {
    override fun registerUser(
        userRegistrationRequest: UserRegistrationRequest,
    ) = SingleResponse
        .of((authService.register(UserCreateDto.of(userRegistrationRequest))))
        .toResponseEntity(HttpStatus.CREATED)

    companion object : KLogging()
}
