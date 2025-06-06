/**
 * k6 성능 테스트용 인증 서비스
 */
import http from "k6/http";
import { config, getHeaders, getUrl } from "../config/environment.js";
import { Logger } from "../utils/logger.js";
import { Assertions } from "../utils/assertions.js";

export class AuthService {
  constructor() {
    this.accessToken = null;
    this.refreshToken = null;
    this.tokenExpiry = null;
  }

  /**
   * 로그인 수행 및 토큰 저장
   */
  login(username = null, password = null) {
    const credentials = {
      username: username || config.TEST_USER.username,
      password: password || config.TEST_USER.password
    };

    Logger.info('로그인 시도 중', { username: credentials.username });

    const response = http.post(
      getUrl(config.API.LOGIN),
      JSON.stringify(credentials),
      {
        headers: getHeaders(),
        timeout: config.TIMEOUTS.LOGIN
      }
    );

    Logger.apiCall('POST', config.API.LOGIN, response.status, response.timings.duration);

    if (Assertions.hasStatus(response, 200, 'Login')) {
      try {
        const body = JSON.parse(response.body);

        // API 응답 구조: { "value": { "accessToken": "...", "refreshToken": "..." } }
        const loginData = body.value || body;

        if (loginData.accessToken && loginData.refreshToken) {
          this.accessToken = loginData.accessToken;
          this.refreshToken = loginData.refreshToken;

          // Calculate expiry time (assume 1 hour if not provided)
          const expiryMinutes = loginData.expiresIn || 60;
          this.tokenExpiry = Date.now() + (expiryMinutes * 60 * 1000);

          Logger.info('로그인 성공', {
            hasAccessToken: !!this.accessToken,
            hasRefreshToken: !!this.refreshToken,
            expiresIn: expiryMinutes + '분'
          });

          return true;
        } else {
          Logger.error('로그인 응답에 토큰이 없음', body);
          return false;
        }
      } catch (e) {
        Logger.error('로그인 응답 파싱 실패', e.message);
        return false;
      }
    } else {
      Logger.error('로그인 실패', {
        status: response.status,
        body: response.body
      });
      return false;
    }
  }

  /**
   * 리프레시 토큰을 사용하여 액세스 토큰 갱신
   */
  refreshAccessToken() {
    if (!this.refreshToken) {
      Logger.error('사용 가능한 리프레시 토큰이 없음');
      return false;
    }

    Logger.info('액세스 토큰 갱신 중');

    const response = http.post(
      getUrl(config.API.REFRESH),
      JSON.stringify({ refreshToken: this.refreshToken }),
      {
        headers: getHeaders(),
        timeout: config.TIMEOUTS.LOGIN
      }
    );

    Logger.apiCall('POST', config.API.REFRESH, response.status, response.timings.duration);

    if (Assertions.hasStatus(response, 200, 'Token Refresh')) {
      try {
        const body = JSON.parse(response.body);

        // API 응답 구조: { "value": { "accessToken": "...", "refreshToken": "..." } }
        const tokenData = body.value || body;

        if (tokenData.accessToken) {
          this.accessToken = tokenData.accessToken;

          // Update refresh token if provided
          if (tokenData.refreshToken) {
            this.refreshToken = tokenData.refreshToken;
          }

          // Update expiry time
          const expiryMinutes = tokenData.expiresIn || 60;
          this.tokenExpiry = Date.now() + (expiryMinutes * 60 * 1000);

          Logger.info('토큰 갱신 성공');
          return true;
        } else {
          Logger.error('토큰 갱신 응답에 액세스 토큰이 없음', body);
          return false;
        }
      } catch (e) {
        Logger.error('토큰 갱신 응답 파싱 실패', e.message);
        return false;
      }
    } else {
      Logger.error('토큰 갱신 실패', {
        status: response.status,
        body: response.body
      });
      return false;
    }
  }

  /**
   * 현재 액세스 토큰 반환, 필요시 갱신
   */
  getAccessToken() {
    // Check if token is about to expire (refresh 5 minutes before expiry)
    if (this.tokenExpiry && Date.now() > (this.tokenExpiry - 5 * 60 * 1000)) {
      Logger.warn('액세스 토큰이 곧 만료됨, 갱신 중...');
      if (!this.refreshAccessToken()) {
        Logger.warn('토큰 갱신 실패, 새로운 로그인 시도 중...');
        if (!this.login()) {
          Logger.error('토큰 갱신 및 재로그인 실패');
          return null;
        }
      }
    }

    return this.accessToken;
  }

  /**
   * 인증이 포함된 헤더 반환
   */
  getAuthHeaders() {
    const token = this.getAccessToken();
    if (!token) {
      Logger.error('유효한 액세스 토큰이 없음');
      return getHeaders();
    }
    return getHeaders(token);
  }

  /**
   * 인증 오류 처리 및 재시도
   */
  handleAuthError(originalRequest) {
    Logger.warn('인증 오류 감지, 토큰 갱신 시도 중');

    if (this.refreshAccessToken()) {
      Logger.info('토큰 갱신 완료, 원래 요청 재시도');
      return true;
    } else {
      Logger.warn('토큰 갱신 실패, 새로운 로그인 시도');
      if (this.login()) {
        Logger.info('재로그인 성공, 원래 요청 재시도');
        return true;
      } else {
        Logger.error('인증 오류 복구 실패');
        return false;
      }
    }
  }

  /**
   * /my 엔드포인트 호출로 인증 테스트
   */
  testAuthentication() {
    const response = http.get(
      getUrl(config.API.MY),
      {
        headers: this.getAuthHeaders(),
        timeout: config.TIMEOUTS.FAST
      }
    );

    Logger.apiCall('GET', config.API.MY, response.status, response.timings.duration);

    if (Assertions.isSuccess(response, 'Authentication Test')) {
      Logger.info('인증 테스트 성공');
      return true;
    } else if (Assertions.isAuthError(response)) {
      Logger.warn('인증 테스트 실패 - 토큰이 유효하지 않음');
      return this.handleAuthError();
    } else {
      Logger.error('인증 테스트 실패', {
        status: response.status,
        body: response.body
      });
      return false;
    }
  }
}

// 전역 인증 인스턴스
export const authService = new AuthService();
