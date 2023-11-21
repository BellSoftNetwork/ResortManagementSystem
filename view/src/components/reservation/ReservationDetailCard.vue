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
          <div class="q-py-sm">
            <div class="text-caption">객실 번호</div>
            <div class="text-body1">
              <div v-if="entity.room">
                <q-btn
                  :to="{ name: 'Room', params: { id: entity.room.id } }"
                  color="primary"
                  flat
                  dense
                >
                  {{ entity.room.number }}
                </q-btn>
              </div>
              <div v-else class="text-grey">미배정</div>
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">누적 결제 금액</div>
            <div class="text-body1">{{ entity.paymentAmount }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 수단</div>
            <div class="text-body1">{{ entity.reservationMethod.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 수단 수수료</div>
            <div class="text-body1">{{ entity.brokerFee }}</div>
          </div>
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <div class="q-py-sm">
            <div class="text-caption">예약자명</div>
            <div class="text-body1">{{ entity.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약자 연락처</div>
            <div class="text-body1"><a :href="'tel:' + entity.phone">{{ entity.phone }}</a></div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약인원</div>
            <div class="text-body1">{{ entity.peopleCount }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 상태</div>
            <div class="text-body1">{{ statusMap[entity.status] }}</div>
          </div>
        </div>
      </div>
      <br />
      <div class="row">
        <div class="col">
          <div class="q-py-sm">
            <div class="text-caption">메모</div>
            <div class="text-body1">{{ entity.note }}</div>
          </div>
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
const statusMap = {
  "PENDING": "예약 대기",
  "NORMAL": "예약 확정",
  "CANCEL": "예약 취소",
  "REFUND": "환불 완료",
}
const stayDateDiff = computed(() => dayjs(entity.value.stayEndAt).diff(dayjs(entity.value.stayStartAt), "day"))

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
