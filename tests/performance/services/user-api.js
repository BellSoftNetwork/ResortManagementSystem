/**
 * k6 성능 테스트용 사용자 API 서비스
 */
import http from "k6/http";
import { config, getUrl } from "../config/environment.js";
import { BaseApiService } from "./base-api.js";
import { Logger } from "../utils/logger.js";
import { authService } from "./auth.js";

export class UserApiService extends BaseApiService {
  constructor() {
    super(config.API.USERS);
  }

  /**
   * 필터를 사용하여 사용자 목록 조회
   */
  getUsers(filters = {}) {
    const defaultFilters = {
      size: 100,
      page: 0
    };

    const queryParams = Object.assign({}, defaultFilters, filters);

    Logger.debug('필터로 사용자 목록 조회', queryParams);

    const response = this.get('', queryParams);
    return this.getSuccessfulResponse(response, '사용자 목록 조회');
  }

  /**
   * ID로 사용자 조회
   */
  getUserById(id) {
    Logger.debug('ID로 사용자 조회', { id });

    const response = this.get(id.toString());
    return this.getSuccessfulResponse(response, `사용자 ${id} 조회`);
  }

  /**
   * 새 사용자 생성
   */
  createUser(userData) {
    Logger.debug('사용자 생성', userData);

    const response = this.post('', userData);
    return this.getSuccessfulResponse(response, '사용자 생성');
  }

  /**
   * 사용자 수정
   */
  updateUser(id, userData) {
    Logger.debug('사용자 수정', Object.assign({ id: id }, userData));

    const response = this.put(id.toString(), userData);
    return this.getSuccessfulResponse(response, `사용자 ${id} 수정`);
  }

  /**
   * 사용자 삭제
   */
  deleteUser(id) {
    Logger.debug('사용자 삭제', { id });

    const response = this.delete(id.toString());
    return response.status >= 200 && response.status < 300;
  }

  /**
   * 현재 사용자 정보 조회
   */
  getMyInfo() {
    Logger.debug('현재 사용자 정보 조회');

    // /api/v1/my 엔드포인트를 직접 호출
    const response = http.get(
      getUrl(config.API.MY),
      {
        headers: authService.getAuthHeaders(),
        timeout: config.TIMEOUTS.DEFAULT
      }
    );

    Logger.apiCall('GET', config.API.MY, response.status, response.timings.duration);
    return this.getSuccessfulResponse(response, '내 정보 조회');
  }

  /**
   * 목록에서 첫 번째 사용자 ID 반환
   */
  getFirstUserId(filters = {}) {
    const users = this.getUsers(filters);

    if (users && users.content && users.content.length > 0) {
      const firstUser = users.content[0];
      Logger.debug('첫 번째 사용자 발견', { id: firstUser.id });
      return firstUser.id;
    }

    Logger.warn('목록에서 사용자를 찾을 수 없음');
    return null;
  }

  /**
   * 목록에서 랜덤 사용자 ID 반환
   */
  getRandomUserId(filters = {}) {
    const users = this.getUsers(filters);

    if (users && users.content && users.content.length > 0) {
      const randomIndex = Math.floor(Math.random() * users.content.length);
      const randomUser = users.content[randomIndex];
      Logger.debug('랜덤 사용자 발견', {
        id: randomUser.id,
        index: randomIndex,
        total: users.content.length
      });
      return randomUser.id;
    }

    Logger.warn('목록에서 사용자를 찾을 수 없음');
    return null;
  }
}

export const userApi = new UserApiService();
