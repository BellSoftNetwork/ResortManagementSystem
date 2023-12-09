package net.bellsoft.rms.component.auth

import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.component.auth.dto.UserWrapper
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.exception.UserNotFoundException
import org.springframework.data.repository.findByIdOrNull
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.Authentication
import org.springframework.security.core.context.DeferredSecurityContext
import org.springframework.security.core.context.SecurityContext
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.web.context.HttpSessionSecurityContextRepository

class SessionSecurityContextRepository(
    private val userRepository: UserRepository,
) : HttpSessionSecurityContextRepository() {
    override fun loadDeferredContext(request: HttpServletRequest?): DeferredSecurityContext {
        val deferredSecurityContext = super.loadDeferredContext(request)

        return object : DeferredSecurityContext {
            private var securityContext: SecurityContext? = null
            private var isGenerated = false

            override fun get(): SecurityContext {
                init()

                return securityContext!!
            }

            override fun isGenerated(): Boolean {
                init()

                return isGenerated
            }

            private fun init() {
                if (securityContext != null)
                    return

                val context = deferredSecurityContext.get()

                isGenerated = deferredSecurityContext.isGenerated
                context.authentication?.let {
                    try {
                        context.authentication = reloadAuthenticationFromDB(it)
                    } catch (ex: UserNotFoundException) {
                        logger.warn("현재 세션의 사용자 정보가 DB에 존재하지 않아 인증 정보 무효화 (ID: ${context.authentication.name})")
                        context.authentication.isAuthenticated = false
                    }
                }

                securityContext = context
            }
        }
    }

    override fun saveContext(context: SecurityContext, request: HttpServletRequest, response: HttpServletResponse) {
        val authentication = context.authentication
        val user = authentication?.principal

        if (user is User)
            context.authentication = UsernamePasswordAuthenticationToken(
                UserWrapper.of(user),
                authentication.credentials,
                authentication.authorities,
            )

        super.saveContext(context, request, response)
    }

    private fun reloadAuthenticationFromDB(authentication: Authentication): Authentication {
        if (authentication is UsernamePasswordAuthenticationToken) {
            val userDetails = createNewUserDetailsFromPrincipal(authentication.principal)

            return UsernamePasswordAuthenticationToken(
                userDetails,
                authentication.getCredentials(),
                userDetails.authorities,
            )
        }

        return authentication
    }

    private fun createNewUserDetailsFromPrincipal(principal: Any): UserDetails {
        val userId = when (principal) {
            is User -> {
                principal.id
            }

            is UserWrapper -> {
                principal.id
            }

            else -> {
                logger.error("알 수 없는 타입 ($principal)")
                throw UserNotFoundException("알 수 없는 유저")
            }
        }

        return userRepository.findByIdOrNull(userId) ?: throw UserNotFoundException()
    }

    companion object : KLogging()
}
