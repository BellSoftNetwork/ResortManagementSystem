/**
 * 혼합 워크로드 테스트 예제
 *
 * 실제 사용자 행동 패턴을 시뮬레이션하여
 * 다양한 API 호출을 혼합한 현실적인 테스트 시나리오.
 */

import { sleep } from "k6";
import { config } from "../config/environment.js";
import { Logger } from "../utils/logger.js";
import { authService } from "../services/auth.js";
import { reservationApi } from "../services/reservation-api.js";
import { roomApi } from "../services/room-api.js";
import { roomGroupApi } from "../services/room-group-api.js";

export const options = {
  scenarios: {
    mixed_workload: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '1m', target: 5 },    // 1분에 걸쳐 5명으로 증가
        { duration: '3m', target: 5 },    // 5명으로 3분간 유지
        { duration: '1m', target: 10 },   // 1분에 걸쳐 10명으로 증가
        { duration: '3m', target: 10 },   // 10명으로 3분간 유지
        { duration: '1m', target: 0 },    // 1분에 걸쳐 0명으로 감소
      ],
    },
  },
  thresholds: config.THRESHOLDS,
};

export function setup() {
  Logger.info('혼합 워크로드 테스트 초기화 중');

  const testAuth = authService.login();
  if (!testAuth) {
    throw new Error('초기화 실패: 인증할 수 없음');
  }

  return { setupComplete: true };
}

export default function() {
  // 각 VU는 한 번 로그인 수행
  if (!authService.accessToken) {
    Logger.info(`VU ${__VU}: 초기 로그인 수행`);

    const loginSuccess = authService.login();
    if (!loginSuccess) {
      Logger.error(`VU ${__VU}: 로그인 실패, 반복 건너뜀`);
      return;
    }
  }

  // 다양한 사용자 행동 패턴 시뮬레이션
  const userBehavior = Math.random();

  if (userBehavior < 0.4) {
    // 40% - 예약 중심 작업
    performReservationHeavyWorkload();
  } else if (userBehavior < 0.7) {
    // 30% - 객실 관리 작업
    performRoomManagementWorkload();
  } else if (userBehavior < 0.9) {
    // 20% - 일반 조회 작업
    performGeneralBrowsingWorkload();
  } else {
    // 10% - 특정 예약 심층 분석
    performDetailedReservationWorkload();
  }
}

function performReservationHeavyWorkload() {
  Logger.debug(`VU ${__VU}: 예약 중심 워크로드 수행`);

  // 다양한 필터로 여러 예약 목록 조회
  const filters = [
    { status: 'NORMAL', type: 'STAY' },
    { status: 'NORMAL', type: 'MONTHLY_RENT' },
    { status: 'CANCELED' }
  ];

  filters.forEach((filter, index) => {
    const reservations = reservationApi.getReservations(filter);
    if (reservations) {
      Logger.debug(`VU ${__VU}: 필터 ${index + 1}로 예약 조회 완료`, {
        count: reservations.content && reservations.content.length
      });
    }
    sleep(0.5);
  });

  // 예약 통계 조회
  const stats = reservationApi.getReservationStats();
  if (stats) {
    Logger.debug(`VU ${__VU}: 예약 통계 조회 완료`);
  }

  sleep(1);
}

function performRoomManagementWorkload() {
  Logger.debug(`VU ${__VU}: 객실 관리 워크로드 수행`);

  // 먼저 객실그룹 조회
  const roomGroups = roomGroupApi.getRoomGroups();
  if (roomGroups && roomGroups.content && roomGroups.content.length > 0) {
    const randomGroupId = roomGroups.content[Math.floor(Math.random() * roomGroups.content.length)].id;
    const groupDetails = roomGroupApi.getRoomGroupById(randomGroupId);
    if (groupDetails) {
      Logger.debug(`VU ${__VU}: 객실그룹 상세 정보 조회 완료`, { id: randomGroupId });
    }
  }

  sleep(0.8);

  // 객실 목록 조회
  const rooms = roomApi.getRooms();
  if (rooms && rooms.content && rooms.content.length > 0) {
    // 2-3개의 랜덤 객실 상세 정보 조회
    const roomCount = Math.min(3, rooms.content.length);
    for (let i = 0; i < roomCount; i++) {
      const randomRoom = rooms.content[Math.floor(Math.random() * rooms.content.length)];
      const roomDetails = roomApi.getRoomById(randomRoom.id);
      if (roomDetails) {
        Logger.debug(`VU ${__VU}: 객실 상세 정보 조회 완료`, { id: randomRoom.id });
      }
      sleep(0.3);
    }
  }

  sleep(1);
}

function performGeneralBrowsingWorkload() {
  Logger.debug(`VU ${__VU}: 일반 조회 워크로드 수행`);

  // 사용자 정보 확인
  const authTest = authService.testAuthentication();
  if (authTest) {
    Logger.debug(`VU ${__VU}: 인증 확인 완료`);
  }

  sleep(0.5);

  // 예약 조회
  const reservations = reservationApi.getReservations();
  if (reservations) {
    Logger.debug(`VU ${__VU}: 예약 조회 완료`, {
      count: reservations.content && reservations.content.length
    });
  }

  sleep(1);

  // 객실 조회
  const rooms = roomApi.getRooms({ size: 50 });
  if (rooms) {
    Logger.debug(`VU ${__VU}: 객실 조회 완료`, {
      count: rooms.content && rooms.content.length
    });
  }

  sleep(1.5);
}

function performDetailedReservationWorkload() {
  Logger.debug(`VU ${__VU}: 상세 예약 분석 워크로드 수행`);

  // 예약 조회
  const reservations = reservationApi.getReservations();
  if (reservations && reservations.content && reservations.content.length > 0) {
    // 심층 분석을 위한 랜덤 예약 선택
    const randomReservation = reservations.content[Math.floor(Math.random() * reservations.content.length)];

    // 상세 정보 조회
    const details = reservationApi.getReservationById(randomReservation.id);
    if (details) {
      Logger.debug(`VU ${__VU}: 예약 심층 분석 완료`, {
        id: details.id,
        status: details.status
      });

      sleep(2); // 사용자가 상세 정보를 읽는 시간 시뮬레이션

      // 관련 객실 정보 조회 (있는 경우)
      if (details.rooms && details.rooms.length > 0) {
        const roomId = details.rooms[0].id;
        const roomDetails = roomApi.getRoomById(roomId);
        if (roomDetails) {
          Logger.debug(`VU ${__VU}: 관련 객실 상세 정보 조회 완료`, { id: roomId });
        }
        sleep(1);
      }
    }
  }

  sleep(1);
}

export function teardown(data) {
  Logger.info('혼합 워크로드 테스트 정리 중');
}
