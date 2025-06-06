/**
 * k6 성능 테스트용 기본 API 서비스
 */
import http from "k6/http";
import { config, getUrl } from "../config/environment.js";
import { Logger } from "../utils/logger.js";
import { Assertions } from "../utils/assertions.js";
import { authService } from "./auth.js";

export class BaseApiService {
  constructor(resourcePath) {
    this.resourcePath = resourcePath;
  }

  /**
   * 인증된 GET 요청 수행
   */
  get(endpoint = '', queryParams = {}, options = {}) {
    const url = this._buildUrl(endpoint, queryParams);
    const requestOptions = this._buildRequestOptions(options);

    const response = http.get(url, requestOptions);

    Logger.apiCall('GET', endpoint || this.resourcePath, response.status, response.timings.duration);

    // Handle authentication errors
    if (Assertions.isAuthError(response)) {
      if (authService.handleAuthError()) {
        Logger.info('인증 복구 후 GET 요청 재시도');
        const retryOptions = this._buildRequestOptions(options);
        const retryResponse = http.get(url, retryOptions);
        Logger.apiCall('GET', endpoint || this.resourcePath + ' (재시도)', retryResponse.status, retryResponse.timings.duration);
        return retryResponse;
      }
    }

    return response;
  }

  /**
   * 인증된 POST 요청 수행
   */
  post(endpoint = '', data = {}, options = {}) {
    const url = this._buildUrl(endpoint);
    const requestOptions = this._buildRequestOptions(options);

    const response = http.post(url, JSON.stringify(data), requestOptions);

    Logger.apiCall('POST', endpoint || this.resourcePath, response.status, response.timings.duration);

    // Handle authentication errors
    if (Assertions.isAuthError(response)) {
      if (authService.handleAuthError()) {
        Logger.info('인증 복구 후 POST 요청 재시도');
        const retryOptions = this._buildRequestOptions(options);
        const retryResponse = http.post(url, JSON.stringify(data), retryOptions);
        Logger.apiCall('POST', endpoint || this.resourcePath + ' (재시도)', retryResponse.status, retryResponse.timings.duration);
        return retryResponse;
      }
    }

    return response;
  }

  /**
   * 인증된 PUT 요청 수행
   */
  put(endpoint = '', data = {}, options = {}) {
    const url = this._buildUrl(endpoint);
    const requestOptions = this._buildRequestOptions(options);

    const response = http.put(url, JSON.stringify(data), requestOptions);

    Logger.apiCall('PUT', endpoint || this.resourcePath, response.status, response.timings.duration);

    // Handle authentication errors
    if (Assertions.isAuthError(response)) {
      if (authService.handleAuthError()) {
        Logger.info('인증 복구 후 PUT 요청 재시도');
        const retryOptions = this._buildRequestOptions(options);
        const retryResponse = http.put(url, JSON.stringify(data), retryOptions);
        Logger.apiCall('PUT', endpoint || this.resourcePath + ' (재시도)', retryResponse.status, retryResponse.timings.duration);
        return retryResponse;
      }
    }

    return response;
  }

  /**
   * 인증된 DELETE 요청 수행
   */
  delete(endpoint = '', options = {}) {
    const url = this._buildUrl(endpoint);
    const requestOptions = this._buildRequestOptions(options);

    const response = http.del(url, null, requestOptions);

    Logger.apiCall('DELETE', endpoint || this.resourcePath, response.status, response.timings.duration);

    // Handle authentication errors
    if (Assertions.isAuthError(response)) {
      if (authService.handleAuthError()) {
        Logger.info('인증 복구 후 DELETE 요청 재시도');
        const retryOptions = this._buildRequestOptions(options);
        const retryResponse = http.del(url, null, retryOptions);
        Logger.apiCall('DELETE', endpoint || this.resourcePath + ' (재시도)', retryResponse.status, retryResponse.timings.duration);
        return retryResponse;
      }
    }

    return response;
  }

  /**
   * 쿼리 파라미터를 포함한 전체 URL 생성
   */
  _buildUrl(endpoint = '', queryParams = {}) {
    let url = getUrl(this.resourcePath);

    if (endpoint) {
      // 앞에 있는 슬래시 제거
      endpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;
      url += `/${endpoint}`;
    }

    // 쿼리 파라미터 추가 (수동으로 구성)
    const paramPairs = [];
    Object.entries(queryParams).forEach(([key, value]) => {
      if (value !== null && value !== undefined) {
        paramPairs.push(encodeURIComponent(key) + '=' + encodeURIComponent(value.toString()));
      }
    });

    if (paramPairs.length > 0) {
      url += '?' + paramPairs.join('&');
    }

    return url;
  }

  /**
   * 인증 헤더를 포함한 요청 옵션 생성
   */
  _buildRequestOptions(customOptions = {}) {
    const defaultOptions = {
      headers: authService.getAuthHeaders(),
      timeout: config.TIMEOUTS.DEFAULT
    };

    // 사용자 정의 옵션 병합
    return Object.assign({}, defaultOptions, customOptions, {
      headers: Object.assign({}, defaultOptions.headers, customOptions.headers || {})
    });
  }

  /**
   * 응답 본문을 안전하게 파싱
   */
  parseResponse(response) {
    try {
      return JSON.parse(response.body);
    } catch (e) {
      Logger.error('응답 본문 파싱 실패', {
        status: response.status,
        body: response.body
      });
      return null;
    }
  }

  /**
   * 응답이 성공인지 확인하고 파싱된 본문 반환
   */
  getSuccessfulResponse(response, operation = '요청') {
    if (Assertions.isSuccessfulApiCall(response, operation)) {
      const parsed = this.parseResponse(response);

      // API 응답 구조에 따라 실제 데이터 추출
      // 단건: { "value": { ... } }
      // 리스트: { "values": [ ... ] }
      if (parsed) {
        if (parsed.value !== undefined) {
          return parsed.value; // 단건 응답
        } else if (parsed.values !== undefined) {
          // 리스트 응답을 페이징 정보와 함께 반환
          return {
            content: parsed.values,
            totalElements: parsed.totalElements || parsed.values.length,
            totalPages: parsed.totalPages || 1,
            number: parsed.number || 0,
            size: parsed.size || parsed.values.length
          };
        } else {
          return parsed; // 기존 구조 그대로 반환
        }
      }

      return parsed;
    }
    return null;
  }
}
