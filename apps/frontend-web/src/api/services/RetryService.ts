import { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { CRITICAL_APIS, EXCLUDE_RETRY_URLS, MAX_RETRY_COUNT, RETRY_DELAY_BASE } from "../constants";
import { networkStatusService } from "./NetworkStatusService";
import { notificationService } from "./NotificationService";

/**
 * 요청 재시도 관련 로직을 관리하는 서비스
 */
class RetryService {
  /**
   * 요청 재시도 처리
   * @param originalRequest 원래 요청
   * @param error 발생한 에러
   * @param api API 인스턴스
   */
  handleRetry(
    originalRequest: InternalAxiosRequestConfig & { retryCount?: number },
    error: AxiosError,
    api: AxiosInstance,
  ): Promise<unknown> {
    // 재시도 횟수 초기화
    if (originalRequest.retryCount === undefined) {
      originalRequest.retryCount = 0;
    }

    // 재시도 가능 여부 확인
    const isTimeout = error.code === "ECONNABORTED" || error.message?.includes("timeout");
    const isNetworkError = !error.response;
    const shouldRetry =
      (isTimeout || isNetworkError) &&
      originalRequest.retryCount < MAX_RETRY_COUNT &&
      originalRequest.method?.toUpperCase() === "GET" &&
      !EXCLUDE_RETRY_URLS.some((url) => originalRequest.url?.includes(url));

    if (!shouldRetry) {
      // 모든 재시도가 실패한 경우 사용자에게 알림
      if ((isTimeout || isNetworkError) && originalRequest.retryCount >= MAX_RETRY_COUNT) {
        const requestUrl = originalRequest.url || "";
        const isCriticalApi = CRITICAL_APIS.some((url) => requestUrl.includes(url));

        if (isCriticalApi) {
          // Critical API 재시도 모두 실패 → 서버 전체 장애
          notificationService.showOfflineNotification(-1); // -1은 모든 재시도 실패를 나타냄
          if (!networkStatusService.isOffline) {
            networkStatusService.setNetworkError();
          }
        } else {
          // 일반 API 재시도 모두 실패 → Toast만
          notificationService.showApiErrorNotification("서버에 연결할 수 없습니다.");
        }
      }

      return Promise.reject(error);
    }

    // 재시도 횟수 증가
    originalRequest.retryCount++;

    // 재시도 간격 계산 (1초, 2초, 3초)
    const retryDelay = originalRequest.retryCount * RETRY_DELAY_BASE;

    // 오프라인 상태 설정 및 알림 표시
    if (!networkStatusService.isOffline) {
      networkStatusService.setNetworkError();
    }

    // 재시도 중임을 통합된 알림으로 표시
    notificationService.showOfflineNotification(originalRequest.retryCount);

    // 지정된 시간 후 재시도
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        api(originalRequest)
          .then((response) => {
            // 재시도 성공 시 온라인 상태로 복귀
            if (networkStatusService.isOffline) {
              networkStatusService.setOnline();
              notificationService.showOnlineNotification();
            }
            resolve(response);
          })
          .catch((err) => reject(err));
      }, retryDelay);
    });
  }
}

// 싱글톤 인스턴스 생성
export const retryService = new RetryService();
