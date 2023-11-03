package net.bellsoft.rms.controller.v1.dto

import net.bellsoft.rms.domain.user.User

data class WhoAmIResponse(val email: String) {
    companion object {
        fun of(user: User) = WhoAmIResponse(user.email)
    }
}
