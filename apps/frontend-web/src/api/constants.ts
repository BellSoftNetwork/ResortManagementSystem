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

// 서버 전체 장애로 판단하는 Critical API 목록
// 이 API들이 5xx 에러를 반환하면 ServerError UI로 전환
// 일반 API 5xx는 Toast 알림만 표시하고 페이지 유지
export const CRITICAL_APIS = ["/api/v1/auth/refresh", "/api/v1/auth/login", "/api/v1/config"];
