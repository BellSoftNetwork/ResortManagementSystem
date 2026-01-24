---
id: frontend-web-room-management
title: "frontend-web 객실 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: medium
effort: medium
---

# frontend-web 객실 관리

> 객실 CRUD 페이지 및 컴포넌트

---

## 1. 페이지

| 페이지 | 경로 | 권한 |
|--------|------|------|
| RoomStatus | /room-status | ADMIN |
| RoomList | /rooms | ADMIN |
| RoomDetail | /rooms/:id | ADMIN |
| RoomEdit | /rooms/:id/edit | ADMIN |
| RoomCreate | /rooms/create | ADMIN |

---

## 2. 컴포넌트

| 컴포넌트 | 역할 |
|----------|------|
| RoomListDynamicTable | 객실 목록 테이블 (필터, 페이지네이션) |
| RoomListTable | 간단한 객실 테이블 |
| RoomDetailCard | 객실 상세 정보 표시 |
| RoomEditCard | 객실 수정 폼 |
| RoomCreateCard | 객실 생성 폼 |
| RoomSelectDynamicTable | 객실 선택 테이블 (예약 생성 시) |
| RoomHistoryDynamicTable | 객실 변경 이력 테이블 |

---

## 3. 객실 현황 페이지

### 3.1 RoomStatus.vue

```vue
<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">객실 현황</div>
    
    <!-- 날짜 범위 선택 -->
    <div class="row q-gutter-sm q-mb-md">
      <q-input v-model="startDate" type="date" label="시작일" />
      <q-input v-model="endDate" type="date" label="종료일" />
    </div>
    
    <!-- 객실 그룹별 현황 -->
    <div v-for="group in roomGroups" :key="group.id">
      <div class="text-h6">{{ group.name }}</div>
      <div class="row q-gutter-sm">
        <q-card v-for="room in group.rooms" :key="room.id" class="room-card">
          <q-card-section :class="getRoomStatusClass(room)">
            <div>{{ room.number }}</div>
            <div v-if="room.currentReservation">
              {{ room.currentReservation.name }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>
```

---

## 4. 객실 목록 페이지

### 4.1 RoomList.vue

```vue
<template>
  <q-page padding>
    <div class="row justify-between items-center q-mb-md">
      <div class="text-h5">객실 관리</div>
      <q-btn color="primary" label="객실 추가" :to="{ name: 'RoomCreate' }" />
    </div>
    
    <RoomListDynamicTable
      @row-click="onRowClick"
    />
  </q-page>
</template>

<script setup lang="ts">
const router = useRouter()

const onRowClick = (room: Room) => {
  router.push({ name: 'RoomDetail', params: { id: room.id } })
}
</script>
```

---

## 5. 동적 테이블 컴포넌트

### 5.1 RoomListDynamicTable.vue

```vue
<template>
  <div>
    <!-- 필터 -->
    <div class="row q-gutter-sm q-mb-md">
      <q-select
        v-model="filter.roomGroupId"
        :options="roomGroups"
        option-value="id"
        option-label="name"
        label="객실 그룹"
        clearable
        emit-value
        map-options
      />
      <q-select
        v-model="filter.status"
        :options="statusOptions"
        label="상태"
        clearable
        emit-value
        map-options
      />
      <q-input v-model="filter.search" label="검색" debounce="300" />
    </div>
    
    <!-- 테이블 -->
    <q-table
      :rows="rooms"
      :columns="columns"
      :loading="loading"
      :pagination="pagination"
      @request="onRequest"
      row-key="id"
      @row-click="(evt, row) => $emit('row-click', row)"
    >
      <template #body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="getStatusColor(props.row.status)">
            {{ getStatusLabel(props.row.status) }}
          </q-badge>
        </q-td>
      </template>
      
      <template #body-cell-actions="props">
        <q-td :props="props">
          <q-btn flat icon="edit" :to="{ name: 'RoomEdit', params: { id: props.row.id } }" />
          <q-btn flat icon="delete" color="negative" @click="onDelete(props.row)" />
        </q-td>
      </template>
    </q-table>
  </div>
</template>

<script setup lang="ts">
import { roomApi } from 'src/api/v1/room'
import { roomGroupApi } from 'src/api/v1/room-group'

const rooms = ref<Room[]>([])
const roomGroups = ref<RoomGroup[]>([])
const loading = ref(false)

const filter = reactive({
  roomGroupId: null as number | null,
  status: null as string | null,
  search: ''
})

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
})

const columns = [
  { name: 'number', label: '객실 번호', field: 'number', sortable: true },
  { name: 'roomGroup', label: '객실 그룹', field: (row: Room) => row.roomGroup?.name },
  { name: 'status', label: '상태', field: 'status' },
  { name: 'note', label: '메모', field: 'note' },
  { name: 'actions', label: '작업', field: 'actions' }
]

const statusOptions = [
  { value: 'NORMAL', label: '정상' },
  { value: 'INACTIVE', label: '비활성' },
  { value: 'CONSTRUCTION', label: '공사중' },
  { value: 'DAMAGED', label: '파손' }
]

const loadRooms = async () => {
  loading.value = true
  try {
    const response = await roomApi.getAll({
      page: pagination.value.page - 1,
      size: pagination.value.rowsPerPage,
      ...filter
    })
    rooms.value = response.data.data
    pagination.value.rowsNumber = response.data.page.totalElements
  } finally {
    loading.value = false
  }
}

const onRequest = (props: { pagination: typeof pagination.value }) => {
  pagination.value = props.pagination
  loadRooms()
}

watch(filter, loadRooms, { deep: true })
onMounted(loadRooms)
</script>
```

---

## 6. 상세/수정/생성 페이지

### 6.1 RoomDetail.vue

```vue
<template>
  <q-page padding>
    <RoomDetailCard :room="room" />
    
    <div class="text-h6 q-mt-lg">변경 이력</div>
    <RoomHistoryDynamicTable :room-id="roomId" />
  </q-page>
</template>
```

### 6.2 RoomEdit.vue

```vue
<template>
  <q-page padding>
    <RoomEditCard :room="room" @success="onSuccess" />
  </q-page>
</template>
```

### 6.3 RoomCreate.vue

```vue
<template>
  <q-page padding>
    <RoomCreateCard @success="onSuccess" />
  </q-page>
</template>
```

---

## 7. 상태 표시

| 상태 | 색상 | 라벨 |
|------|------|------|
| NORMAL | green | 정상 |
| INACTIVE | grey | 비활성 |
| CONSTRUCTION | orange | 공사중 |
| DAMAGED | red | 파손 |
