<template>
  <q-stepper
    v-model="formModel.status.step"
    ref="stepper"
    color="primary"
    animated
    vertical
  >
    <q-step
      :name="1"
      :caption="stayDateDiff <= 0 ? null : `${formModel.value.stayDate.from} ~ ${formModel.value.stayDate.to}`"
      :error="stayDateDiff <= 0"
      :done="formModel.status.step > 1"
      title="숙박 기간"
      icon="settings"
    >
      <q-date
        v-model="formModel.value.stayDate"
        :color="stayDateDiff > 0 ? 'primary' : 'red'"
        :title="formatSubTitle(stayDateDiff)"
        subtitle="숙박 기간"
        range
        mask="YYYY-MM-DD"
      />

      <q-stepper-navigation>
        <q-btn
          @click="$refs.stepper.next()"
          label="다음"
          color="primary"
        />
      </q-stepper-navigation>
    </q-step>

    <q-step
      :name="2"
      :done="formModel.status.step > 2"
      title="객실 배정"
      :caption="(selectedRoom[0] && Object.keys(selectedRoom[0]).includes('number')) ? selectedRoom[0].number : '추후 배정'"
      icon="create_new_folder"
    >
      <RoomSelectTable
        v-model:selected="selectedRoom"
        :first-value="entity.room"
        :stay-start-at="formModel.value.stayDate.from"
        :stay-end-at="formModel.value.stayDate.to"
      />

      <q-stepper-navigation>
        <q-btn
          @click="$refs.stepper.next()"
          label="다음"
          color="primary"
        />
        <q-btn
          @click="$refs.stepper.previous()"
          color="primary"
          label="이전"
          class="q-ml-sm"
          flat
        />
      </q-stepper-navigation>
    </q-step>

    <q-step
      :name="3"
      :done="formModel.status.step > 3"
      title="금액 확인"
      :caption="formatPrice(formModel.value.price)"
      icon="create_new_folder"
    >
      <q-form @submit="formModel.status.step = 4">
        <q-input
          v-model="formModel.value.price"
          :rules="rules.price"
          @update:model-value="changePrice()"
          label="판매 금액"
          placeholder="100000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-input
          v-model="formModel.value.paymentAmount"
          :rules="rules.paymentAmount"
          label="누적 결제 금액"
          placeholder="80000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-select
          v-model="formModel.value.reservationMethod"
          @update:model-value="changePrice()"
          :loading="reservationMethods.status.isLoading"
          :disable="!reservationMethods.status.isLoaded"
          :options="reservationMethods.values"
          option-label="name"
          label="예약 수단"
          required
          map-options
        ></q-select>

        <q-input
          v-model="formModel.value.brokerFee"
          :rules="rules.brokerFee"
          :readonly="true"
          label="예약 수단 수수료"
          placeholder="5000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-stepper-navigation>
          <q-btn
            type="submit"
            label="다음"
            color="primary"
          />
          <q-btn
            @click="$refs.stepper.previous()"
            color="primary"
            label="이전"
            class="q-ml-sm"
            flat
          />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step
      :name="4"
      :done="formModel.status.step > 4"
      title="예약자 정보"
      :caption="`${formModel.value.name} (${formModel.value.phone}) / ${formModel.value.peopleCount}명`"
      icon="create_new_folder"
    >
      <q-form @submit="formModel.status.step = 5">
        <q-input
          v-model="formModel.value.name"
          :rules="rules.name"
          label="예약자명"
          placeholder="홍길동"
          required
        ></q-input>

        <q-input
          v-model="formModel.value.phone"
          :rules="rules.phone"
          label="예약자 연락처"
          placeholder="010-0000-0000"
          required
        ></q-input>

        <q-input
          v-model="formModel.value.peopleCount"
          :rules="rules.peopleCount"
          label="예약인원"
          placeholder="4"
          type="number"
          min="0"
          max="1000"
          required
        ></q-input>

        <q-select
          v-model="formModel.value.status"
          :options="options.status"
          label="예약 상태"
          required
          emit-value
          map-options
        ></q-select>

        <q-input
          v-model="formModel.value.note"
          :rules="rules.note"
          type="textarea"
          label="메모"
          placeholder="밤 늦게 입실 예정"
        ></q-input>

        <q-stepper-navigation>
          <q-btn
            label="다음"
            type="submit"
            color="primary"
          />
          <q-btn
            @click="$refs.stepper.previous()"
            color="primary"
            label="이전"
            class="q-ml-sm"
            flat
          />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step
      :name="5"
      title="수정"
      icon="add_comment"
    >
      <div>{{ Object.keys(patchedData()).length }}개 항목이 변경되었습니다.</div>

      <q-stepper-navigation>
        <q-btn
          @click="update"
          label="수정"
          color="primary"
        />
        <q-btn
          @click="$refs.stepper.previous()"
          color="primary"
          label="이전"
          class="q-ml-sm"
          flat
        />
      </q-stepper-navigation>
    </q-step>
  </q-stepper>
</template>

<script setup>
import { computed, onBeforeMount, ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"
import dayjs from "dayjs"
import RoomSelectTable from "components/room/RoomSelectTable.vue"

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const props = defineProps({
  id: Number,
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
  peopleCount: 4,
  stayStartAt: "",
  stayEndAt: "",
  checkInAt: "",
  checkOutAt: "",
  price: "",
  paymentAmount: "",
  refundAmount: "",
  brokerFee: "",
  note: "",
  canceledAt: "",
  status: "NORMAL",
  createdAt: "",
  createdBy: "",
  updatedAt: "",
  updatedBy: "",
})
const formModel = ref({
  status: {
    step: 1,
    isProgress: false,
  },
  value: {
    id: props.id,
    reservationMethod: {
      id: -1,
      name: "네이버",
      commissionRate: 0.1,
      createdAt: "2021-01-01T00:00:00.000Z",
      updatedAt: "2021-01-01T00:00:00.000Z",
    },
    name: "",
    phone: "",
    peopleCount: 4,
    stayDate: { from: "", to: "" },
    price: 0,
    paymentAmount: 0,
    brokerFee: 0,
    note: "",
    status: "PENDING",
  },
})
const selectedRoom = ref([])
const status = ref({
  isProgress: false,
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
    value => (value <= formModel.value.value.price) || "판매 금액보다 클 수 없습니다",
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
const stayDateDiff = computed(() => {
    try {
      return dayjs(formModel.value.value.stayDate.to).diff(dayjs(formModel.value.value.stayDate.from), "day")
    } catch (e) {
      return 0
    }
  },
)
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
    .then(res => {
      entity.value = res.data.value
    })
    .catch(error => {
      if (error.response.status === 404)
        router.push({ name: "ErrorNotFound" })

      console.log(error)
    }).finally(() => {
      status.value.isProgress = false
    })
}

function loadReservationMethods() {
  reservationMethods.value.status.isLoading = true
  reservationMethods.value.status.isLoaded = false
  reservationMethods.value.values = []

  return api.get(`/api/v1/reservation-methods?sort=name`)
    .then(response => {
      const values = response.data.values

      reservationMethods.value.values = values

      reservationMethods.value.status.isLoaded = true
    }).finally(() => {
      reservationMethods.value.status.isLoading = false
    })
}

function update() {
  status.value.isProgress = true

  api.patch(`/api/v1/reservations/${entity.value.id}`, patchedData())
    .then(() => {
      router.push({ name: "Reservation", params: { id: entity.value.id } })

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

function patchedData() {
  const patchData = {}

  if (entity.value.reservationMethod.id !== formModel.value.value.reservationMethod.id)
    patchData.reservationMethodId = formModel.value.value.reservationMethod.id
  const roomId = (selectedRoom.value[0] && Object.keys(selectedRoom.value[0]).includes("id")) ? selectedRoom.value[0].id : null
  if (entity.value.room && entity.value.room.id !== roomId)
    patchData.roomId = roomId
  if (entity.value.name !== formModel.value.value.name)
    patchData.name = formModel.value.value.name
  if (entity.value.phone !== formModel.value.value.phone)
    patchData.phone = formModel.value.value.phone
  if (entity.value.peopleCount !== formModel.value.value.peopleCount)
    patchData.peopleCount = formModel.value.value.peopleCount
  if (dayjs(entity.value.stayStartAt).format("YYYY-MM-DD") !== dayjs(formModel.value.value.stayDate.from).format("YYYY-MM-DD"))
    patchData.stayStartAt = formModel.value.value.stayDate.from
  if (dayjs(entity.value.stayEndAt).format("YYYY-MM-DD") !== dayjs(formModel.value.value.stayDate.to).format("YYYY-MM-DD"))
    patchData.stayEndAt = formModel.value.value.stayDate.to
  if (entity.value.price !== formModel.value.value.price)
    patchData.price = formModel.value.value.price
  if (entity.value.paymentAmount !== formModel.value.value.paymentAmount)
    patchData.paymentAmount = formModel.value.value.paymentAmount
  if (entity.value.brokerFee !== formModel.value.value.brokerFee)
    patchData.brokerFee = formModel.value.value.brokerFee
  if (entity.value.note !== formModel.value.value.note)
    patchData.note = formModel.value.value.note
  if (entity.value.status !== formModel.value.value.status)
    patchData.status = formModel.value.value.status

  return patchData
}

function formatSubTitle(dateDiff) {
  return `${dateDiff}박 ${dateDiff + 1}일`
}

function formatPrice(value) {
  return new Intl.NumberFormat("ko-KR", {
    style: "currency",
    currency: "KRW",
  }).format(value)
}

function changePrice() {
  formModel.value.value.brokerFee = formModel.value.value.price * formModel.value.value.reservationMethod.commissionRate
}

function resetForm() {
  formModel.value.value.name = entity.value.name
  formModel.value.value.phone = entity.value.phone
  formModel.value.value.peopleCount = entity.value.peopleCount
  formModel.value.value.stayDate.from = dayjs(entity.value.stayStartAt).format("YYYY-MM-DD")
  formModel.value.value.stayDate.to = dayjs(entity.value.stayEndAt).format("YYYY-MM-DD")
  formModel.value.value.price = entity.value.price
  formModel.value.value.paymentAmount = entity.value.paymentAmount
  formModel.value.value.brokerFee = entity.value.brokerFee
  formModel.value.value.note = entity.value.note
  formModel.value.value.status = entity.value.status
}

onBeforeMount(() => {
  fetchData().then(() => {
    resetForm()
    if (entity.value.room)
      selectedRoom.value = [entity.value.room]
    loadReservationMethods().then(() => {
      formModel.value.value.reservationMethod = reservationMethods.value.values.find(item => item.id === entity.value.reservationMethod.id)
    })
  })
})
</script>
