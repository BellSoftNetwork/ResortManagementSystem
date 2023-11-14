<template>
  <q-card>
    <q-card-section class="text-h6">
      객실 추가
    </q-card-section>

    <q-form @submit="create">
      <q-card-section>
        <q-input
          v-model="formData.number"
          :rules="rules.number"
          label="번호"
          placeholder="101호"
          required
        ></q-input>

        <q-input
          v-model="formData.peekPrice"
          :rules="rules.peekPrice"
          label="성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.offPeekPrice"
          :rules="rules.offPeekPrice"
          label="비성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.description"
          :rules="rules.description"
          type="textarea"
          label="설명"
          placeholder="와이파이 사용 가능"
        ></q-input>

        <q-input
          v-model="formData.note"
          :rules="rules.note"
          type="textarea"
          label="메모 (관리용)"
          placeholder="문고리 고장"
        ></q-input>

        <q-select
          v-model="formData.status"
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
          label="추가"
          flat
        />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup>
import { onBeforeMount, ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const dialog = ref({
  isOpen: false,
})
const status = ref({
  isProgress: false,
})
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

function create() {
  status.value.isProgress = true

  const data = {
    number: formData.value.number,
    peekPrice: formData.value.peekPrice,
    offPeekPrice: formData.value.offPeekPrice,
    description: formData.value.description,
    note: formData.value.note,
    status: formData.value.status,
  }

  api.post("/api/v1/rooms", data)
    .then(() => {
      router.push({ name: "Rooms" })

      resetForm()
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

function resetForm() {
  formData.value.name = ""
  formData.value.number = ""
  formData.value.peekPrice = 0
  formData.value.offPeekPrice = 0
  formData.value.description = ""
  formData.value.staus = "NORMAL"
}

onBeforeMount(() => {
  resetForm()
})
</script>
