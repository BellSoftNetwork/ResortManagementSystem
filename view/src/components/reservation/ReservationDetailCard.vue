<template>
  <q-card flat bordered>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6">
      예약 정보
    </q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col-12 col-md-auto q-px-sm">
          <q-date
            :model-value="{from: entity.stayStartAt, to: entity.stayEndAt}"
            :title="formatSubTitle(stayDateDiff)"
            subtitle="숙박 기간"
            range
            mask="YYYY-MM-DD"
            :readonly="true"
          />
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <q-input
            v-model="entity.room.number"
            :readonly="true"
            label="객실 번호"
          />

          <q-input
            v-model="entity.price"
            :readonly="true"
            label="판매 금액"
            type="number"
          ></q-input>

          <q-input
            v-model="entity.reservationFee"
            :readonly="true"
            label="예약 선입금액"
          ></q-input>

          <q-input
            v-model="entity.paymentAmount"
            :readonly="true"
            label="누적 결제 금액"
            type="number"
          ></q-input>

          <q-select
            v-model="entity.reservationMethod.name"
            :readonly="true"
            label="예약 수단"
          ></q-select>

          <q-input
            v-model="entity.brokerFee"
            :readonly="true"
            label="예약 수단 수수료"
            type="number"
          ></q-input>
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <q-input
            v-model="entity.name"
            :readonly="true"
            label="예약자명"
          ></q-input>

          <q-input
            v-model="entity.phone"
            :readonly="true"
            label="예약자 연락처"
          ></q-input>

          <q-input
            v-model="entity.peopleCount"
            :readonly="true"
            label="예약인원"
          ></q-input>

          <q-select
            v-model="entity.status"
            :readonly="true"
            :options="options.status"
            label="예약 상태"
            emit-value
            map-options
          ></q-select>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <q-input
            v-model="entity.note"
            :readonly="true"
            type="textarea"
            label="메모"
          ></q-input>
        </div>
      </div>
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
        :to="{ name: 'EditReservation', params: { id: entity.id } }"
        color="primary"
        label="수정"
        flat
      />
    </q-card-actions>
  </q-card>
</template>

<script setup>
import { computed, onBeforeMount, ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"
import dayjs from "dayjs"

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
const entity = ref({
  id: props.id,
  reservationMethod: {
    id: 0,
    name: "",
  },
  room: {
    id: 0,
    number: "",
  },
  name: "",
  phone: "",
  peopleCount: 0,
  stayStartAt: "",
  stayEndAt: "",
  checkInAt: "",
  checkOutAt: "",
  price: "",
  paymentAmount: "",
  refundAmount: "",
  reservationFee: "",
  brokerFee: "",
  note: "",
  canceledAt: "",
  status: "NORMAL",
  createdAt: "",
  createdBy: "",
  updatedAt: "",
  updatedBy: "",
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 30) || "2~30 글자가 필요합니다"],
  phone: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  peopleCount: [value => (value >= 0 && value <= 1000) || "1000 명 이하만 입실 가능합니다"],
  stayStartAt: [value => (/^-?[\d]+-[0-1]\d-[0-3]\d$/.test(value)) || "####-##-## 형태의 날짜만 입력 가능합니다."],
  stayEndAt: [value => (/^-?[\d]+-[0-1]\d-[0-3]\d$/.test(value)) || "####-##-## 형태의 날짜만 입력 가능합니다."],
  price: [value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"],
  paymentAmount: [
    value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다",
    value => (value <= entity.value.price) || "판매 금액보다 클 수 없습니다",
    value => (value >= entity.value.reservationFee) || "선입금액보다 작을 수 없습니다",
  ],
  reservationFee: [
    value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다",
    value => (value <= entity.value.price) || "판매 금액보다 클 수 없습니다",
  ],
  brokerFee: [value => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"],
  note: [value => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다"],
}
const options = {
  status: [
    { label: "예약 대기", value: "PENDING" },
    { label: "예약 확정", value: "NORMAL" },
    { label: "예약 취소", value: "CANCEL" },
    { label: "환불 완료", value: "REFUND" },
  ],
}
const stayDateDiff = computed(() => dayjs(entity.value.stayEndAt).diff(dayjs(entity.value.stayStartAt), "day"))
const reservationMethods = ref({
  status: {
    isLoading: false,
    isLoaded: false,
  },
  values: [
    {
      id: -1,
      name: "네이버",
      commissionRate: 0.1,
      createdAt: "2021-01-01T00:00:00.000Z",
      updatedAt: "2021-01-01T00:00:00.000Z",
    },
  ],
})


function fetchData() {
  status.value.isProgress = true

  return api.get(`/api/v1/reservations/${entity.value.id}`)
    .then(response => {
      entity.value = response.data.value
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
  const itemId = entity.value.id
  const itemName = entity.value.name

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 ${formatSubTitle(stayDateDiff.value)} 예약 정보를 삭제하시겠습니까?`,
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
    api.delete(`/api/v1/reservations/${itemId}`)
      .then(() => {
        router.push({ name: "ReservationList" })
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

function formatSubTitle(dateDiff) {
  return `${dateDiff}박 ${dateDiff + 1}일`
}

onBeforeMount(() => {
  fetchData()
})
</script>
