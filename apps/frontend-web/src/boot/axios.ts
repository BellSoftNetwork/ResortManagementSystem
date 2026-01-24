import { boot } from "quasar/wrappers";
import axios, { AxiosError, AxiosResponse } from "axios";
import { TIMEOUT_DURATION } from "src/api/constants";
import { networkStatusService } from "src/api/services/NetworkStatusService";
import { notificationService } from "src/api/services/NotificationService";
import { authInterceptorService } from "src/api/services/AuthInterceptorService";
import { retryService } from "src/api/services/RetryService";

/**
 * API 인스턴스 생성
 */
const api = axios.create({
  timeout: TIMEOUT_DURATION,
});

/**
 * 요청 인터셉터 설정
 */
function setupInterceptors(): void {
  // 요청 인터셉터 설정
  authInterceptorService.setupRequestInterceptor(api);

  // 응답 인터셉터 설정
  api.interceptors.response.use(
    (response: AxiosResponse) => {
      // 온라인 상태 복구 처리
      if (networkStatusService.isOffline) {
        networkStatusService.setOnline();
        notificationService.showOnlineNotification();
      }

      return response;
    },
    async (error: AxiosError) => {
      const originalRequest = error.config || {};

      // 401 에러인 경우 토큰 갱신 시도
      if (error.response?.status === 401 && !originalRequest._retry) {
        // refresh token API 호출인 경우 재시도하지 않음
        if (originalRequest.url?.includes("/auth/refresh")) {
          return Promise.reject(error);
        }
        return authInterceptorService.handleTokenRefresh(originalRequest, error, api);
      }

      // 네트워크 오류나 타임아웃인 경우 재시도
      return retryService.handleRetry(originalRequest, error, api);
    },
  );
}

// 인터셉터 설정 실행
setupInterceptors();

export default boot(({ app }) => {
  // Vue 컴포넌트에서 사용할 수 있도록 전역 속성으로 등록
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;
});

export { api };
