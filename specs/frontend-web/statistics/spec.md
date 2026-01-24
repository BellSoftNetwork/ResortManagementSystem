---
id: frontend-web-statistics
title: "frontend-web 통계"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: low
effort: medium
---

# frontend-web 통계

> 예약 통계 대시보드

---

## 1. 페이지

| 페이지 | 경로 | 권한 |
|--------|------|------|
| StatsPage | /stats | ADMIN |

---

## 2. 컴포넌트

| 컴포넌트 | 역할 |
|----------|------|
| MonthSelector | 월 선택기 |
| YearlyReservationChart | 연간 예약 차트 |
| YearlyRevenueChart | 연간 매출 차트 |
| YearlyReservationCountChart | 연간 예약 건수 차트 |
| YearlyPeopleCountChart | 연간 인원 차트 |
| YearlyRoomCountChart | 연간 객실 사용 차트 |
| MonthlyGuestsCard | 월별 투숙객 카드 |
| MonthlySalesCard | 월별 매출 카드 |
| MonthlyReservationCountCard | 월별 예약 건수 카드 |
| RoomGroupOccupancyTable | 객실 그룹별 점유율 테이블 |
| RoomAllocationTable | 객실 배정 테이블 |

---

## 3. 통계 페이지 구성

### 3.1 StatsPage.vue

```vue
<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">통계</div>
    
    <!-- 월 선택 -->
    <MonthSelector v-model="selectedMonth" />
    
    <!-- 연간 차트 -->
    <div class="row q-gutter-md q-mt-md">
      <div class="col-12 col-md-6">
        <YearlyReservationChart :year="selectedYear" />
      </div>
      <div class="col-12 col-md-6">
        <YearlyRevenueChart :year="selectedYear" />
      </div>
    </div>
    
    <!-- 월별 요약 카드 -->
    <div class="row q-gutter-md q-mt-md">
      <MonthlySalesCard :month="selectedMonth" class="col" />
      <MonthlyGuestsCard :month="selectedMonth" class="col" />
      <MonthlyReservationCountCard :month="selectedMonth" class="col" />
    </div>
    
    <!-- 상세 테이블 -->
    <div class="row q-gutter-md q-mt-md">
      <div class="col-12 col-md-6">
        <RoomGroupOccupancyTable :month="selectedMonth" />
      </div>
      <div class="col-12 col-md-6">
        <RoomAllocationTable :month="selectedMonth" />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
const selectedMonth = ref(new Date().toISOString().slice(0, 7)) // YYYY-MM
const selectedYear = computed(() => selectedMonth.value.slice(0, 4))
</script>
```

---

## 4. 차트 컴포넌트

### 4.1 YearlyReservationChart.vue

```vue
<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">연간 예약 현황</div>
      <apexchart
        type="bar"
        :options="chartOptions"
        :series="series"
        height="300"
      />
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { reservationApi } from 'src/api/v1/reservation'

const props = defineProps<{ year: string }>()

const stats = ref<StatisticsData[]>([])

const chartOptions = computed(() => ({
  chart: { type: 'bar' },
  xaxis: {
    categories: stats.value.map(s => s.period)
  },
  colors: ['#1976D2']
}))

const series = computed(() => [{
  name: '예약 건수',
  data: stats.value.map(s => s.totalReservations)
}])

const loadStats = async () => {
  const response = await reservationApi.getStatistics({
    startDate: `${props.year}-01-01`,
    endDate: `${props.year}-12-31`,
    periodType: 'MONTHLY'
  })
  stats.value = response.data.stats
}

watch(() => props.year, loadStats, { immediate: true })
</script>
```

### 4.2 YearlyRevenueChart.vue

```vue
<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">연간 매출 현황</div>
      <apexchart
        type="line"
        :options="chartOptions"
        :series="series"
        height="300"
      />
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
const series = computed(() => [{
  name: '매출',
  data: stats.value.map(s => s.totalSales)
}])
</script>
```

---

## 5. 요약 카드 컴포넌트

### 5.1 MonthlySalesCard.vue

```vue
<template>
  <q-card>
    <q-card-section>
      <div class="text-subtitle2">{{ month }} 매출</div>
      <div class="text-h4">{{ formatCurrency(totalSales) }}</div>
      <div class="text-caption" :class="changeClass">
        {{ changeText }}
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
const props = defineProps<{ month: string }>()

const totalSales = ref(0)
const previousSales = ref(0)

const change = computed(() => {
  if (previousSales.value === 0) return 0
  return ((totalSales.value - previousSales.value) / previousSales.value) * 100
})

const changeClass = computed(() => change.value >= 0 ? 'text-positive' : 'text-negative')
const changeText = computed(() => `전월 대비 ${change.value >= 0 ? '+' : ''}${change.value.toFixed(1)}%`)
</script>
```

---

## 6. ApexCharts 설정

```typescript
// boot/apexcharts.ts
import VueApexCharts from 'vue3-apexcharts'

export default boot(({ app }) => {
  app.use(VueApexCharts)
  app.component('apexchart', VueApexCharts)
})
```

---

## 7. 통계 데이터 타입

```typescript
interface StatisticsData {
  period: string          // "2026-01"
  totalSales: number      // 총 매출
  totalReservations: number // 총 예약 건수
  totalGuests: number     // 총 투숙객
}
```
