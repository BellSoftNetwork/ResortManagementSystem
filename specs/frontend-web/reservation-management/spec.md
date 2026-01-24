---
id: frontend-web-reservation-management
title: "frontend-web 예약 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: high
effort: medium
---

# frontend-web 예약 관리

> 예약 CRUD 페이지 및 컴포넌트

---

## 1. 페이지

| 페이지 | 경로 | 권한 | 설명 |
|--------|------|------|------|
| ReservationList | /reservations | ADMIN | 숙박 예약 목록 |
| ReservationDetail | /reservations/:id | ADMIN | 예약 상세 |
| ReservationEdit | /reservations/:id/edit | ADMIN | 예약 수정 |
| ReservationCreate | /reservations/create | ADMIN | 예약 생성 |

---

## 2. 컴포넌트

| 컴포넌트 | 역할 |
|----------|------|
| ReservationListDynamicTable | 예약 목록 테이블 (필터, 페이지네이션) |
| ReservationEditor | 예약 생성/수정 폼 |
| ReservationHistoryDynamicTable | 예약 변경 이력 테이블 |

---

## 3. 예약 목록 페이지

### 3.1 ReservationList.vue

```vue
<template>
  <q-page padding>
    <div class="row justify-between items-center q-mb-md">
      <div class="text-h5">예약 관리</div>
      <q-btn color="primary" label="예약 추가" :to="{ name: 'ReservationCreate' }" />
    </div>
    
    <ReservationListDynamicTable
      :type="'STAY'"
      @row-click="onRowClick"
    />
  </q-page>
</template>
```

---

## 4. 동적 테이블

### 4.1 ReservationListDynamicTable.vue

```vue
<template>
  <div>
    <!-- 필터 -->
    <div class="row q-gutter-sm q-mb-md">
      <q-input v-model="filter.stayStartAt" type="date" label="입실일 시작" />
      <q-input v-model="filter.stayEndAt" type="date" label="입실일 종료" />
      <q-select
        v-model="filter.status"
        :options="statusOptions"
        label="상태"
        clearable
        emit-value
        map-options
      />
      <q-select
        v-model="filter.roomId"
        :options="rooms"
        option-value="id"
        option-label="number"
        label="객실"
        clearable
        emit-value
        map-options
      />
      <q-input v-model="filter.search" label="검색 (이름, 전화번호)" debounce="300" />
    </div>
    
    <!-- 테이블 -->
    <q-table
      :rows="reservations"
      :columns="columns"
      :loading="loading"
      :pagination="pagination"
      @request="onRequest"
      row-key="id"
    >
      <template #body-cell-rooms="props">
        <q-td :props="props">
          {{ props.row.rooms.map(r => r.number).join(', ') }}
        </q-td>
      </template>
      
      <template #body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="getStatusColor(props.row.status)">
            {{ getStatusLabel(props.row.status) }}
          </q-badge>
        </q-td>
      </template>
      
      <template #body-cell-price="props">
        <q-td :props="props">
          {{ formatCurrency(props.row.price) }}
        </q-td>
      </template>
    </q-table>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  type?: 'STAY' | 'MONTHLY_RENT'
}>()

const filter = reactive({
  stayStartAt: null as string | null,
  stayEndAt: null as string | null,
  status: null as string | null,
  roomId: null as number | null,
  search: '',
  type: props.type
})

const columns = [
  { name: 'name', label: '예약자', field: 'name', sortable: true },
  { name: 'phone', label: '연락처', field: 'phone' },
  { name: 'rooms', label: '객실', field: 'rooms' },
  { name: 'stayStartAt', label: '입실일', field: 'stayStartAt', sortable: true },
  { name: 'stayEndAt', label: '퇴실일', field: 'stayEndAt' },
  { name: 'peopleCount', label: '인원', field: 'peopleCount' },
  { name: 'price', label: '가격', field: 'price' },
  { name: 'status', label: '상태', field: 'status' },
  { name: 'actions', label: '작업', field: 'actions' }
]

const statusOptions = [
  { value: 'NORMAL', label: '정상' },
  { value: 'PENDING', label: '대기' },
  { value: 'CANCEL', label: '취소' },
  { value: 'REFUND', label: '환불' }
]
</script>
```

---

## 5. 예약 에디터

### 5.1 ReservationEditor.vue

```vue
<template>
  <q-card>
    <q-card-section>
      <q-form @submit="onSubmit">
        <!-- 예약자 정보 -->
        <div class="text-h6">예약자 정보</div>
        <q-input v-model="form.name" label="이름" :rules="[required]" />
        <q-input v-model="form.phone" label="전화번호" :rules="[required]" />
        <q-input v-model.number="form.peopleCount" label="인원" type="number" />
        
        <!-- 숙박 정보 -->
        <div class="text-h6 q-mt-md">숙박 정보</div>
        <div class="row q-gutter-sm">
          <q-input v-model="form.stayStartAt" type="date" label="입실일" :rules="[required]" />
          <q-input v-model="form.stayEndAt" type="date" label="퇴실일" :rules="[required]" />
        </div>
        
        <!-- 객실 선택 -->
        <div class="text-h6 q-mt-md">객실 선택</div>
        <RoomSelectDynamicTable v-model="form.roomIds" :stay-start-at="form.stayStartAt" :stay-end-at="form.stayEndAt" />
        
        <!-- 결제 정보 -->
        <div class="text-h6 q-mt-md">결제 정보</div>
        <q-select
          v-model="form.paymentMethodId"
          :options="paymentMethods"
          option-value="id"
          option-label="name"
          label="결제 수단"
          emit-value
          map-options
          :rules="[required]"
        />
        <q-input v-model.number="form.price" label="총 가격" type="number" />
        <q-input v-model.number="form.deposit" label="보증금" type="number" />
        <q-input v-model.number="form.paymentAmount" label="결제액" type="number" />
        
        <!-- 메모 -->
        <q-input v-model="form.note" label="메모" type="textarea" />
        
        <div class="q-mt-md">
          <q-btn type="submit" color="primary" :label="isEdit ? '수정' : '생성'" :loading="loading" />
          <q-btn flat label="취소" @click="$router.back()" />
        </div>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
const props = defineProps<{
  reservation?: Reservation
}>()

const isEdit = computed(() => !!props.reservation)

const form = reactive({
  name: props.reservation?.name ?? '',
  phone: props.reservation?.phone ?? '',
  peopleCount: props.reservation?.peopleCount ?? 1,
  stayStartAt: props.reservation?.stayStartAt ?? '',
  stayEndAt: props.reservation?.stayEndAt ?? '',
  roomIds: props.reservation?.rooms.map(r => r.id) ?? [],
  paymentMethodId: props.reservation?.paymentMethodId ?? null,
  price: props.reservation?.price ?? 0,
  deposit: props.reservation?.deposit ?? 0,
  paymentAmount: props.reservation?.paymentAmount ?? 0,
  note: props.reservation?.note ?? '',
  type: props.reservation?.type ?? 'STAY'
})

const onSubmit = async () => {
  loading.value = true
  try {
    if (isEdit.value) {
      await reservationApi.update(props.reservation!.id, form)
    } else {
      await reservationApi.create(form)
    }
    emit('success')
  } finally {
    loading.value = false
  }
}
</script>
```

---

## 6. 상태 표시

| 상태 | 색상 | 라벨 |
|------|------|------|
| NORMAL | green | 정상 |
| PENDING | orange | 대기 |
| CANCEL | grey | 취소 |
| REFUND | red | 환불 |

---

## 7. 가격 포맷

```typescript
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('ko-KR', {
    style: 'currency',
    currency: 'KRW'
  }).format(value)
}
```
