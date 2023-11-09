<template>
  <v-dialog
    v-model="status.isDialogActive"
    :persistent="status.isProgress"
    width="500"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        variant="text"
      >
        <slot>예약 수단 수정</slot>
      </v-btn>
    </template>

    <template v-slot:default>
      <v-card
        title="예약 수단 수정"
        :loading="status.isProgress"
        :disabled="status.isProgress"
      >
        <v-card-text>
          <v-form
            v-model="status.isValid"
            @submit.prevent
            fast-fail
            ref="form"
          >
            <v-text-field
              v-model="reservationMethod.name"
              :rules="rules.name"
              label="이름"
            ></v-text-field>

            <v-text-field
              v-model="reservationMethod.commissionRatePercent"
              label="수수료율"
              :rules="rules.commissionRatePercent"
              append-inner-icon="fa-solid fa-percent"
              type="number"
              min="0"
              max="100"
              required
            ></v-text-field>
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            text="취소"
            @click="status.isDialogActive = false"
          ></v-btn>

          <v-btn
            text="수정"
            color="primary"
            @click="updateReservationMethod"
            :disabled="!(status.isValid && isChanged())"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

  <v-snackbar
    v-model="status.isError"
  >
    예약 수단 수정 실패 ({{ status.errorMessage }})

    <template v-slot:actions>
      <v-btn
        color="pink"
        variant="text"
        @click="status.isError = false"
      >
        닫기
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/store/auth.js"
import axios from "@/modules/axios-wrapper"

const router = useRouter()
const authStore = useAuthStore()

const emit = defineEmits(["complete"])
const props = defineProps({
  reservationMethod: Object,
})
const status = ref({
  isDialogActive: false,
  isValid: false,
  isProgress: false,
  isError: false,
  errorMessage: null,
})
const reservationMethod = ref({
  name: props.reservationMethod.name,
  commissionRatePercent: props.reservationMethod.commissionRate * 100,
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  commissionRatePercent: [value => (value >= 0 && value <= 100) || "수수료율이 유효하지 않습니다."],
}


function updateReservationMethod() {
  status.value.isProgress = true

  axios.patch(`/api/v1/reservation-methods/${props.reservationMethod.id}`, patchedData())
    .then(() => {
      emit("complete")
      status.value.isDialogActive = false
    })
    .catch((error) => {
      status.value.errorMessage = error.response.data.message
      status.value.isError = true
    })
    .finally(() => {
      status.value.isProgress = false
    })
}

function isChanged() {
  return reservationMethod.value.name !== props.reservationMethod.name ||
    reservationMethod.value.commissionRatePercent / 100 !== props.reservationMethod.commissionRate
}

function patchedData() {
  const patchData = {}

  if (props.reservationMethod.name !== reservationMethod.value.name)
    patchData.name = reservationMethod.value.name
  if (props.reservationMethod.commissionRate !== reservationMethod.value.commissionRatePercent / 100)
    patchData.commissionRate = reservationMethod.value.commissionRatePercent / 100

  return patchData
}
</script>
