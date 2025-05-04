package net.bellsoft.rms.authentication.util

/**
 * 인증 관련 테스트에서 공통으로 사용하는 상수 모음
 */
object AuthTestConstants {
    /**
     * 기본 Windows 사용자 에이전트
     */
    const val DEFAULT_WINDOWS_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
        "(KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

    /**
     * 기본 Android 사용자 에이전트
     */
    const val DEFAULT_ANDROID_USER_AGENT = "Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 " +
        "(KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36"

    /**
     * 기본 iOS 사용자 에이전트
     */
    const val DEFAULT_IOS_USER_AGENT = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 " +
        "(KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"

    /**
     * 기본 한국어 언어 설정
     */
    const val DEFAULT_KOREAN_LANGUAGE = "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7"

    /**
     * 기본 영어 언어 설정
     */
    const val DEFAULT_ENGLISH_LANGUAGE = "en-US,en;q=0.9"
}
