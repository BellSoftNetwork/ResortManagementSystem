import { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "src/stores/auth";

/**
 * 인증 관련 인터셉터 로직을 관리하는 서비스
 */
class AuthInterceptorService {
  private refreshPromise: Promise<any> | null = null;
  private refreshAttempts = 0;
  private readonly MAX_REFRESH_ATTEMPTS = 3;
  private readonly REFRESH_RETRY_DELAY = 1000; // 1초
  private lastRefreshAttempt = 0;
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

    // refresh 토큰 요청 자체인 경우 재시도하지 않음
    if (originalRequest.url?.includes("/auth/refresh")) {
      this.resetRefreshAttempts();
      return Promise.reject(error);
    }

    // 재시도 횟수 초과 확인
    if (this.refreshAttempts >= this.MAX_REFRESH_ATTEMPTS) {
      console.error("Refresh token retry limit exceeded");
      // 네트워크 오류인 경우 토큰을 유지하기 위해 handleRefreshFailure()를 호출하지 않음
      // 실제 토큰 제거는 catch 블록에서 인증 오류(401/403)인 경우에만 수행
      this.resetRefreshAttempts();
      return Promise.reject(error);
    }

    // 최소 재시도 간격 확인
    const now = Date.now();
    const timeSinceLastAttempt = now - this.lastRefreshAttempt;
    if (timeSinceLastAttempt < this.REFRESH_RETRY_DELAY) {
      const delay = this.REFRESH_RETRY_DELAY - timeSinceLastAttempt;
      await new Promise((resolve) => setTimeout(resolve, delay));
    }

    originalRequest._retry = true;

    // 이미 진행 중인 refresh 요청이 있으면 기다림
    if (this.refreshPromise) {
      try {
        await this.refreshPromise;
        originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`;
        return api(originalRequest);
      } catch (error) {
        return Promise.reject(error);
      }
    }

    // 새로운 refresh 요청 시작
    this.refreshAttempts++;
    this.lastRefreshAttempt = Date.now();

    this.refreshPromise = authStore
      .refreshAccessToken()
      .then(() => {
        // 갱신 성공 시 재시도 카운터 리셋
        this.resetRefreshAttempts();

        // 갱신된 토큰으로 원래 요청 재시도
        originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`;
        return api(originalRequest);
      })
      .catch((refreshError) => {
        // 401/403 인증 오류에서만 로그아웃 처리, 네트워크 오류는 토큰 유지
        const status = refreshError?.response?.status;
        if (status === 401 || status === 403) {
          this.handleRefreshFailure(authStore);
        }
        return Promise.reject(refreshError);
      })
      .finally(() => {
        this.refreshPromise = null;
      });

    return this.refreshPromise;
  }

  /**
   * 리프레시 실패 처리
   */
  private handleRefreshFailure(authStore: any): void {
    // 재시도 카운터 리셋
    this.resetRefreshAttempts();

    // 토큰 제거
    authStore.clearTokens();

    // 로그인 페이지로 리다이렉트
    if (authStore.router) {
      authStore.router.push({ name: "Login" });
    }
  }

  /**
   * 재시도 카운터 리셋
   */
  private resetRefreshAttempts(): void {
    this.refreshAttempts = 0;
    this.lastRefreshAttempt = 0;
    this.refreshPromise = null;
  }
}

// 싱글톤 인스턴스 생성
export const authInterceptorService = new AuthInterceptorService();
