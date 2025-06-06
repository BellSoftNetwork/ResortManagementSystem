/**
 * 예약 정보 조회 시나리오
 *
 * 사용자가 예약 정보를 조회하는 과정을 시뮬레이션하는 테스트 시나리오:
 * 1. 로그인
 * 2. /my 엔드포인트로 인증 확인
 * 3. 예약 목록 조회
 * 4. 0.5초 대기
 * 5. 예약 목록 재조회
 * 6. 2초 대기
 * 7. 목록의 첫 번째 예약 상세 정보 조회
 * 8. 조회 사이클 반복
 */

import { sleep } from "k6";
import { config } from "../config/environment.js";
import { Logger } from "../utils/logger.js";
import { authService } from "../services/auth.js";
import { reservationApi } from "../services/reservation-api.js";

export const options = {
  scenarios: {
    reservation_browsing: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '30s', target: 5 },   // 30초에 걸쳐 5명으로 증가
        { duration: '2m', target: 5 },    // 5명으로 2분간 유지
        { duration: '30s', target: 10 },  // 30초에 걸쳐 10명으로 증가
        { duration: '2m', target: 10 },   // 10명으로 2분간 유지
        { duration: '30s', target: 0 },   // 30초에 걸쳐 0명으로 감소
      ],
    },
  },
  thresholds: config.THRESHOLDS,
};

export function setup() {
  Logger.info('예약 조회 테스트 초기화 중');

  // 초기화 검증을 위한 로그인 수행
  const testAuth = authService.login();
  if (!testAuth) {
    throw new Error('초기화 실패: 인증할 수 없음');
  }

  // 인증 테스트
  const authTest = authService.testAuthentication();
  if (!authTest) {
    throw new Error('초기화 실패: 인증 테스트 실패');
  }

  Logger.info('초기화 성공적으로 완료됨');
  return { setupComplete: true };
}

export default function() {
  // 각 VU는 시작 시 한 번 로그인 수행
  if (!authService.accessToken) {
    Logger.info(`VU ${__VU}: 초기 로그인 수행`);

    const loginSuccess = authService.login();
    if (!loginSuccess) {
      Logger.error(`VU ${__VU}: 로그인 실패, 반복 건너뜀`);
      return;
    }

    // 인증이 작동하는지 확인
    const authTest = authService.testAuthentication();
    if (!authTest) {
      Logger.error(`VU ${__VU}: 인증 테스트 실패, 반복 건너뜀`);
      return;
    }

    Logger.info(`VU ${__VU}: 로그인 및 인증 테스트 성공`);
  }

  // 예약 조회 사이클 수행
  performReservationBrowsingCycle();
}

function performReservationBrowsingCycle() {
  Logger.debug(`VU ${__VU}: 예약 조회 사이클 시작`);

  // 1단계: 예약 목록 조회 (첫 번째)
  Logger.debug(`VU ${__VU}: 초기 예약 목록 조회`);
  const reservationFilters = {
    stayStartAt: '2025-05-26',
    stayEndAt: '2025-07-06',
    status: 'NORMAL',
    type: 'STAY',
    size: 200
  };

  const firstReservationList = reservationApi.getReservations(reservationFilters);
  if (!firstReservationList) {
    Logger.warn(`VU ${__VU}: 초기 예약 목록 조회 실패`);
    return;
  }

  // 2단계: 0.5초 대기
  Logger.debug(`VU ${__VU}: 0.5초 대기`);
  sleep(0.5);

  // 3단계: 예약 목록 조회 (두 번째)
  Logger.debug(`VU ${__VU}: 예약 목록 재조회`);
  const secondReservationList = reservationApi.getReservations(reservationFilters);
  if (!secondReservationList) {
    Logger.warn(`VU ${__VU}: 두 번째 예약 목록 조회 실패`);
    return;
  }

  // 4단계: 2초 대기
  Logger.debug(`VU ${__VU}: 2초 대기`);
  sleep(2);

  // 5단계: 첫 번째 예약 상세 정보 조회
  if (firstReservationList.content && firstReservationList.content.length > 0) {
    const firstReservationId = firstReservationList.content[0].id;
    Logger.debug(`VU ${__VU}: 예약 ${firstReservationId} 상세 정보 조회`);

    const reservationDetails = reservationApi.getReservationById(firstReservationId);
    if (reservationDetails) {
      Logger.debug(`VU ${__VU}: 예약 상세 정보 조회 성공`, {
        id: reservationDetails.id,
        status: reservationDetails.status
      });
    } else {
      Logger.warn(`VU ${__VU}: 예약 ID ${firstReservationId} 상세 정보 조회 실패`);
    }
  } else {
    Logger.warn(`VU ${__VU}: 상세 정보를 조회할 예약이 목록에 없음`);
  }

  Logger.debug(`VU ${__VU}: 예약 조회 사이클 완료`);
}

export function teardown(data) {
  Logger.info('예약 조회 테스트 정리 중');
}
