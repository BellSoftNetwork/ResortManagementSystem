---
id: frontend-web-authentication
title: "frontend-web 인증"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: high
effort: medium
---

# frontend-web 인증

> JWT 기반 로그인, 회원가입, 토큰 관리

---

## 1. 관련 파일

| 파일 | 역할 |
|------|------|
| `pages/auth/MainLogin.vue` | 로그인 페이지 |
| `pages/auth/MainRegister.vue` | 회원가입 페이지 |
| `components/auth/LoginCard.vue` | 로그인 폼 |
| `components/auth/RegisterCard.vue` | 회원가입 폼 |
| `components/auth/LogoutDialog.vue` | 로그아웃 확인 |
| `stores/auth.ts` | 인증 상태 관리 |
| `api/v1/auth.ts` | 인증 API 함수 |
| `api/services/AuthInterceptorService.ts` | JWT 인터셉터 |

---

## 2. 페이지

### 2.1 로그인 (/login)

```vue
<template>
  <q-page class="flex flex-center">
    <LoginCard @success="onLoginSuccess" />
  </q-page>
</template>

<script setup lang="ts">
const router = useRouter()
const authStore = useAuthStore()

const onLoginSuccess = () => {
  router.push('/')
}
</script>
```

### 2.2 회원가입 (/register)

- 앱 설정에서 회원가입 가능 여부 확인
- 불가능 시 로그인 페이지로 리다이렉트

---

## 3. Auth Store

### 3.1 상태

```typescript
interface AuthState {
  status: {
    isFirstRequest: boolean
    isRefreshingToken: boolean
  }
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  accessTokenExpiresIn: number | null
  tokenRefreshTimer: number | null
  refreshAttempts: number
  lastRefreshAttempt: number
}
```

### 3.2 Getters

```typescript
getters: {
  isLoggedIn: (state) => !!state.accessToken,
  isNormalRole: (state) => state.user?.role === 'NORMAL',
  isAdminRole: (state) => ['ADMIN', 'SUPER_ADMIN'].includes(state.user?.role ?? ''),
  isSuperAdminRole: (state) => state.user?.role === 'SUPER_ADMIN'
}
```

### 3.3 Actions

```typescript
actions: {
  async login(credentials: LoginRequest) {
    const response = await authApi.login(credentials)
    this.accessToken = response.data.accessToken
    this.refreshToken = response.data.refreshToken
    this.accessTokenExpiresIn = response.data.accessTokenExpiresIn
    this.startTokenRefreshTimer()
    await this.loadAccountInfo()
  },
  
  async logout() {
    await authApi.logout()
    this.stopTokenRefreshTimer()
    this.$reset()
  },
  
  async refreshAccessToken() {
    if (!this.refreshToken) return
    
    const response = await authApi.refresh({ refreshToken: this.refreshToken })
    this.accessToken = response.data.accessToken
    this.refreshToken = response.data.refreshToken
    this.startTokenRefreshTimer()
  },
  
  async loadAccountInfo() {
    const response = await mainApi.getMy()
    this.user = response.data
  },
  
  startTokenRefreshTimer() {
    // 만료 5분 전에 갱신
    const refreshTime = (this.accessTokenExpiresIn! - 300) * 1000
    this.tokenRefreshTimer = setTimeout(() => {
      this.refreshAccessToken()
    }, refreshTime)
  },
  
  hydrate() {
    // 페이지 새로고침 시 localStorage에서 복원
    const refreshToken = localStorage.getItem('refreshToken')
    if (refreshToken) {
      this.refreshToken = refreshToken
      this.refreshAccessToken()
    }
  }
}
```

---

## 4. API 함수

```typescript
// api/v1/auth.ts
export const authApi = {
  register: (data: RegisterRequest) =>
    api.post<SingleResponse<User>>('/api/v1/auth/register', data),
    
  login: (data: LoginRequest) =>
    api.post<SingleResponse<TokenResponse>>('/api/v1/auth/login', data),
    
  logout: () =>
    api.post('/api/v1/auth/logout'),
    
  refresh: (data: RefreshTokenRequest) =>
    api.post<SingleResponse<TokenResponse>>('/api/v1/auth/refresh', data)
}
```

---

## 5. Axios 인터셉터

### 5.1 요청 인터셉터

```typescript
api.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  
  // refresh 요청에는 토큰 추가하지 않음
  if (config.url?.includes('/auth/refresh')) {
    return config
  }
  
  if (authStore.accessToken) {
    config.headers.Authorization = `Bearer ${authStore.accessToken}`
  }
  
  return config
})
```

### 5.2 응답 인터셉터

```typescript
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const authStore = useAuthStore()
    
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true
      
      try {
        await authStore.refreshAccessToken()
        return api.request(error.config)
      } catch {
        authStore.logout()
        router.push('/login')
      }
    }
    
    return Promise.reject(error)
  }
)
```

---

## 6. 라우트 가드

### 6.1 authenticate-guard

```typescript
export const authenticateGuard: NavigationGuard = async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 공개 페이지는 통과
  if (to.meta.public) {
    return next()
  }
  
  // 첫 요청 시 hydrate
  if (authStore.status.isFirstRequest) {
    authStore.hydrate()
    authStore.status.isFirstRequest = false
  }
  
  // 인증 확인
  if (!authStore.isLoggedIn) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }
  
  next()
}
```

### 6.2 role-guard

```typescript
export const roleGuard = (allowedRoles: string[]): NavigationGuard => {
  return (to, from, next) => {
    const authStore = useAuthStore()
    
    if (!authStore.user || !allowedRoles.includes(authStore.user.role)) {
      return next('/error/403')
    }
    
    next()
  }
}
```

---

## 7. 컴포넌트

### 7.1 LoginCard

```vue
<template>
  <q-card class="login-card">
    <q-card-section>
      <q-form @submit="onSubmit">
        <q-input v-model="form.username" label="아이디" />
        <q-input v-model="form.password" label="비밀번호" type="password" />
        <q-btn type="submit" label="로그인" :loading="loading" />
      </q-form>
    </q-card-section>
  </q-card>
</template>
```

### 7.2 RegisterCard

```vue
<template>
  <q-card class="register-card">
    <q-card-section>
      <q-form @submit="onSubmit">
        <q-input v-model="form.userId" label="아이디" />
        <q-input v-model="form.email" label="이메일" />
        <q-input v-model="form.name" label="이름" />
        <q-input v-model="form.password" label="비밀번호" type="password" />
        <q-input v-model="form.passwordConfirm" label="비밀번호 확인" type="password" />
        <q-btn type="submit" label="회원가입" :loading="loading" />
      </q-form>
    </q-card-section>
  </q-card>
</template>
```

---

## 8. 토큰 저장

| 토큰 | 저장 위치 | 이유 |
|------|-----------|------|
| Access Token | 메모리 (Pinia) | XSS 방지 |
| Refresh Token | localStorage | 페이지 새로고침 시 유지 |
