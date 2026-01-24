---
id: frontend-web-state-management
title: "frontend-web 상태 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: medium
effort: small
---

# frontend-web 상태 관리

> Pinia 기반 전역 상태 관리

---

## 1. 스토어 목록

| 스토어 | 파일 | 역할 |
|--------|------|------|
| useAuthStore | `stores/auth.ts` | 인증, 토큰, 사용자 정보 |
| useMenuStore | `stores/menu.ts` | 역할 기반 메뉴 |
| useAppConfigStore | `stores/app-config.ts` | 앱 설정 |

---

## 2. Auth Store

### 2.1 상태

```typescript
state: () => ({
  status: {
    isFirstRequest: true,
    isRefreshingToken: false
  },
  user: null as User | null,
  accessToken: null as string | null,
  refreshToken: null as string | null,
  accessTokenExpiresIn: null as number | null,
  tokenRefreshTimer: null as number | null,
  refreshAttempts: 0,
  lastRefreshAttempt: 0
})
```

### 2.2 Getters

```typescript
getters: {
  isFirstRequest: (state) => state.status.isFirstRequest,
  isRefreshingToken: (state) => state.status.isRefreshingToken,
  isLoggedIn: (state) => !!state.accessToken,
  isNormalRole: (state) => state.user?.role === 'NORMAL',
  isAdminRole: (state) => ['ADMIN', 'SUPER_ADMIN'].includes(state.user?.role ?? ''),
  isSuperAdminRole: (state) => state.user?.role === 'SUPER_ADMIN'
}
```

### 2.3 Actions

| Action | 역할 |
|--------|------|
| `login(credentials)` | 로그인, 토큰 저장, 타이머 시작 |
| `logout()` | 로그아웃, 상태 초기화 |
| `refreshAccessToken()` | 토큰 갱신 |
| `loadAccountInfo()` | 사용자 정보 로드 |
| `startTokenRefreshTimer()` | 자동 갱신 타이머 |
| `stopTokenRefreshTimer()` | 타이머 정지 |
| `hydrate()` | 페이지 새로고침 시 복원 |

---

## 3. Menu Store

### 3.1 역할 기반 메뉴

```typescript
getters: {
  allLinks: () => {
    const authStore = useAuthStore()
    const links: MenuLink[] = [
      { title: '대시보드', icon: 'dashboard', to: '/', gnb: true }
    ]
    
    if (authStore.isAdminRole) {
      links.push(
        { title: '예약 관리', icon: 'event', to: '/reservations', gnb: true },
        { title: '월세 관리', icon: 'apartment', to: '/monthly-rents', gnb: true },
        { title: '객실 현황', icon: 'bed', to: '/room-status', gnb: true },
        { title: '통계', icon: 'bar_chart', to: '/stats', gnb: true },
        { title: '객실 관리', icon: 'meeting_room', to: '/rooms', gnb: false },
        { title: '객실 그룹', icon: 'category', to: '/room-groups', gnb: false },
        { title: '결제 수단', icon: 'payment', to: '/payment-methods', gnb: false },
        { title: '계정 관리', icon: 'manage_accounts', to: '/admin/accounts', gnb: false }
      )
    }
    
    if (authStore.isSuperAdminRole) {
      links.push(
        { title: '개발 테스트', icon: 'science', to: '/admin/dev-test', gnb: false }
      )
    }
    
    return links
  },
  
  tabLinks: (state) => state.allLinks.filter(link => link.gnb)
}
```

### 3.2 MenuLink 타입

```typescript
interface MenuLink {
  title: string
  icon: string
  to: string
  gnb: boolean  // 상단 탭에 표시 여부
}
```

---

## 4. App Config Store

### 4.1 상태

```typescript
state: () => ({
  status: {
    isLoading: false,
    isLoaded: false
  },
  config: {
    isAvailableRegistration: false
  }
})
```

### 4.2 Actions

```typescript
actions: {
  async loadAppConfig() {
    if (this.status.isLoaded) return
    
    this.status.isLoading = true
    try {
      const response = await mainApi.getConfig()
      this.config = response.data
      this.status.isLoaded = true
    } finally {
      this.status.isLoading = false
    }
  }
}
```

### 4.3 사용처

- 회원가입 페이지에서 가입 가능 여부 확인
- 가입 불가 시 로그인 페이지로 리다이렉트

---

## 5. Store 초기화

### 5.1 Pinia 설정

```typescript
// stores/index.ts
import { createPinia } from 'pinia'

export const pinia = createPinia()

// 플러그인: auth store hydration
pinia.use(({ store }) => {
  if (store.$id === 'auth') {
    store.hydrate()
  }
})
```

### 5.2 Quasar 부트

```typescript
// boot/pinia.ts
import { boot } from 'quasar/wrappers'
import { pinia } from 'src/stores'

export default boot(({ app }) => {
  app.use(pinia)
})
```

---

## 6. 반응형 사용

```vue
<script setup lang="ts">
import { useAuthStore } from 'src/stores/auth'
import { storeToRefs } from 'pinia'

const authStore = useAuthStore()
const { user, isLoggedIn, isAdminRole } = storeToRefs(authStore)

// 액션 호출
const handleLogout = () => {
  authStore.logout()
}
</script>

<template>
  <div v-if="isLoggedIn">
    <span>{{ user?.name }}</span>
    <q-btn v-if="isAdminRole" label="관리자 메뉴" />
  </div>
</template>
```
