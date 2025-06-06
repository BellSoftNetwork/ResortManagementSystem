# Grafana k6를 이용한 성능 테스트

이 디렉토리는 리조트 관리 시스템 백엔드 API를 위한 Grafana k6 성능 테스트 시나리오를 포함

## 프로젝트 구조

```
tests/performance/
├── config/
│   └── environment.js          # 환경 설정 및 설정값
├── utils/
│   ├── logger.js              # 로깅 유틸리티
│   └── assertions.js          # 테스트 검증 헬퍼
├── services/
│   ├── auth.js                # 토큰 관리가 포함된 인증 서비스
│   ├── base-api.js            # 공통 기능이 포함된 기본 API 서비스
│   ├── reservation-api.js     # 예약 API 서비스
│   ├── room-api.js            # 객실 API 서비스
│   ├── room-group-api.js      # 객실그룹 API 서비스
│   ├── payment-method-api.js  # 결제수단 API 서비스
│   └── user-api.js            # 사용자 API 서비스
├── scenarios/
│   ├── reservation-browsing.js    # 예약 조회 시나리오
│   └── reservation-load-test.js   # 고강도 예약 부하 테스트
├── examples/
│   ├── simple-api-test.js     # 기본 API 테스트 예제
│   └── mixed-workload.js      # 혼합 현실적 워크로드 예제
└── README.md
```

## 설정

### 기본 URL 설정

기본 URL: `https://staging.rms.bellsoft.net/`

```bash
# 환경 변수로 설정
export BASE_URL=https://your-api-server.com

# 또는 k6에 직접 전달
k6 run -e BASE_URL=https://your-api-server.com scenario.js
```

### 테스트 사용자 인증정보

환경 변수로 테스트 인증정보 설정

```bash
export TEST_USERNAME=your-test-email@example.com
export TEST_PASSWORD=your-test-password
```

### 디버그 모드

상세 로깅 활성화

```bash
export DEBUG=true
```

## 사용 가능한 테스트 시나리오

### 1. 예약 조회 시나리오

**파일**: `scenarios/reservation-browsing.js`

현실적인 사용자 조회 행동 시뮬레이션
- VU당 1회 로그인
- 예약 목록 조회
- 0.5초 대기
- 예약 목록 재조회  
- 2초 대기
- 첫 번째 예약 상세 정보 조회
- 사이클 반복

**사용법**
```bash
k6 run scenarios/reservation-browsing.js
```

**부하 프로파일**
- 30초에 걸쳐 5명으로 증가
- 5명으로 2분간 유지
- 30초에 걸쳐 10명으로 증가
- 10명으로 2분간 유지
- 30초에 걸쳐 0명으로 감소

### 2. 예약 부하 테스트

**파일**: `scenarios/reservation-load-test.js`

최대 처리량 측정을 위한 고강도 부하 테스트
- 20명의 동시 사용자
- 1분간 연속 API 호출
- 요청 간 지연 없음
- 초당 최대 요청 수 측정

**사용법**
```bash
k6 run scenarios/reservation-load-test.js
```

### 3. 간단한 API 테스트

**파일**: `examples/simple-api-test.js`

모든 주요 API를 다루는 기본 기능 테스트
- 인증 테스트
- 모든 리소스 API 테스트 (예약, 객실, 객실그룹, 결제수단, 사용자)
- API 기능 검증을 위한 단일 반복

**사용법**
```bash
k6 run examples/simple-api-test.js
```

### 4. 혼합 워크로드 테스트

**파일**: `examples/mixed-workload.js`

다양한 사용자 행동 패턴을 가진 현실적인 혼합 워크로드
- 40% 예약 중심 사용자
- 30% 객실 관리 사용자  
- 20% 일반 조회 사용자
- 10% 상세 분석 사용자

**사용법**
```bash
k6 run examples/mixed-workload.js
```

## API 서비스

### 인증 기능

- **자동 로그인**: 각 가상 사용자는 시작 시 1회 로그인
- **토큰 관리**: 만료 시 자동 액세스 토큰 갱신
- **인증 오류 복구**: 401/403 오류 시 자동 재인증
- **토큰 만료 추적**: 만료 전 미리 토큰 갱신

### 사용 가능한 API 서비스

모든 API 서비스는 미리 설정되어 바로 사용 가능

```javascript
import { reservationApi } from '../services/reservation-api.js';
import { roomApi } from '../services/room-api.js';
import { roomGroupApi } from '../services/room-group-api.js';
import { paymentMethodApi } from '../services/payment-method-api.js';
import { userApi } from '../services/user-api.js';

// 예제
const reservations = reservationApi.getReservations();
const reservation = reservationApi.getReservationById(123);
const rooms = roomApi.getRooms();
const room = roomApi.getRoomById(456);
```

### CRUD 작업

각 서비스는 전체 CRUD 작업을 지원

```javascript
// GET 작업
const items = api.getItems(filters);
const item = api.getItemById(id);

// POST 작업  
const newItem = api.createItem(data);

// PUT 작업
const updatedItem = api.updateItem(id, data);

// DELETE 작업
const success = api.deleteItem(id);
```

## 테스트 실행

### 사전 요구사항

1. k6 설치: https://k6.io/docs/getting-started/installation/
2. 백엔드 API가 실행 중이고 접근 가능한지 확인
3. 테스트 인증정보 설정

### 기본 사용법

```bash
# 특정 시나리오 실행
k6 run scenarios/reservation-browsing.js

# 사용자 정의 기본 URL로 실행
k6 run -e BASE_URL=https://your-server.com scenarios/reservation-load-test.js

# 사용자 정의 인증정보로 실행
k6 run -e TEST_USERNAME=user@example.com -e TEST_PASSWORD=password scenarios/reservation-browsing.js

# 디버그 로깅으로 실행
k6 run -e DEBUG=true examples/simple-api-test.js

# 사용자 정의 VU 수와 실행 시간으로 실행
k6 run --vus 50 --duration 30s scenarios/reservation-load-test.js
```

### HTML 리포트 생성

```bash
# HTML 리포트 생성
k6 run --out html=report.html scenarios/reservation-browsing.js

# JSON 출력 생성
k6 run --out json=results.json scenarios/reservation-load-test.js
```

## 사용자 정의 시나리오 만들기

### 기본 템플릿

```javascript
import { config } from '../config/environment.js';
import { Logger } from '../utils/logger.js';
import { authService } from '../services/auth.js';
import { reservationApi } from '../services/reservation-api.js';

export const options = {
  scenarios: {
    my_test: {
      executor: 'constant-vus',
      vus: 10,
      duration: '30s',
    },
  },
  thresholds: config.THRESHOLDS,
};

export default function() {
  // VU당 1회 로그인
  if (!authService.accessToken) {
    const loginSuccess = authService.login();
    if (!loginSuccess) {
      return;
    }
  }
  
  // 여기에 테스트 로직 작성
  const reservations = reservationApi.getReservations();
  if (reservations) {
    Logger.info('테스트 성공');
  }
}
```

### 모범 사례

1. **인증**: API 호출 전 항상 `authService.accessToken` 확인
2. **오류 처리**: 일관된 오류 확인을 위해 제공된 검증 유틸리티 사용
3. **로깅**: 일관된 로그 형식을 위해 Logger 유틸리티 사용
4. **임계값**: 표준 성능 기준을 위해 `config.THRESHOLDS` 사용
5. **리소스 정리**: 복잡한 시나리오에서는 teardown 함수 구현

## 문제 해결

### 일반적인 문제

1. **로그인 실패**: 테스트 인증정보와 API 엔드포인트 확인
2. **토큰 만료**: 인증 서비스가 자동으로 처리
3. **API 타임아웃**: 필요시 `config.TIMEOUTS` 조정
4. **높은 오류율**: API 서버 용량과 네트워크 연결 확인

### 디버그 모드

상세한 API 호출 정보를 보려면 디버그 로깅 활성화

```bash
k6 run -e DEBUG=true your-scenario.js
```

다음 정보가 표시됨
- 상세한 요청/응답 로깅
- 토큰 갱신 이벤트
- API 호출 타이밍 정보
- 인증 상태 업데이트
