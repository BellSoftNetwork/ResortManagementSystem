---
id: frontend-web-dashboard
title: "frontend-web 대시보드"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: low
effort: small
---

# frontend-web 대시보드

> 메인 대시보드 페이지

---

## 1. 관련 파일

| 파일 | 역할 |
|------|------|
| `pages/DefaultHome.vue` | 대시보드 페이지 |
| `components/dashboard/ReservationSummary.vue` | 예약 요약 컴포넌트 |

---

## 2. 페이지 구성

### 2.1 DefaultHome.vue

```vue
<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">대시보드</div>
    <ReservationSummary />
  </q-page>
</template>

<script setup lang="ts">
import ReservationSummary from 'src/components/dashboard/ReservationSummary.vue'
</script>
```

---

## 3. 예약 요약 컴포넌트

### 3.1 표시 정보

| 항목 | 설명 |
|------|------|
| 오늘 체크인 | 오늘 입실 예정 예약 수 |
| 오늘 체크아웃 | 오늘 퇴실 예정 예약 수 |
| 이번 주 예약 | 이번 주 예약 수 |
| 이번 달 예약 | 이번 달 예약 수 |

### 3.2 구현

```vue
<template>
  <div class="row q-gutter-md">
    <q-card class="col">
      <q-card-section>
        <div class="text-h6">오늘 체크인</div>
        <div class="text-h3">{{ todayCheckIn }}</div>
      </q-card-section>
    </q-card>
    
    <q-card class="col">
      <q-card-section>
        <div class="text-h6">오늘 체크아웃</div>
        <div class="text-h3">{{ todayCheckOut }}</div>
      </q-card-section>
    </q-card>
    
    <q-card class="col">
      <q-card-section>
        <div class="text-h6">이번 주 예약</div>
        <div class="text-h3">{{ weeklyReservations }}</div>
      </q-card-section>
    </q-card>
    
    <q-card class="col">
      <q-card-section>
        <div class="text-h6">이번 달 예약</div>
        <div class="text-h3">{{ monthlyReservations }}</div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { reservationApi } from 'src/api/v1/reservation'

const todayCheckIn = ref(0)
const todayCheckOut = ref(0)
const weeklyReservations = ref(0)
const monthlyReservations = ref(0)

const loadSummary = async () => {
  const today = new Date().toISOString().split('T')[0]
  
  // 오늘 체크인
  const checkInResponse = await reservationApi.getAll({
    stayStartAt: today,
    stayEndAt: today,
    status: 'NORMAL'
  })
  todayCheckIn.value = checkInResponse.data.page.totalElements
  
  // 오늘 체크아웃
  const checkOutResponse = await reservationApi.getAll({
    stayEndAt: today,
    status: 'NORMAL'
  })
  todayCheckOut.value = checkOutResponse.data.page.totalElements
  
  // 이번 주/달 계산...
}

onMounted(loadSummary)
</script>
```

---

## 4. 권한

- 경로: `/`
- 권한: 모든 인증된 사용자
- NORMAL 사용자도 대시보드 조회 가능
