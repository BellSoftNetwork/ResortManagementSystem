package net.bellsoft.rms.common.controller

import org.springframework.boot.web.servlet.error.ErrorController
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping

@Controller
class ErrorController : ErrorController {
    @Suppress("SpringMVCViewInspection")
    @GetMapping(ERROR_PATH)
    fun redirectRoot() = "index.html"

    companion object {
        private const val ERROR_PATH = "/error"
    }
}
