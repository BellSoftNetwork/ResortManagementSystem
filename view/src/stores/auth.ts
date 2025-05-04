// Utilities
import { defineStore } from "pinia";
import { User } from "src/schema/user";
import { LoginParams, postLogin, postLogout, postRefreshToken } from "src/api/v1/auth";
import { getMy } from "src/api/v1/main";

interface State {
  status: {
    isFirstRequest: boolean;
    isRefreshingToken: boolean;
  };
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  accessTokenExpiresIn: number | null;
  tokenRefreshTimer: number | null;
}

// 로컬 스토리지 키
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const TOKEN_EXPIRES_KEY = "token_expires";

export const useAuthStore = defineStore("auth", {
  state: (): State => ({
    status: {
      isFirstRequest: true,
      isRefreshingToken: false,
    },
    user: null,
    accessToken: localStorage.getItem(ACCESS_TOKEN_KEY),
    refreshToken: localStorage.getItem(REFRESH_TOKEN_KEY),
    accessTokenExpiresIn: localStorage.getItem(TOKEN_EXPIRES_KEY)
      ? Number(localStorage.getItem(TOKEN_EXPIRES_KEY))
      : null,
    tokenRefreshTimer: null,
  }),

  getters: {
    isFirstRequest: (state) => state.status.isFirstRequest,
    isRefreshingToken: (state) => state.status.isRefreshingToken,
    isLoggedIn: (state) => state.user !== null && state.accessToken !== null,
    isNormalRole: (state) => ["NORMAL", "ADMIN", "SUPER_ADMIN"].includes(state.user?.role),
    isAdminRole: (state) => ["ADMIN", "SUPER_ADMIN"].includes(state.user?.role),
    isSuperAdminRole: (state) => ["SUPER_ADMIN"].includes(state.user?.role),
  },

  actions: {
    loadAccountInfo() {
      // 토큰이 없으면 로드하지 않음
      if (!this.accessToken) {
        this.status.isFirstRequest = false;
        return Promise.resolve();
      }

      return getMy()
        .then((response) => {
          this.user = response.value;
        })
        .catch(() => {
          this.user = null;
          this.clearTokens();
        })
        .finally(() => {
          this.status.isFirstRequest = false;
        });
    },

    login(params: LoginParams) {
      this.user = null;
      this.clearTokens();

      return postLogin(params).then((response) => {
        // 사용자 정보 저장
        this.user = response.value.user;

        // 토큰 정보 저장
        this.setTokens(response.value.accessToken, response.value.refreshToken, response.value.accessTokenExpiresIn);

        // 토큰 자동 갱신 타이머 시작
        this.startTokenRefreshTimer();
      });
    },

    logout() {
      return postLogout().then(() => {
        this.user = null;
        this.clearTokens();
        // 토큰 자동 갱신 타이머 중지
        this.stopTokenRefreshTimer();
      });
    },

    // 토큰이 만료되었는지 확인
    isTokenExpired() {
      if (!this.accessTokenExpiresIn) return true;
      return Date.now() >= this.accessTokenExpiresIn;
    },

    // 액세스 토큰 갱신
    async refreshAccessToken(forceLoading = false) {
      if (!this.refreshToken) {
        throw new Error("리프레시 토큰이 없습니다.");
      }

      // 토큰이 만료되었거나 강제 로딩이 필요한 경우에만 로딩 상태 표시
      const isExpired = this.isTokenExpired();
      if (isExpired || forceLoading) {
        this.status.isRefreshingToken = true;
      }

      try {
        const response = await postRefreshToken({
          refreshToken: this.refreshToken,
        });

        this.setTokens(response.value.accessToken, response.value.refreshToken, response.value.accessTokenExpiresIn);

        if (!this.isLoggedIn) {
          await this.loadAccountInfo();
        }

        // 토큰 갱신 후 타이머 재시작
        this.startTokenRefreshTimer();

        return response;
      } catch (error) {
        // 네트워크 오류가 아닌 인증 오류인 경우에만 토큰 제거
        // 네트워크 오류는 서버 응답이 없는 경우로, 토큰이 실제로 만료되지 않았을 수 있음
        const isNetworkError = !error.response;
        const isAuthError = error.response?.status === 401 || error.response?.status === 403;

        if (!isNetworkError && isAuthError) {
          // 인증 오류인 경우에만 토큰 제거
          this.clearTokens();
        }

        throw error;
      } finally {
        // 토큰이 만료되었거나 강제 로딩이 필요한 경우에만 로딩 상태 해제
        if (isExpired || forceLoading) {
          this.status.isRefreshingToken = false;
        }
      }
    },

    // 토큰 저장 (액세스 토큰은 메모리에만, 리프레시 토큰은 로컬 스토리지에도 저장)
    setTokens(accessToken: string, refreshToken: string, expiresIn: number) {
      this.accessToken = accessToken;
      this.refreshToken = refreshToken;
      this.accessTokenExpiresIn = expiresIn;

      // 액세스 토큰은 메모리에만 저장 (로컬 스토리지에 저장하지 않음)
      localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
      localStorage.setItem(TOKEN_EXPIRES_KEY, String(expiresIn));
    },

    // 토큰 제거
    clearTokens() {
      this.accessToken = null;
      this.refreshToken = null;
      this.accessTokenExpiresIn = null;

      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(TOKEN_EXPIRES_KEY);

      // 토큰 자동 갱신 타이머 중지
      this.stopTokenRefreshTimer();
    },

    // 토큰이 만료되기 5분 전인지 확인
    isTokenNearExpiration() {
      if (!this.accessTokenExpiresIn) return false;

      const currentTime = Date.now();
      const fiveMinutesInMs = 5 * 60 * 1000; // 5분을 밀리초로 변환

      // 현재 시간 + 5분이 토큰 만료 시간보다 크거나 같으면 곧 만료됨
      return currentTime + fiveMinutesInMs >= this.accessTokenExpiresIn;
    },

    // 토큰 자동 갱신 타이머 시작
    startTokenRefreshTimer() {
      // 이미 타이머가 실행 중이면 중지
      this.stopTokenRefreshTimer();

      // 1분마다 토큰 상태 확인
      this.tokenRefreshTimer = window.setInterval(() => {
        // 로그인 상태가 아니면 타이머 중지
        if (!this.isLoggedIn) {
          this.stopTokenRefreshTimer();
          return;
        }

        // 토큰이 곧 만료되면 자동 갱신
        if (this.isTokenNearExpiration()) {
          this.refreshAccessToken().catch(() => {
            // 갱신 실패 시 타이머 중지
            this.stopTokenRefreshTimer();
          });
        }
      }, 60000); // 1분마다 체크
    },

    // 토큰 자동 갱신 타이머 중지
    stopTokenRefreshTimer() {
      if (this.tokenRefreshTimer !== null) {
        window.clearInterval(this.tokenRefreshTimer);
        this.tokenRefreshTimer = null;
      }
    },

    // 스토어 초기화 시 자동 실행
    async hydrate(storeState) {
      // 액세스 토큰이 없지만 리프레시 토큰이 있는 경우 (페이지 새로고침 후)
      if (!storeState.accessToken && storeState.refreshToken) {
        try {
          // 액세스 토큰 갱신 시도 (페이지 새로고침 후에는 forceLoading=true로 설정)
          await this.refreshAccessToken(true);
          // 갱신 성공 시 여기서 리턴 (아래 코드는 실행되지 않음)
          return;
        } catch (error) {
          // 리프레시 토큰이 만료된 경우에만 로그인 페이지로 리다이렉트
          if (this.router) {
            this.router.push({ name: "Login" });
          }
          return;
        }
      }

      // 토큰이 있지만 사용자 정보가 없는 경우 사용자 정보 로드
      if (storeState.accessToken && !storeState.user) {
        this.loadAccountInfo().catch(() => {
          this.clearTokens();
        });
      }

      // 로그인 상태면 토큰 자동 갱신 타이머 시작
      if (storeState.accessToken) {
        this.startTokenRefreshTimer();
      }
    },
  },
});
