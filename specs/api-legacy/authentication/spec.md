---
id: api-legacy-authentication
title: "api-legacy 인증"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: high
effort: medium
---

# api-legacy 인증

> Spring Security + JWT 기반 인증

---

## 1. 개요

### 1.1 인증 방식

- Spring Security 기반 인증
- JWT (JSON Web Token)
- Access Token (15분) + Refresh Token (7일)

### 1.2 관련 파일

| 파일 | 역할 |
|------|------|
| `AuthController.kt` | 인증 엔드포인트 |
| `AuthService.kt` | 인증 서비스 인터페이스 |
| `AuthServiceImpl.kt` | 인증 서비스 구현 |
| `JwtAuthenticationFilter.kt` | JWT 필터 |
| `JsonAuthenticationFilter.kt` | JSON 로그인 필터 |
| `LoginAttemptService.kt` | 로그인 시도 추적 |
| `SecurityConfig.kt` | 보안 설정 |

---

## 2. 엔드포인트

| Method | Path | 설명 |
|--------|------|------|
| POST | `/api/v1/auth/register` | 회원가입 |
| POST | `/api/v1/auth/login` | 로그인 |
| POST | `/api/v1/auth/refresh` | 토큰 갱신 |

---

## 3. Spring Security 설정

```kotlin
@Configuration
@EnableWebSecurity
@EnableMethodSecurity
class SecurityConfig(
    private val jwtAuthenticationFilter: JwtAuthenticationFilter
) {
    @Bean
    fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
        http
            .csrf { it.disable() }
            .sessionManagement { it.sessionCreationPolicy(STATELESS) }
            .authorizeHttpRequests {
                it.requestMatchers("/api/v1/auth/**").permitAll()
                it.requestMatchers("/actuator/**").permitAll()
                it.anyRequest().authenticated()
            }
            .addFilterBefore(jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter::class.java)
        return http.build()
    }
}
```

---

## 4. JWT 필터

```kotlin
class JwtAuthenticationFilter(
    private val jwtTokenProvider: JwtTokenProvider,
    private val userDetailsService: UserDetailsService
) : OncePerRequestFilter() {
    
    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        chain: FilterChain
    ) {
        val token = resolveToken(request)
        if (token != null && jwtTokenProvider.validateToken(token)) {
            val username = jwtTokenProvider.getUsername(token)
            val userDetails = userDetailsService.loadUserByUsername(username)
            val auth = UsernamePasswordAuthenticationToken(
                userDetails, null, userDetails.authorities
            )
            SecurityContextHolder.getContext().authentication = auth
        }
        chain.doFilter(request, response)
    }
}
```

---

## 5. 로그인 시도 추적

### 5.1 LoginAttemptService

```kotlin
@Service
class LoginAttemptService(
    private val loginAttemptRepository: LoginAttemptRepository
) {
    fun recordAttempt(username: String, ipAddress: String, successful: Boolean, deviceInfo: DeviceInfoDto?)
    fun isBlocked(username: String, ipAddress: String): Boolean
    fun getRecentFailedAttempts(username: String, ipAddress: String): Int
}
```

### 5.2 브루트포스 방지

| 항목 | 값 |
|------|-----|
| 시간 윈도우 | 15분 |
| 최대 실패 횟수 | 5회 |
| 차단 응답 | 429 Too Many Requests |

---

## 6. DTO

### 6.1 요청

```kotlin
data class UserRegistrationRequest(
    val userId: String,
    val email: String?,
    val name: String,
    val password: String
)

data class LoginRequest(
    val username: String,
    val password: String,
    val device: DeviceInfoDto?
)

data class RefreshTokenRequest(
    val refreshToken: String
)
```

### 6.2 응답

```kotlin
data class TokenDto(
    val accessToken: String,
    val refreshToken: String,
    val accessTokenExpiresIn: Long
)
```

---

## 7. 비밀번호 저장

Spring Security 호환 BCrypt 형식:

```
{bcrypt}$2a$10$...
```

---

## 8. 예외 처리

| 예외 | 상황 |
|------|------|
| AuthenticationException | 인증 실패 |
| TokenExpiredException | 토큰 만료 |
| TooManyRequestsException | 브루트포스 차단 |
