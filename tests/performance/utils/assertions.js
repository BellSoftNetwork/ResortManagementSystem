/**
 * k6 테스트용 검증 유틸리티
 */
import { check } from "k6";
import { Logger } from "./logger.js";

export class Assertions {
  /**
   * 응답이 성공인지 확인 (2xx 상태코드)
   */
  static isSuccess(response, operation = '요청') {
    return check(response, {
      [`${operation}: 상태코드가 2xx`]: (r) => r.status >= 200 && r.status < 300,
    });
  }

  /**
   * 특정 상태코드 확인
   */
  static hasStatus(response, expectedStatus, operation = '요청') {
    return check(response, {
      [`${operation}: 상태코드가 ${expectedStatus}`]: (r) => r.status === expectedStatus,
    });
  }

  /**
   * 응답시간 확인
   */
  static isResponseTimeFast(response, maxMs = 2000, operation = '요청') {
    return check(response, {
      [`${operation}: 응답시간 < ${maxMs}ms`]: (r) => r.timings.duration < maxMs,
    });
  }

  /**
   * 응답에 필수 필드가 있는지 확인
   */
  static hasRequiredFields(response, requiredFields, operation = '요청') {
    const checks = {};

    try {
      const body = JSON.parse(response.body);

      requiredFields.forEach(field => {
        checks[`${operation}: '${field}' 필드 존재`] = () => {
          return body.hasOwnProperty(field);
        };
      });

      return check(response, checks);
    } catch (e) {
      Logger.error(`${operation}의 응답 본문 파싱 실패`, e.message);
      return false;
    }
  }

  /**
   * 인증 오류 확인
   */
  static isAuthError(response) {
    return response.status === 401 || response.status === 403;
  }

  /**
   * 성공적인 API 호출에 대한 통합 검증
   */
  static isSuccessfulApiCall(response, operation = 'API 호출', maxMs = 2000) {
    const success = this.isSuccess(response, operation);
    const fast = this.isResponseTimeFast(response, maxMs, operation);

    Logger.apiCall('API', operation, response.status, response.timings.duration);

    return success && fast;
  }
}
