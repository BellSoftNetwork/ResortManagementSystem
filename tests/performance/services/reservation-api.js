/**
 * k6 성능 테스트용 예약 API 서비스
 */
import { config } from "../config/environment.js";
import { BaseApiService } from "./base-api.js";
import { Logger } from "../utils/logger.js";

export class ReservationApiService extends BaseApiService {
  constructor() {
    super(config.API.RESERVATIONS);
  }

  /**
   * 필터를 사용하여 예약 목록 조회
   */
  getReservations(filters = {}) {
    const defaultFilters = {
      stayStartAt: '2025-05-26',
      stayEndAt: '2025-07-06',
      status: 'NORMAL',
      type: 'STAY',
      size: 200,
      page: 0
    };

    const queryParams = Object.assign({}, defaultFilters, filters);

    Logger.debug('필터로 예약 목록 조회', queryParams);

    const response = this.get('', queryParams);
    return this.getSuccessfulResponse(response, '예약 목록 조회');
  }

  /**
   * ID로 예약 조회
   */
  getReservationById(id) {
    Logger.debug('ID로 예약 조회', { id });

    const response = this.get(id.toString());
    return this.getSuccessfulResponse(response, `예약 ${id} 조회`);
  }

  /**
   * 새 예약 생성
   */
  createReservation(reservationData) {
    Logger.debug('예약 생성', reservationData);

    const response = this.post('', reservationData);
    return this.getSuccessfulResponse(response, '예약 생성');
  }

  /**
   * 예약 수정
   */
  updateReservation(id, reservationData) {
    Logger.debug('예약 수정', Object.assign({ id: id }, reservationData));

    const response = this.put(id.toString(), reservationData);
    return this.getSuccessfulResponse(response, `예약 ${id} 수정`);
  }

  /**
   * 예약 삭제
   */
  deleteReservation(id) {
    Logger.debug('예약 삭제', { id });

    const response = this.delete(id.toString());
    return response.status >= 200 && response.status < 300;
  }

  /**
   * 예약 통계 조회
   */
  getReservationStats(filters = {}) {
    Logger.debug('예약 통계 조회', filters);

    const response = this.get('statistics', filters);
    return this.getSuccessfulResponse(response, '예약 통계 조회');
  }

  /**
   * 목록에서 첫 번째 예약 ID 반환
   */
  getFirstReservationId(filters = {}) {
    const reservations = this.getReservations(filters);

    if (reservations && reservations.content && reservations.content.length > 0) {
      const firstReservation = reservations.content[0];
      Logger.debug('첫 번째 예약 발견', { id: firstReservation.id });
      return firstReservation.id;
    }

    Logger.warn('목록에서 예약을 찾을 수 없음');
    return null;
  }

  /**
   * 목록에서 랜덤 예약 ID 반환
   */
  getRandomReservationId(filters = {}) {
    const reservations = this.getReservations(filters);

    if (reservations && reservations.content && reservations.content.length > 0) {
      const randomIndex = Math.floor(Math.random() * reservations.content.length);
      const randomReservation = reservations.content[randomIndex];
      Logger.debug('랜덤 예약 발견', {
        id: randomReservation.id,
        index: randomIndex,
        total: reservations.content.length
      });
      return randomReservation.id;
    }

    Logger.warn('목록에서 예약을 찾을 수 없음');
    return null;
  }
}

export const reservationApi = new ReservationApiService();
