/**
 * 예약 부하 테스트 시나리오
 *
 * 예약 목록 API를 지속적으로 호출하여
 * 1분 동안 최대 처리량을 측정하는 고강도 부하 테스트.
 *
 * 특징:
 * - VU당 1회 로그인
 * - 지연 없는 연속 API 호출
 * - 인증 오류 시 자동 토큰 갱신
 * - 최대 처리량 측정
 */

import { Logger } from "../utils/logger.js";
import { authService } from "../services/auth.js";
import { reservationApi } from "../services/reservation-api.js";

export const options = {
  scenarios: {
    reservation_load_test: {
      executor: 'constant-vus',
      vus: 20,                    // 20명의 동시 사용자
      duration: '1m',             // 1분간 실행
      gracefulStop: '10s',        // 요청 완료를 위해 10초 대기
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<5000'],     // 95%의 요청이 5초 이내에 완료되어야 함
    http_req_failed: ['rate<0.2'],         // 오류율이 20% 이하여야 함
    http_reqs: ['rate>50'],                // 초당 50개 이상의 요청을 처리해야 함
    iteration_duration: ['p(95)<6000'],    // 95%의 반복이 6초 이내에 완료되어야 함
  },
};

export function setup() {
  Logger.info('예약 부하 테스트 초기화 중');

  // 초기화 검증을 위한 로그인 수행
  const testAuth = authService.login();
  if (!testAuth) {
    throw new Error('초기화 실패: 인증할 수 없음');
  }

  // 예약 API 테스트
  const testReservations = reservationApi.getReservations();
  if (!testReservations) {
    throw new Error('초기화 실패: 예약 정보를 가져올 수 없음');
  }

  Logger.info('초기화 성공적으로 완료됨', {
    reservationCount: testReservations.totalElements || (testReservations.content && testReservations.content.length) || 0
  });

  return { setupComplete: true };
}

export default function() {
  // 각 VU는 시작 시 한 번 로그인 수행
  if (!authService.accessToken) {
    Logger.info(`VU ${__VU}: 부하 테스트를 위한 초기 로그인 수행`);

    const loginSuccess = authService.login();
    if (!loginSuccess) {
      Logger.error(`VU ${__VU}: 로그인 실패, VU 중단`);
      return;
    }

    Logger.info(`VU ${__VU}: 로그인 성공, 부하 테스트 시작`);
  }

  // 지연 없는 연속 API 호출
  performHighFrequencyReservationCall();
}

function performHighFrequencyReservationCall() {
  // 일관된 테스트를 위한 최적화된 필터
  const reservationFilters = {
    stayStartAt: '2025-05-26',
    stayEndAt: '2025-07-06',
    status: 'NORMAL',
    type: 'STAY',
    size: 200,
    page: 0
  };

  Logger.debug(`VU ${__VU}: 고빈도 예약 API 호출 수행`);

  const reservations = reservationApi.getReservations(reservationFilters);

  if (reservations) {
    Logger.debug(`VU ${__VU}: 부하 테스트 호출 성공`, {
      totalElements: reservations.totalElements,
      contentLength: reservations.content && reservations.content.length,
      currentPage: reservations.number,
      totalPages: reservations.totalPages
    });
  } else {
    Logger.warn(`VU ${__VU}: 부하 테스트 호출 실패`);
  }

  // sleep 없음 - 최대 처리량을 위한 연속 호출
}

export function teardown(data) {
  Logger.info('예약 부하 테스트 정리 중');

  // 통계가 있다면 최종 로그 출력
  if (data && data.setupComplete) {
    Logger.info('부하 테스트 성공적으로 완료됨');
  }
}
