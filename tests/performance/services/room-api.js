/**
 * k6 성능 테스트용 객실 API 서비스
 */
import { config } from "../config/environment.js";
import { BaseApiService } from "./base-api.js";
import { Logger } from "../utils/logger.js";

export class RoomApiService extends BaseApiService {
  constructor() {
    super(config.API.ROOMS);
  }

  /**
   * 필터를 사용하여 객실 목록 조회
   */
  getRooms(filters = {}) {
    const defaultFilters = {
      size: 100,
      page: 0
    };

    const queryParams = Object.assign({}, defaultFilters, filters);

    Logger.debug('필터로 객실 목록 조회', queryParams);

    const response = this.get('', queryParams);
    return this.getSuccessfulResponse(response, '객실 목록 조회');
  }

  /**
   * ID로 객실 조회
   */
  getRoomById(id) {
    Logger.debug('ID로 객실 조회', { id });

    const response = this.get(id.toString());
    return this.getSuccessfulResponse(response, `객실 ${id} 조회`);
  }

  /**
   * 새 객실 생성
   */
  createRoom(roomData) {
    Logger.debug('객실 생성', roomData);

    const response = this.post('', roomData);
    return this.getSuccessfulResponse(response, '객실 생성');
  }

  /**
   * 객실 수정
   */
  updateRoom(id, roomData) {
    Logger.debug('객실 수정', Object.assign({ id: id }, roomData));

    const response = this.put(id.toString(), roomData);
    return this.getSuccessfulResponse(response, `객실 ${id} 수정`);
  }

  /**
   * 객실 삭제
   */
  deleteRoom(id) {
    Logger.debug('객실 삭제', { id });

    const response = this.delete(id.toString());
    return response.status >= 200 && response.status < 300;
  }

  /**
   * 객실 가용성 조회
   */
  getRoomAvailability(filters = {}) {
    Logger.debug('객실 가용성 조회', filters);

    const response = this.get('availability', filters);
    return this.getSuccessfulResponse(response, '객실 가용성 조회');
  }

  /**
   * 목록에서 첫 번째 객실 ID 반환
   */
  getFirstRoomId(filters = {}) {
    const rooms = this.getRooms(filters);

    if (rooms && rooms.content && rooms.content.length > 0) {
      const firstRoom = rooms.content[0];
      Logger.debug('첫 번째 객실 발견', { id: firstRoom.id });
      return firstRoom.id;
    }

    Logger.warn('목록에서 객실을 찾을 수 없음');
    return null;
  }

  /**
   * 목록에서 랜덤 객실 ID 반환
   */
  getRandomRoomId(filters = {}) {
    const rooms = this.getRooms(filters);

    if (rooms && rooms.content && rooms.content.length > 0) {
      const randomIndex = Math.floor(Math.random() * rooms.content.length);
      const randomRoom = rooms.content[randomIndex];
      Logger.debug('랜덤 객실 발견', {
        id: randomRoom.id,
        index: randomIndex,
        total: rooms.content.length
      });
      return randomRoom.id;
    }

    Logger.warn('목록에서 객실을 찾을 수 없음');
    return null;
  }
}

export const roomApi = new RoomApiService();
