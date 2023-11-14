<template>
  <q-card>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6">
      객실 수정
    </q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-input
          v-model="formData.number"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="rules.number"
          label="번호"
          placeholder="101호"
          required
        ></q-input>

        <q-input
          v-model="formData.peekPrice"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="rules.peekPrice"
          label="성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.offPeekPrice"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="rules.offPeekPrice"
          label="비성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.description"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="rules.description"
          type="textarea"
          label="설명"
          placeholder="와이파이 사용 가능"
        ></q-input>

        <q-input
          v-model="formData.note"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="rules.note"
          type="textarea"
          label="메모 (관리용)"
          placeholder="문고리 고장"
        ></q-input>

        <q-select
          v-model="formData.status"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :options="options.status"
          label="상태"
          required
          emit-value
          map-options
        ></q-select>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          :disable="status.isProgress"
          :to="{ name: 'Rooms' }"
          color="primary"
          label="취소"
          flat
        />
        <q-btn
          :loading="status.isProgress"
          type="submit"
          color="red"
          label="수정"
          flat
        />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup>
import { onBeforeMount, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const $q = useQuasar()
const dialog = ref({
  isOpen: false,
})
const status = ref({
  isProgress: false,
})
const entity = {
  id: route.params.id,
  number: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
  note: "",
  status: "NORMAL",
}
const formData = ref({
  number: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
  note: "",
  status: "NORMAL",
})
const rules = {
  number: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  peekPrice: [value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"],
  offPeekPrice: [value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"],
  description: [value => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다"],
  note: [value => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다"],
}
const options = {
  status: [
    { label: "정상", value: "NORMAL" },
    { label: "이용불가", value: "INACTIVE" },
    { label: "파손", value: "DAMAGED" },
    { label: "공사 중", value: "CONSTRUCTION" },
  ],
}

function fetchData() {
  status.value.isProgress = true

  return api.get(`/api/v1/rooms/${entity.id}`)
    .then(res => {
      entity.id = res.data.value.id
      entity.number = res.data.value.number
      entity.peekPrice = res.data.value.peekPrice
      entity.offPeekPrice = res.data.value.offPeekPrice
      entity.description = res.data.value.description
      entity.note = res.data.value.note
      entity.status = res.data.value.status
    })
    .catch(error => {
      if (error.response.status === 404)
        router.push({ name: "ErrorNotFound" })

      console.log(error)
    }).finally(() => {
      status.value.isProgress = false
    })
}


function update() {
  if (!isChanged()) {
    $q.notify({
      message: "수정된 항목이 없습니다.",
      type: "info",
      actions: [
        {
          icon: "close", color: "white", round: true,
        },
      ],
    })

    return
  }

  status.value.isProgress = true

  api.patch(`/api/v1/rooms/${entity.id}`, patchedData())
    .then(() => {
      router.push({ name: "Room", params: { id: entity.id } })
    })
    .catch((error) => {
      $q.notify({
        message: error.response.data.message,
        type: "negative",
        actions: [
          {
            icon: "close", color: "white", round: true,
          },
        ],
      })
    })
    .finally(() => {
      status.value.isProgress = false
    })
}

function isChanged() {
  const checkKeys = Object.keys(formData.value)

  for (const key of checkKeys) {
    if (formData.value[key] !== entity[key]) {
      return true
    }
  }

  return false
}

function patchedData() {
  const checkKeys = Object.keys(formData.value)
  const patchData = {}

  checkKeys.forEach((key) => {
    if (formData.value[key] !== entity[key]) {
      patchData[key] = formData.value[key]
    }
  })

  return patchData
}

function resetForm() {
  formData.value.name = entity.name
  formData.value.number = entity.number
  formData.value.peekPrice = entity.peekPrice
  formData.value.offPeekPrice = entity.offPeekPrice
  formData.value.description = entity.description
  formData.value.note = entity.note
  formData.value.status = entity.status
}

onBeforeMount(() => {
  fetchData().then(() => {
    resetForm()
  })
})
</script>
