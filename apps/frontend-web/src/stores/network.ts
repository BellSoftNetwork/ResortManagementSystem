import { defineStore } from "pinia";

/**
 * 에러 타입 정의
 * - NONE: 정상 상태
 * - NETWORK_ERROR: 네트워크 연결 불가 (서버 응답 없음)
 * - SERVER_ERROR: 서버 오류 (5xx 응답)
 */
export type ErrorType = "NONE" | "NETWORK_ERROR" | "SERVER_ERROR";

interface NetworkState {
  isOffline: boolean;
  errorType: ErrorType;
  lastError: string | null;
  statusCode: number | null;
  retryCount: number;
  autoRetryEnabled: boolean;
  autoRetryCount: number;
}

/**
 * 에러 타입별 기본 메시지
 */
const ERROR_MESSAGES: Record<ErrorType, string> = {
  NONE: "",
  NETWORK_ERROR: "서버에 연결할 수 없습니다",
  SERVER_ERROR: "서버에 문제가 발생했습니다",
};

/**
 * 자동 재시도 설정 상수
 */
const MAX_AUTO_RETRY = 10;
export const RETRY_INTERVAL_MS = 10000; // 10초

export const useNetworkStore = defineStore("network", {
  state: (): NetworkState => ({
    isOffline: false,
    errorType: "NONE",
    lastError: null,
    statusCode: null,
    retryCount: 0,
    autoRetryEnabled: false,
    autoRetryCount: 0,
  }),

  getters: {
    isOnline: (state) => !state.isOffline,
    isNetworkError: (state) => state.errorType === "NETWORK_ERROR",
    isServerError: (state) => state.errorType === "SERVER_ERROR",
    /**
     * 에러 메시지 반환 (커스텀 메시지 또는 기본 메시지)
     */
    errorMessage: (state) => state.lastError || ERROR_MESSAGES[state.errorType],
    /**
     * 자동 재시도 가능 여부 (활성화 상태 && 최대 횟수 미만)
     */
    canAutoRetry: (state) => state.autoRetryEnabled && state.autoRetryCount < MAX_AUTO_RETRY,
    /**
     * 재시도 진행 상태 문자열 ("3/10")
     */
    retryProgress: (state) => `${state.autoRetryCount}/${MAX_AUTO_RETRY}`,
  },

  actions: {
    /**
     * 네트워크 오류 설정 (서버 응답 없음)
     */
    setNetworkError(error?: string) {
      this.isOffline = true;
      this.errorType = "NETWORK_ERROR";
      this.lastError = error || null;
      this.statusCode = null;
    },

    /**
     * 서버 오류 설정 (5xx 응답)
     */
    setServerError(statusCode: number, error?: string) {
      this.isOffline = true;
      this.errorType = "SERVER_ERROR";
      this.lastError = error || null;
      this.statusCode = statusCode;
    },

    /**
     * 오프라인 상태 설정 (하위 호환성 유지)
     * @deprecated setNetworkError() 또는 setServerError() 사용 권장
     */
    setOffline(error?: string) {
      this.setNetworkError(error);
    },

    /**
     * 온라인 상태로 복구
     */
    setOnline() {
      this.isOffline = false;
      this.errorType = "NONE";
      this.lastError = null;
      this.statusCode = null;
      this.retryCount = 0;
      this.resetRetryState();
    },

    /**
     * 재시도 횟수 증가
     */
    incrementRetryCount() {
      this.retryCount++;
    },

    /**
     * 자동 재시도 시작
     */
    startAutoRetry() {
      this.autoRetryEnabled = true;
      this.autoRetryCount = 0;
    },

    /**
     * 자동 재시도 횟수 증가
     */
    incrementAutoRetry() {
      this.autoRetryCount++;
    },

    /**
     * 자동 재시도 중단
     */
    stopAutoRetry() {
      this.autoRetryEnabled = false;
    },

    /**
     * 재시도 상태 초기화 (복구 시)
     */
    resetRetryState() {
      this.autoRetryEnabled = false;
      this.autoRetryCount = 0;
    },
  },
});
