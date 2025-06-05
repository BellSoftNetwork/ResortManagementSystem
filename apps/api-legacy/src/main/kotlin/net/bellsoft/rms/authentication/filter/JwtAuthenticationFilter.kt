package net.bellsoft.rms.authentication.filter

import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import mu.KLogging
import net.bellsoft.rms.authentication.component.JwtTokenProvider
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.util.StringUtils
import org.springframework.web.filter.OncePerRequestFilter

/**
 * JWT 인증 필터
 *
 * 요청 헤더에서 JWT 토큰을 추출하고 유효성을 검사한 후 인증 정보를 SecurityContext에 저장한다.
 */
class JwtAuthenticationFilter(
    private val jwtTokenProvider: JwtTokenProvider,
) : OncePerRequestFilter() {
    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        filterChain: FilterChain,
    ) {
        try {
            // 토큰 갱신 엔드포인트는 토큰 검증을 건너뛴다
            if (request.requestURI == "/api/v1/auth/refresh" || request.requestURI == "/api/v1/auth/login") {
                filterChain.doFilter(request, response)
                return
            }

            val accessToken = resolveToken(request)

            // 토큰이 유효한 경우 인증 정보 설정
            if (
                accessToken != null &&
                StringUtils.hasText(accessToken) &&
                jwtTokenProvider.validateToken(accessToken)
            ) {
                val authentication = jwtTokenProvider.getAuthentication(accessToken)
                SecurityContextHolder.getContext().authentication = authentication
                Companion.logger.debug { "인증 정보 저장 완료 (URI: ${request.requestURI})" }
            }

            filterChain.doFilter(request, response)
        } catch (ex: Exception) {
            Companion.logger.error { "JWT 인증 필터 오류: ${ex.message}" }
            // 예외 처리는 ExceptionController에게 위임
            throw ex
        }
    }

    /**
     * 요청 헤더에서 JWT 토큰을 추출한다.
     *
     * @param request HTTP 요청
     * @return JWT 토큰 (없으면 null)
     */
    private fun resolveToken(request: HttpServletRequest): String? {
        val bearerToken = request.getHeader(AUTHORIZATION_HEADER)

        return if (StringUtils.hasText(bearerToken) && bearerToken.startsWith(BEARER_PREFIX)) {
            bearerToken.substring(BEARER_PREFIX.length)
        } else {
            null
        }
    }

    companion object : KLogging() {
        private const val AUTHORIZATION_HEADER = "Authorization"
        private const val BEARER_PREFIX = "Bearer "
    }
}
