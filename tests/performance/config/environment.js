/**
 * k6 성능 테스트를 위한 환경 설정
 */

export const config = {
  BASE_URL: __ENV.BASE_URL || 'https://staging.rms.bellsoft.net',

  // 테스트 사용자 인증정보
  TEST_USER: {
    username: __ENV.TEST_USERNAME || '@@@-PLEASE-SET-USERNAME-@@@',
    password: __ENV.TEST_PASSWORD || '@@@-PLEASE-SET-PASSWORD-@@@'
  },

  // API 엔드포인트
  API: {
    LOGIN: '/api/v1/auth/login',
    REFRESH: '/api/v1/auth/refresh',
    MY: '/api/v1/my',
    RESERVATIONS: '/api/v1/reservations',
    ROOMS: '/api/v1/rooms',
    ROOM_GROUPS: '/api/v1/room-groups',
    PAYMENT_METHODS: '/api/v1/payment-methods',
    USERS: '/api/v1/admin/accounts'
  },

  // 요청 타임아웃 설정
  TIMEOUTS: {
    DEFAULT: '30s',
    LOGIN: '10s',
    FAST: '5s'
  },

  // 테스트 임계값 설정
  THRESHOLDS: {
    http_req_duration: ['p(95)<2000'], // 95%의 요청이 2초 이내에 완료되어야 함
    http_req_failed: ['rate<0.1'],     // 오류율이 10% 이하여야 함
    http_reqs: ['rate>10']             // 초당 요청수가 10개 이상이어야 함
  }
};

/**
 * 엔드포인트에 대한 전체 URL 반환
 */
export function getUrl(endpoint) {
  return `${config.BASE_URL}${endpoint}`;
}

/**
 * 기본 요청 헤더 반환
 */
export function getHeaders(accessToken = null) {
  const headers = {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  };

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }

  return headers;
}
