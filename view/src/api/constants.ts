// 네트워크 및 API 관련 상수

// 타임아웃 설정
export const TIMEOUT_DURATION = 10000; // 10초

// 재시도 관련 설정
export const MAX_RETRY_COUNT = 3;
export const RETRY_DELAY_BASE = 1000; // 1초

// 알림 설정
export const NOTIFICATION_TIMEOUT = 5000; // 5초

// 재시도 제외 URL 목록
export const EXCLUDE_RETRY_URLS = ["/api/v1/auth/logout"];
