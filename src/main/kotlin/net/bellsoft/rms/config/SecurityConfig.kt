package net.bellsoft.rms.config

import com.fasterxml.jackson.databind.ObjectMapper
import net.bellsoft.rms.component.auth.JsonAuthenticationFailureHandler
import net.bellsoft.rms.component.auth.JsonAuthenticationSuccessHandler
import net.bellsoft.rms.component.auth.SessionSecurityContextRepository
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.filter.JsonAuthenticationFilter
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpMethod
import org.springframework.http.HttpStatus
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.crypto.factory.PasswordEncoderFactories
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.security.web.SecurityFilterChain
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter
import org.springframework.security.web.authentication.logout.HttpStatusReturningLogoutSuccessHandler
import org.springframework.security.web.csrf.CookieCsrfTokenRepository
import org.springframework.security.web.csrf.CsrfTokenRequestAttributeHandler
import org.springframework.security.web.util.matcher.AntPathRequestMatcher

@Configuration
@EnableMethodSecurity(securedEnabled = true)
@EnableWebSecurity(debug = false)
class SecurityConfig(
    @Value("\${spring.profiles.active}") private val activeProfile: String,
    private val authenticationConfiguration: AuthenticationConfiguration,
    private val jsonAuthenticationSuccessHandler: JsonAuthenticationSuccessHandler,
    private val jsonAuthenticationFailureHandler: JsonAuthenticationFailureHandler,
    private val userRepository: UserRepository,
    private val objectMapper: ObjectMapper,
) {
    @Bean
    fun filterChain(http: HttpSecurity): SecurityFilterChain {
        return http.run {
            addFilterBefore(jsonAuthenticationFilter(), UsernamePasswordAuthenticationFilter::class.java)
            securityContext {
                it.securityContextRepository(securityContextRepository())
            }
            csrf {
                it.csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse())
                it.csrfTokenRequestHandler(CsrfTokenRequestAttributeHandler())

                if (isLocalMode())
                    it.ignoringRequestMatchers("/h2-console/**")
            }
            formLogin {
                it.disable()
            }
            logout {
                it.logoutUrl("/api/v1/auth/logout")
                it.logoutSuccessHandler(HttpStatusReturningLogoutSuccessHandler(HttpStatus.OK))
                it.permitAll()
            }
            authorizeHttpRequests {
                if (isLocalMode())
                    it.requestMatchers("/h2-console/**").permitAll()

                it.requestMatchers("/api/*/auth/**").permitAll()
                it.requestMatchers("/api/v1/env").permitAll()
                it.requestMatchers("/api/v1/config").permitAll()
                it.requestMatchers("/api/**").authenticated()
                it.requestMatchers("/api/*/admin/**").hasAnyRole("ADMIN", "SUPER_ADMIN")
                it.requestMatchers("/docs/**").permitAll()
                it.anyRequest().permitAll()
            }
        }.build()
    }

    @Bean
    fun jsonAuthenticationFilter(): JsonAuthenticationFilter {
        return JsonAuthenticationFilter(
            AntPathRequestMatcher("/api/v1/auth/login", HttpMethod.POST.name()),
            objectMapper,
        ).apply {
            setSecurityContextRepository(securityContextRepository())
            setAuthenticationManager(authenticationConfiguration.authenticationManager)
            setAuthenticationSuccessHandler(jsonAuthenticationSuccessHandler)
            setAuthenticationFailureHandler(jsonAuthenticationFailureHandler)
        }
    }

    @Bean
    fun securityContextRepository() = SessionSecurityContextRepository(userRepository)

    @Bean
    fun passwordEncoder(): PasswordEncoder = PasswordEncoderFactories.createDelegatingPasswordEncoder()

    private fun isLocalMode() = activeProfile == "local"
}
