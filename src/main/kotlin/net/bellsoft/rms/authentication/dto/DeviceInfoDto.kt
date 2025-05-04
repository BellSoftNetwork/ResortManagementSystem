package net.bellsoft.rms.authentication.dto

import jakarta.servlet.http.HttpServletRequest
import java.security.MessageDigest

/**
 * 디바이스 정보 DTO
 *
 * 클라이언트의 디바이스 정보를 담는 데이터 클래스
 */
data class DeviceInfoDto(
    /**
     * 클라이언트 IP 주소
     */
    val ipAddress: String,

    /**
     * 운영체제 정보
     */
    val osInfo: String,

    /**
     * 언어 설정 정보
     */
    val languageInfo: String,

    /**
     * 사용자 에이전트 정보
     */
    val userAgent: String,

    /**
     * 디바이스 핑거프린트 해시
     */
    val deviceFingerprint: String,
) {
    companion object {
        /**
         * HttpServletRequest에서 디바이스 정보를 추출하여 DeviceInfo 객체를 생성한다.
         *
         * @param request HTTP 요청
         * @return DeviceInfo 객체
         */
        fun fromRequest(request: HttpServletRequest): DeviceInfoDto {
            val ipAddress = getClientIp(request)
            val userAgent = request.getHeader("User-Agent") ?: ""
            val acceptLanguage = request.getHeader("Accept-Language") ?: ""
            val osInfo = extractPlatformFromUserAgent(userAgent)
            val languageInfo = acceptLanguage.split(',').firstOrNull()?.trim() ?: ""

            // OS 정보를 해싱하여 디바이스 핑거프린트 생성
            val deviceFingerprint = generateFingerprintHash(osInfo)

            return DeviceInfoDto(
                ipAddress = ipAddress,
                osInfo = osInfo,
                languageInfo = languageInfo,
                userAgent = userAgent,
                deviceFingerprint = deviceFingerprint,
            )
        }

        /**
         * 클라이언트 IP 주소를 가져온다.
         *
         * @param request HTTP 요청
         * @return 클라이언트 IP 주소
         */
        private fun getClientIp(request: HttpServletRequest): String {
            var ip = request.getHeader("X-Forwarded-For")
            if (ip.isNullOrEmpty() || "unknown".equals(ip, ignoreCase = true)) {
                ip = request.getHeader("Proxy-Client-IP")
            }
            if (ip.isNullOrEmpty() || "unknown".equals(ip, ignoreCase = true)) {
                ip = request.getHeader("WL-Proxy-Client-IP")
            }
            if (ip.isNullOrEmpty() || "unknown".equals(ip, ignoreCase = true)) {
                ip = request.getHeader("HTTP_CLIENT_IP")
            }
            if (ip.isNullOrEmpty() || "unknown".equals(ip, ignoreCase = true)) {
                ip = request.getHeader("HTTP_X_FORWARDED_FOR")
            }
            if (ip.isNullOrEmpty() || "unknown".equals(ip, ignoreCase = true)) {
                ip = request.remoteAddr
            }
            return ip
        }

        /**
         * User-Agent 헤더에서 운영체제 정보를 추출한다.
         *
         * @param userAgent User-Agent 헤더 값
         * @return 운영체제 정보
         */
        private fun extractPlatformFromUserAgent(userAgent: String): String {
            return when {
                userAgent.contains("Windows") -> "Windows"
                userAgent.contains("Mac OS X") -> "macOS"
                userAgent.contains("iPhone") || userAgent.contains("iPad") -> "iOS"
                userAgent.contains("Android") -> "Android"
                userAgent.contains("Linux") -> "Linux"
                else -> "Unknown"
            }
        }

        /**
         * 문자열을 SHA-256으로 해싱한다.
         *
         * @param input 해싱할 문자열
         * @return 해시 문자열
         */
        private fun generateFingerprintHash(input: String): String {
            if (input.isBlank()) return ""

            val bytes = MessageDigest.getInstance("SHA-256").digest(input.toByteArray())
            return bytes.joinToString("") { "%02x".format(it) }
        }
    }
}
