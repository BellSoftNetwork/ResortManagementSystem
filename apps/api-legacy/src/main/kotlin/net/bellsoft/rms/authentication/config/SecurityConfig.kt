package net.bellsoft.rms.authentication.config

import net.bellsoft.rms.authentication.component.JwtTokenProvider
import net.bellsoft.rms.authentication.filter.JwtAuthenticationFilter
import net.bellsoft.rms.authentication.handler.JwtAuthenticationEntryPoint
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpStatus
import org.springframework.security.authentication.AuthenticationManager
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.http.SessionCreationPolicy
import org.springframework.security.crypto.factory.PasswordEncoderFactories
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.security.web.SecurityFilterChain
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter
import org.springframework.security.web.authentication.logout.HttpStatusReturningLogoutSuccessHandler

@Configuration
@EnableMethodSecurity(securedEnabled = true)
@EnableWebSecurity(debug = false)
class SecurityConfig(
    @Value("\${spring.profiles.active}") private val activeProfile: String,
    private val jwtTokenProvider: JwtTokenProvider,
    private val jwtAuthenticationEntryPoint: JwtAuthenticationEntryPoint,
) {
    @Bean
    fun filterChain(http: HttpSecurity): SecurityFilterChain {
        return http.run {
            // JWT 필터 추가
            addFilterBefore(
                JwtAuthenticationFilter(jwtTokenProvider),
                UsernamePasswordAuthenticationFilter::class.java,
            )

            // 세션 사용 안함 (JWT 사용)
            sessionManagement {
                it.sessionCreationPolicy(SessionCreationPolicy.STATELESS)
            }

            csrf {
                it.disable()
            }
            formLogin {
                it.disable()
            }
            exceptionHandling {
                it.authenticationEntryPoint(jwtAuthenticationEntryPoint)
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
    fun authenticationManager(authenticationConfiguration: AuthenticationConfiguration): AuthenticationManager {
        return authenticationConfiguration.authenticationManager
    }

    @Bean
    fun passwordEncoder(): PasswordEncoder = PasswordEncoderFactories.createDelegatingPasswordEncoder()

    private fun isLocalMode() = activeProfile == "local"
}
