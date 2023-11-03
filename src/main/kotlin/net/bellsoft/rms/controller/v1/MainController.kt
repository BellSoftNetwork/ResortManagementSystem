package net.bellsoft.rms.controller.v1

import net.bellsoft.rms.controller.v1.dto.EnvResponse
import net.bellsoft.rms.controller.v1.dto.WhoAmIResponse
import net.bellsoft.rms.domain.user.User
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/api/v1")
class MainController {
    @GetMapping("/env")
    fun displayEnv() = EnvResponse.of("Resort Management System", "RMS", "v1.0")

    @GetMapping("/whoami")
    fun displayMySelf(@AuthenticationPrincipal user: User) = WhoAmIResponse.of(user)
}
