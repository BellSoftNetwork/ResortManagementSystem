/**
 * k6 성능 테스트용 객실그룹 API 서비스
 */
import { config } from "../config/environment.js";
import { BaseApiService } from "./base-api.js";
import { Logger } from "../utils/logger.js";

export class RoomGroupApiService extends BaseApiService {
  constructor() {
    super(config.API.ROOM_GROUPS);
  }

  /**
   * 필터를 사용하여 객실그룹 목록 조회
   */
  getRoomGroups(filters = {}) {
    const defaultFilters = {
      size: 100,
      page: 0
    };

    const queryParams = Object.assign({}, defaultFilters, filters);

    Logger.debug('필터로 객실그룹 목록 조회', queryParams);

    const response = this.get('', queryParams);
    return this.getSuccessfulResponse(response, '객실그룹 목록 조회');
  }

  /**
   * ID로 객실그룹 조회
   */
  getRoomGroupById(id) {
    Logger.debug('ID로 객실그룹 조회', { id });

    const response = this.get(id.toString());
    return this.getSuccessfulResponse(response, `객실그룹 ${id} 조회`);
  }

  /**
   * 새 객실그룹 생성
   */
  createRoomGroup(roomGroupData) {
    Logger.debug('객실그룹 생성', roomGroupData);

    const response = this.post('', roomGroupData);
    return this.getSuccessfulResponse(response, '객실그룹 생성');
  }

  /**
   * 객실그룹 수정
   */
  updateRoomGroup(id, roomGroupData) {
    Logger.debug('객실그룹 수정', Object.assign({ id: id }, roomGroupData));

    const response = this.put(id.toString(), roomGroupData);
    return this.getSuccessfulResponse(response, `객실그룹 ${id} 수정`);
  }

  /**
   * 객실그룹 삭제
   */
  deleteRoomGroup(id) {
    Logger.debug('객실그룹 삭제', { id });

    const response = this.delete(id.toString());
    return response.status >= 200 && response.status < 300;
  }

  /**
   * 목록에서 첫 번째 객실그룹 ID 반환
   */
  getFirstRoomGroupId(filters = {}) {
    const roomGroups = this.getRoomGroups(filters);

    if (roomGroups && roomGroups.content && roomGroups.content.length > 0) {
      const firstRoomGroup = roomGroups.content[0];
      Logger.debug('첫 번째 객실그룹 발견', { id: firstRoomGroup.id });
      return firstRoomGroup.id;
    }

    Logger.warn('목록에서 객실그룹을 찾을 수 없음');
    return null;
  }

  /**
   * 목록에서 랜덤 객실그룹 ID 반환
   */
  getRandomRoomGroupId(filters = {}) {
    const roomGroups = this.getRoomGroups(filters);

    if (roomGroups && roomGroups.content && roomGroups.content.length > 0) {
      const randomIndex = Math.floor(Math.random() * roomGroups.content.length);
      const randomRoomGroup = roomGroups.content[randomIndex];
      Logger.debug('랜덤 객실그룹 발견', {
        id: randomRoomGroup.id,
        index: randomIndex,
        total: roomGroups.content.length
      });
      return randomRoomGroup.id;
    }

    Logger.warn('목록에서 객실그룹을 찾을 수 없음');
    return null;
  }
}

export const roomGroupApi = new RoomGroupApiService();
