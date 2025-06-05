import { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "src/stores/auth";

/**
 * 인증 관련 인터셉터 로직을 관리하는 서비스
 */
class AuthInterceptorService {
  /**
   * 요청 인터셉터 설정: JWT 토큰을 요청 헤더에 추가
   * @param instance Axios 인스턴스
   */
  setupRequestInterceptor(instance: AxiosInstance): void {
    instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const authStore = useAuthStore();
        const token = authStore.accessToken;

        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
      },
      (error) => {
        return Promise.reject(error);
      },
    );
  }

  /**
   * 토큰 갱신 처리
   * @param originalRequest 원래 요청
   * @param error 발생한 에러
   * @param api API 인스턴스
   */
  async handleTokenRefresh(
    originalRequest: InternalAxiosRequestConfig & { _retry?: boolean },
    error: AxiosError,
    api: AxiosInstance,
  ): Promise<unknown> {
    const authStore = useAuthStore();

    // 이미 재시도 중인 경우 또는 리프레시 토큰이 없는 경우
    if (originalRequest._retry || !authStore.refreshToken) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    try {
      // 토큰 갱신 시도
      await authStore.refreshAccessToken();

      // 갱신된 토큰으로 원래 요청 재시도
      originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`;
      return api(originalRequest);
    } catch (refreshError) {
      // 리프레시 토큰도 만료된 경우 로그아웃
      authStore.clearTokens();

      // 로그인 페이지로 리다이렉트
      if (authStore.router) {
        authStore.router.push({ name: "Login" });
      }

      return Promise.reject(refreshError);
    }
  }
}

// 싱글톤 인스턴스 생성
export const authInterceptorService = new AuthInterceptorService();
