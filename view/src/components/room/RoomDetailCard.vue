<template>
  <q-card flat bordered>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6">
      {{ entity.number }}
    </q-card-section>

    <q-card-section>
      <q-input
        v-model="entity.peekPrice"
        :loading="status.isProgress"
        :readonly="true"
        label="성수기 예약금"
        type="number"
      ></q-input>

      <q-input
        v-model="entity.offPeekPrice"
        :loading="status.isProgress"
        :readonly="true"
        label="비성수기 예약금"
        type="number"
      ></q-input>

      <q-input
        v-model="entity.description"
        :loading="status.isProgress"
        :readonly="true"
        type="textarea"
        label="설명"
      ></q-input>

      <q-input
        v-model="entity.note"
        :loading="status.isProgress"
        :readonly="true"
        type="textarea"
        label="메모 (관리용)"
      ></q-input>

      <q-select
        v-model="entity.status"
        :loading="status.isProgress"
        :readonly="true"
        :options="options.status"
        label="상태"
        emit-value
        map-options
      ></q-select>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn
        @click="deleteItem()"
        color="red"
        label="삭제"
        dense
        flat
      ></q-btn>
      <q-btn
        :disable="status.isProgress"
        :to="{ name: 'EditRoom', params: { id: entity.id } }"
        color="primary"
        label="수정"
        flat
      />
    </q-card-actions>
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
const props = defineProps({
  id: Number,
})
const dialog = ref({
  isOpen: false,
})
const status = ref({
  isProgress: false,
})
const entity = {
  id: props.id,
  number: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
  note: "",
  status: "NORMAL",
}
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

function deleteItem() {
  const itemId = entity.id
  const itemName = entity.number

  $q.dialog({
    title: "삭제",
    message: `정말로 '${itemName}'을 삭제하시겠습니까?`,
    ok: {
      label: "삭제",
      flat: true,
      color: "negative",
    },
    cancel: {
      label: "유지",
      flat: true,
    },
    focus: "cancel",
  }).onOk(() => {
    api.delete(`/api/v1/rooms/${itemId}`)
      .then(() => {
        router.push({ name: "Rooms" })
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
  })
}

onBeforeMount(() => {
  fetchData()
})
</script>
