import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useAuthStore } from "../auth";
import * as authApi from "src/api/v1/auth";
import * as mainApi from "src/api/v1/main";
import { createMockApiResponse, createMockApiError } from "test/vitest/helpers";
import type { User } from "src/schema/user";
import type { LoginResponse, RefreshTokenResponse } from "src/api/v1/auth";

vi.mock("src/api/v1/auth");
vi.mock("src/api/v1/main");

const createMockUser = (role: "NORMAL" | "ADMIN" | "SUPER_ADMIN" = "NORMAL"): User => ({
  id: 1,
  name: "테스트 사용자",
  userId: "testuser",
  email: "test@example.com",
  role,
  status: "ACTIVE",
  profileImageUrl: "",
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
  createdBy: "admin",
  updatedBy: "admin",
});

const createMockLoginResponse = (role: "NORMAL" | "ADMIN" | "SUPER_ADMIN" = "NORMAL"): LoginResponse => ({
  user: createMockUser(role),
  accessToken: "mock-access-token",
  refreshToken: "mock-refresh-token",
  accessTokenExpiresIn: Date.now() + 3600000,
});

const createMockRefreshResponse = (): RefreshTokenResponse => ({
  accessToken: "new-access-token",
  refreshToken: "new-refresh-token",
  accessTokenExpiresIn: Date.now() + 3600000,
});

describe("useAuthStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    localStorage.clear();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe("초기 상태", () => {
    it("기본 상태값이 올바르게 설정된다", () => {
      const store = useAuthStore();

      expect(store.status.isFirstRequest).toBe(true);
      expect(store.status.isRefreshingToken).toBe(false);
      expect(store.user).toBeNull();
      expect(store.accessToken).toBeNull();
      expect(store.refreshToken).toBeNull();
      expect(store.refreshAttempts).toBe(0);
    });
  });

  describe("getters", () => {
    it("isLoggedIn은 user와 accessToken이 모두 있을 때 true를 반환한다", () => {
      const store = useAuthStore();

      expect(store.isLoggedIn).toBe(false);

      store.user = createMockUser();
      store.accessToken = "token";

      expect(store.isLoggedIn).toBe(true);
    });

    it("isNormalRole은 NORMAL, ADMIN, SUPER_ADMIN 역할에서 true를 반환한다", () => {
      const store = useAuthStore();

      store.user = createMockUser("NORMAL");
      expect(store.isNormalRole).toBe(true);

      store.user = createMockUser("ADMIN");
      expect(store.isNormalRole).toBe(true);

      store.user = createMockUser("SUPER_ADMIN");
      expect(store.isNormalRole).toBe(true);
    });

    it("isAdminRole은 ADMIN, SUPER_ADMIN 역할에서만 true를 반환한다", () => {
      const store = useAuthStore();

      store.user = createMockUser("NORMAL");
      expect(store.isAdminRole).toBe(false);

      store.user = createMockUser("ADMIN");
      expect(store.isAdminRole).toBe(true);

      store.user = createMockUser("SUPER_ADMIN");
      expect(store.isAdminRole).toBe(true);
    });

    it("isSuperAdminRole은 SUPER_ADMIN 역할에서만 true를 반환한다", () => {
      const store = useAuthStore();

      store.user = createMockUser("NORMAL");
      expect(store.isSuperAdminRole).toBe(false);

      store.user = createMockUser("ADMIN");
      expect(store.isSuperAdminRole).toBe(false);

      store.user = createMockUser("SUPER_ADMIN");
      expect(store.isSuperAdminRole).toBe(true);
    });
  });

  describe("actions", () => {
    describe("loadAccountInfo", () => {
      it("토큰이 없으면 사용자 정보를 로드하지 않는다", async () => {
        const store = useAuthStore();

        await store.loadAccountInfo();

        expect(mainApi.getMy).not.toHaveBeenCalled();
        expect(store.status.isFirstRequest).toBe(false);
      });

      it("토큰이 있으면 사용자 정보를 로드한다", async () => {
        const mockUser = createMockUser();
        vi.mocked(mainApi.getMy).mockResolvedValue(createMockApiResponse(mockUser));

        const store = useAuthStore();
        store.accessToken = "valid-token";

        await store.loadAccountInfo();

        expect(mainApi.getMy).toHaveBeenCalled();
        expect(store.user).toEqual(mockUser);
        expect(store.status.isFirstRequest).toBe(false);
      });

      it("401 에러 시 토큰을 제거한다", async () => {
        const error = { response: { status: 401 } };
        vi.mocked(mainApi.getMy).mockRejectedValue(error);

        const store = useAuthStore();
        store.accessToken = "invalid-token";
        store.refreshToken = "refresh-token";

        await store.loadAccountInfo();

        expect(store.user).toBeNull();
        expect(store.accessToken).toBeNull();
        expect(store.refreshToken).toBeNull();
      });
    });

    describe("login", () => {
      it("로그인 성공 시 사용자 정보와 토큰을 저장한다", async () => {
        const mockResponse = createMockLoginResponse();
        vi.mocked(authApi.postLogin).mockResolvedValue(createMockApiResponse(mockResponse));

        const store = useAuthStore();
        await store.login({ username: "testuser", password: "password" });

        expect(store.user).toEqual(mockResponse.user);
        expect(store.accessToken).toBe(mockResponse.accessToken);
        expect(store.refreshToken).toBe(mockResponse.refreshToken);
        expect(store.accessTokenExpiresIn).toBe(mockResponse.accessTokenExpiresIn);
      });

      it("로그인 전 기존 토큰을 제거한다", async () => {
        const mockResponse = createMockLoginResponse();
        vi.mocked(authApi.postLogin).mockResolvedValue(createMockApiResponse(mockResponse));

        const store = useAuthStore();
        store.accessToken = "old-token";
        store.user = createMockUser();

        await store.login({ username: "testuser", password: "password" });

        expect(authApi.postLogin).toHaveBeenCalledWith({ username: "testuser", password: "password" });
      });
    });

    describe("logout", () => {
      it("로그아웃 시 사용자 정보와 토큰을 제거한다", async () => {
        vi.mocked(authApi.postLogout).mockResolvedValue({});

        const store = useAuthStore();
        store.user = createMockUser();
        store.accessToken = "token";
        store.refreshToken = "refresh";

        await store.logout();

        expect(store.user).toBeNull();
        expect(store.accessToken).toBeNull();
        expect(store.refreshToken).toBeNull();
      });
    });

    describe("isTokenExpired", () => {
      it("accessTokenExpiresIn이 없으면 true를 반환한다", () => {
        const store = useAuthStore();

        expect(store.isTokenExpired()).toBe(true);
      });

      it("현재 시간이 만료 시간을 지났으면 true를 반환한다", () => {
        const store = useAuthStore();
        store.accessTokenExpiresIn = Date.now() - 1000;

        expect(store.isTokenExpired()).toBe(true);
      });

      it("현재 시간이 만료 시간 전이면 false를 반환한다", () => {
        const store = useAuthStore();
        store.accessTokenExpiresIn = Date.now() + 3600000;

        expect(store.isTokenExpired()).toBe(false);
      });
    });

    describe("isTokenNearExpiration", () => {
      it("accessTokenExpiresIn이 없으면 false를 반환한다", () => {
        const store = useAuthStore();

        expect(store.isTokenNearExpiration()).toBe(false);
      });

      it("만료까지 5분 미만이면 true를 반환한다", () => {
        const store = useAuthStore();
        store.accessTokenExpiresIn = Date.now() + 4 * 60 * 1000;

        expect(store.isTokenNearExpiration()).toBe(true);
      });

      it("만료까지 5분 이상이면 false를 반환한다", () => {
        const store = useAuthStore();
        store.accessTokenExpiresIn = Date.now() + 10 * 60 * 1000;

        expect(store.isTokenNearExpiration()).toBe(false);
      });
    });

    describe("refreshAccessToken", () => {
      it("리프레시 토큰이 없으면 에러를 발생시킨다", async () => {
        const store = useAuthStore();

        await expect(store.refreshAccessToken()).rejects.toThrow("리프레시 토큰이 없습니다.");
      });

      it("토큰 갱신 성공 시 새 토큰을 저장한다", async () => {
        const mockResponse = createMockRefreshResponse();
        vi.mocked(authApi.postRefreshToken).mockResolvedValue(createMockApiResponse(mockResponse));
        vi.mocked(mainApi.getMy).mockResolvedValue(createMockApiResponse(createMockUser()));

        const store = useAuthStore();
        store.refreshToken = "old-refresh-token";
        store.accessTokenExpiresIn = Date.now() - 1000;

        await store.refreshAccessToken();

        expect(store.accessToken).toBe(mockResponse.accessToken);
        expect(store.refreshToken).toBe(mockResponse.refreshToken);
        expect(store.refreshAttempts).toBe(0);
      });

      it("최대 재시도 횟수 초과 시 에러를 발생시킨다", async () => {
        const store = useAuthStore();
        store.refreshToken = "token";
        store.refreshAttempts = 3;

        await expect(store.refreshAccessToken()).rejects.toThrow("리프래시 토큰 재시도 횟수 초과");
      });
    });

    describe("setTokens", () => {
      it("토큰을 상태와 localStorage에 저장한다", () => {
        const store = useAuthStore();
        const expiresIn = Date.now() + 3600000;

        store.setTokens("access", "refresh", expiresIn);

        expect(store.accessToken).toBe("access");
        expect(store.refreshToken).toBe("refresh");
        expect(store.accessTokenExpiresIn).toBe(expiresIn);
        expect(localStorage.getItem("refresh_token")).toBe("refresh");
        expect(localStorage.getItem("token_expires")).toBe(String(expiresIn));
      });
    });

    describe("clearTokens", () => {
      it("모든 토큰과 관련 상태를 제거한다", () => {
        const store = useAuthStore();
        store.accessToken = "access";
        store.refreshToken = "refresh";
        store.accessTokenExpiresIn = Date.now() + 3600000;
        store.refreshAttempts = 2;
        store.lastRefreshAttempt = Date.now();
        localStorage.setItem("refresh_token", "refresh");
        localStorage.setItem("token_expires", "123456");

        store.clearTokens();

        expect(store.accessToken).toBeNull();
        expect(store.refreshToken).toBeNull();
        expect(store.accessTokenExpiresIn).toBeNull();
        expect(store.refreshAttempts).toBe(0);
        expect(store.lastRefreshAttempt).toBe(0);
        expect(localStorage.getItem("refresh_token")).toBeNull();
        expect(localStorage.getItem("token_expires")).toBeNull();
      });
    });

    describe("토큰 자동 갱신 타이머", () => {
      it("startTokenRefreshTimer는 1분마다 토큰 상태를 확인한다", () => {
        const store = useAuthStore();
        store.user = createMockUser();
        store.accessToken = "token";
        store.refreshToken = "refresh";
        store.accessTokenExpiresIn = Date.now() + 4 * 60 * 1000;

        vi.mocked(authApi.postRefreshToken).mockResolvedValue(createMockApiResponse(createMockRefreshResponse()));

        store.startTokenRefreshTimer();

        expect(store.tokenRefreshTimer).not.toBeNull();

        vi.advanceTimersByTime(60000);

        expect(authApi.postRefreshToken).toHaveBeenCalled();
      });

      it("stopTokenRefreshTimer는 타이머를 중지한다", () => {
        const store = useAuthStore();
        store.startTokenRefreshTimer();

        expect(store.tokenRefreshTimer).not.toBeNull();

        store.stopTokenRefreshTimer();

        expect(store.tokenRefreshTimer).toBeNull();
      });
    });
  });
});
