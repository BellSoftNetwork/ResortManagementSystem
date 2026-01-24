---
id: frontend-web-admin-management
title: "frontend-web 관리자 기능"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: frontend
risk: medium
effort: small
---

# frontend-web 관리자 기능

> 계정 관리 및 개발 도구

---

## 1. 페이지

| 페이지 | 경로 | 권한 |
|--------|------|------|
| AccountList | /admin/accounts | ADMIN |
| DevTest | /admin/dev-test | SUPER_ADMIN |
| DevDebug | /debug | SUPER_ADMIN |

---

## 2. 계정 관리

### 2.1 컴포넌트

| 컴포넌트 | 역할 |
|----------|------|
| AccountListDynamicTable | 계정 목록 테이블 |
| AccountCreateDialog | 계정 생성 다이얼로그 |
| AccountEditDialog | 계정 수정 다이얼로그 |

### 2.2 AccountList.vue

```vue
<template>
  <q-page padding>
    <div class="row justify-between items-center q-mb-md">
      <div class="text-h5">계정 관리</div>
      <q-btn color="primary" label="계정 추가" @click="showCreateDialog = true" />
    </div>
    
    <AccountListDynamicTable @edit="onEdit" />
    
    <AccountCreateDialog v-model="showCreateDialog" @success="onCreateSuccess" />
    <AccountEditDialog v-model="showEditDialog" :user="selectedUser" @success="onEditSuccess" />
  </q-page>
</template>
```

### 2.3 AccountListDynamicTable.vue

```vue
<template>
  <q-table
    :rows="users"
    :columns="columns"
    :loading="loading"
    :pagination="pagination"
    @request="onRequest"
  >
    <template #body-cell-role="props">
      <q-td :props="props">
        <q-badge :color="getRoleColor(props.row.role)">
          {{ getRoleLabel(props.row.role) }}
        </q-badge>
      </q-td>
    </template>
    
    <template #body-cell-status="props">
      <q-td :props="props">
        <q-badge :color="props.row.status === 'ACTIVE' ? 'green' : 'grey'">
          {{ props.row.status === 'ACTIVE' ? '활성' : '비활성' }}
        </q-badge>
      </q-td>
    </template>
    
    <template #body-cell-actions="props">
      <q-td :props="props">
        <q-btn flat icon="edit" @click="$emit('edit', props.row)" />
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
const columns = [
  { name: 'userId', label: '아이디', field: 'userId', sortable: true },
  { name: 'name', label: '이름', field: 'name', sortable: true },
  { name: 'email', label: '이메일', field: 'email' },
  { name: 'role', label: '역할', field: 'role' },
  { name: 'status', label: '상태', field: 'status' },
  { name: 'createdAt', label: '생성일', field: 'createdAt' },
  { name: 'actions', label: '작업', field: 'actions' }
]

const getRoleColor = (role: string) => {
  switch (role) {
    case 'SUPER_ADMIN': return 'red'
    case 'ADMIN': return 'orange'
    default: return 'blue'
  }
}

const getRoleLabel = (role: string) => {
  switch (role) {
    case 'SUPER_ADMIN': return '슈퍼 관리자'
    case 'ADMIN': return '관리자'
    default: return '일반'
  }
}
</script>
```

### 2.4 AccountCreateDialog.vue

```vue
<template>
  <q-dialog v-model="show">
    <q-card style="width: 400px">
      <q-card-section>
        <div class="text-h6">계정 생성</div>
      </q-card-section>
      
      <q-card-section>
        <q-form @submit="onSubmit">
          <q-input v-model="form.userId" label="아이디" :rules="[required]" />
          <q-input v-model="form.name" label="이름" :rules="[required]" />
          <q-input v-model="form.email" label="이메일" />
          <q-input v-model="form.password" label="비밀번호" type="password" :rules="[required, minLength(8)]" />
          <q-select v-model="form.role" :options="roleOptions" label="역할" emit-value map-options />
          <q-select v-model="form.status" :options="statusOptions" label="상태" emit-value map-options />
          
          <div class="q-mt-md">
            <q-btn type="submit" color="primary" label="생성" :loading="loading" />
            <q-btn flat label="취소" v-close-popup />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>
```

---

## 3. 개발 테스트

### 3.1 DevTest.vue

```vue
<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">개발 테스트</div>
    
    <DevTestComponent />
  </q-page>
</template>
```

### 3.2 DevTestComponent.vue

```vue
<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">테스트 데이터 생성</div>
      
      <q-form @submit="onSubmit">
        <q-select
          v-model="form.type"
          :options="typeOptions"
          label="생성 타입"
          emit-value
          map-options
        />
        
        <template v-if="form.type === 'reservation' || form.type === 'all'">
          <div class="text-subtitle2 q-mt-md">예약 옵션</div>
          <q-input v-model="form.reservationOptions.startDate" type="date" label="시작일" />
          <q-input v-model="form.reservationOptions.endDate" type="date" label="종료일" />
          <q-input v-model.number="form.reservationOptions.regularReservations" label="숙박 예약 수" type="number" />
          <q-input v-model.number="form.reservationOptions.monthlyReservations" label="월세 예약 수" type="number" />
        </template>
        
        <q-btn type="submit" color="primary" label="생성" :loading="loading" class="q-mt-md" />
      </q-form>
    </q-card-section>
    
    <q-card-section v-if="result">
      <div class="text-h6">결과</div>
      <pre>{{ JSON.stringify(result, null, 2) }}</pre>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { devTestApi } from 'src/services/dev-test'

const typeOptions = [
  { value: 'essential', label: '필수 데이터만' },
  { value: 'reservation', label: '예약 데이터만' },
  { value: 'all', label: '모든 데이터' }
]

const form = reactive({
  type: 'all',
  reservationOptions: {
    startDate: null,
    endDate: null,
    regularReservations: 50,
    monthlyReservations: 5
  }
})

const result = ref(null)
const loading = ref(false)

const onSubmit = async () => {
  loading.value = true
  try {
    const response = await devTestApi.generateTestData(form)
    result.value = response.data
  } finally {
    loading.value = false
  }
}
</script>
```

---

## 4. 디버그 페이지

### 4.1 DevDebug.vue

```vue
<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">디버그</div>
    
    <div class="row q-gutter-md">
      <ApiTestCard class="col" />
      <UserInfoCard class="col" />
    </div>
  </q-page>
</template>
```

### 4.2 ApiTestCard.vue

- API 호출 테스트
- 응답 시간 측정
- 에러 확인

### 4.3 UserInfoCard.vue

- 현재 로그인 사용자 정보
- 토큰 정보
- 권한 정보
