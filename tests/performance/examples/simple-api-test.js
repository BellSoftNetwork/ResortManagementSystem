/**
 * 간단한 API 테스트 예제
 *
 * 개별 API들의 기본 기능을 테스트하는
 * 인증 및 오류 처리가 포함된 기본 예제.
 */

import { config } from "../config/environment.js";
import { Logger } from "../utils/logger.js";
import { authService } from "../services/auth.js";
import { reservationApi } from "../services/reservation-api.js";
import { roomApi } from "../services/room-api.js";
import { roomGroupApi } from "../services/room-group-api.js";
import { paymentMethodApi } from "../services/payment-method-api.js";
import { userApi } from "../services/user-api.js";

export const options = {
  scenarios: {
    simple_api_test: {
      executor: 'shared-iterations',
      vus: 1,
      iterations: 1,
    },
  },
  thresholds: config.THRESHOLDS,
};

export function setup() {
  Logger.info('간단한 API 테스트 초기화 중');
  return { setupComplete: true };
}

export default function() {
  Logger.info('간단한 API 테스트 시작');

  // 1단계: 로그인
  Logger.info('1단계: 로그인 수행');
  const loginSuccess = authService.login();
  if (!loginSuccess) {
    Logger.error('로그인 실패, 테스트 중단');
    return;
  }

  // 2단계: 인증 테스트
  Logger.info('2단계: 인증 테스트');
  const authTest = authService.testAuthentication();
  if (!authTest) {
    Logger.error('인증 테스트 실패');
    return;
  }

  // 3단계: 예약 API 테스트
  Logger.info('3단계: 예약 API 테스트');
  testReservationApis();

  // 4단계: 객실 API 테스트
  Logger.info('4단계: 객실 API 테스트');
  testRoomApis();

  // 5단계: 객실그룹 API 테스트
  Logger.info('5단계: 객실그룹 API 테스트');
  testRoomGroupApis();

  // 6단계: 결제수단 API 테스트
  Logger.info('6단계: 결제수단 API 테스트');
  testPaymentMethodApis();

  // 7단계: 사용자 API 테스트
  Logger.info('7단계: 사용자 API 테스트');
  testUserApis();

  Logger.info('간단한 API 테스트 완료');
}

function testReservationApis() {
  // 예약 목록 조회
  const reservations = reservationApi.getReservations();
  if (reservations) {
    Logger.info('✓ 예약 목록 조회 성공', {
      totalElements: reservations.totalElements,
      contentLength: reservations.content && reservations.content.length
    });

    // 첫 번째 예약 상세 정보 조회 (있는 경우)
    if (reservations.content && reservations.content.length > 0) {
      const firstId = reservations.content[0].id;
      const reservationDetails = reservationApi.getReservationById(firstId);
      if (reservationDetails) {
        Logger.info('✓ 예약 상세 정보 조회 성공', { id: firstId });
      } else {
        Logger.warn('✗ 예약 상세 정보 조회 실패', { id: firstId });
      }
    }
  } else {
    Logger.warn('✗ 예약 목록 조회 실패');
  }

  // 예약 통계 조회 (필터 매개변수가 필요할 수 있음)
  const stats = reservationApi.getReservationStats({
    stayStartAt: '2025-05-26',
    stayEndAt: '2025-07-06'
  });
  if (stats) {
    Logger.info('✓ 예약 통계 조회 성공');
  } else {
    Logger.warn('✗ 예약 통계 조회 실패 (필터 매개변수 필요할 수 있음)');
  }
}

function testRoomApis() {
  // 객실 목록 조회
  const rooms = roomApi.getRooms();
  if (rooms) {
    Logger.info('✓ 객실 목록 조회 성공', {
      totalElements: rooms.totalElements,
      contentLength: rooms.content && rooms.content.length
    });

    // 첫 번째 객실 상세 정보 조회 (있는 경우)
    if (rooms.content && rooms.content.length > 0) {
      const firstId = rooms.content[0].id;
      const roomDetails = roomApi.getRoomById(firstId);
      if (roomDetails) {
        Logger.info('✓ 객실 상세 정보 조회 성공', { id: firstId });
      } else {
        Logger.warn('✗ 객실 상세 정보 조회 실패', { id: firstId });
      }
    }
  } else {
    Logger.warn('✗ 객실 목록 조회 실패');
  }
}

function testRoomGroupApis() {
  // 객실그룹 목록 조회
  const roomGroups = roomGroupApi.getRoomGroups();
  if (roomGroups) {
    Logger.info('✓ 객실그룹 목록 조회 성공', {
      totalElements: roomGroups.totalElements,
      contentLength: roomGroups.content && roomGroups.content.length
    });

    // 첫 번째 객실그룹 상세 정보 조회 (있는 경우)
    if (roomGroups.content && roomGroups.content.length > 0) {
      const firstId = roomGroups.content[0].id;
      const roomGroupDetails = roomGroupApi.getRoomGroupById(firstId);
      if (roomGroupDetails) {
        Logger.info('✓ 객실그룹 상세 정보 조회 성공', { id: firstId });
      } else {
        Logger.warn('✗ 객실그룹 상세 정보 조회 실패', { id: firstId });
      }
    }
  } else {
    Logger.warn('✗ 객실그룹 목록 조회 실패');
  }
}

function testPaymentMethodApis() {
  // 결제수단 목록 조회
  const paymentMethods = paymentMethodApi.getPaymentMethods();
  if (paymentMethods) {
    Logger.info('✓ 결제수단 목록 조회 성공', {
      totalElements: paymentMethods.totalElements,
      contentLength: paymentMethods.content && paymentMethods.content.length
    });

    // 첫 번째 결제수단 상세 정보 조회 (있는 경우)
    if (paymentMethods.content && paymentMethods.content.length > 0) {
      const firstId = paymentMethods.content[0].id;
      const paymentMethodDetails = paymentMethodApi.getPaymentMethodById(firstId);
      if (paymentMethodDetails) {
        Logger.info('✓ 결제수단 상세 정보 조회 성공', { id: firstId });
      } else {
        Logger.warn('✗ 결제수단 상세 정보 조회 실패', { id: firstId });
      }
    }
  } else {
    Logger.warn('✗ 결제수단 목록 조회 실패');
  }
}

function testUserApis() {
  // 사용자 목록 조회는 관리자 권한이 필요하므로 건너뛰기
  Logger.info('사용자 API 테스트는 관리자 권한이 필요하여 건너뛰기');

  // 대신 내 정보 조회 테스트
  try {
    const myInfo = userApi.getMyInfo();
    if (myInfo) {
      Logger.info('✓ 내 정보 조회 성공', { id: myInfo.id || 'unknown' });
    } else {
      Logger.warn('✗ 내 정보 조회 실패');
    }
  } catch (error) {
    Logger.warn('✗ 내 정보 조회 오류', error.message);
  }
}

export function teardown(data) {
  Logger.info('간단한 API 테스트 정리 중');
}
