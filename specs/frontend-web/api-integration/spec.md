---
id: frontend-web-api-integration
title: "frontend-web API 통신"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: medium
effort: small
---

# frontend-web API 통신

> Axios 기반 API 통신 레이어

---

## 1. 파일 구조

```
src/api/
├── v1/                     # API 엔드포인트
│   ├── auth.ts             # 인증
│   ├── main.ts             # 프로필, 설정
│   ├── room.ts             # 객실
│   ├── room-group.ts       # 객실 그룹
│   ├── reservation.ts      # 예약
│   ├── payment-method.ts   # 결제 수단
│   └── admin/
│       └── account.ts      # 관리자 계정
├── services/
│   ├── AuthInterceptorService.ts  # JWT 인터셉터
│   ├── RetryService.ts            # 재시도 로직
│   ├── NotificationService.ts     # 알림
│   └── NetworkStatusService.ts    # 네트워크 상태
└── constants.ts
```

---

## 2. Axios 설정

### 2.1 인스턴스 생성

```typescript
// boot/axios.ts
import axios from 'axios'

const api = axios.create({
  baseURL: '/',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export { api }
```

### 2.2 인터셉터

```typescript
// 요청 인터셉터 - JWT 토큰 추가
api.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.accessToken && !config.url?.includes('/auth/refresh')) {
    config.headers.Authorization = `Bearer ${authStore.accessToken}`
  }
  return config
})

// 응답 인터셉터 - 401 처리
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true
      await authStore.refreshAccessToken()
      return api.request(error.config)
    }
    return Promise.reject(error)
  }
)
```

---

## 3. API 함수 패턴

### 3.1 CRUD 패턴

```typescript
// api/v1/room.ts
export const roomApi = {
  // 목록 (페이지네이션 + 필터)
  getAll: (params?: RoomFilter & PaginationParams) =>
    api.get<ListResponse<Room>>('/api/v1/rooms', { params }),
  
  // 단일 조회
  getById: (id: number) =>
    api.get<SingleResponse<Room>>(`/api/v1/rooms/${id}`),
  
  // 생성
  create: (data: CreateRoomRequest) =>
    api.post<SingleResponse<Room>>('/api/v1/rooms', data),
  
  // 수정
  update: (id: number, data: UpdateRoomRequest) =>
    api.patch<SingleResponse<Room>>(`/api/v1/rooms/${id}`, data),
  
  // 삭제
  delete: (id: number) =>
    api.delete(`/api/v1/rooms/${id}`),
  
  // 히스토리
  getHistories: (id: number, params?: PaginationParams) =>
    api.get<ListResponse<Revision<Room>>>(`/api/v1/rooms/${id}/histories`, { params })
}
```

### 3.2 응답 타입

```typescript
// schema/response.ts
interface SingleResponse<T> {
  data: T
}

interface ListResponse<T> {
  data: T[]
  page: PageInfo
}

interface PageInfo {
  size: number
  number: number
  totalElements: number
  totalPages: number
}
```

---

## 4. API 목록

### 4.1 인증 (auth.ts)

| 함수 | Method | Path |
|------|--------|------|
| register | POST | /api/v1/auth/register |
| login | POST | /api/v1/auth/login |
| logout | POST | /api/v1/auth/logout |
| refresh | POST | /api/v1/auth/refresh |

### 4.2 메인 (main.ts)

| 함수 | Method | Path |
|------|--------|------|
| getMy | GET | /api/v1/my |
| updateMy | PATCH | /api/v1/my |
| getConfig | GET | /api/v1/config |
| getEnv | GET | /api/v1/env |

### 4.3 객실 (room.ts)

| 함수 | Method | Path |
|------|--------|------|
| getAll | GET | /api/v1/rooms |
| getById | GET | /api/v1/rooms/:id |
| create | POST | /api/v1/rooms |
| update | PATCH | /api/v1/rooms/:id |
| delete | DELETE | /api/v1/rooms/:id |
| getHistories | GET | /api/v1/rooms/:id/histories |

### 4.4 객실 그룹 (room-group.ts)

| 함수 | Method | Path |
|------|--------|------|
| getAll | GET | /api/v1/room-groups |
| getById | GET | /api/v1/room-groups/:id |
| create | POST | /api/v1/room-groups |
| update | PATCH | /api/v1/room-groups/:id |
| delete | DELETE | /api/v1/room-groups/:id |

### 4.5 예약 (reservation.ts)

| 함수 | Method | Path |
|------|--------|------|
| getAll | GET | /api/v1/reservations |
| getById | GET | /api/v1/reservations/:id |
| create | POST | /api/v1/reservations |
| update | PATCH | /api/v1/reservations/:id |
| delete | DELETE | /api/v1/reservations/:id |
| getHistories | GET | /api/v1/reservations/:id/histories |
| getStatistics | GET | /api/v1/reservation-statistics |

### 4.6 결제 수단 (payment-method.ts)

| 함수 | Method | Path |
|------|--------|------|
| getAll | GET | /api/v1/payment-methods |
| getById | GET | /api/v1/payment-methods/:id |
| create | POST | /api/v1/payment-methods |
| update | PATCH | /api/v1/payment-methods/:id |
| delete | DELETE | /api/v1/payment-methods/:id |

### 4.7 관리자 계정 (admin/account.ts)

| 함수 | Method | Path |
|------|--------|------|
| getAll | GET | /api/v1/admin/accounts |
| create | POST | /api/v1/admin/accounts |
| update | PATCH | /api/v1/admin/accounts/:id |

---

## 5. 서비스

### 5.1 RetryService

```typescript
// 네트워크 실패 시 자동 재시도
const MAX_RETRIES = 3
const RETRY_DELAY = 1000

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const retries = error.config._retries || 0
    
    if (isNetworkError(error) && retries < MAX_RETRIES) {
      error.config._retries = retries + 1
      await delay(RETRY_DELAY * (retries + 1))
      return api.request(error.config)
    }
    
    return Promise.reject(error)
  }
)
```

### 5.2 NotificationService

```typescript
// 사용자 알림
export const notifyError = (message: string) => {
  Notify.create({
    type: 'negative',
    message,
    position: 'top'
  })
}

export const notifySuccess = (message: string) => {
  Notify.create({
    type: 'positive',
    message,
    position: 'top'
  })
}
```

### 5.3 NetworkStatusService

```typescript
// 네트워크 연결 상태 모니터링
export const useNetworkStatus = () => {
  const isOnline = ref(navigator.onLine)
  
  onMounted(() => {
    window.addEventListener('online', () => isOnline.value = true)
    window.addEventListener('offline', () => isOnline.value = false)
  })
  
  return { isOnline }
}
```

---

## 6. 컴포넌트에서 사용

```vue
<script setup lang="ts">
import { roomApi } from 'src/api/v1/room'
import type { Room } from 'src/schema/room'

const rooms = ref<Room[]>([])
const loading = ref(false)

const loadRooms = async () => {
  loading.value = true
  try {
    const response = await roomApi.getAll({ page: 0, size: 20 })
    rooms.value = response.data.data
  } catch (error) {
    notifyError('객실 목록을 불러올 수 없습니다')
  } finally {
    loading.value = false
  }
}

onMounted(loadRooms)
</script>
```
