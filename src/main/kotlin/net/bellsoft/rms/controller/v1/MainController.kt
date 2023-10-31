package net.bellsoft.rms.controller.v1

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/api/v1")
class MainController {
    @GetMapping("/")
    fun displayIndex(): String {
        return "Resort Management System (RMS) v1.0"
    }
}
